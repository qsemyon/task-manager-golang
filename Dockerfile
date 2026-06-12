FROM node:24-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM golang:1.26.4-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

COPY --from=frontend-builder /static ./static
RUN go build -o main .

FROM alpine:latest
WORKDIR /root/
COPY --from=backend-builder /app/main .
COPY --from=backend-builder /app/static ./static

EXPOSE 8080
CMD ["./main"]