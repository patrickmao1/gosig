#!/bin/bash

CONTAINER_NAME="gosig-node1-1"
INTERVAL=1

prev_sent=0
prev_recv=0

convert_to_bytes() {
    local value=$1
    local number=$(echo $value | sed -E 's/([0-9.]+).*/\1/')
    local unit=$(echo $value | sed -E 's/[0-9.]+(.*)/\1/')

    case "$unit" in
        B)   multiplier=1 ;;
        kB)  multiplier=1024 ;;
        MB)  multiplier=$((1024 * 1024)) ;;
        GB)  multiplier=$((1024 * 1024 * 1024)) ;;
        TB)  multiplier=$((1024 * 1024 * 1024 * 1024)) ;;
        *)   multiplier=1 ;;  # fallback
    esac

    # bc handles floating point multiplication
    bytes=$(echo "$number * $multiplier" | bc)
    # Remove decimal part if any
    echo "${bytes%.*}"
}

while true; do
    stats=$(docker stats --no-stream --format "{{.Name}} {{.NetIO}}" | grep "$CONTAINER_NAME")

    sent_raw=$(echo $stats | awk '{print $2}')
    recv_raw=$(echo $stats | awk '{print $4}')

    sent_bytes=$(convert_to_bytes "$sent_raw")
    recv_bytes=$(convert_to_bytes "$recv_raw")

    if [ "$prev_sent" -ne 0 ]; then
        sent_diff=$((sent_bytes - prev_sent))
        recv_diff=$((recv_bytes - prev_recv))

        echo "$(date '+%H:%M:%S') - Sent: ${sent_diff} B/s, Received: ${recv_diff} B/s"
    fi

    prev_sent=$sent_bytes
    prev_recv=$recv_bytes

    sleep $INTERVAL
done