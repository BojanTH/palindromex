# This is a local Dockerfile used for development

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
    cd ./web && npm install && npm start && cd ..

RUN go get -d -v ./...
RUN go test ./...
RUN go install palindromex

EXPOSE 8080

CMD cd /app && palindromex --config=".env.docker.yaml" serve