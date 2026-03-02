#!/bin/bash

# Check Data Collection Status

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
LOG_DIR="$SCRIPT_DIR/logs"
STATUS_FILE="$LOG_DIR/collector_status.json"
PID_FILE="$LOG_DIR/collector.pid"

echo "=== Data Collection Status ==="
echo ""

if [ ! -f "$STATUS_FILE" ]; then
    echo "❌ No active collection found"
    echo ""
    echo "To start collection:"
    echo "  cd $SCRIPT_DIR"
    echo "  ./run-data-collection.sh"
    exit 0
fi

# Read status
if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    
    if ps -p "$PID" > /dev/null 2>&1; then
        echo "✅ Collection RUNNING"
        echo ""
        
        # Parse status file
        START_TIME=$(grep '"start_time"' "$STATUS_FILE" | cut -d'"' -f4)
        END_TIME=$(grep '"end_time"' "$STATUS_FILE" | cut -d'"' -f4)
        LOG_FILE=$(grep '"log_file"' "$STATUS_FILE" | cut -d'"' -f4)
        DATA_DIR=$(grep '"data_dir"' "$STATUS_FILE" | cut -d'"' -f4)
        
        echo "PID: $PID"
        echo "Start Time: $START_TIME"
        echo "Expected End: $END_TIME"
        echo ""
        
        # Calculate progress
        START_EPOCH=$(date -d "$START_TIME" +%s 2>/dev/null || echo 0)
        END_EPOCH=$(date -d "$END_TIME" +%s 2>/dev/null || echo 0)
        NOW_EPOCH=$(date +%s)
        
        if [ $START_EPOCH -gt 0 ] && [ $END_EPOCH -gt 0 ]; then
            TOTAL_DURATION=$((END_EPOCH - START_EPOCH))
            ELAPSED=$((NOW_EPOCH - START_EPOCH))
            REMAINING=$((END_EPOCH - NOW_EPOCH))
            
            if [ $TOTAL_DURATION -gt 0 ]; then
                PROGRESS=$((ELAPSED * 100 / TOTAL_DURATION))
                echo "Progress: ${PROGRESS}%"
                echo "Elapsed: $(($ELAPSED / 3600))h $(($ELAPSED % 3600 / 60))m"
                echo "Remaining: $(($REMAINING / 3600))h $(($REMAINING % 3600 / 60))m"
            fi
        fi
        
        echo ""
        echo "Commands:"
        echo "  View log: tail -f $LOG_FILE"
        echo "  View data: ls -lh $DATA_DIR"
        echo "  Stop: kill $PID"
        
    else
        echo "⚠️  Collection STOPPED (PID $PID no longer running)"
        echo ""
        echo "Check results:"
        DATA_DIR=$(grep '"data_dir"' "$STATUS_FILE" | cut -d'"' -f4)
        echo "  ls -lh $DATA_DIR"
        
        # Clean up
        rm -f "$PID_FILE"
    fi
else
    echo "⚠️  Unknown status"
fi

echo ""
