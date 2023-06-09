FROM golang:1.20

WORKDIR /app
COPY go.* .
COPY main.go .

RUN go build -o /long_running_app

CMD [ "/long_running_app" ]
