#!/bin/bash

echo "🚀 Setting up P-WVC development environment..."

# Update package lists
sudo apt-get update

# Install additional tools
echo "📦 Installing additional development tools..."
sudo apt-get install -y \
    curl \
    wget \
    unzip \
    jq \
    tree \
    htop \
    postgresql-client \
    make \
    build-essential

# Install Go tools
echo "🔧 Installing Go development tools..."
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/ramya-rao-a/go-outline@latest
go install -v github.com/cweill/gotests/gotests@latest
go install -v github.com/fatih/gomodifytags@latest
go install -v github.com/josharian/impl@latest
go install -v github.com/haya14busa/goplay/cmd/goplay@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
go install -v honnef.co/go/tools/cmd/staticcheck@latest

# Set up Go environment
echo "🏗️ Configuring Go environment..."
echo 'export PATH=$PATH:/go/bin:$HOME/go/bin' >> ~/.bashrc
echo 'export GOPATH=/go' >> ~/.bashrc
echo 'export GOROOT=/usr/local/go' >> ~/.bashrc

# Install npm global packages for frontend development
echo "📱 Installing Node.js development tools..."
npm install -g \
    @types/node \
    typescript \
    @typescript-eslint/eslint-plugin \
    @typescript-eslint/parser \
    prettier \
    create-react-app \
    @vitejs/create-app

# Set up project-specific aliases and functions
echo "⚡ Setting up development aliases..."
cat >> ~/.bashrc << 'EOF'

# P-WVC Development Aliases
alias ll='ls -alF'
alias la='ls -A'
alias l='ls -CF'
alias ..='cd ..'
alias ...='cd ../..'

# Go shortcuts
alias gor='go run .'
alias gob='go build .'
alias got='go test ./...'
alias gof='go fmt ./...'
alias gov='go vet ./...'
alias gom='go mod tidy'

# Git shortcuts
alias gs='git status'
alias ga='git add'
alias gc='git commit'
alias gp='git push'
alias gl='git log --oneline'
alias gd='git diff'

# P-WVC specific shortcuts
alias start-db='docker run --name pwvc-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=pwvc -p 5432:5432 -d postgres:15-alpine'
alias stop-db='docker stop pwvc-postgres && docker rm pwvc-postgres'
alias logs-db='docker logs pwvc-postgres'

# Development functions
pwvc-setup() {
    echo "🏗️ Setting up P-WVC development environment..."
    
    # Initialize Go module if not exists
    if [ ! -f "go.mod" ]; then
        go mod init pwvc
    fi
    
    # Create basic project structure
    mkdir -p {cmd/server,internal/{api,domain,repository,service,websocket},pkg,migrations,web/src}
    
    echo "✅ P-WVC project structure created!"
}

pwvc-run() {
    echo "🚀 Starting P-WVC application..."
    go run ./cmd/server
}

pwvc-test() {
    echo "🧪 Running P-WVC tests..."
    go test -v ./...
}

pwvc-build() {
    echo "📦 Building P-WVC application..."
    go build -o pwvc ./cmd/server
}

pwvc-dev() {
    echo "🔄 Starting P-WVC in development mode..."
    # You can add auto-reload functionality here later
    go run ./cmd/server
}

EOF

# Create development database setup script
echo "🗄️ Creating database setup script..."
cat > /workspaces/pwvc/setup-db.sh << 'EOF'
#!/bin/bash

echo "🗄️ Setting up P-WVC PostgreSQL database..."

# Start PostgreSQL container
docker run --name pwvc-postgres \
    -e POSTGRES_PASSWORD=password \
    -e POSTGRES_DB=pwvc \
    -e POSTGRES_USER=pwvc \
    -p 5432:5432 \
    -d postgres:15-alpine

echo "⏳ Waiting for database to start..."
sleep 5

# Check if database is ready
while ! docker exec pwvc-postgres pg_isready -U pwvc; do
    echo "Waiting for database connection..."
    sleep 2
done

echo "✅ Database is ready!"
echo "📍 Connection string: postgres://pwvc:password@localhost:5432/pwvc?sslmode=disable"

EOF

chmod +x /workspaces/pwvc/setup-db.sh

# Create a useful development README
echo "📖 Creating development README..."
cat > /workspaces/pwvc/README.dev.md << 'EOF'
# P-WVC Development Environment

## Quick Start

### 1. Setup Database
```bash
./setup-db.sh
```

### 2. Initialize Project Structure
```bash
pwvc-setup
```

### 3. Run Application
```bash
pwvc-run
```

## Useful Commands

### Go Development
- `gor` - Run the application
- `gob` - Build the application  
- `got` - Run tests
- `gof` - Format code
- `gov` - Vet code
- `gom` - Tidy modules

### Git Shortcuts
- `gs` - Git status
- `ga .` - Add all files
- `gc -m "message"` - Commit with message
- `gp` - Push to origin

### Database Management
- `start-db` - Start PostgreSQL container
- `stop-db` - Stop and remove container
- `logs-db` - View database logs

### P-WVC Functions
- `pwvc-setup` - Initialize project structure
- `pwvc-run` - Start the application
- `pwvc-test` - Run all tests
- `pwvc-build` - Build executable
- `pwvc-dev` - Development mode with auto-reload

## Development Ports
- **8080** - Go Backend API
- **3000** - React Frontend
- **5432** - PostgreSQL Database

## Environment Variables
```bash
DATABASE_URL=postgres://pwvc:password@localhost:5432/pwvc?sslmode=disable
PORT=8080
GIN_MODE=debug
```

EOF

echo "✅ P-WVC development environment setup complete!"
echo "🎉 You can now use bash commands and P-WVC shortcuts!"
echo ""
echo "Next steps:"
echo "1. Reload your bash profile: source ~/.bashrc"
echo "2. Setup database: ./setup-db.sh"
echo "3. Initialize project: pwvc-setup"
echo "4. Start developing! 🚀"