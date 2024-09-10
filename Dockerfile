# базовый образ
FROM golang:1.22

# скопировать все файлы
COPY . /var/app

#рабочая директория
WORKDIR /var/app

# установить значения переменным окружения CGO_ENABLED и GOOS, выполнить команду go build -o /docker-gs-ping
# Build
RUN go build -o /zadanie ./cmd

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run
CMD ["/zadanie"]

