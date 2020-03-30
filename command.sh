#!/bin/sh
./main & ./gateway/reflex -r '\.json$' -s -- sh -c "./gateway/krakend run -c ./gateway/config.json -p 8080"