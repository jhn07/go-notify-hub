# NotifyHub - Notification Service

A microservice for sending notifications through multiple channels (Telegram, Email, SMS, WebPush) with PostgreSQL storage and RabbitMQ queue processing.

## 🚀 Quick Start

### Requirements
- Docker & Docker Compose
- Go 1.21+ (for local development)
- PostgreSQL 15+ (included in docker-compose)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/jhn07/go-notify-hub.git
cd go-notify-hub
```

2. Create `.env` file:
```bash
cp .env.example .env
```

3. Edit `.env` with your settings

4. Start services:
```bash
docker-compose up -d --build
```

## 📡 API Endpoints

### Health Check
```bash
GET http://localhost:8080/healthz
```

### Send Notification
```bash
POST http://localhost:8080/send
Content-Type: application/json

{
  "user_id": "user123",
  "message": "Hello from NotifyHub!",
  "channels": ["telegram", "email"]
}
```

**Response:**
```json
{
  "status": "queued",
  "message_id": "msg_a1b2c3d4...",
  "channels": ["telegram", "email"]
}
```

### Check Notification Status
```bash
GET http://localhost:8080/status/{message_id}
```

**Response:**
```json
{
  "status": "sent"
}
```

**Available statuses:**
- `queued` - Notification queued for processing
- `sending` - Currently being sent
- `sent` - Successfully sent to all channels
- `failed` - Failed to send to all channels
- `partial` - Sent to some channels, failed on others

## 📂 Project Structure

```
go-notify-hub/
├── cmd/
│   ├── server/              # API Server
│   │   └── main.go
│   └── worker/              # Queue Worker
│       └── main.go
├── internal/
│   ├── api/                 # HTTP handlers and middleware
│   │   ├── handlers.go
│   │   └── middleware.go
│   ├── channels/            # Notification channel implementations
│   │   ├── channel.go       # Channel interface and factory
│   │   ├── telegram.go      # Telegram channel
│   │   └── email.go         # Email channel
│   ├── db/                  # Database layer
│   │   └── db.go            # PostgreSQL connection and schema
│   ├── models/              # Data models
│   │   ├── notification.go  # Notification request/response
│   │   └── status.go        # Status constants
│   ├── queue/               # RabbitMQ integration
│   │   └── rabbitmq.go
│   └── service/             # Business logic
│       └── notifier.go
├── .env                     # Environment variables (not in git)
├── .env.example             # Environment variables template
├── .gitignore
├── .dockerignore
├── docker-compose.yml       # Docker Compose configuration
├── Dockerfile               # Multi-stage Docker build
├── go.mod                   # Go modules
├── go.sum
└── README.md
```

## 🏗️ Architecture

```
┌──────────┐      ┌───────────┐      ┌──────────┐      ┌────────┐      ┌──────────┐
│  Client  │─────▶│    API    │─────▶│PostgreSQL│      │RabbitMQ│─────▶│  Worker  │
└──────────┘      │  Server   │      │    DB    │      │ Queue  │      └─────┬────┘
                  └─────┬─────┘      └──────────┘      └────────┘            │
                        │                  ▲                                 │
                        └──────────────────┴─────────────────────────────────┘
                                           │
                                           ▼
                                     ┌─────────────┐
                                     │  Channels   │
                                     ├─────────────┤
                                     │  • Telegram │
                                     │  • Email    │
                                     │  • SMS      │
                                     │  • WebPush  │
                                     └─────────────┘
```

**Flow:**
1. Client sends notification request to API
2. API validates request and saves to PostgreSQL with status `queued`
3. API publishes message to RabbitMQ queue and returns 202 Accepted
4. Worker consumes message from queue
5. Worker updates status to `sending` in database
6. Worker sends notification through specified channels
7. Worker updates final status (`sent`, `failed`, or `partial`) in database
8. Client can check status via GET /status/{message_id}

## 🔧 Configuration

All settings are managed through `.env` file:

### API Configuration
- `API_PORT` - API server port (default: 8080)
- `API_ADDR` - API server address (default: :8080)
- `READ_TIMEOUT` - HTTP read timeout
- `WRITE_TIMEOUT` - HTTP write timeout
- `IDLE_TIMEOUT` - HTTP idle timeout

### RabbitMQ Configuration
- `RABBITMQ_URL` - RabbitMQ connection URL
- `RABBITMQ_HOST` - RabbitMQ host
- `RABBITMQ_PORT` - RabbitMQ port (default: 5672)
- `RABBITMQ_USER` - RabbitMQ username
- `RABBITMQ_PASSWORD` - RabbitMQ password

### PostgreSQL Configuration
- `DATABASE_URL` - PostgreSQL connection URL
- `POSTGRES_USER` - PostgreSQL username (default: postgres)
- `POSTGRES_PASSWORD` - PostgreSQL password (default: notifyhub)
- `POSTGRES_DB` - Database name (default: notifyhub)
- `POSTGRES_PORT` - PostgreSQL port (default: 5432)

### Worker Configuration
- `WORKER_COUNT` - Number of workers (default: 1)
- `QUEUE_NAME` - Queue name (default: notifyhub_queue)

### Channel Configuration
- `TELEGRAM_BOT_TOKEN` - Telegram bot token
- `EMAIL_SMTP_HOST` - SMTP server host
- `EMAIL_SMTP_PORT` - SMTP server port
- And more...

## 📊 Monitoring

### RabbitMQ Management UI
**URL:** http://localhost:15672  
**Default credentials:** guest / guest

### PostgreSQL Database
**Access database:**
```bash
docker exec -it notifyhub_postgres psql -U postgres -d notifyhub
```

**Useful queries:**
```sql
-- View all notifications
SELECT id, user_id, status, channels, created_at 
FROM notifications 
ORDER BY created_at DESC 
LIMIT 10;

