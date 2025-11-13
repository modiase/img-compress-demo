#!/usr/bin/env -S bash -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "$SCRIPT_DIR/lib.sh"

PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
mkdir -p "$PROJECT_ROOT/logs"

cleanup() {
    log_info "Shutting down development servers..."

    if [ ! -z "${FRONTEND_PID:-}" ] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
        log_info "Stopping frontend server (PID: $FRONTEND_PID)"
        kill -TERM "$FRONTEND_PID" 2>/dev/null || true
    fi

    if [ ! -z "${BACKEND_PID:-}" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        log_info "Stopping backend server (PID: $BACKEND_PID)"
        kill -TERM "$BACKEND_PID" 2>/dev/null || true
    fi

    local count=0
    while [ $count -lt 10 ]; do
        local still_running=false

        if [ ! -z "${FRONTEND_PID:-}" ] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
            still_running=true
        fi

        if [ ! -z "${BACKEND_PID:-}" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
            still_running=true
        fi

        if [ "$still_running" = false ]; then
            break
        fi

        sleep 0.5
        count=$((count + 1))
    done

    if [ ! -z "${FRONTEND_PID:-}" ] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
        log_warn "Force killing frontend server"
        kill -KILL "$FRONTEND_PID" 2>/dev/null || true
    fi

    if [ ! -z "${BACKEND_PID:-}" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        log_warn "Force killing backend server"
        kill -KILL "$BACKEND_PID" 2>/dev/null || true
    fi

    log_success "Development servers stopped"
    exit 0
}

trap cleanup SIGINT SIGTERM EXIT

log_info "Starting development servers..."

log_info "Starting backend server on port 8080..."
(cd "$PROJECT_ROOT/backend" && go run main.go 2>&1 | while IFS= read -r line; do
    padded_label="$(_pad_center "backend" 20)"
    padded_level="$(_pad_center "info" 7)"
    echo -e "$(date '+%H:%M:%S') | ${PURPLE}${padded_label}${NC} | ${WHITE}${padded_level}${NC} | $line"
done | tee "$PROJECT_ROOT/logs/backend.log") &
BACKEND_PID=$!

log_info "Waiting for backend server to be ready..."
while true; do
    if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
        log_error "Backend server process died"
        if [ -f "$PROJECT_ROOT/logs/backend.log" ]; then
            cat "$PROJECT_ROOT/logs/backend.log"
        fi
        exit 1
    fi

    if curl -s http://localhost:8080/api/compress > /dev/null 2>&1; then
        log_success "Backend server is ready"
        break
    fi

    sleep 0.5
done

log_info "Starting frontend server on port 5173..."
(cd "$PROJECT_ROOT/frontend" && pnpm exec vite dev 2>&1 | while IFS= read -r line; do
    padded_label="$(_pad_center "frontend" 20)"
    padded_level="$(_pad_center "info" 7)"
    clean_line="$(echo "$line" | sed 's/^[0-9][0-9]:[0-9][0-9]:[0-9][0-9] [AP]M //')"
    echo -e "$(date '+%H:%M:%S') | ${CYAN}${padded_label}${NC} | ${WHITE}${padded_level}${NC} | $clean_line"
done | tee "$PROJECT_ROOT/logs/frontend.log") &
FRONTEND_PID=$!

log_info "Waiting for frontend server to be ready..."
while true; do
    if ! kill -0 "$FRONTEND_PID" 2>/dev/null; then
        log_error "Frontend server process died"
        if [ -f "$PROJECT_ROOT/logs/frontend.log" ]; then
            cat "$PROJECT_ROOT/logs/frontend.log"
        fi
        exit 1
    fi

    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        log_success "Frontend server is ready"
        break
    fi

    sleep 0.5
done

log_success "Development servers are running!"
log_info "Frontend: http://localhost:5173"
log_info "Backend: http://localhost:8080"
log_info "Logs: logs/frontend.log and logs/backend.log"
log_info "Press Ctrl+C to stop all servers"

wait
