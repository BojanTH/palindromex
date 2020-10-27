# This is a live Dockerfile used by Cloud Run

FROM golang:1.14

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .
RUN apt-get -y update && \
    apt-get install -y apt-utils && \
    apt-get install -y apt-transport-https curl software-properties-common && \
    curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs && \
    cd ./web && npm install && npm build && cd ..

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN go get -d -v ./...
RUN go test ./...
RUN go install palindromex

CMD cd /app && palindromex serve