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

RUN DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# install everything instead of one by one like below
RUN go get -d -v ./...
RUN go install palindromex

#EXPOSE 8080

CMD cd /app && palindromex serve