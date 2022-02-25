##
##Build
##
FROM golang:1.17.6-alpine AS builder

WORKDIR /eKYC
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o eKYC.o ./cmd/ekyc_api

##
##Deploy
##
FROM alpine
COPY --from=builder /eKYC/eKYC.o .
COPY --from=builder /eKYC/.env .
EXPOSE 8080

ENTRYPOINT ["./eKYC.o" ]