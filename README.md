# Golang Inventory Management APIs
    REST services used to maintain the product catalogue of an ecommerce company.

## Setup Postgres Database

    $ cd migration
    $ goose postgres "<POSTGRES_DB_URL>" up
    $ cd ..

## Run Development Environment

    Rename env.example to development.env in /config
    Set the necessary environment variables in development.env
    
    $ source config/development.env
    $ go run main.go
