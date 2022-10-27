FROM golang:1.19-alpine


WORKDIR /usr/local/app/cmd/server

COPY go.mod /usr/local/app
COPY go.sum /usr/local/app

RUN go mod download

COPY . /usr/local/app

#WORKDIR /usr/local/app/cmd/server
RUN go build .

EXPOSE 8080

CMD [ "./server" ]