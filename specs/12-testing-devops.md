# 12. Testing and DevOps

## Backend tests

Покрыть:

- services;
- repositories;
- validation logic;
- auth;
- risk rules;
- HTTP handlers.

## SQL tests

Создать:

```text
database/tests/test_functions.sql
database/tests/test_triggers.sql
```

Проверить:

- функции возвращают строки;
- wallet owner trigger;
- negative balance trigger;
- status flow trigger;
- audit trigger.

## Frontend tests

MVP:

- ручной QA;
- screenshots;
- smoke tests.

Pet project:

- Vitest;
- React Testing Library;
- Playwright demo flow.

## Docker Compose

Сервисы:

- postgres;
- redis;
- rabbitmq;
- backend;
- web.

## CI

Этапы:

1. Go tests.
2. Go vet/gofmt.
3. Frontend build.
4. Migration check.
5. Optional Docker build.

## Defense readiness checklist

- DBeaver подключается.
- Схема открывается.
- Web запускается.
- Таблицы показываются.
- Функции выполняются.
- Графики строятся.
- Demo payment работает.
- Триггерная ошибка показывается.
- Скриншоты готовы.
