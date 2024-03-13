FROM golang:1.22

WORKDIR /price_management

ADD . .

RUN go mod tidy
RUN go build

CMD ["./price_management", "-config_path=../config/price_management.yaml"]