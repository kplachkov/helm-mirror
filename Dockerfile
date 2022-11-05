FROM opensuse/leap:15.1

RUN zypper -n in \
		git \
		go1.12 \
		golang-github-cpuguy83-go-md2man \
		make \
		tar \
		gzip

ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH
RUN go get -u golang.org/x/lint/golint && \
	go get -u github.com/vbatts/git-validation && type git-validation

VOLUME ["/go/src/github.com/kplachkov/helm-mirror"]
WORKDIR /go/src/github.com/kplachkov/helm-mirror
COPY . /go/src/github.com/kplachkov/helm-mirror