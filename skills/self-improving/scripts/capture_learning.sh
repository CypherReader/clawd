#!/bin/bash
# Capture a learning entry

LEARNING_DIR="/root/clawd/memory/learnings"
TODAY=$(date +%Y-%m-%d)
LEARNING_FILE="$LEARNING_DIR/$TODAY.md"

# Create directory if needed
mkdir -p "$LEARNING_DIR"

# Create or append to today's learning file
if [ ! -f "$LEARNING_FILE" ]; then
    cat > "$LEARNING_FILE" << EOF
# Learnings - $TODAY

EOF
fi

# Add timestamp and separator
echo -e "\n---\n## $(date +%H:%M:%S) - Learning Entry\n" >> "$LEARNING_FILE"

# Interactive prompt or accept parameters
if [ $# -eq 0 ]; then
    echo "=== Capture Learning ==="
    echo ""
    read -p "Task Description: " task
    read -p "Approach: " approach
    read -p "Outcome (success/failure/partial): " outcome
    read -p "Key Lesson: " lesson
    read -p "Apply To (future scenarios): " apply
    read -p "Confidence (high/medium/low): " confidence
    
    cat >> "$LEARNING_FILE" << EOF
**Task**: $task
**Approach**: $approach
**Outcome**: $outcome
**Lesson**: $lesson
**Apply To**: $apply
**Confidence**: $confidence
EOF
else
    # Accept as arguments for programmatic use
    cat >> "$LEARNING_FILE" << EOF
**Task**: $1
**Approach**: $2
**Outcome**: $3
**Lesson**: $4
**Apply To**: ${5:-Future similar tasks}
**Confidence**: ${6:-medium}
EOF
fi

echo ""
echo "✅ Learning captured in $LEARNING_FILE"
