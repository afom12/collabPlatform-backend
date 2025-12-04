# Quick Start Guide

## Prerequisites

- Docker Desktop (includes Docker Compose) OR Docker + docker-compose
- Go 1.21+ (optional, for local development without Docker)

## 🚀 Running with Docker (Easiest Way)

### Windows (PowerShell)
```powershell
.\run.ps1
```

### Linux/Mac (Bash)
```bash
chmod +x run.sh
./run.sh
```

### Manual Docker Commands
```bash
# Start all services
docker compose up -d

# Or if you have docker-compose (with hyphen)
docker-compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

## ✅ Verify Installation

1. **Check health endpoint:**
   ```bash
   curl http://localhost:8080/health
   ```
   Should return: `{"status":"ok"}`

2. **Check services:**
   - PostgreSQL: `docker compose ps postgres`
   - Redis: `docker compose ps redis`
   - Backend: `docker compose ps backend`

## 🧪 Test the API

### 1. Register a User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "testuser",
    "password": "password123"
  }'
```

### 2. Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

Save the `token` from the response.

### 3. Create a Document
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Document",
    "type": "text"
  }'
```

Save the `id` from the response.

### 4. Connect via WebSocket

Use a WebSocket client or browser console:

```javascript
const token = 'YOUR_TOKEN_HERE';
const docId = 'DOCUMENT_ID_HERE';

const ws = new WebSocket(`ws://localhost:8080/api/v1/ws?token=${token}&document_id=${docId}`);

ws.onopen = () => console.log('Connected!');
ws.onmessage = (event) => {
  console.log('Received:', JSON.parse(event.data));
};

// Send an operation
ws.send(JSON.stringify({
  type: 'operation',
  operation: {
    type: 'insert',
    position: 0,
    length: 0,
    content: 'Hello World!'
  },
  document_id: docId
}));
```

## 🔧 Local Development (Without Docker)

1. **Install PostgreSQL and Redis locally**

2. **Set environment variables:**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=collab_platform
   export REDIS_ADDR=localhost:6379
   export JWT_SECRET=your-secret-key
   export SERVER_PORT=8080
   ```

3. **Create database:**
   ```bash
   createdb collab_platform
   ```

4. **Run the server:**
   ```bash
   go run cmd/server/main.go
   ```

## 🐛 Troubleshooting

### Port Already in Use
- Change ports in `docker-compose.yml` or stop conflicting services

### Database Connection Error
- Wait a few seconds for PostgreSQL to fully start
- Check: `docker compose logs postgres`

### Redis Connection Error
- Check: `docker compose logs redis`
- Verify Redis is running: `docker compose ps redis`

### WebSocket Connection Failed
- Ensure you're using `ws://` (not `http://`) for WebSocket connections
- Check token is valid and not expired
- Verify document_id exists

## 📚 Next Steps

- Read the full [README.md](README.md) for detailed API documentation
- Explore the codebase structure
- Customize for your needs
- Add additional features (comments, presence, etc.)

