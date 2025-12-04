# Real-Time Collaboration Platform Backend

A production-grade real-time collaboration platform backend built in Go, similar to Google Docs. This project demonstrates advanced backend engineering concepts including real-time synchronization, conflict resolution, WebSockets, concurrency, and distributed systems.

## 🚀 Features

### Core Features
- **Real-Time Editing**: WebSocket-based real-time collaboration with instant updates
- **Conflict Resolution**: CRDT-based conflict resolution for concurrent edits
- **Role-Based Access Control (RBAC)**: Owner, Editor, and Viewer roles
- **Document Management**: Create, update, share documents with different types (text, notes, whiteboards, tasks)
- **Version History**: Track document versions and changes
- **Activity Feed**: Monitor who edited what and when
- **Secure Sharing**: Share documents with secure tokens and permissions
- **Offline Support**: Queue operations for offline users
- **JWT Authentication**: Secure token-based authentication
- **API Documentation**: Interactive Swagger/OpenAPI documentation

### Architecture
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Microservices Ready**: Modular design that can be split into microservices
- **PostgreSQL**: Reliable relational database for data persistence
- **Redis**: Pub/Sub for distributed message broadcasting
- **WebSocket Hub**: Goroutine-based WebSocket management for real-time communication
- **Docker**: Containerized deployment with Docker Compose

## 🛠 Tech Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **WebSocket**: Gorilla WebSocket
- **Database**: PostgreSQL 15
- **Cache/PubSub**: Redis 7
- **ORM**: GORM
- **Authentication**: JWT
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized setup)
- PostgreSQL 15+ (if running locally without Docker)
- Redis 7+ (if running locally without Docker)

## 🚀 Quick Start

### Option 1: Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/afom12/collabPlatform-backend.git
   cd collabPlatform-backend
   ```

2. **Start services**
   ```bash
   docker-compose up -d
   ```

   This will start:
   - PostgreSQL on port 5432
   - Redis on port 6379
   - Backend API on port 8080

3. **Check health**
   ```bash
   curl http://localhost:8080/health
   ```

4. **Access Swagger Documentation**
   ```
   http://localhost:8080/swagger/index.html
   ```

### Option 2: Local Development

1. **Install dependencies**
   ```bash
   go mod download
   ```

2. **Set up PostgreSQL and Redis**
   - Start PostgreSQL and create database `collab_platform`
   - Start Redis server

3. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Generate Swagger documentation** (optional)
   ```bash
   make swagger
   ```

5. **Run the server**
   ```bash
   go run cmd/server/main.go
   ```

## 📡 API Endpoints

### Interactive API Documentation

Visit **http://localhost:8080/swagger/index.html** for interactive Swagger UI documentation where you can:
- Browse all endpoints
- Test API calls directly
- View request/response schemas
- Authenticate with JWT tokens

### Authentication

- `POST /api/v1/auth/register` - Register a new user
  ```json
  {
    "email": "user@example.com",
    "username": "username",
    "password": "password123"
  }
  ```

- `POST /api/v1/auth/login` - Login and get JWT token
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- `GET /api/v1/auth/profile` - Get current user profile (requires auth)

### Documents

- `POST /api/v1/documents` - Create a new document (requires auth)
  ```json
  {
    "title": "My Document",
    "type": "text"
  }
  ```

- `GET /api/v1/documents` - List user's documents (requires auth)
- `GET /api/v1/documents/:id` - Get document by ID (requires auth)
- `PUT /api/v1/documents/:id` - Update document (requires auth)
  ```json
  {
    "title": "Updated Title",
    "content": "Document content"
  }
  ```

- `POST /api/v1/documents/:id/share` - Share document with user (requires auth)
  ```json
  {
    "user_id": "user-uuid",
    "role": "editor"
  }
  ```

- `GET /api/v1/documents/:id/versions` - Get document version history
- `GET /api/v1/documents/:id/activities` - Get document activity feed

### WebSocket

- `GET /api/v1/ws?token=<jwt_token>&document_id=<doc_id>` - Connect to WebSocket for real-time collaboration

**WebSocket Message Format:**
```json
{
  "type": "operation",
  "operation": {
    "type": "insert",
    "position": 10,
    "length": 0,
    "content": "Hello"
  },
  "document_id": "doc-uuid"
}
```

## 📚 Swagger Documentation

### Setup Swagger

1. **Install Swagger CLI**
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate Documentation**
   ```bash
   make swagger
   # Or manually: swag init -g cmd/server/main.go -o docs
   ```

3. **Access Swagger UI**
   - Start the server
   - Visit: `http://localhost:8080/swagger/index.html`

