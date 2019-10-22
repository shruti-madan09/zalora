FROM golang:latest

LABEL maintainer="Shruti Madan <shruti.madan09@gmail.com>"

WORKDIR /workspace/zalora

ENV GOPATH=/workspace/zalora/vendor:/workspace/zalora \
    PATH=/usr/local/go/bin:${GOPATH}:${PATH}

ENV Mode=release

COPY ./ /workspace/zalora

RUN make

EXPOSE 8080

CMD ["./bin/zalora"]
