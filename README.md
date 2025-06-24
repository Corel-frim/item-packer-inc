# Item Packer Inc

A web application for packing items into standard sized packs.

## Notes

- Obviously, I used AI, but only to help with some routine tasks to save time, like: add frontend (since it's not the main focus for the task), basic unit tests, write documentation, check that my implementation covers all the requirements.
- What I did not use AI for: CI/CD (get this one from my pet project), fiber usage, main.go setup (again, got basic settings from a pet project), project structure and business logic implementation (but used help to optimize it for readability and comments).
- Code is written with moving to more appropriate storage in mind (like Postgres or MongoDB) in the future, so it is not tightly coupled with the in-memory storage implementation.
- I used `sync.Mutex` to make the in-memory storage thread-safe, but it is not a production-ready solution (though, should be good enough for the test task).
- Swagger wasn't a part of the requirements, but I used to have that in every project (if there's no other communication channel like gRPC), so I added it here as well.
- Also, I added a few notes to my implementation in index.html ('Orders' tab) thinking you might want to visit the website first and see that, but if not - now you know where to find it.

## Overview

This application provides a RESTful API and web interface for managing packs and creating orders with optimal packing. It calculates the most efficient way to pack items into standard-sized packs, minimizing waste.

## Features

- Define and manage pack sizes (add, update, delete, list)
- Calculate optimal packing for orders
- RESTful API with JSON responses
- Web-based user interface
- Swagger API documentation

## Live Demo

You can try the application at: https://item-packer-inc.fly.dev/

## Getting Started

### Prerequisites

- Go 1.24 or higher
- Docker (optional, for containerized deployment)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/corel-frim/item-packer-inc.git
   cd item-packer-inc
   ```

2. Install required tools:
   ```bash
   make install
   ```
   This installs Swagger for API documentation and golangci-lint for code quality.

3. Generate Swagger documentation if annotations changed (optional):
   ```bash
   make swagger
   ```

### Running the Application

#### Local Development

Run the application locally:
```bash
make run
```

#### Using Docker

Build and run using Docker:
```bash
make run-docker
```

This builds a Docker image and runs it, exposing the application on port 8080.

### Testing

Run the test suite:
```bash
make test
```

Run the linter:
```bash
make linter
```

Fix linting issues automatically:
```bash
make lint-fix
```

### Access the Application

Once running, you can access:
- Web UI: http://localhost:8080
- API Documentation: http://localhost:8080/swagger/index.html

## API Endpoints

### Packs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/packs` | Get all available packs |
| POST | `/packs/{amount}` | Add a new pack with specified amount |
| PUT | `/packs/{oldAmount}/{newAmount}` | Update a pack's amount |
| DELETE | `/packs/{amount}` | Delete a pack |

### Orders

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/orders/items/{amount}` | Create an order with specified number of items |
| GET | `/orders` | Get all orders |

## Storage

The application uses an in-memory storage implementation:

- Data is not persisted across application restarts
- Both packs and orders are stored in memory
- There's a soft limit of 20 items for both packs and orders
- Thread-safe implementation using mutexes

In a production environment, you might want to replace this with a database implementation.

## Project Structure

```
item-packer-inc/
├── api/              # API implementation
│   ├── handlers/     # Request handlers
│   └── api.go        # API setup
├── cmd/              # Application entry points
│   └── main.go       # Main application
├── docs/             # API documentation
│   └── swagger/      # Swagger definitions
├── frontend/         # Web UI
│   ├── css/          # Stylesheets
│   ├── js/           # JavaScript files
│   └── index.html    # Main HTML file
├── internal/         # Internal packages
│   ├── models/       # Data models
│   └── storage/      # Data storage
├── Dockerfile        # Docker configuration
├── Makefile          # Build and run commands
└── README.md         # This file
```

## Deployment

The application is configured for deployment to [fly.io](https://fly.io) using GitHub Actions. When code is pushed to the main branch, it's automatically deployed to the test environment.

## Example Usage

### Using the Web Interface

1. Open your browser and navigate to `http://localhost:8080`
2. Use the "Packs" tab to manage available pack sizes
3. Use the "Orders" tab to create new orders and view results

### Using the API

#### Add a new pack

```bash
curl -X POST http://localhost:8080/packs/100
```

#### Get all packs

```bash
curl -X GET http://localhost:8080/packs
```

#### Create an order

```bash
curl -X POST http://localhost:8080/orders/items/1234
```

## Data Models

### Pack

```json
{
  "amount": 250
}
```

### Order

```json
{
  "requestedItems": 1234,
  "overpackedItems": 16,
  "totalItems": 1250,
  "packs": [
    {
      "quantity": 2,
      "pack": {
        "amount": 500
      }
    },
    {
      "quantity": 1,
      "pack": {
        "amount": 250
      }
    }
  ]
}
```

---

For any questions or issues, please open an issue on the GitHub repository.
