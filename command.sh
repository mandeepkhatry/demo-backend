./gateway/reflex -r '\.json$' -s -- sh -c "./gateway/krakend run -c ./gateway/config.json -p 8080" & ./main &