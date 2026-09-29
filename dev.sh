#!/bin/bash
# Care4u Local Development Orchestrator
# Part of the Infortts Swarm Ecosystem.

# --- Configuration ---
GREEN='\x1b[32m'
BLUE='\x1b[34m'
RED='\x1b[31m'
CYAN='\x1b[36m'
YELLOW='\x1b[33m'
MAGENTA='\x1b[35m'
NC='\x1b[0m' # No Color

PROJECT_NAME="Care4u"

# Ensure we run from the project root regardless of the caller's CWD
cd "$(dirname "$0")" || exit 1

SERVICES="auth booking payment location notifications provider"

# --- Helpers ---

# load_env applies .env for every variable that is NOT already exported, so a
# one-off override (e.g. `JWT_SECRET=... ./dev.sh`) still wins over the file.
load_env() {
    if [ ! -f .env ]; then
        echo -e "${YELLOW}⚠️  No .env found — using built-in dev defaults (see .env.example)${NC}"
        return 0
    fi
    while IFS= read -r line || [ -n "$line" ]; do
        case "$line" in
            ''|\#*) continue ;;
        esac
        key="${line%%=*}"
        [ -z "$key" ] && continue
        [ -n "${!key:-}" ] && continue
        export "$key=${line#*=}"
    done < .env
}

# dsn builds a host-reachable Postgres DSN. The service defaults point at the
# compose network name `db`, which does not resolve for a host-run `go run`.
dsn() {
    echo "postgres://${POSTGRES_USER:-user}:${POSTGRES_PASSWORD:-pass}@127.0.0.1:${POSTGRES_PORT:-5432}/$1?sslmode=disable"
}

require_go() {
    if ! command -v go > /dev/null 2>&1; then
        echo -e "${RED}❌ Go toolchain not found — install Go 1.25+ to run the backend.${NC}"
        return 1
    fi
    return 0
}

require_docker() {
    if ! docker compose version > /dev/null 2>&1; then
        echo -e "${RED}❌ docker compose v2 is required (infra: Postgres, Redis, NATS).${NC}"
        return 1
    fi
    return 0
}

# --- Command Parsing ---
if [ "$1" = "clean" ]; then
    echo -e "${YELLOW}🧹 Cleaning up Care4u...${NC}"
    require_docker || exit 1
    docker compose down --remove-orphans
    for port in 8080 8081 8082 8083 8084 8085; do
        kill_port $port
    done
    pkill -f "go run" 2>/dev/null
    echo -e "${GREEN}✅ Care4u stopped.${NC}"
    exit 0
fi

if [ "$1" = "install" ]; then
    echo -e "${YELLOW}📦 Installing Care4u dependencies...${NC}"
    require_go || exit 1
    require_docker || exit 1
    for svc in $SERVICES; do
        echo -e "  → services/${svc}"
        (cd "services/${svc}" && go mod download) || {
            echo -e "${RED}❌ go mod download failed in services/${svc}${NC}"
            exit 1
        }
    done
    docker compose pull db redis nats
    echo -e "${GREEN}✅ Dependencies installed.${NC}"
    exit 0
fi


# --- Functions ---
kill_port() {
    local port=$1
    local pids=""
    if command -v lsof > /dev/null 2>&1; then
        pids=$(lsof -ti :$port 2>/dev/null)
    elif command -v fuser > /dev/null 2>&1; then
        pids=$(fuser $port/tcp 2>/dev/null)
    fi
    if [ ! -z "$pids" ]; then
        echo -e " -> Clearing port $port (PIDs: $pids)"
        echo "$pids" | xargs kill -9 2>/dev/null
        sleep 1
    fi
}

cleanup() {
    echo -e "\n${RED}🛑 Shutting down Care4u services...${NC}"
    docker compose down
    # Kill go run processes
    pkill -f "go run" 2>/dev/null
    exit
}

trap cleanup SIGINT SIGTERM

# --- Initialization ---
echo -e "${CYAN}🚀 Manifesting ${PROJECT_NAME} Environment...${NC}"

require_go || exit 1
require_docker || exit 1
load_env

# The services below run on the HOST, so they must reach the infra over the
# loopback-mapped compose ports (see docker-compose.yml), not the compose DNS.
export NATS_URL="${NATS_URL:-nats://127.0.0.1:4222}"
export REDIS_URL="${REDIS_URL:-127.0.0.1:6379}"

# Free ports
echo -e "${YELLOW}🧹 Cleaning up old processes and freeing ports...${NC}"
for port in 8080 8081 8082 8083 8084 8085; do
    kill_port $port
done

# 1. Start Docker Services (Postgres, Redis, NATS)
echo -e "${CYAN}[DOCKER]${NC} Starting infrastructure (DB, Redis, NATS)..."
# We only want the infra, not the build services if we run them locally
docker compose up -d db redis nats

# 2. Start Go Services
echo -e "${MAGENTA}[SERVICES]${NC} Launching Go Microservices..."

# Auth (entrypoint lives in cmd/api, the rest sit at the module root)
(cd services/auth && DATABASE_URL="$(dsn auth_db)" PORT=8080 go run cmd/api/main.go 2>&1 | sed "s/^/${BLUE}[AUTH]${NC} /") &
# Booking
(cd services/booking && DATABASE_URL="$(dsn booking_db)" PORT=8081 go run main.go 2>&1 | sed "s/^/${GREEN}[BOOKING]${NC} /") &
# Payment
(cd services/payment && DATABASE_URL="$(dsn payment_db)" PORT=8082 go run main.go 2>&1 | sed "s/^/${YELLOW}[PAYMENT]${NC} /") &
# Location
(cd services/location && PORT=8083 go run main.go 2>&1 | sed "s/^/${CYAN}[LOCATION]${NC} /") &
# Notifications
(cd services/notifications && PORT=8084 go run main.go 2>&1 | sed "s/^/${MAGENTA}[NOTIF]${NC} /") &
# Provider
(cd services/provider && DATABASE_URL="$(dsn provider_db)" PORT=8085 go run main.go 2>&1 | sed "s/^/${BLUE}[PROVIDER]${NC} /") &

echo -e "${BLUE}⌨️  Press Ctrl+C to stop all services.${NC}"

# Keep script running
while true; do
    sleep 1
done
