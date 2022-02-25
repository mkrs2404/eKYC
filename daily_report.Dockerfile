##
##Build
##
FROM golang:1.17.6-alpine AS builder

WORKDIR /eKYC
COPY go.mod ./
COPY go.sum ./

COPY . ./

RUN go build -o daily_report.o ./cmd/daily_report

##
##Deploy
##
FROM alpine
COPY --from=builder /eKYC/daily_report.o .
COPY --from=builder /eKYC/.env .

ENTRYPOINT ["./daily_report.o" ]