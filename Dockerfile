FROM golang:latest AS builder
WORKDIR /src
COPY ./ /src
RUN CGO_ENABLED=0 GOOS=linux go build -o ./dist/app .

FROM scratch AS runtime
COPY --from=builder /src/dist/app /app
ENTRYPOINT /app
