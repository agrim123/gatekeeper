FROM golang:1.15.6-alpine as builder

WORKDIR /src

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o gatekeeper .

FROM alpine:3.12

WORKDIR /opt/gatekeeper

COPY --from=builder /src/gatekeeper /opt/gatekeeper/gatekeeper

CMD ["sleep", "10"]
