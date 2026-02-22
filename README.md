
# GoL7-LB: A Lightweight Layer 7 Load Balancer

A high-performance Layer 7 Load Balancer implemented in Go, featuring active health checks, reverse proxying, and graceful shutdown capabilities. This project demonstrates the practical application of Go's concurrency primitives and networking standard library.

## 🚀 Features

* **Layer 7 Reverse Proxying**: Utilizes `httputil.ReverseProxy` to forward HTTP requests to backend servers.
* **Round-Robin Scheduling**: Evenly distributes incoming traffic across healthy backend instances.
* **Active Health Checks**: Periodically monitors backend availability via TCP dialing to ensure high service uptime.
* **Thread-Safe State Management**: Uses `sync.RWMutex` to allow high-frequency reads of backend status while ensuring safe periodic updates.
* **Graceful Shutdown**: Listens for system interrupts to clean up resources and stop background workers without dropping active connections.
* **Containerized Deployment**: Includes `Dockerfile` and `docker-compose` for rapid local orchestration.

## 🛠 Architecture

The system consists of a **Load Balancer** acting as the entry point and multiple **Backend Servers** that respond with their specific port information.

* **Server Pool**: Manages the list of backends and the current rotation index.
* **Health Check Worker**: A background goroutine that pings backends every 10 seconds.
* **Context Control**: Uses `context.Context` to prevent goroutine leaks during shutdown.

## 📦 Getting Started

### Prerequisites

* Go 1.21+
* Docker & Docker Compose (optional)

### Running with Docker (Recommended)

The easiest way to see the Load Balancer in action is using Docker Compose:

```bash
docker-compose up --build

```

This will spin up one Load Balancer at `:8080` and three backend instances.

### Running Manually

1. **Start Backend Servers**:
```bash
go run cmd/backend/main.go 8081
go run cmd/backend/main.go 8082

```


2. **Start the Load Balancer**:
```bash
go run cmd/lb/main.go

```



## 🧪 Technical Deep Dive

### Why Mutex over Channels for State?

In this project, the "Alive" status of a backend is a state frequently read by the proxy logic but rarely updated by the health checker. Using `sync.RWMutex` provides better performance than channels for this "broadcast" style state sharing, as it allows multiple readers to access the status simultaneously without blocking.

### Health Check Strategy

The current implementation uses **Active Health Checks** via TCP Dialing. While HTTP-based checks can verify application logic, TCP checks provide a lightweight way to ensure the process and network stack are responsive with minimal overhead.

## 📈 Future Roadmap

* [ ] Implement Weighted Round-Robin for heterogeneous server capacities.
* [ ] Add Prometheus metrics for monitoring request rates and latency.
* [ ] Support Dynamic Configuration reloading (hot-swapping backends without restart).

---

Would you like me to help you write the **Unit Tests** for the `GetNextPeer` logic next?
