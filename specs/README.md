# DePay Dev Pack

Набор markdown-документов для продолжения `galetaa/DePay-server` как монорепозитория: Go backend, PostgreSQL, Redis, RabbitMQ, веб-интерфейс и SQL-часть для курсовой.

## Цель

Не переписывать текущий API-прототип с нуля, а довести его до сильного pet project и одновременно закрыть требования курсовой по БД.

## Главные направления

1. Подключить PostgreSQL и миграции.
2. Расширить предметную область: магазины, invoices, NFC-сессии, KYC, риск-контроль, RPC-мониторинг.
3. Заменить in-memory repositories на PostgreSQL repositories.
4. Добавить веб-интерфейс: таблицы, SQL-функции, графики, demo-flow.
5. Реализовать 12 сложных SQL-функций и минимум 5 разных триггеров.
6. Подготовить пояснительную записку и сценарий защиты.

## Файлы

- `00-current-state-audit.md` — что уже есть и что нужно исправить.
- `01-product-scope.md` — продуктовая концепция и объем.
- `02-architecture.md` — целевая архитектура.
- `03-repository-structure.md` — структура монорепо.
- `04-database-design.md` — расширенная БД.
- `05-coursework-sql-pack.md` — функции, процедуры, триггеры.
- `06-backend-plan.md` — backend roadmap.
- `07-web-plan.md` — веб-интерфейс.
- `08-api-contracts.md` — REST API.
- `09-security-compliance.md` — безопасность, KYC, аудит.
- `10-roadmap.md` — этапы разработки.
- `11-feature-backlog.md` — backlog по фичам.
- `12-testing-devops.md` — тестирование, Docker, CI/CD.
- `13-coursework-map.md` — соответствие требованиям курсовой.
- `14-defense-demo-script.md` — сценарий защиты.
- `15-coding-standards.md` — правила разработки.
