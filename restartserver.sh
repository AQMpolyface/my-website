#!/bin/bash
kill $(ps aux | grep go | awk '{print $2}')

nginx -t

nohup go run backend.go & disown
