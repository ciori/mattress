# Builder

FROM golang:1.24.1-bookworm AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mattress .

# Application

FROM golang:1.24.1-bookworm
WORKDIR /app
RUN adduser --shell /bin/bash mattress
COPY --from=builder /src/mattress .
RUN chown mattress:mattress -R /app
USER mattress
EXPOSE 2121
CMD ["./mattress"]
