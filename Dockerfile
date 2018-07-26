FROM golang as build

ADD . /go/src/certwatcher
WORKDIR /go/src/certwatcher

RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/alecthomas/kingpin.v2
RUN go get k8s.io/apimachinery/pkg/apis/meta/v1
RUN go get k8s.io/client-go/kubernetes
RUN go get k8s.io/client-go/tools/clientcmd

RUN go install
RUN go build

ENTRYPOINT ["/go/bin/certwatcher"]