# go-repair-center

Household appliance repair center management system built with Go REST API and Vue 3.

## Local development admin

- Username: `admin`
- Password: `Admin123!`
- This account is for local development acceptance only. Change it immediately in production.

## Start with Docker

```bash
docker compose -f docker-compose.yml up --build
```

The frontend is exposed at `http://localhost:8080`.

## Additional features

- CSV import preview and apply for customers, parts and devices
- CSV export for repair orders, parts, feedbacks and audit logs
- Session management for current and target users
- Operations center for backlog, SLA alerts, warranty reminders and health snapshot
- Dispatch center for batch assignment, time conflict checks and rework-chain lookup
- Approval center for quotation and warranty review queues
- Inventory control for low-stock monitoring, restock suggestions and stock timelines
