#!/bin/bash

#useless whn pushing new html & css. html and css automatically refreshes, cuz served on command.
kill $(ps aux | grep go | awk '{print $2}')

nginx -t

rm main

go build main.go

chmod +x main

nohup ./main & disown
