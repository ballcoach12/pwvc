# PairWise Application Deployment Guide

## Overview

This guide provides comprehensive instructions for deploying the PairWise (Pairwise-Weighted Value/Complexity) application in various environments.

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Local Development](#local-development)
3. [Docker Deployment](#docker-deployment)
4. [Production Deployment](#production-deployment)
5. [Environment Configuration](#environment-configuration)
6. [Database Setup](#database-setup)
7. [Monitoring and Logging](#monitoring-and-logging)
8. [Troubleshooting](#troubleshooting)

## Prerequisites

### System Requirements

- **Go**: Version 1.23.3 or later
- **Node.js**: Version 18 or later
- **PostgreSQL**: Version 13 or later
- **Redis**: Version 6 or later (optional for production scaling)
- **Docker**: Version 20 or later (for containerized deployment)
- **Docker Compose**: Version 2.0 or later

### Hardware Requirements

- **Minimum**: 2 CPU cores, 4GB RAM, 20GB storage
- **Recommended**: 4 CPU cores, 8GB RAM, 50GB storage
- **Production**: 8 CPU cores, 16GB RAM, 100GB storage

## Local Development

### 1. Clone and Setup

```bash
git clone <repository-url>
cd pairwise
```

### 2. Backend Setup

```bash
# Install Go dependencies
go mod download
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your local configuration

# Run database migrations
go run cmd/migrate/main.go
```

### 3. Frontend Setup

```bash
cd web
npm install
# or
yarn install
```

### 4. Start Development Servers

#### Backend

```bash
# Start Go backend
go run cmd/server/main.go
# Server runs on http://localhost:8080
```

#### Frontend

```bash
cd web
npm run dev
# Frontend runs on http://localhost:5173
```

### 5. Run Tests

```bash
# Backend tests
go test ./...

# Frontend tests
cd web
npm test
```

## Docker Deployment

### 1. Quick Start with Docker Compose

```bash
# Build and start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### 2. Production Docker Compose

```bash
# Use production compose file
docker-compose -f docker-compose.prod.yml up -d
```

### 3. Individual Container Deployment

#### Backend Container

```bash
# Build backend image
docker build -t pairwise-backend .

# Run backend container
docker run -d \
  --name pairwise-backend \
  -p 8080:8080 \
  -e DATABASE_URL=postgres://user:pass@host:5432/pairwise \
  -e REDIS_URL=redis://redis:6379 \
  pairwise-backend
```

#### Frontend Container

```bash
# Build frontend image
docker build -t pairwise-frontend ./web

# Run frontend container
docker run -d \
  --name pairwise-frontend \
  -p 80:80 \
  pairwise-frontend
```

## Production Deployment

### 1. AWS ECS Deployment

```yaml
# ecs-task-definition.json
{
  "family": "pairwise-app",
  "networkMode": "awsvpc",
  "requiresCompatibilities": ["FARGATE"],
  "cpu": "1024",
  "memory": "2048",
  "containerDefinitions":
    [
      {
        "name": "pairwise-backend",
        "image": "your-registry/pairwise-backend:latest",
        "portMappings": [{ "containerPort": 8080, "protocol": "tcp" }],
        "environment":
          [
            {
              "name": "DATABASE_URL",
              "value": "postgres://user:pass@rds-endpoint:5432/pairwise",
            },
          ],
      },
    ],
}
```

### 2. Kubernetes Deployment

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pairwise-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pairwise-backend
  template:
    metadata:
      labels:
        app: pairwise-backend
    spec:
      containers:
        - name: pairwise-backend
          image: pairwise-backend:latest
          ports:
            - containerPort: 8080
          env:
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: pairwise-secrets
                  key: database-url
---
apiVersion: v1
kind: Service
metadata:
  name: pairwise-backend-service
spec:
  selector:
    app: pairwise-backend
  ports:
    - port: 80
      targetPort: 8080
  type: LoadBalancer
```

### 3. Traditional Server Deployment

```bash
# 1. Install dependencies on server
sudo apt update
sudo apt install postgresql redis-server nginx

# 2. Create application user
sudo useradd -r -s /bin/false pairwise
sudo mkdir /opt/pairwise
sudo chown pairwise:pairwise /opt/pairwise

# 3. Deploy application
sudo -u pairwise cp pairwise-binary /opt/pairwise/
sudo -u pairwise cp -r web/dist /opt/pairwise/static/

# 4. Configure systemd service
sudo cp scripts/pairwise.service /etc/systemd/system/
sudo systemctl enable pairwise
sudo systemctl start pairwise

# 5. Configure nginx
sudo cp scripts/nginx.conf /etc/nginx/sites-available/pairwise
sudo ln -s /etc/nginx/sites-available/pairwise /etc/nginx/sites-enabled/
sudo systemctl reload nginx
```

## Environment Configuration

### Backend Environment Variables

```env
# Database Configuration
DATABASE_URL=postgres://username:password@localhost:5432/pairwise
DATABASE_MAX_CONNECTIONS=25
DATABASE_MAX_IDLE_CONNECTIONS=10

# Redis Configuration (Optional)
REDIS_URL=redis://localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

# Server Configuration
PORT=8080
GIN_MODE=release
LOG_LEVEL=info
LOG_FORMAT=json

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:5173,https://yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# WebSocket Configuration
WS_READ_BUFFER_SIZE=1024
WS_WRITE_BUFFER_SIZE=1024
WS_HEARTBEAT_INTERVAL=30s

# Monitoring
PROMETHEUS_ENABLED=true
PROMETHEUS_PORT=9090
HEALTH_CHECK_TIMEOUT=30s

# Security
JWT_SECRET=your-jwt-secret-here
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

### Frontend Environment Variables

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api
VITE_WS_BASE_URL=ws://localhost:8080/ws

# Application Configuration
VITE_APP_TITLE=PairWise Application
VITE_APP_VERSION=1.0.0
VITE_APP_ENVIRONMENT=production

# Feature Flags
VITE_ENABLE_ANALYTICS=true
VITE_ENABLE_DEBUGGING=false
```

## Database Setup

### 1. PostgreSQL Installation and Configuration

```sql
-- Create database and user
CREATE DATABASE pairwise;
CREATE USER pairwise_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE pairwise TO pairwise_user;

-- Connect to pairwise database
\c pairwise

-- Grant schema permissions
GRANT ALL ON SCHEMA public TO pairwise_user;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO pairwise_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO pairwise_user;
```

### 2. Run Migrations

```bash
# Run database migrations
go run cmd/migrate/main.go

# Or using make command
make migrate-up
```

### 3. Seed Data (Optional)

```bash
# Load sample data
psql -h localhost -U pairwise_user -d pairwise -f scripts/sample-data.sql
```

### 4. Database Backup and Restore

```bash
# Backup database
pg_dump -h localhost -U pairwise_user -d pairwise > backup.sql

# Restore database
psql -h localhost -U pairwise_user -d pairwise < backup.sql
```

## Monitoring and Logging

### 1. Prometheus Configuration

```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: "pairwise-backend"
    static_configs:
      - targets: ["localhost:8080"]
    metrics_path: /metrics
    scrape_interval: 5s
```

### 2. Grafana Dashboard

- Import dashboard from `monitoring/grafana-dashboard.json`
- Configure data source: Prometheus server URL
- Set up alerts for critical metrics

### 3. Log Management

```bash
# View application logs
docker-compose logs -f pairwise-backend

# Or for systemd
sudo journalctl -u pairwise -f

# Log rotation configuration
sudo cat > /etc/logrotate.d/pairwise << EOF
/var/log/pairwise/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 pairwise pairwise
}
EOF
```

## Troubleshooting

### Common Issues

#### 1. Database Connection Issues

```bash
# Check database connectivity
pg_isready -h localhost -p 5432

# Test connection with credentials
psql -h localhost -U pairwise_user -d pairwise -c "SELECT 1;"

# Check database logs
sudo tail -f /var/log/postgresql/postgresql-*.log
```

#### 2. Port Conflicts

```bash
# Check port usage
sudo netstat -tulpn | grep :8080
sudo lsof -i :8080

# Kill process using port
sudo kill -9 $(sudo lsof -t -i:8080)
```

#### 3. Frontend Build Issues

```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Check for build errors
npm run build 2>&1 | tee build.log
```

#### 4. Docker Issues

```bash
# Check container logs
docker logs pairwise-backend

# Check container resource usage
docker stats

# Clean up Docker resources
docker system prune -a
```

### Performance Tuning

#### Database Optimization

```sql
-- Add indexes for frequently queried columns
CREATE INDEX idx_projects_created_at ON projects(created_at);
CREATE INDEX idx_features_project_id ON features(project_id);
CREATE INDEX idx_attendees_project_id ON attendees(project_id);
CREATE INDEX idx_pairwise_sessions_project_id ON pairwise_sessions(project_id);

-- Analyze tables for query optimization
ANALYZE projects;
ANALYZE features;
ANALYZE attendees;
```

#### Application Optimization

```bash
# Enable Go runtime optimizations
export GOMAXPROCS=$(nproc)
export GOGC=100

# Configure connection pooling
export DATABASE_MAX_CONNECTIONS=50
export DATABASE_MAX_IDLE_CONNECTIONS=25
```

### Health Checks

#### Backend Health Check

```bash
curl -f http://localhost:8080/health || exit 1
```

#### Database Health Check

```bash
pg_isready -h localhost -p 5432 -U pairwise_user
```

#### Frontend Health Check

```bash
curl -f http://localhost:80/ || exit 1
```

### Backup and Recovery

#### Automated Backup Script

```bash
#!/bin/bash
# scripts/backup.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/pairwise/backups"
DB_NAME="pairwise"
DB_USER="pairwise_user"

# Create backup directory
mkdir -p $BACKUP_DIR

# Database backup
pg_dump -h localhost -U $DB_USER -d $DB_NAME > $BACKUP_DIR/pairwise_$DATE.sql

# Compress backup
gzip $BACKUP_DIR/pairwise_$DATE.sql

# Remove backups older than 30 days
find $BACKUP_DIR -name "*.sql.gz" -mtime +30 -delete

echo "Backup completed: pairwise_$DATE.sql.gz"
```

#### Schedule Backups with Cron

```bash
# Add to crontab
0 2 * * * /opt/pairwise/scripts/backup.sh
```

## Security Considerations

### 1. SSL/TLS Configuration

```nginx
server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /path/to/certificate.crt;
    ssl_certificate_key /path/to/private.key;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 2. Firewall Configuration

```bash
# UFW configuration
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw deny 8080/tcp   # Block direct backend access
sudo ufw enable
```

### 3. Database Security

```sql
-- Disable unnecessary extensions
-- Review and limit user permissions
REVOKE ALL ON SCHEMA public FROM public;
GRANT USAGE ON SCHEMA public TO pairwise_user;
```

## Support and Maintenance

### Regular Maintenance Tasks

1. **Weekly**: Review logs, check disk space, update dependencies
2. **Monthly**: Database maintenance, backup verification, security updates
3. **Quarterly**: Performance review, capacity planning, disaster recovery testing

### Monitoring Checklist

- [ ] Application response time < 500ms
- [ ] Database connection pool < 80% utilization
- [ ] Disk usage < 80%
- [ ] Memory usage < 90%
- [ ] Error rate < 1%
- [ ] Backup completion success

### Support Resources

- **Documentation**: `/docs` directory
- **Issue Tracking**: GitHub Issues
- **Monitoring**: Grafana dashboards
- **Logs**: Application and system logs
- **Health Checks**: `/health` endpoint
