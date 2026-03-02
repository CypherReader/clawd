#!/bin/bash

# 6-Hour Data Collection Script
# Runs in background and saves all output

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

LOG_DIR="$SCRIPT_DIR/logs"
mkdir -p "$LOG_DIR"

TIMESTAMP=$(date +%Y-%m-%d_%H-%M-%S)
LOG_FILE="$LOG_DIR/collection_${TIMESTAMP}.log"
PID_FILE="$LOG_DIR/collector.pid"
STATUS_FILE="$LOG_DIR/collector_status.json"

echo "=== Starting Cross-Platform Data Collection ===" | tee "$LOG_FILE"
echo "Start Time: $(date)" | tee -a "$LOG_FILE"
echo "Log File: $LOG_FILE" | tee -a "$LOG_FILE"
echo "Duration: 6 hours" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"

# Check if already running
if [ -f "$PID_FILE" ]; then
    OLD_PID=$(cat "$PID_FILE")
    if ps -p "$OLD_PID" > /dev/null 2>&1; then
        echo "⚠️  Data collector already running (PID: $OLD_PID)" | tee -a "$LOG_FILE"
        echo "To stop: kill $OLD_PID" | tee -a "$LOG_FILE"
        exit 1
    else
        echo "Removing stale PID file..." | tee -a "$LOG_FILE"
        rm "$PID_FILE"
    fi
fi

# Start collector in background
cd "$SCRIPT_DIR/cmd/data-collector"
nohup go run main.go >> "$LOG_FILE" 2>&1 &
COLLECTOR_PID=$!

echo $COLLECTOR_PID > "$PID_FILE"

# Write status file
cat > "$STATUS_FILE" << EOF
{
  "status": "running",
  "pid": $COLLECTOR_PID,
  "start_time": "$(date -Iseconds)",
  "end_time": "$(date -Iseconds -d '+6 hours')",
  "log_file": "$LOG_FILE",
  "data_dir": "$SCRIPT_DIR/cmd/data-collector/data"
}
EOF

echo "✅ Data collector started!" | tee -a "$LOG_FILE"
echo "PID: $COLLECTOR_PID" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"
echo "Monitoring commands:" | tee -a "$LOG_FILE"
echo "  Check status: cat $STATUS_FILE" | tee -a "$LOG_FILE"
echo "  View progress: tail -f $LOG_FILE" | tee -a "$LOG_FILE"
echo "  Stop collection: kill $COLLECTOR_PID" | tee -a "$LOG_FILE"
echo "" | tee -a "$LOG_FILE"
echo "Results will be available in ~6 hours at:" | tee -a "$LOG_FILE"
echo "  $SCRIPT_DIR/cmd/data-collector/data/" | tee -a "$LOG_FILE"

# Monitor process for a few seconds
sleep 3

if ps -p "$COLLECTOR_PID" > /dev/null 2>&1; then
    echo "✅ Process confirmed running" | tee -a "$LOG_FILE"
else
    echo "❌ Process failed to start. Check log:" | tee -a "$LOG_FILE"
    echo "   tail $LOG_FILE" | tee -a "$LOG_FILE"
    rm "$PID_FILE"
    exit 1
fi

echo "" | tee -a "$LOG_FILE"
echo "Collection running in background. Safe to close terminal." | tee -a "$LOG_FILE"
