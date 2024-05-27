# Go REST API

A simple RESTful API built with Go, GORM, and PostgreSQL for performing CRUD operations in my journey of learning Go and getting comfortable with it.

## Tech Stack

- **Go** - The programming language used for building the API.
- **GORM** - An ORM library for Go, used for interacting with the PostgreSQL database.
- **PostgreSQL** - The relational database management system used for storing and managing note data.

## Features

- **Create Note**: Create a new note by providing title, content, category, and publication status.
- **Read Notes**: Retrieve a list of all available notes or a specific note by ID.
- **Update Note**: Update the details (title, content, category, publication status) of an existing note.
- **Delete Note**: Delete an existing note by providing its ID.

## Getting Started

1. Clone the repository:

    ```bash
    git clone https://github.com/your-username/go-rest-api.git
    ```

2. Copy and set up the PostgreSQL database and update the connection details in the `.env-example` file:

    ```bash
    cp .env-example .env
    ```

3. Install the dependencies:

    ```bash
    go mod download
    ```

4. Run the API:

    ```bash
    go run main.go
    ```

The API will be available at `http://localhost:8750`.

## Endpoints
- `POST /api/auth/register`: Register new user
- `POST /api/auth/login`: Login user
- `POST /api/auth/logout`: Logout user

- `POST /api/notes`: Create a new note
- `GET /api/notes`: Retrieve a list of all notes
- `GET /api/notes/:id`: Retrieve a specific note by ID
- `PUT /api/notes/:id`: Update an existing note by ID
- `DELETE /api/notes/:id`: Delete an existing note by ID

## Todo

- [x] Add authentication feature for securing the API endpoints.
- [ ] Add Role-based access control.
- [ ] Implement rate limiting to prevent abuse.
- [ ] Add pagination to the `GET /api/notes` endpoint.
- [ ] Write unit and integration tests.
- [ ] Implement a search functionality for notes.
- [ ] Add Swagger documentation for the API.
- [ ] Set up a CI/CD pipeline for automated testing and deployment.

