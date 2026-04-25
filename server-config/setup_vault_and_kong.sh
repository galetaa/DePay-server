#!/bin/bash
set -euo pipefail

# -----------------------------
# Настройка переменных окружения
# -----------------------------
export VAULT_ADDR="${VAULT_ADDR:-http://vault.vault.svc.cluster.local:8200}"
export KONG_ADMIN_URL="${KONG_ADMIN_URL:-http://kong-admin.kong.svc.cluster.local:8001}"

# Имена компонентов
CONSUMER_NAME="${KONG_CONSUMER_NAME:-depay-platform}"
KONG_SERVICES=(
  "user-service:1000"
  "merchant-service:1000"
  "wallet-service:2000"
  "transaction-core-service:5000"
  "transaction-validation-service:5000"
  "gas-info-service:2000"
  "kyc-service:1000"
  "admin-service:500"
)

# Путь в Vault для хранения JWT ключей
JWT_SECRET_PATH="${JWT_SECRET_PATH:-secret/jwt}"
APP_SECRET_PATH="${APP_SECRET_PATH:-secret/depay/app}"

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

echo "Проверка прикладных секретов DePay в Vault по пути: ${APP_SECRET_PATH}"
if ! vault kv get "$APP_SECRET_PATH" >/dev/null 2>&1; then
  vault kv put "$APP_SECRET_PATH" \
    database_url="${DATABASE_URL:-}" \
    jwt_secret="${JWT_SECRET:-}" \
    rabbitmq_url="${RABBITMQ_URL:-}" \
    blockchain_rpc_url="${BLOCKCHAIN_RPC_URL:-}" \
    kyc_provider_url="${KYC_PROVIDER_URL:-}" \
    kyc_provider_api_key="${KYC_PROVIDER_API_KEY:-}"
  echo "Прикладные секреты загружены."
else
  echo "Прикладные секреты уже присутствуют."
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

# 2.3. Настройка Rate Limiting для всех backend-сервисов
for service_spec in "${KONG_SERVICES[@]}"; do
  service="${service_spec%%:*}"
  limit="${service_spec##*:}"
  echo "Настраиваем Rate Limiting для сервиса '$service'..."
  HTTP_CODE=$(kong_post "$KONG_ADMIN_URL/services/$service/plugins" "--data name=rate-limiting --data config.minute=$limit --data config.policy=local")
  if [[ "$HTTP_CODE" != "201" && "$HTTP_CODE" != "409" ]]; then
    echo "Ошибка при настройке Rate Limiting для $service, HTTP-код: $HTTP_CODE"
    exit 1
  fi
done

echo "Автоматизация для Vault и конфигурации Kong завершена успешно."
