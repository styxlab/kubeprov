#!/bin/sh
while true; do
  go build
  $@ &
  PID=$!
  inotifywait -r -e modify .
  #kill $PID
done
