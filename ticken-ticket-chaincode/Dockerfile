FROM golang:1.18-alpine

WORKDIR /chaincode/ticken-event
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 9999
CMD ["ticken-event"]