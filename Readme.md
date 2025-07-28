# Shared Drawboard

This is a shared drawboard application.

## Architecture

The application follows a layered architecture:

* **`cmd/main.go`**: Entry point of the application.
* **`internal/database/`**: Database configuration and connection logic (MongoDB).
* **`internal/handler/`**: HTTP handler logic and routing (`gorilla/mux`).
* **`internal/middleware/`**: Authentication middleware (JWT).
* **`internal/models/`**: Data models.
* **`internal/service/`**: Business logic.
* **`web/`**: Frontend files (login and drawboard).

## Functionality

The application provides the following functionality:

* User signup and signin.
* Authentication using JWT.
* Shared drawboard functionality (under development).
