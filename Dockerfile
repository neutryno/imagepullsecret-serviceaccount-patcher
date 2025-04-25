FROM docker.io/library/golang:1.24.2
WORKDIR /src
COPY ./ /src
RUN GOOS=linux go build -o ./dist/app .

FROM docker.io/library/debian
COPY --from=0 /src/dist/app /app
ENTRYPOINT /app
