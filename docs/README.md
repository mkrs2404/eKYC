# eKYC Client

## Features

- Signup a new API client to get access and secret keys
- Upload images with metadata
- Perform face matching between two image ids to get the face match score
- Perform optical character recognition (OCR) on images
- Reporting - generate client wise daily/weekly/monthly reports for billing purposes via cron jobs

## Architecture

![](/app/assets/architecture.png)

## Assumptions

- Actural ML models not used for OCR and face match.
- Fake KYC data is generated and used

## Requirements

- go version >= 1.17
- docker
- docker-compose

## Usage

Clone the repository using:

```bash
git clone github.com/mkrs2404/eKYC.git
```

You can create an .env file as per the [template](#env-template), or pass all the parameters in the 'make run' command as follows. You can also override .env file's variables by passing as flags

Run following commands in the root directory:

- Run API

```
make host=<hostname> db=<db_name> user=<username> pwd=<password> port=<db_port> server=<IP:Port> minio_server=<IP:Port> minio_pwd=<minio_pwd> minio_user=<minio_user> redis_server=<IP:Port> redis_pwd= <redis_pwd> rabbitmq_server=<rabbitmq_server> rabbitmq_user=<rabbitmq_user> rabbitmq_pwd=<rabbitmq_pwd> face_worker_queue=<face_worker_queue> run
```

- Generate Daily Report
  
```
make run-daily-report
```

- Generate Monthly Report for a client
  
```
make client=<clientID> run-monthly-report
```

- Test

```
make test
```

- Clean workspace

```
make clean
```

## env Template

Create an ".env" file in the root directory with the following template

```
DB_HOST=
DB_NAME=
DB_USER=
DB_PASSWORD=
DB_PORT=
SERVER_ADDR=

#JWT Secret Key
SECRET_KEY=
#Token delay after mentioned hours
TOKEN_EXPIRY_DELAY=

#Minio Client
MINIO_SERVER=
MINIO_USER=
MINIO_PWD=

#Redis
REDIS_SERVER=
REDIS_PASSWORD=

#Test Environment
TEST_DB_HOST=
TEST_DB_NAME=
TEST_DB_USER=
TEST_DB_PASSWORD=
TEST_DB_PORT=
TEST_MINIO_SERVER=
TEST_MINIO_USER=
TEST_MINIO_PWD=

#RabbitMQ Client
RABBITMQ_SERVER=
RABBITMQ_USER=
RABBITMQ_PWD=
FACE_WORKER_QUEUE=
```
