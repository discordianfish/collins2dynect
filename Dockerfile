FROM       ubuntu
MAINTAINER Johannes 'fish' Ziemke <fish@docker.com> (@discordianfish)

ADD        . /collins2dynect
WORKDIR    /collins2dynect
ENTRYPOINT [ "./looper.sh" ]
RUN        chmod a+x looper.sh collins2dynect
