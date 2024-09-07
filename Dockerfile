FROM golang:1.23.1-alpine AS build

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod verify

COPY . .
RUN make build

RUN adduser -D -g '' -s /bin/false -h /mukesh mukesh

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /etc/passwd /etc/passwd

COPY --from=build /app/db/migration /migration
COPY --from=build /app/bin/marketing-api /bin/marketing-api

USER mukesh

ENTRYPOINT ["/bin/marketing-api"]