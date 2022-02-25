##  
##Build
##
FROM golang:1.17.6-alpine AS builder

WORKDIR /eKYC
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o migration.o ./cmd/migrations

##
##Deploy
##
FROM alpine
COPY --from=builder /eKYC/migration.o .
COPY --from=builder /eKYC/.env .

CMD ["./migration.o" ]