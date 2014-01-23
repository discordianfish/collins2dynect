FROM       ubuntu
MAINTAINER Johannes 'fish' Ziemke <fish@docker.com> (@discordianfish)

ADD        . /collins2dynect
WORKDIR    /collins2dynect
ENTRYPOINT [ "./looper.sh" ]
RUN        apt-get update && \
           apt-get -y install ca-certificates && \
	   chmod a+x looper.sh collins2dynect
