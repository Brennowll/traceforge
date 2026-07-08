FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -o /out/traceforge ./cmd/traceforge

FROM alpine:3.20
RUN adduser -D -g '' traceforge
USER traceforge
WORKDIR /work
COPY --from=build /out/traceforge /usr/local/bin/traceforge
ENTRYPOINT ["traceforge"]
CMD ["version"]
