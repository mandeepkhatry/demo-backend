FROM golang:latest

ENV GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY . .

# COPY ./main /usr/bin/main

# COPY ./gateway/reflex /usr/bin/reflex

# COPY ./gateway/krakend /usr/bin/krakend

RUN ["chmod", "+x", "./command.sh"]

EXPOSE 8080

CMD ["./command.sh"]
