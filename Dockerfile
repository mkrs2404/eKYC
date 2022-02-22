FROM golang:1.17.6

WORKDIR /eKYC
COPY go.mod /
COPY go.sum /

COPY . /eKYC

RUN go build -o eKYC.o ./cmd/ekyc_api

EXPOSE 8080

CMD [ "source ./.env", "./eKYC.o" ]