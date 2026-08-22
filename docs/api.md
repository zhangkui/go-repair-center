# API Overview

Base path: `/api/v1`

## Authentication

- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `GET /me`
- `POST /auth/password`

## Dashboard and reports

- `GET /dashboard`
- `GET /reports/revenue`
- `GET /reports/revenue.csv`
- `GET /ops/backlog`
- `GET /ops/sla-alerts`
- `GET /ops/warranty-expiring`
- `GET /ops/health`
- `GET /sessions/me`
- `GET /sessions/users/{userID}`
- `POST /sessions/{id}/revoke`
- `POST /sessions/revoke-all`
- `POST /sessions/cleanup`
- `POST /dispatch/conflicts`
- `POST /dispatch/batch`
- `GET /dispatch/rework-chain/{id}`
- `GET /approvals/overview`
- `GET /approvals/quotations`
- `POST /approvals/quotations/{id}/decision`
- `GET /approvals/warranties`
- `POST /approvals/warranties/{id}/decision`
- `GET /inventory/summary`
- `GET /inventory/low-stock`
- `GET /inventory/restock-suggestions`
- `GET /inventory/parts/{id}/timeline`
- `POST /inventory/restock`
- `GET /tools/import/templates/customers`
- `GET /tools/import/templates/parts`
- `GET /tools/import/templates/devices`
- `POST /tools/import/customers/preview`
- `POST /tools/import/customers/apply`
- `POST /tools/import/parts/preview`
- `POST /tools/import/parts/apply`
- `POST /tools/import/devices/preview`
- `POST /tools/import/devices/apply`
- `GET /tools/export/repair-orders.csv`
- `GET /tools/export/parts.csv`
- `GET /tools/export/feedbacks.csv`
- `GET /tools/export/audit-logs.csv`

## Core resources

- `GET|POST /users`
- `GET|PUT /users/{id}`
- `PUT /users/{id}/enable`
- `PUT /users/{id}/disable`
- `POST /users/{id}/reset-password`
- `GET|POST /roles`
- `GET|POST /permissions`
- `GET|POST /customers`
- `GET|POST /devices`
- `GET|POST /repair-orders`
- `GET|POST /quotations`
- `GET|POST /repair-executions`
- `GET|POST /parts`
- `GET|POST /warranties`
- `GET|POST /feedbacks`
- `GET|POST /audit-logs`
- `GET|POST /system-configs`
- `GET|POST /fault-codes`
