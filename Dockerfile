FROM alpine
RUN apk add --no-cache git tzdata
COPY build/scm-status-linux-amd64 /usr/bin/scm-status
WORKDIR /home
RUN ["ls"]
ENTRYPOINT ["scm-status"]