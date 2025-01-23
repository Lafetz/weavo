# Introduction

A basic CRUD API built with the net/http package, following Hexagonal Architecture principles.
[**LIVE SWAGGER**](https://weavo.onrender.com/api/swagger/index.html)

- Unit and Integration Tests
- Continuous Integration (CI) using github actions
- API Documentation using Swagger
- net/http Package
- Docker

## How to Run

- Make sure to set the following environment variables:

```sh
OPEN_URL="https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s"
OPEN_KEY=YOUR_API_KEY_HERE
PORT=8080
LOG_LEVEL=info
ENV=development
```

### Using Docker

1. Clone the Repository
2. Build the Docker Image:

```sh
 docker build -t weavo-api .
 Run the Docker Container:
docker run --env-file .env -p weavo-api
```

Access the API:
The API will be available at /api/v1/persons.

View Swagger Documentation:
Access the Swagger UI at /swagger/index.html.

Run Tests:

 ```sh
    make test
 ```

Coverage:

 ```sh
   make coverage 
 ```

### Without Docker

Clone the Repository:

Install Dependencies: Make sure you have Go and make installed and set up. Then, install any required dependencies:

```sh
go mod tidy

```

Run the Application:

```sh

make run

```
