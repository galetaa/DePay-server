# 15. Coding Standards

## Go

Слои:

```text
controller -> service -> repository -> database
```

Правила:

- `context.Context` первым параметром;
- `gofmt` перед commit;
- не возвращать сырые ошибки БД пользователю;
- in-memory repositories только для тестов;
- бизнес-логика не должна быть в controller.

## SQL

Naming:

- tables: plural snake_case;
- columns: snake_case;
- indexes: `idx_table_column`;
- triggers: `trg_table_purpose`;
- functions: `get_report_name`.

Типы:

- crypto/деньги: `NUMERIC(30,8)`;
- timestamps: `TIMESTAMPTZ`;
- audit payloads: `JSONB`;
- email: `VARCHAR(255)`.

## API

Все через `/api`.

Ошибка:

```json
{
  "error": {
    "code": "INVALID_INPUT",
    "message": "..."
  }
}
```

## Frontend

- components: PascalCase;
- hooks: useSomething;
- server state через TanStack Query;
- формы показывают ошибки;
- не использовать Redux на MVP.

## Git

Branches:

- `feature/db-migrations`;
- `feature/web-admin-functions`;
- `fix/user-postgres-repo`;
- `coursework/sql-triggers`.

Commits:

- `feat(db): add payment transaction triggers`;
- `fix(user): replace in-memory repository`;
- `feat(web): add SQL function runner`.

## Definition of Done

Фича готова, если:

- код написан;
- тесты проходят;
- миграции есть;
- OpenAPI обновлен;
- ошибки обработаны;
- можно показать в demo.
