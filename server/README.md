# Develapar Blog Server

A robust, production-ready REST API server for a blog application built with Go, Gin framework, and PostgreSQL. This server provides comprehensive blog functionality with advanced features like authentication, rate limiting, metrics collection, and structured logging.

## 🚀 Features

### Core Functionality

- **User Management**: Registration, authentication, and user profiles
- **Article Management**: CRUD operations with slug-based URLs
- **Category & Tag System**: Organize content with categories and tags
- **Comment System**: User comments on articles
- **Like System**: Article likes and user engagement
- **Bookmark System**: Save articles for later reading

### Advanced Features

- **JWT Authentication**: Secure token-based authentication
- **Rate Limiting**: Configurable rate limiting with sliding window algorithm
- **Metrics Collection**: Comprehensive application and system metrics
- **Structured Logging**: JSON-based logging with context correlation
- **Database Connection Pooling**: Optimized database connections
- **Health Checks**: Application and database health monitoring
- **API Documentation**: Auto-generated Swagger documentation
- **CORS Support**: Cross-origin resource sharing configuration

## 🏗️ Architecture

### Clean Architecture Pattern

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Controllers   │────│    Services     │────│  Repositories   │
│   (HTTP Layer)  │    │ (Business Logic)│    │  (Data Layer)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Middleware    │
                    │ (Cross-cutting) │
                    └─────────────────┘
```

### Middleware Chain

The server uses a carefully ordered middleware chain for optimal performance and security:

1. **Recovery Middleware** - Panic recovery and error handling
2. **CORS Middleware** - Cross-origin resource sharing
3. **Context Middleware** - Request ID and context injection
4. **Request Logging Middleware** - Structured request/response logging
5. **Rate Limiting Middleware** - Request rate limiting
6. **Metrics Middleware** - Performance and usage metrics
7. **Error Handling Middleware** - Centralized error processing

## 📁 Project Structure

```
server/
├── config/                 # Configuration management
│   ├── config.go          # Main configuration
│   └── database_pool.go   # Database connection pooling
├── controller/            # HTTP request handlers
│   ├── article_controller.go
│   ├── user_controller.go
│   ├── health_controller.go
│   └── metrics_controller.go
├── middleware/            # HTTP middleware
│   ├── auth_middleware.go
│   ├── context_middleware.go
│   ├── rate_limiter.go
│   ├── metrics_middleware.go
│   └── request_logging_middleware.go
├── model/                 # Data models
│   ├── dto/              # Data transfer objects
│   ├── user.go
│   ├── article.go
│   └── category.go
├── repository/            # Data access layer
│   ├── user_repository.go
│   ├── article_repository.go
│   └── category_repository.go
├── service/               # Business logic layer
│   ├── user_service.go
│   ├── article_service.go
│   ├── jwt_service.go
│   ├── metrics_service.go
│   └── validation_service.go
├── utils/                 # Utility functions
│   ├── logger.go
│   ├── password_hash.go
│   ├── errors.go
│   └── slug_generator.go
├── docs/                  # API documentation
│   ├── swagger.json
│   └── swagger.yaml
├── main.go               # Application entry point
├── server.go             # Server setup and configuration
├── ddl.sql              # Database schema
└── .env                 # Environment configuration
```

## 🛠️ Installation & Setup

### Prerequisites

- Go 1.24.0 or higher
- PostgreSQL 12 or higher
- Git

### 1. Clone the Repository

```bash
git clone <repository-url>
cd server
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Database Setup

```bash
# Create PostgreSQL database
createdb develapar_blog_db

# Run the DDL script
psql -d develapar_blog_db -f ddl.sql
```

### 4. Environment Configuration

Copy the `.env` file and configure your environment variables:

```bash
cp .env.example .env
```

Edit `.env` with your configuration:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=develapar_blog_db
DB_USER=postgres
DB_PASSWORD=your_password

# Application Configuration
PORT_APP=:4300

# JWT Configuration
JWT_KEY=your_secret_key
JWT_LIFE_TIME=1
JWT_ISSUER_NAME=develapar

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_AUTHENTICATED_RPM=120
RATE_LIMIT_ANONYMOUS_RPM=30
```

### 5. Run the Server

```bash
# Development mode
go run main.go

# Build and run
go build -o develapar-server .
./develapar-server
```

The server will start on `http://localhost:4300`

