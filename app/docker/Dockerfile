ARG BuildImage=go:1.21.0-bullseye
ARG ProductionImage=debian:bullseye

FROM --platform=linux/x86_64 $BuildImage AS build

WORKDIR /app

COPY . .

RUN go mod tidy && \
    go build -o main

FROM --platform=linux/x86_64 $ProductionImage AS production
COPY --from=build /app/main /main

CMD ["/main"]