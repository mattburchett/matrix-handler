FROM golang:alpine

RUN go build -o /root/matrix-handler cmd/svr/main.go

WORKDIR /root/
CMD ["./matrix-handler"]
