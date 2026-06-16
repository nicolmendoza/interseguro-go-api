FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o go-api .

FROM alpine:3.20

WORKDIR /app
COPY --from=build /app/go-api /app/go-api
EXPOSE 3000
CMD ["/app/go-api"]
