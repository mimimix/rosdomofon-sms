FROM golang:1.24.4-alpine3.22 AS builder

# Setup base software for building an app.
RUN apk update && apk add ca-certificates git

# Build argument for the application directory
ARG APP_DIR
#ENV APP_DIR=$APP_DIR

WORKDIR /opt/build

# Download dependencies.
COPY ./apps/$APP_DIR/go.mod ./apps/$APP_DIR/go.sum ./apps/$APP_DIR/
COPY ./pkg/ ./pkg/

WORKDIR /opt/build/apps/$APP_DIR
RUN go mod download && go mod verify

# Copy application source.
COPY ./apps/$APP_DIR/ .

# Build the application.
RUN go build -o /opt/bin/application ./cmd/main.go

# Prepare executor image.
FROM alpine:3.21 AS runner

RUN apk update && apk add ca-certificates bash && rm -rf /var/cache/apk/*

WORKDIR /app

#COPY ./migrations migrations

COPY --from=builder /opt/bin/application ./

# Run the application.
CMD ["./application"]
