FROM golang:1.25-alpine AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o wishlist ./cmd/

ENV CGO_ENABLED=0

FROM alpine
WORKDIR /app
COPY --from=build /app/wishlist .
CMD ["./wishlist"]