See [SWAGGER_SETUP.md](SWAGGER_SETUP.md) for detailed instructions.

## 🏗 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── domain/                  # Domain entities and business rules
│   │   ├── user.go
│   │   ├── document.go
│   │   └── collaboration.go
│   ├── usecase/                 # Business logic
│   │   ├── auth.go
│   │   ├── document.go
│   │   └── collaboration.go
│   ├── infrastructure/          # External dependencies
│   │   ├── database/
│   │   ├── redis/
│   │   └── repository/
│   └── delivery/                # Delivery mechanisms
│       ├── http/
│       │   ├── handlers/
│       │   └── middleware/
│       └── websocket/
├── docs/                        # Swagger documentation
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## 🔐 Security Considerations

- **JWT Secret**: Change `JWT_SECRET` in production
- **CORS**: Configure CORS properly for production
- **Password Hashing**: Uses bcrypt with default cost
- **SQL Injection**: Protected by GORM parameterized queries
- **WebSocket Origin**: Currently allows all origins - restrict in production

## 🧪 Testing the Platform

### 1. Register a user
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

Save the token from the response.

### 3. Create a document
```bash
curl -X POST http://localhost:8080/api/v1/documents \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{
    "title": "My First Document",
    "type": "text"
  }'
```

### 4. Test with Swagger UI
- Visit `http://localhost:8080/swagger/index.html`
- Click "Authorize" and enter your JWT token
- Test endpoints directly from the UI

### 5. Connect via WebSocket
Use a WebSocket client or browser console:
```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws?token=<your-token>&document_id=<doc-id>');
ws.onmessage = (event) => console.log('Received:', JSON.parse(event.data));
ws.send(JSON.stringify({
  type: 'operation',
  operation: {
    type: 'insert',
    position: 0,
    length: 0,
    content: 'Hello World'
  },
  document_id: '<doc-id>'
}));
```

## 🚀 Production Deployment

1. **Environment Variables**: Set all required environment variables
2. **Database**: Use managed PostgreSQL service
3. **Redis**: Use managed Redis service
4. **Load Balancer**: Add load balancer for multiple instances
5. **Monitoring**: Add logging, metrics, and tracing
6. **SSL/TLS**: Enable HTTPS
7. **Rate Limiting**: Add rate limiting middleware
8. **Scaling**: Consider horizontal scaling with Redis pub/sub

## 📝 License

This project is for educational and demonstration purposes.

## 🤝 Contributing

This is a showcase project demonstrating advanced Go backend engineering concepts.

## 📚 Advanced Features (Future Enhancements)

- [ ] Kafka integration for event streaming
- [ ] NATS for messaging
- [ ] MongoDB for document storage
- [ ] gRPC microservices
- [ ] Advanced CRDT implementation (Yjs, Automerge)
- [ ] Presence indicators (who's online)
- [ ] Comments and suggestions
- [ ] File attachments
- [ ] Search functionality
- [ ] Export/Import features

## 📖 Documentation

- [README.md](README.md) - This file
- [QUICKSTART.md](QUICKSTART.md) - Quick start guide
- [SWAGGER_SETUP.md](SWAGGER_SETUP.md) - Swagger documentation setup
- [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md) - Architecture overview
