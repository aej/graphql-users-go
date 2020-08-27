# Graphql Users

A production-ready user service exposing a GraphqlAPI api. Follows [12 Factor App](https://12factor.net/) best practices.

# Features

* User account creation with email and password
* Email confirmation 
* User login and session persistence (cookie-based with db-backed session storage)
* User logout
* Password reset


The stack:

- Go
- Postgres
- Docker 


## Development

Configuration is managed by the environment. 

Create a `.env` file with the following:

```
APP_ENV=dev
DATABASE_URL=postgres://postgres@localhost:5432/graphql-users?sslmode=disable
SECRET_KEY=andyjones1111111  # should be 16, 24 or 32 bytes
HASH_KEY=andyjones1111111    # should be 16, 24 or 32 bytes
HOST=localhost
```

The application uses [Air](https://github.com/cosmtrek/air) for hot reloading while developing locally. 

Simply use

`air` 

to run the app.

The application exposes the following endpoints:

- `/query` - The GraphQL endpoint
- `/graphql` - [GraphQL Playground](https://github.com/prisma-labs/graphql-playground)
- `/health` - Health check endpoint which just returns HTTP 200

Behind the scenes the application uses [gqlgen](https://github.com/99designs/gqlgen). Check out the gqlgen documentation for information on extending the Graphql endpoints.

## Migrations

A separate migration entrypoint manages the migrations. To run migrations in the development environment use

```bash
go run cmd/migrate/migrate.go up
```

In production, you can build the `migrate` binary and run migrations

```bash
go build -o /path/to/binary/migrate cmd/migrate/migrate.go
/bin/migrate up
```

## Production

Build and run the docker image

```
docker build . -t graphql-users
```

```
docker run -p 8000:8000 \
    -e SECRET_KEY=wPSZOmn5Ykp9yacq \
    -e HASH_KEY=wPSZOmn5Ykp9yacq \
    -e HOST=localhost \
    -e DATABASE_URL=postgres://postgres@localhost:5432/graphql-users?sslmode=disable \
    graphql-users
```
