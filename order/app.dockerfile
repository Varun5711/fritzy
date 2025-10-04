FROM golang:1.24-alpine AS build

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /go/bin/app ./order/cmd/order

FROM alpine:3.11

WORKDIR /usr/bin
COPY --from=build /go/bin .

EXPOSE 8080

CMD ["app"]