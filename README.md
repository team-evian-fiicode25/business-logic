# Fiicode25 Business Logic
Core logic library of the [RideMe](https://github.com/team-evian-fiicode25) project.

## Database schema
![diagram](https://team-evian-fiicode25.github.io/business-logic/DB-Schema.svg)

## Migrating database
```bash
POSTGRES_CONNECTION='<CONNECTION_STRING>' go run github.com/team-evian-fiicode25/business-logic/cmd/migrate@latest
```
