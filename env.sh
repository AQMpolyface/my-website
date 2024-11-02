#!/bin/bash

kill $(ps aux | grep ./main | awk '{print $2}')

export DB_USER=polyface
export DB_PASSWORD=tUring#2007#
export DB_NAME=auth
export DB_HOST=localhost
export protUsername=uwupolyface

sudo systemctl restart nginx

rm main

go build main.go

chmod +x main

nohup ./main & disown