## 📚 API Documentation

### Swagger UI

Access the interactive API documentation at:

```
http://localhost:4300/swagger/index.html
```

### Main Endpoints

#### Authentication

- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh JWT token

#### Users

- `GET /api/v1/users/` - Get all users
- `GET /api/v1/users/paginated` - Get users with pagination
- `GET /api/v1/users/:user_id` - Get user by ID

#### Articles

- `GET /api/v1/article/` - Get all articles
- `GET /api/v1/article/paginated` - Get articles with pagination
- `GET /api/v1/article/:slug` - Get article by slug
- `POST /api/v1/article/` - Create new article (auth required)
- `PUT /api/v1/article/:article_id` - Update article (auth required)
- `DELETE /api/v1/article/:article_id` - Delete article (auth required)

#### Categories

- `GET /api/v1/category/` - Get all categories
- `POST /api/v1/category/` - Create category (auth required)
- `PUT /api/v1/category/:cat_id` - Update category (auth required)
- `DELETE /api/v1/category/:cat_id` - Delete category (auth required)

#### Health & Monitoring

- `GET /api/v1/health` - Basic health check
- `GET /api/v1/health/detailed` - Detailed health information
- `GET /api/v1/metrics` - Application metrics
- `GET /api/v1/metrics/summary` - Metrics summary

## ⚙️ Configuration

### Environment Variables

#### Database Configuration

```env
DB_HOST=localhost                    # Database host
DB_PORT=5432                        # Database port
DB_NAME=develapar_blog_db           # Database name
DB_USER=postgres                    # Database user
DB_PASSWORD=your_password           # Database password
DB_DRIVER=postgres                  # Database driver

# Connection Pool Settings
DB_MAX_OPEN_CONNS=5      # Maximum open connections
DB_MAX_IDLE_CONNS=2               # Maximum idle connections
DB_CONN_MAX_LIFETIME=5m           # Connection maximum lifetime
DB_CONN_MAX_IDLE_TIME=5m          # Connection maximum idle time
DB_CONNECT_TIMEOUT=10s             # Connection timeout
DB_QUERY_TIMEOUT=10s               # Query timeout
```

#### Application Configuration

```env
PORT_APP=:4300                     # Application port

# JWT Configuration
JWT_KEY=your_secret_key            # JWT signing key
JWT_LIFE_TIME=1                    # JWT lifetime in hours
JWT_ISSUER_NAME=develapar          # JWT issuer name
```

#### Context Configuration

```env
CONTEXT_REQUEST_TIMEOUT=30s        # Maximum request processing time
CONTEXT_DATABASE_TIMEOUT=15s       # Maximum database operation time
CONTEXT_VALIDATION_TIMEOUT=5s      # Maximum validation processing time
CONTEXT_LOGGING_TIMEOUT=2s         # Maximum logging operation time
```

#### Logging Configuration

```env
LOG_LEVEL=info                     # Logging level (debug, info, warn, error)
LOG_FORMAT=json                    # Log format (json, text)
LOG_OUTPUT_PATH=stdout             # Log output path
LOG_ERROR_OUTPUT_PATH=stderr       # Error log output path
LOG_MAX_SIZE=100                   # Maximum log file size (MB)
LOG_MAX_BACKUPS=3                  # Number of backup log files
LOG_MAX_AGE=28                     # Maximum age of log files (days)
LOG_COMPRESS=true                  # Compress rotated logs
```

#### Rate Limiting Configuration

```env
RATE_LIMIT_ENABLED=true            # Enable/disable rate limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=60  # Default requests per minute
RATE_LIMIT_BURST_SIZE=10           # Burst size for rate limiting
RATE_LIMIT_CLEANUP_INTERVAL=5m     # Cleanup interval for rate limiter
RATE_LIMIT_WINDOW_SIZE=1m          # Sliding window size
RATE_LIMIT_AUTHENTICATED_RPM=120   # Rate limit for authenticated users
RATE_LIMIT_ANONYMOUS_RPM=30        # Rate limit for anonymous users
```

## 🔒 Security Features

### Authentication & Authorization

- **JWT Tokens**: Secure token-based authentication
- **Password Hashing**: Bcrypt password hashing
- **Role-based Access**: User role management
- **Token Refresh**: Secure token refresh mechanism

