FROM golang:1.23-alpine AS build
RUN apk add --no-cache alpine-sdk curl


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o tailwindcss && \
    chmod +x tailwindcss && \
	./tailwindcss -i internal/assets/input.css -o public/assets/css/output.css  -m -c internal/assets/tailwind.config.js

RUN go build -o main cmd/main.go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main cmd/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
# Copy the templates from the build stage to the production image
COPY --from=build /app/internal/templates /app/internal/templates
COPY --from=build /app/public /app/public
EXPOSE ${PORT}
CMD ./main --port=${PORT}  # Use the shell form here
