#!/bin/bash
set -euo pipefail

# -----------------------------
# Настройка переменных окружения
# -----------------------------
export VAULT_ADDR="http://vault.vault.svc.cluster.local:8200"
export KONG_ADMIN_URL="http://kong-admin.kong.svc.cluster.local:8001"

# Имена компонентов
CONSUMER_NAME="example-consumer"
SERVICE_USER="user-service"
SERVICE_TRANSACTION="transaction-service"
SERVICE_WALLET="wallet-service"

# Путь в Vault для хранения JWT ключей
JWT_SECRET_PATH="secret/jwt"

# Пути к файлам с ключами (убедитесь, что файлы находятся в той же директории, где запускается скрипт)
JWT_PUBLIC_FILE="./jwt_public_key.pem"
JWT_PRIVATE_FILE="./jwt_private_key.pem"

# -----------------------------
# 1. Автоматизация для Vault
# -----------------------------
echo "Проверка наличия JWT ключей в Vault по пути: ${JWT_SECRET_PATH}"
if ! vault kv get -field=public_key "$JWT_SECRET_PATH" >/dev/null 2>&1; then
  echo "JWT ключи не найдены. Загружаем их в Vault..."
  vault kv put "$JWT_SECRET_PATH" public_key="$(cat "$JWT_PUBLIC_FILE")" private_key="$(cat "$JWT_PRIVATE_FILE")"
  echo "Ключи успешно загружены."
else
  echo "JWT ключи уже присутствуют в Vault."
fi

echo "Получение публичного ключа из Vault..."
JWT_PUBLIC_KEY=$(vault kv get -field=public_key "$JWT_SECRET_PATH")
if [ -z "$JWT_PUBLIC_KEY" ]; then
  echo "Ошибка: публичный ключ JWT не найден!"
  exit 1
fi

# -----------------------------
# 2. Конфигурация Kong через Admin API
# -----------------------------

# Функция для отправки запросов и базовой обработки ошибок
function kong_post() {
  local url="$1"
  shift
  local data="$@"
  curl -s -o /dev/null -w "%{http_code}" -X POST "$url" $data
}

# 2.1. Создание потребителя (consumer)
echo "Создаём consumer '$CONSUMER_NAME' в Kong..."
HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/consumers" "--data username=$CONSUMER_NAME")
if [[ "$HTTP_CODE" != "201" && "$HTTP_CODE" != "409" ]]; then
  echo "Ошибка при создании consumer, HTTP-код: $HTTP_CODE"
  exit 1
fi

# 2.2. Конфигурация плагина JWT для потребителя
echo "Настраиваем плагин JWT для consumer '$CONSUMER_NAME'..."
HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/consumers/$CONSUMER_NAME/jwt" "--data key=$JWT_PUBLIC_KEY --data algorithm=RS256")
if [[ "$HTTP_CODE" != "201" && "$HTTP_CODE" != "409" ]]; then
  echo "Ошибка при настройке JWT плагина, HTTP-код: $HTTP_CODE"
  exit 1
fi

# 2.3. Настройка Rate Limiting для user-service
echo "Настраиваем плагин Rate Limiting для сервиса '$SERVICE_USER'..."
HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/services/$SERVICE_USER/plugins" "--data name=rate-limiting --data config.minute=1000")
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "Ошибка при настройке Rate Limiting для $SERVICE_USER, HTTP-код: $HTTP_CODE"
  exit 1
fi

# 2.4. Настройка Rate Limiting для transaction-service
echo "Настраиваем плагин Rate Limiting для сервиса '$SERVICE_TRANSACTION'..."
HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/services/$SERVICE_TRANSACTION/plugins" "--data name=rate-limiting --data config.minute=5000")
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "Ошибка при настройке Rate Limiting для $SERVICE_TRANSACTION, HTTP-код: $HTTP_CODE"
  exit 1
fi

# 2.5. Настройка плагина кэширования для wallet-service
# Примечание: здесь приведён пример использования плагина, который может обращаться к Redis.
echo "Настраиваем плагин кэширования для сервиса '$SERVICE_WALLET'..."
HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/services/$SERVICE_WALLET/plugins" "--data name=response-caching --data config.redis_host=redis.database.svc.cluster.local --data config.redis_port=6379 --data config.ttl=300")
if [[ "$HTTP_CODE" != "201" ]]; then
  echo "Ошибка при настройке кэширования для $SERVICE_WALLET, HTTP-код: $HTTP_CODE"
  exit 1
fi

echo "Автоматизация для Vault и конфигурации Kong завершена успешно."
