# 07. Web Frontend Plan

## Стек

Рекомендуется:

- React;
- Vite;
- TypeScript;
- React Router;
- TanStack Query;
- Chart.js/Recharts;
- Tailwind или Bootstrap.

Для очень быстрого MVP можно Go templates + Bootstrap + Chart.js.

## Минимум для курсовой

### `/admin/tables`

- выбрать таблицу;
- показать данные;
- pagination/limit;
- обработка ошибок.

### `/admin/functions`

- выбрать SQL-функцию;
- ввести параметры;
- выполнить;
- показать результат таблицей;
- показать ошибку неверных параметров.

### `/admin/analytics`

Графики:

1. оборот магазинов;
2. статусы транзакций;
3. динамика платежей;
4. ошибки по сетям/RPC.

### `/admin/demo`

Demo-flow:

1. создать invoice;
2. выбрать пользователя и кошелек;
3. отправить транзакцию;
4. показать статус;
5. показать запись в таблице и графиках.

## Личные кабинеты

### User

- профиль;
- KYC;
- кошельки;
- балансы;
- история;
- оплата invoice.

### Merchant

- профиль магазина;
- verification;
- кошельки;
- invoices;
- terminals;
- analytics.

### Compliance

- KYC queue;
- merchant verification queue;
- risk alerts;
- blacklist.

### Admin

- tables;
- functions;
- analytics;
- audit;
- rpc nodes.

## Компоненты

- `Layout`;
- `Sidebar`;
- `DataTable`;
- `FunctionRunner`;
- `DateRangePicker`;
- `StatusBadge`;
- `ErrorAlert`;
- `ChartCard`;
- `DemoFlow`.

## Обязательная валидация

- required fields;
- date format `YYYY-MM-DD`;
- amount > 0;
- понятные ошибки API;
- disabled button while loading.