-- Count by status
SELECT status, COUNT(*) 
FROM notifications 
GROUP BY status;

-- Recent failed notifications
SELECT * 
FROM notifications 
WHERE status = 'failed' 
ORDER BY created_at DESC 
LIMIT 5;
```

## 🛠️ Docker Commands

### Start services
```bash
docker-compose up -d --build
```

### Stop services
```bash
docker-compose down
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api
docker-compose logs -f worker
docker-compose logs -f rabbitmq
docker-compose logs -f postgres
```

### Check status
```bash
docker-compose ps
```

### Restart services
```bash
docker-compose restart
```

### Scale workers
```bash
# Run 3 workers
docker-compose up -d --scale worker=3

# Return to 1 worker
docker-compose up -d --scale worker=1
```

### Clean up
```bash
# Stop and remove containers
docker-compose down

# Stop and remove containers + volumes
docker-compose down -v
```

## 🧪 Testing

### Test the API
```bash
# Health check
curl http://localhost:8080/healthz

# Send notification
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "test_user",
    "message": "Test notification",
    "channels": ["telegram", "email"]
  }'

# Check notification status (replace with actual message_id)
curl http://localhost:8080/status/msg_abc123...
```

### Load testing
```bash
# Send 10 notifications simultaneously
for i in {1..10}; do
  curl -X POST http://localhost:8080/send \
    -H "Content-Type: application/json" \
    -d "{
      \"user_id\": \"user$i\",
      \"message\": \"Test message $i\",
      \"channels\": [\"telegram\", \"email\"]
    }" &
done
wait
```

## 🔌 Supported Channels

### Currently Implemented
- ✅ **Telegram** - Send messages via Telegram Bot API
- ✅ **Email** - Send emails via SMTP

### Planned
- 🔄 **SMS** - Send SMS messages
- 🔄 **WebPush** - Browser push notifications

## 🚀 Adding New Channels

1. Create new file in `internal/channels/`:
```go
// internal/channels/sms.go
package channels

type SMSChannel struct{}

func (c *SMSChannel) Send(userID, message string) error {
    // Implementation here
    return nil
}
```

2. Register in factory (`internal/channels/channel.go`):
```go
case "sms":
    return &SMSChannel{}, nil
```

3. Add to allowed channels in validation (`internal/api/handlers.go`)

## 🐛 Troubleshooting

### RabbitMQ connection issues
```bash
# Check if RabbitMQ is running
docker-compose ps rabbitmq

# Check RabbitMQ logs
docker-compose logs rabbitmq
```

### API not responding
```bash
# Check API logs
docker-compose logs api

# Restart API
docker-compose restart api
```

### Worker not processing messages
```bash
# Check worker logs
docker-compose logs worker

# Check RabbitMQ queue
# Visit http://localhost:15672 and check "Queues" tab
```

### Database connection issues
```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Test connection
docker exec -it notifyhub_postgres pg_isready -U postgres

# Restart PostgreSQL
docker-compose restart postgres
```

## 📝 Development

### Local development without Docker
```bash
# Start RabbitMQ and PostgreSQL
docker-compose up -d rabbitmq postgres

# Set environment variables
export DATABASE_URL="postgres://postgres:notifyhub@localhost:5432/notifyhub?sslmode=disable"
export RABBITMQ_URL="amqp://guest:guest@localhost:5672/"

# Run API server
go run cmd/server/main.go

# Run worker (in another terminal)
go run cmd/worker/main.go
```

### Run tests
```bash
go test -v ./...
```

### Build binaries
```bash
# Build server
go build -o bin/server cmd/server/main.go

# Build worker
go build -o bin/worker cmd/worker/main.go
```

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!
