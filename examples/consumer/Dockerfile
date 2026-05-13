FROM golang:1.24 AS builder

ARG SKAFFOLD_GO_GCFLAGS

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -gcflags="${SKAFFOLD_GO_GCFLAGS}" -o /bin/consumer ./examples/consumer/

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/consumer /bin/consumer

ENTRYPOINT ["/bin/consumer"]
