##
##Build
##
FROM golang:1.17.6-alpine AS builder

WORKDIR /eKYC
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o seeder.o ./cmd/seeder

##
##Deploy
##
FROM alpine
COPY --from=builder /eKYC/seeder.o .
COPY --from=builder /eKYC/.env .

CMD ["./seeder.o" ]