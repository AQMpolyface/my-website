#!/bin/sh

VISIT_FILE="visits.txt"

if [ ! -f "$VISIT_FILE" ]; then
    echo 0 > "$VISIT_FILE"
fi

current_visits=$(cat "$VISIT_FILE")

new_visits=$((current_visits + 1))

echo "$new_visits" > "$VISIT_FILE"

