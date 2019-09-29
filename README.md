# SQL Migration

This repository contains an example of how to use the [`golang-migrate`] to automatically migrate
a SQL database to a desired schema.

```bash
# Start a Postgres database that our example will connect to.
docker run --name postgres-migration -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres

# Run the application.
go run .

# Clean up the database.
docker stop postgres-migration && docker rm postgres-migration
```

[`golang-migrate`]: https://github.com/golang-migrate/migrate
