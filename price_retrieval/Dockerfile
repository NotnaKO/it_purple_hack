FROM golang:1.22

RUN export PATH=/usr/pgsql-16.2/bin:$PATH

ADD . /price_retrieval/


WORKDIR /price_retrieval

RUN go mod tidy
RUN go build

CMD ["./price_retrieval", "-config_path=../config/price_retrieval.yaml"]