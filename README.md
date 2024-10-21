# FxDemo2

FxDemo2 is a simple HTTP server implemented in Go using the [Fx](https://github.com/uber/fx) framework for dependency injection, [Gorilla Mux](https://github.com/gorilla/mux) for routing, and [Zap](https://github.com/uber-go/zap) for logging. This application demonstrates how to set up an HTTP server with lifecycle management using Fx.

## Features

- HTTP server running on port 3333
- Simple routing using Gorilla Mux
- Structured logging with Zap
- Graceful startup and shutdown using Fx lifecycle hooks

## Prerequisites

- Go 1.16 or later
- The following Go packages:
  - `github.com/gorilla/mux`
  - `go.uber.org/fx`
  - `go.uber.org/zap`

You can install the required packages using:

```bash
go get github.com/gorilla/mux
go get go.uber.org/fx
go get go.uber.org/zap
```

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/fxdemo2.git
   cd fxdemo2
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

To run the application, use the following command:

```bash
go run main.go
```

Once the server is running, you can access the home endpoint by navigating to `http://localhost:3333` in your web browser or using a tool like `curl`:

```bash
curl http://localhost:3333
```

You should see the response:

```
Hello from FxDemo2 Server!
```

## Code Overview

### Main Components

- **MyService**: The core service that handles HTTP requests and logging.
- **Router**: A `mux.Router` instance used for routing HTTP requests.
- **Logger**: A `zap.Logger` instance used for logging application events.

### Key Functions

- `NewMyService`: Initializes a new instance of `MyService`.
- `newMuxRouter`: Creates and returns a new `mux.Router`.
- `newZapLogger`: Creates and returns a new Zap logger in development mode.
- `IndexHandler`: Handles requests to the root URL and logs the request.
- `Start`: Starts the HTTP server.
- `Stop`: Logs the stopping of the HTTP server.
- `runService`: Registers the `IndexHandler` and sets up lifecycle hooks for starting and stopping the server.

### Lifecycle Management

The server starts and stops gracefully using Fx's lifecycle hooks. The `OnStart` hook initializes the route and starts the server in a goroutine, while the `OnStop` hook handles cleanup.

## Contributing

If you have suggestions for improvements or find bugs, please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

Feel free to modify any section as needed to fit your specific use case or project structure!