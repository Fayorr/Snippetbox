# Snippetbox

## Run the Go app locally with MySQL in Docker

1. Create your local configuration:

   ```sh
   cp .env.example .env
   ```

2. Start only MySQL:

   ```sh
   docker compose up -d db
   ```

3. Run the application on your Mac:

   ```sh
   go run ./cmd/web
   ```

The Docker MySQL port is published as `3309` on your Mac. The `.env.example`
file deliberately uses `DB_PORT=3309` for this mode. Inside the `web`
container, Compose supplies the correct internal address (`db:3306`) instead.

## Run with a locally installed MySQL server

Keep the same `DB_USER`, `DB_PASSWORD`, and `DB_NAME` values in `.env`, but
set `DB_PORT=3306` (or the port used by your server). Create the database and
user with matching credentials, for example:

```sql
CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER 'web'@'localhost' IDENTIFIED BY 'pass';
GRANT ALL PRIVILEGES ON snippetbox.* TO 'web'@'localhost';
```

Then run `go run ./cmd/web`. Do not commit `.env`; it can contain real
credentials and is already ignored by Git.

`using password: YES` in a MySQL error means the application did send a
password. It does not mean a password was missing; it means the server did not
accept the supplied user/password combination.
