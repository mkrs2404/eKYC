##
##Build
##
FROM golang:1.17.6-alpine AS builder

WORKDIR /eKYC
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o monthly_report.o ./cmd/monthly_report

##
##Deploy
##
FROM alpine
COPY --from=builder /eKYC/monthly_report.o .
COPY --from=builder /eKYC/.env .

ENTRYPOINT ["./monthly_report.o" ]