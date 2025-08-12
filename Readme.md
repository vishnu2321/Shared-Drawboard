# Shared Drawboard

Shared Drawboard is a real-time, collaborative web application that allows users to sign up, log in, and draw together on a shared canvas.

## Features

### User Authentication
*   **Secure Signup & Login**: A secure process for user registration and authentication.
*   **JWT-based Sessions**: User sessions are managed using JSON Web Tokens (JWT), with automated token retrieval and refresh to maintain a seamless user experience.

### Collaborative Drawboard
*   **Real-time Collaboration**: Users can join a shared drawing session and see updates from other participants instantly.
*   **Comprehensive Drawing Tools**:
    *   ✏️ **Freehand Drawing**: Draw freely on the canvas.
    *   📏 **Line**: Create straight lines.
    *   ⬜ **Rectangle**: Draw rectangles and squares.
    *   ⭕ **Circle**: Draw circles and ovals.
    *   🔤 **Text**: Add text annotations to the drawing.
    *   🧽 **Eraser**: Remove parts of the drawing.
    *   ↖️ **Selection & Resizing**: Select, move, and resize existing shapes.
*   **Tool Customization**:
    *   Select custom colors for shapes and text.
    *   Adjust the stroke thickness for drawing tools.
    *   Change the size of the eraser.
*   **Board Management**:
    *   Clear the entire drawing board with a single click.

## Tech Stack

*   **Backend**: Go
*   **API & Routing**: `gorilla/mux`
*   **Database**: MongoDB
*   **Authentication**: JSON Web Tokens (JWT)
*   **Real-time Communication**: WebSockets
*   **Frontend**: HTML, CSS, JavaScript

## Project Structure

The application follows a standard layered architecture to separate concerns.


## Getting Started

### Prerequisites
*   [Go](https://golang.org/doc/install) (version 1.18 or newer)
*   [MongoDB](https://www.mongodb.com/try/download/community)
*   A running MongoDB instance.

### Installation & Running

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/shared-drawboard.git
    cd shared-drawboard
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

3.  **Configure Environment Variables:**
    Create a `.env` file in the root of the project and add the necessary configuration.
    ```env
    MONGO_URI="mongodb://localhost:27017"
    JWT_SECRET="your-strong-jwt-secret"
    PORT="8080"
    ```

4.  **Run the application:**
    ```sh
    go run cmd/main.go
    ```

5.  **Access the application:**
    Open your web browser and navigate to `http://localhost:8080`.