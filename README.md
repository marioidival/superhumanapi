# Superhuman API

This project is a backend service for the Superhuman application, designed to manage user data and interactions. It's built using Go, with PostgreSQL for data persistence.

## Prerequisites

- Go 1.22+
- Docker and Docker Compose
- PostgreSQL

## Getting Started

1. **Clone the repository**

    ```sh
    git clone https://github.com/marioidival/superhumanapi.git
    cd superhumanapi
    ```

2. **Start the PostgreSQL database**

    Use Docker Compose to start a local PostgreSQL instance:

    ```sh
    docker-compose -p superhuman -f docker-compose.yml up -d postgres
    ```

3. **Run database migrations**

    Migrate the database schema:

    ```sh
    make db/migrate
    ```

4. **Generate SQL queries**

    Generate SQL queries using `sqlc`:

    ```sh
    make generate/queries
    ```

5. **Build the API server**

    Compile the API server binary:

    ```sh
    make $(BIN)/api
    ```


6. **Start the API server builded locally**

    Start the server:

    ```sh
    ./bin/api
    ```

7. **Start the API server by docker compose**

    Start the server:

    ```sh
    docker-compose -p superhuman -f docker-compose.yml up app
    ```

## Development

- **Adding new migrations**

    To add new database migrations, place your SQL migration files in `internal/db/schema/migrations` and run `make db/migrate`.

- **Generating client code**

    If you make changes to the database schema or queries, regenerate the client code:

    ```sh
    make generate/queries
    ```

## Contributing

Contributions are welcome! Please feel free to submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.