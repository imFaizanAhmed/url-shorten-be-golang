# URL Shortener with Redis and Docker

A high-performance URL shortening service built with Go, Beego framework, and Redis for storage, fully containerized with Docker.

## 🚀 Features

- **Fast URL Shortening**: Generate short codes for long URLs
- **Redis Storage**: Persistent storage with configurable TTL
- **Docker Ready**: Full Docker and docker-compose support
- **Health Checks**: Built-in health monitoring
- **Redis UI**: Optional Redis Commander for database management

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Client        │    │   Go App        │    │   Redis         │
│                 │───▶│   (Beego)       │───▶│   (Storage)     │
│   HTTP Requests │    │   Port: 8080    │    │   Port: 6378    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                                ▼
                       ┌─────────────────┐
                       │   Redis UI      │
                       │                 │
                       │   Port: 8081    │
                       └─────────────────┘
```

## 🐳 Docker Configuration Explained

### Docker Compose Services

1. **Redis Service**
   - **Image**: `redis:7-alpine` (lightweight, production-ready)
   - **Purpose**: Primary data storage for URL mappings
   - **Persistence**: Data stored in Docker volume `redis_data`
   - **Health Check**: Ensures Redis is ready before starting the app
   - **Configuration**: Custom `redis.conf` for optimized performance

2. **Go Application Service**
   - **Build**: Multi-stage Dockerfile for optimized image size
   - **Dependencies**: Waits for Redis health check to pass
   - **Environment**: Configurable Redis connection settings
   - **Health Check**: Validates app is responding to HTTP requests

3. **Redis Commander** (Optional)
   - **Purpose**: Web-based Redis management interface
   - **Access**: http://localhost:8081 (admin/admin123)
   - **Use Case**: Development and debugging

### Redis Configuration Deep Dive

Our Redis setup includes several important configurations:

#### **Persistence**
```redis
# RDB Snapshots (point-in-time backups)
save 900 1      # Save if at least 1 key changed in 900 seconds
save 300 10     # Save if at least 10 keys changed in 300 seconds  
save 60 10000   # Save if at least 10000 keys changed in 60 seconds

# AOF (Append Only File - logs every write operation)
appendonly yes
appendfsync everysec  # Sync to disk every second
```

#### **Memory Management**
```redis
# Evict keys using the Least Recently Used policy when memory is full
maxmemory-policy allkeys-lru
```

#### **Performance**
```redis
# Connection settings
maxclients 10000
tcp-keepalive 300
tcp-backlog 511
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### 1. Clone and Setup
```bash
git clone <repository-url>
cd url-shorten-be
cp env.example .env  # Copy and customize environment variables
```

### 2. Start with Docker Compose
```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Check service status
docker-compose ps
```

### 3. Test the API
```bash
# Create a short URL
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://www.example.com/very-long-url"}'

# Access short URL (replace {shortCode} with actual code)
curl -L http://localhost:8080/{shortCode}
```

## 🔧 Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `REDIS_HOST` | `localhost` | Redis server hostname |
| `REDIS_PORT` | `6378` | Redis server port |
| `REDIS_PASSWORD` | `` | Redis password (if auth enabled) |
| `REDIS_DB` | `0` | Redis database number |
| `BEE_RUNMODE` | `dev` | Beego run mode (dev/prod) |

## 🐛 Development

### Local Development (without Docker)
```bash
# Start Redis locally
redis-server

# Set environment variables
export REDIS_HOST=localhost
export REDIS_PORT=6378

# Run the application
cd core-api
go run main.go
```

### Building Docker Image Manually
```bash
cd core-api
docker build -t url-shortener .
```

## 📊 Monitoring

### Health Checks
- **App Health**: http://localhost:8080/
- **Redis Health**: `docker-compose exec redis redis-cli ping`

### Redis Management
- **Redis Commander**: http://localhost:8081
- **Direct Redis CLI**: `docker-compose exec redis redis-cli`

### Useful Redis Commands
```bash
# View all URL keys
KEYS url:*

# Get a specific URL
GET url:abc123

# Check TTL for a key
TTL url:abc123

# View Redis info
INFO memory
```

## 🔒 Production Considerations

### Security
1. **Enable Redis Authentication**
   ```redis
   requirepass your_strong_password
   ```

2. **Network Security**
   - Use internal Docker networks
   - Don't expose Redis port externally
   - Use HTTPS for the application

### Performance
1. **Redis Optimization**
   - Monitor memory usage
   - Adjust `maxmemory-policy` based on use case
   - Use Redis clustering for high availability

2. **Application Scaling**
   - Run multiple app instances behind a load balancer
   - Use Redis Sentinel for Redis high availability

### Backup
```bash
# Backup Redis data
docker-compose exec redis redis-cli BGSAVE

# Copy backup file
docker cp url_shortener_redis:/data/dump.rdb ./backup-$(date +%Y%m%d).rdb
```

## 🐳 Docker Best Practices Implemented

1. **Multi-stage builds**: Smaller production images
2. **Non-root user**: Security best practice
3. **Health checks**: Ensure services are ready
4. **Named volumes**: Persistent data storage
5. **Custom networks**: Service isolation
6. **Resource limits**: Prevent resource exhaustion
7. **Restart policies**: Automatic recovery

## 📝 API Endpoints

| Method | Endpoint | Description | Body |
|--------|----------|-------------|------|
| `GET` | `/` | API info | - |
| `POST` | `/shorten` | Create short URL | `{"long_url": "https://example.com"}` |
| `GET` | `/{shortCode}` | Redirect to long URL | - |

## ❓ Interview Questions & Answers

### Redis Configuration
**Q: Why use Redis over in-memory storage?**
A: Redis provides persistence, scalability, and can be shared across multiple application instances.

**Q: Explain RDB vs AOF persistence.**
A: RDB creates point-in-time snapshots (good for backups), AOF logs every write (better durability).

**Q: What's the TTL strategy?**
A: 24-hour TTL prevents unlimited growth while allowing reasonable usage periods.

### Docker
**Q: Why multi-stage builds?**
A: Separates build dependencies from runtime, resulting in smaller, more secure production images.

**Q: What's the health check purpose?**
A: Ensures services are actually ready to handle requests, not just started.

**Q: Why named volumes for Redis?**
A: Provides persistent storage that survives container recreations.

## 🛠️ Troubleshooting

### Common Issues
1. **Port conflicts**: Change ports in docker-compose.yml
2. **Redis connection failed**: Check Redis service status
3. **Build failures**: Ensure Docker has internet access for dependencies

### Debug Commands
```bash
# Check service logs
docker-compose logs redis
docker-compose logs app

# Access container shells
docker-compose exec app sh
docker-compose exec redis sh

# Restart specific service
docker-compose restart app
```
