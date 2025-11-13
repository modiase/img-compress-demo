#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
NC='\033[0m'

log_info() {
    echo -e "$(date '+%H:%M:%S') | ${BLUE}INFO${NC}  | $*"
}

log_success() {
    echo -e "$(date '+%H:%M:%S') | ${GREEN}OK${NC}    | $*"
}

log_warn() {
    echo -e "$(date '+%H:%M:%S') | ${YELLOW}WARN${NC}  | $*"
}

log_error() {
    echo -e "$(date '+%H:%M:%S') | ${RED}ERROR${NC} | $*"
}

_pad_center() {
    local str="$1"
    local width="$2"
    local str_len=${#str}
    local pad_total=$((width - str_len))

    if [ $pad_total -le 0 ]; then
        echo "$str"
        return
    fi

    local pad_left=$((pad_total / 2))
    local pad_right=$((pad_total - pad_left))

    printf "%*s%s%*s" $pad_left "" "$str" $pad_right ""
}