### Rate Limiting

- **Sliding Window Algorithm**: Advanced rate limiting
- **User-based Limits**: Different limits for authenticated/anonymous users
- **Path Exclusions**: Skip rate limiting for health checks and metrics
- **Automatic Cleanup**: Periodic cleanup of rate limit data

### Security Headers & CORS

- **CORS Configuration**: Configurable cross-origin resource sharing
- **Security Headers**: Standard security headers
- **Request Validation**: Input validation and sanitization

## 📊 Monitoring & Observability

### Structured Logging

- **JSON Format**: Machine-readable log format
- **Context Correlation**: Request ID and user ID tracking
- **Log Levels**: Configurable log levels (debug, info, warn, error)
- **Log Rotation**: Automatic log file rotation

### Metrics Collection

- **Request Metrics**: Response times, status codes, request counts
- **System Metrics**: Memory usage, CPU usage, goroutine counts
- **Database Metrics**: Connection pool statistics
- **Custom Metrics**: Business-specific metrics

### Health Checks

- **Basic Health**: Simple application health status
- **Detailed Health**: Comprehensive system information
- **Database Health**: Database connection and query testing
- **Dependency Health**: External service health checks

## 🧪 Testing

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Structure

- **Unit Tests**: Individual component testing
- **Integration Tests**: Component interaction testing
- **Mock Services**: Isolated testing with mocks
- **Test Coverage**: Comprehensive test coverage reporting

## 🚀 Deployment

### Build for Production

```bash
# Build optimized binary
go build -ldflags="-w -s" -o develapar-server .

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o develapar-server-linux .
```

### Docker Deployment

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o develapar-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/develapar-server .
CMD ["./develapar-server"]
```

### Environment-specific Configuration

- **Development**: Debug logging, detailed error messages
- **Staging**: Production-like settings with debug capabilities
- **Production**: Optimized performance, minimal logging

## 🔧 Development

### Code Style & Standards

- **Go Conventions**: Follow standard Go conventions
- **Clean Architecture**: Maintain separation of concerns
- **Error Handling**: Comprehensive error handling
- **Documentation**: Inline code documentation

### Adding New Features

1. **Model**: Define data structures in `model/`
2. **Repository**: Implement data access in `repository/`
3. **Service**: Add business logic in `service/`
4. **Controller**: Create HTTP handlers in `controller/`
5. **Routes**: Register routes in `server.go`
6. **Tests**: Add comprehensive tests
7. **Documentation**: Update API documentation

### Database Migrations

```bash
# Add new migration
echo "ALTER TABLE users ADD COLUMN phone VARCHAR(20);" >> migrations/001_add_phone.sql

# Apply migrations
psql -d develapar_blog_db -f migrations/001_add_phone.sql
```

## 📈 Performance Optimization

### Database Optimization

- **Connection Pooling**: Optimized connection pool settings
- **Query Optimization**: Efficient database queries
- **Indexing**: Proper database indexing
- **Pagination**: Efficient data pagination

### Application Optimization

- **Middleware Ordering**: Optimized middleware chain
- **Memory Management**: Efficient memory usage
- **Goroutine Management**: Proper concurrency handling
- **Caching**: Strategic caching implementation

## 🐛 Troubleshooting

### Common Issues

#### Database Connection Issues

```bash
# Check database connectivity
psql -h localhost -p 5432 -U postgres -d develapar_blog_db

# Verify environment variables
echo $DB_HOST $DB_PORT $DB_NAME
```

#### Port Already in Use

```bash
# Find process using port 4300
lsof -i :4300

# Kill process
kill -9 <PID>
```

#### JWT Token Issues

- Verify JWT_KEY is set correctly
- Check token expiration time
- Ensure proper token format in Authorization header

### Logging & Debugging

- Check application logs for detailed error information
- Use health check endpoints to verify system status
- Monitor metrics for performance issues
- Enable debug logging for development

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and best practices
- Write comprehensive tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting PR

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/) - HTTP web framework
- [PostgreSQL](https://www.postgresql.org/) - Database system
- [JWT-Go](https://github.com/golang-jwt/jwt) - JWT implementation
- [Swagger](https://swagger.io/) - API documentation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

---

**Develapar Blog Server** - Built with ❤️ using Go
