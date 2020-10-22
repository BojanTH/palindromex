FROM golang:1.14

WORKDIR /go/src/palindromex
COPY . .
RUN apt-get -y update && \
    apt-get install -y apt-transport-https curl software-properties-common && \
    curl -sL https://deb.nodesource.com/setup_14.x | bash - && \
    apt-get install -y nodejs && \
    cd ./web && npm install && npm start && cd ..

# install everything instead of one by one like below
RUN go get -d -v ./...
RUN go install palindromex

EXPOSE 8080

CMD palindromex serve