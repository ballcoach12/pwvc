#!/bin/bash

# P-WVC Deployment Script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="${PROJECT_ROOT}/.env"
BACKUP_DIR="${PROJECT_ROOT}/backups"

# Functions
log() {
    echo -e "${BLUE}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

success() {
    echo -e "${GREEN}[SUCCESS] $1${NC}"
}

warn() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# Check if Docker and Docker Compose are installed
check_prerequisites() {
    log "Checking prerequisites..."
    
    if ! command -v docker &> /dev/null; then
        error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    
    success "Prerequisites check passed"
}

# Setup environment file
setup_environment() {
    log "Setting up environment configuration..."
    
    if [ ! -f "$ENV_FILE" ]; then
        if [ -f "${PROJECT_ROOT}/.env.docker" ]; then
            cp "${PROJECT_ROOT}/.env.docker" "$ENV_FILE"
            warn "Created .env file from .env.docker template"
            warn "Please review and update the configuration in .env file"
        else
            error ".env file not found and no template available"
            exit 1
        fi
    else
        success "Environment file exists"
    fi
}

# Create necessary directories
setup_directories() {
    log "Creating necessary directories..."
    
    mkdir -p "$BACKUP_DIR"
    mkdir -p "${PROJECT_ROOT}/logs"
    mkdir -p "${PROJECT_ROOT}/uploads"
    
    success "Directories created"
}

# Build Docker images
build_images() {
    log "Building Docker images..."
    
    cd "$PROJECT_ROOT"
    
    # Build backend
    log "Building backend image..."
    docker build -t pwvc-backend .
    
    # Build frontend
    log "Building frontend image..."
    docker build -t pwvc-frontend ./web
    
    success "Docker images built successfully"
}

# Start services
start_services() {
    log "Starting services..."
    
    cd "$PROJECT_ROOT"
    
    # Start infrastructure services first
    log "Starting database and cache services..."
    docker-compose up -d postgres redis
    
    # Wait for database to be ready
    log "Waiting for database to be ready..."
    sleep 10
    
    # Run migrations
    log "Running database migrations..."
    docker-compose run --rm migrate
    
    # Start application services
    log "Starting application services..."
    docker-compose up -d backend frontend
    
    success "All services started successfully"
}

# Check service health
check_health() {
    log "Checking service health..."
    
    cd "$PROJECT_ROOT"
    
    # Wait a bit for services to start
    sleep 15
    
    # Check backend health
    if curl -f http://localhost:8080/health &> /dev/null; then
        success "Backend is healthy"
    else
        error "Backend health check failed"
        return 1
    fi
    
    # Check frontend health
    if curl -f http://localhost:3000/health &> /dev/null; then
        success "Frontend is healthy"
    else
        error "Frontend health check failed"
        return 1
    fi
    
    success "All health checks passed"
}

# Create backup
create_backup() {
    log "Creating database backup..."
    
    cd "$PROJECT_ROOT"
    
    TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
    BACKUP_FILE="${BACKUP_DIR}/pwvc_backup_${TIMESTAMP}.sql"
    
    docker-compose exec -T postgres pg_dump -U pwvc pwvc > "$BACKUP_FILE"
    
    if [ -f "$BACKUP_FILE" ]; then
        success "Backup created: $BACKUP_FILE"
    else
        error "Backup creation failed"
        return 1
    fi
}

# Display status
show_status() {
    log "Service Status:"
    cd "$PROJECT_ROOT"
    docker-compose ps
    
    echo ""
    log "Application URLs:"
    echo "  Frontend: http://localhost:3000"
    echo "  Backend API: http://localhost:8080/api"
    echo "  Health Check: http://localhost:8080/health"
    
    if docker-compose ps | grep -q prometheus; then
        echo "  Prometheus: http://localhost:9090"
    fi
    
    if docker-compose ps | grep -q grafana; then
        echo "  Grafana: http://localhost:3001"
    fi
}

# Main deployment function
deploy() {
    log "Starting P-WVC deployment..."
    
    check_prerequisites
    setup_environment
    setup_directories
    build_images
    start_services
    
    if check_health; then
        success "Deployment completed successfully!"
        show_status
    else
        error "Deployment failed during health checks"
        log "Checking logs..."
        cd "$PROJECT_ROOT"
        docker-compose logs --tail=50
        exit 1
    fi
}

# Stop services
stop_services() {
    log "Stopping services..."
    cd "$PROJECT_ROOT"
    docker-compose down
    success "Services stopped"
}

# Restart services
restart_services() {
    log "Restarting services..."
    stop_services
    start_services
    check_health
    success "Services restarted successfully"
}

# Show logs
show_logs() {
    cd "$PROJECT_ROOT"
    if [ -n "$1" ]; then
        docker-compose logs -f "$1"
    else
        docker-compose logs -f
    fi
}

# Update application
update() {
    log "Updating application..."
    
    # Create backup before update
    create_backup
    
    # Pull latest changes (if this is a git deployment)
    if [ -d "${PROJECT_ROOT}/.git" ]; then
        log "Pulling latest changes..."
        cd "$PROJECT_ROOT"
        git pull
    fi
    
    # Rebuild and restart
    build_images
    restart_services
    
    success "Update completed successfully"
}

# Main script logic
case "${1:-deploy}" in
    "deploy")
        deploy
        ;;
    "stop")
        stop_services
        ;;
    "start")
        start_services
        check_health
        ;;
    "restart")
        restart_services
        ;;
    "status")
        show_status
        ;;
    "logs")
        show_logs "$2"
        ;;
    "backup")
        create_backup
        ;;
    "update")
        update
        ;;
    "build")
        build_images
        ;;
    "health")
        check_health
        ;;
    *)
        echo "Usage: $0 {deploy|start|stop|restart|status|logs|backup|update|build|health}"
        echo ""
        echo "Commands:"
        echo "  deploy  - Full deployment (default)"
        echo "  start   - Start services"
        echo "  stop    - Stop services"
        echo "  restart - Restart services"
        echo "  status  - Show service status"
        echo "  logs    - Show logs (optionally for specific service)"
        echo "  backup  - Create database backup"
        echo "  update  - Update and restart application"
        echo "  build   - Build Docker images"
        echo "  health  - Check service health"
        echo ""
        echo "Examples:"
        echo "  $0 deploy          # Full deployment"
        echo "  $0 logs backend    # Show backend logs"
        echo "  $0 restart         # Restart all services"
        exit 1
        ;;
esac