FROM alpine:3.24.1

COPY ghlabels /usr/local/bin/ghlabels
RUN chmod +x /usr/local/bin/ghlabels

RUN mkdir /workdir
WORKDIR /workdir

ENTRYPOINT [ "/usr/local/bin/ghlabels" ]