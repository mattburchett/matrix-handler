FROM golang:alpine AS build

WORKDIR /root/
COPY * /root/
RUN go build -o /root/matrix-handler cmd/svr/main.go

FROM scratch
COPY --from=build /root/matrix-handler /bin/matrix-handler
ENTRYPOINT ["/bin/matrix-handler"]