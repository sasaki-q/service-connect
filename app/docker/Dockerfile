FROM --platform=linux/x86_64 golang:1.21.0-bullseye AS build

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o main

FROM --platform=linux/x86_64 debian:bullseye AS production
COPY --from=build /app/main /main
RUN apt-get update

CMD ["/main"]