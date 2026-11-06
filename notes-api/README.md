# Notes and Bookmarks API

A Go REST API for managing notes, bookmarks, and tags. It uses PostgreSQL, sqlc-generated queries, Chi routing, request IDs, request logging, and graceful HTTP shutdown.

## Run locally with Docker Compose

Requirements: Go 1.27.1 or later and Docker Compose.

From this directory, create your local environment file and start the services:

```powershell
Copy-Item .env.example .env
docker compose up --build -d
docker compose ps
```

The API listens on `http://localhost:8080`; PostgreSQL is available locally on port `5433`. On the first API startup, Goose applies the embedded SQL migrations. The database data is stored in the `postgres_data` volume.

The checked-in `.env.example` contains development-only credentials. Change them in `.env` for your local environment; do not commit `.env`. If the password contains URL-reserved characters, URL-encode it in `DATABASE_URL`.

Stop the services with `docker compose down`. To also permanently delete the database volume, use `docker compose down -v`.

## API

All responses use JSON except for empty `204 No Content` responses. `GET /health` is a liveness check; it does not verify database readiness.

| Method | Path | Purpose |
|---|---|---|
| `GET` | `/health` | Liveness check |
| `GET`, `POST` | `/notes` | List/search or create notes |
| `GET`, `PUT`, `DELETE` | `/notes/{id}` | Read, update, or delete a note |
| `POST` | `/notes/{id}/tags` | Add tags to a note |
| `GET`, `POST` | `/bookmarks` | List/search or create bookmarks |
| `GET`, `PUT`, `DELETE` | `/bookmarks/{id}` | Read, update, or delete a bookmark |
| `POST` | `/bookmarks/{id}/tags` | Add tags to a bookmark |

List endpoints accept optional `search` and `tag` query parameters. For example:

```text
GET /notes?search=learning&tag=go
GET /bookmarks?search=docs&tag=reference
```

Create a note with:

```json
{"title":"Learning Go","content":"Notes go here"}
```

Create a bookmark with:

```json
{"url":"https://go.dev","title":"The Go Programming Language","description":"Official Go site"}
```

## Develop and test

Run these commands from this directory:

```powershell
go test ./...
go vet ./...
sqlc generate
```

The API applies outstanding migrations automatically at startup. Migration files are in `db/migrations`; query definitions are in `db/queries`, and generated Go code is in `internal/store/dbgen`.

## Local database backups

Create a PostgreSQL custom-format dump while Compose is running:

```powershell
docker compose exec -T db pg_dump -U app -d notes_bookmarks -Fc -f /tmp/notes-bookmarks.dump
docker compose cp db:/tmp/notes-bookmarks.dump .\notes-bookmarks.dump
```

To restore into a fresh database, copy the dump into the container and restore it:

```powershell
docker compose cp .\notes-bookmarks.dump db:/tmp/notes-bookmarks.dump
docker compose exec -T db pg_restore -U app -d notes_bookmarks /tmp/notes-bookmarks.dump
```

Keep backups outside the repository and verify restores periodically. These commands create manual local backups; they do not configure automated retention or off-site storage.

## Deployment and security

This Compose configuration is for local development: both published ports bind to loopback, and PostgreSQL TLS is disabled for the private Compose network. The API currently has no authentication or authorization, and `GET /health` is only a liveness check. Do not expose this setup directly to the public internet.

Before a public deployment, use managed secrets rather than a checked-in or broadly accessible environment file, configure TLS and network access appropriate to the deployment, automate encrypted and tested backups with retention/off-site storage, and design and implement authentication and authorization for the API. Terminate external HTTPS at a properly configured reverse proxy or load balancer.
