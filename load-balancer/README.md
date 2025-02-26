# Go Load Balancer

This project implements a simple load balancer in Go that distributes incoming requests using goroutines and implements least-connections load balancing.

## Features

- Distributes incoming requests using goroutines
- Implements least-connections load balancing

## Getting Started

### Prerequisites

- Go 1.16 or later

### Running the Load Balancer

1. Clone the repository:
    ```sh
    git clone https://github.com/yourusername/load-balancer.git
    cd load-balancer
    ```

2. Start the backend servers (for example, using Python's `http.server`):
    ```sh
    python3 -m http.server 8081
    python3 -m http.server 8082
    python3 -m http.server 8083
    ```

3. Run the load balancer:
    ```sh
    go run main.go
    ```

4. Send requests to the load balancer:
    ```sh
    curl http://localhost:8080
    ```

### Running the Tests

Run the unit tests using the following command:
```sh
go test
```

## License

This project is licensed under the MIT License.
