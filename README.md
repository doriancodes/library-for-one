# Library for One

## Running the project locally

This project uses **Docker**, **PostgreSQL**, and **golang-migrate** for
database migrations.

### Requirements

You need the following installed:

-   Docker
-   Docker Compose
-   Go (for installing the migration CLI)

Install the migration tool:

``` bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

Make sure `$GOPATH/bin` is in your `PATH` so the `migrate` command is
available.

------------------------------------------------------------------------

## 1. Start the database

Start the PostgreSQL container:

``` bash
docker compose up -d db
```

The database will be available at:

    host: localhost
    port: 5432
    database: library
    user: postgres
    password: postgres

------------------------------------------------------------------------

## 2. Run database migrations

Apply all migrations:

``` bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/library?sslmode=disable" up
```

This will create the required tables (for example the `books` table).

You can verify the tables with:

``` bash
docker exec -it library-db psql -U postgres -d library
```

Then inside psql:

``` sql
\dt
```

------------------------------------------------------------------------

## 3. Start the application

Run the application container:

``` bash
docker compose up --build app
```

The web server will start at:

    http://localhost:8080

------------------------------------------------------------------------

## 4. Stopping the project

Stop containers:

``` bash
docker compose down
```

To also remove the database data:

``` bash
docker compose down -v
```

------------------------------------------------------------------------

## Project structure

    .
    ├── docker-compose.yml
    ├── migrations/
    │   ├── 0001_create_books.up.sql
    │   └── 0001_create_books.down.sql
    ├── static/
    │   └── page.html
    ├── main.go
    ├── db.go
    ├── models.go
    └── README.md

------------------------------------------------------------------------

## Development notes

-   The application reads database configuration from environment
    variables.
-   Inside Docker the application connects to PostgreSQL using the
    service name `db`.
-   External tools (e.g. DBeaver) should connect to `localhost:5432`.
