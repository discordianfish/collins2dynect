FROM       ubuntu
MAINTAINER Johannes 'fish' Ziemke <fish@docker.com> (@discordianfish)

RUN        apt-get update && apt-get install -yq curl git ca-certificates
RUN        curl -s https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV        PATH    /usr/local/go/bin:$PATH
ENV        GOPATH  /go

ADD        . /collins2dynect
WORKDIR    /collins2dynect
RUN        go get -d && go build && chmod a+x looper.sh
ENTRYPOINT [ "./looper.sh" ]
