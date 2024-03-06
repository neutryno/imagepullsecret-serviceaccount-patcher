FROM golang:latest
WORKDIR /src
COPY ./ /src
RUN GOOS=linux go build -o ./dist/app .

FROM debian
COPY --from=0 /src/dist/app /app
ENTRYPOINT /app