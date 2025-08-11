# Shared Drawboard

This is a shared drawboard application that allows users to sign up, log in, and collaborate on a drawing in real-time.

## Architecture

The application follows a layered architecture:

***`cmd/main.go`**: Entry point of the application.
***`internal/database/`**: Database configuration and connection logic (MongoDB).
***`internal/handler/`**: HTTP handler logic and routing (using `gorilla/mux`).
***`internal/middleware/`**: Authentication middleware (JWT).
***`internal/models/`**: Data models.
***`internal/service/`**: Business logic.
***`web/`**: Frontend files for user interface (login and drawboard).

## Functionality

***User Authentication**:
    *Secure user signup and signin process.
    *Authentication is handled using JSON Web Tokens (JWT).
    *The login page handles token retrieval and refresh to maintain user sessions.
***Shared Drawboard**:
    *Real-time collaborative drawing experience.
    *Users can join a shared drawing session and see updates from other participants in real-time.
    ***Available Tools**:
        *Freehand Drawing (✏️)
        *Line (📏)
        *Rectangle (⬜)
        *Circle (⭕)
        *Text (🔤)
        *Eraser (🧽)
        *Selection and Resizing (↖️)
    ***Customization**:
        *Color selection for drawing tools.
        *Adjustable stroke thickness for drawing tools.
        *Adjustable size for the eraser.
    ***Board Management**:
        *Option to clear the entire drawing board.

## Project Structure

***`cmd/`**: Contains the main application executable.
***`internal/`**: Houses the core application logic, including database interactions, request handlers, middleware, models, and services.
***`pkg/`**: Contains reusable packages, such as authentication utilities and helper functions.
***`web/`**: Contains the frontend static assets, including:
    *`web/login/`: Files for the user login and signup interface.
    *`web/drawboard/`: Files for the main shared drawboard interface, including HTML, CSS, and JavaScript.
