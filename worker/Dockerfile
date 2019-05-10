FROM golang:1.11
WORKDIR /usr/src/saas
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["saas"]