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

# --- Functions ---
kill_port() {
    local port=$1
    local pids=$(lsof -ti :$port)
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

# Auth
(cd services/auth && go run cmd/api/main.go 2>&1 | sed "s/^/${BLUE}[AUTH]${NC} /") &
# Booking
(cd services/booking && go run main.go 2>&1 | sed "s/^/${GREEN}[BOOKING]${NC} /") &
# Payment
(cd services/payment && go run main.go 2>&1 | sed "s/^/${YELLOW}[PAYMENT]${NC} /") &
# Location
(cd services/location && go run main.go 2>&1 | sed "s/^/${CYAN}[LOCATION]${NC} /") &
# Notifications
(cd services/notifications && go run main.go 2>&1 | sed "s/^/${MAGENTA}[NOTIF]${NC} /") &
# Provider
(cd services/provider && go run main.go 2>&1 | sed "s/^/${BLUE}[PROVIDER]${NC} /") &

echo -e "${BLUE}⌨️  Press Ctrl+C to stop all services.${NC}"

# Keep script running
while true; do
    sleep 1
done
