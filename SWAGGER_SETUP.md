# Swagger/OpenAPI Documentation Setup

## 📚 Overview

This project includes Swagger/OpenAPI documentation for all API endpoints. The interactive API documentation is available at `/swagger/index.html` when the server is running.

## 🚀 Quick Start

### 1. Install Swagger CLI Tool

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 2. Generate Swagger Documentation

```bash
# Using Makefile (recommended)
make swagger

# Or manually
swag init -g cmd/server/main.go -o docs
```

### 3. Start the Server

```bash
go run cmd/server/main.go
```

### 4. Access Swagger UI

Open your browser and navigate to:
```
http://localhost:8080/swagger/index.html
```

## 📖 Using Swagger UI

1. **View All Endpoints**: Browse all available API endpoints organized by tags
2. **Test Endpoints**: Click "Try it out" on any endpoint to test it directly
3. **Authentication**: 
   - Click the "Authorize" button at the top
   - Enter your JWT token (without "Bearer" prefix)
   - Click "Authorize" to authenticate all requests
4. **View Schemas**: See request/response models and examples

## 🔧 Endpoints Documented

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `GET /api/v1/auth/profile` - Get user profile (protected)

### Documents
- `POST /api/v1/documents` - Create document (protected)
- `GET /api/v1/documents` - List user documents (protected)
- `GET /api/v1/documents/{id}` - Get document by ID (protected)
- `PUT /api/v1/documents/{id}` - Update document (protected)
- `POST /api/v1/documents/{id}/share` - Share document (protected)
- `GET /api/v1/documents/{id}/versions` - Get version history
- `GET /api/v1/documents/{id}/activities` - Get activity feed

### Health
- `GET /health` - Health check

## 📝 Adding New Endpoints

To document a new endpoint, add Swagger annotations above the handler function:

```go
// MyHandler godoc
// @Summary      Brief description
// @Description  Detailed description
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Document ID"
// @Param        request  body      MyRequest  true  "Request body"
// @Success      200  {object}  MyResponse
// @Failure      400   {object}  ErrorResponse
// @Router       /my-endpoint/{id} [post]
func (h *MyHandler) MyHandler(c *gin.Context) {
    // Handler implementation
}
```

Then regenerate Swagger docs:
```bash
make swagger
```

## 🔄 Updating Documentation

After making changes to handlers or adding new endpoints:

1. Update Swagger annotations in handler files
2. Run `make swagger` to regenerate documentation
3. Restart the server
4. Refresh Swagger UI to see changes

## 📦 Dependencies

The following packages are used for Swagger:

- `github.com/swaggo/swag` - Swagger generator
- `github.com/swaggo/gin-swagger` - Gin middleware for Swagger
- `github.com/swaggo/files` - Static file server for Swagger UI

## 🎯 Features

- ✅ Interactive API testing
- ✅ Request/response examples
- ✅ Authentication support
- ✅ Schema definitions
- ✅ Error response documentation
- ✅ Parameter validation examples

## 📚 Resources

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [OpenAPI Specification](https://swagger.io/specification/)
- [Swagger UI](https://swagger.io/tools/swagger-ui/)

