FROM golang:latest

ENV WORK_DIR "$GOPATH/src/github.com/dikhan/terraform-provider-openapi/examples/goa/api"

RUN mkdir -p $WORK_DIR
ADD . $WORK_DIR
WORKDIR $WORK_DIR

COPY swagger/swagger.json /opt/goa/swagger/
COPY swagger/swagger.yaml /opt/goa/swagger/

RUN go get
RUN go build -o goa-service-provider .

EXPOSE 9090

CMD ["./goa-service-provider"]