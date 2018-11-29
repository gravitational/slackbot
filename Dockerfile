FROM alpine

RUN apk update

RUN apk add nodejs npm

ADD . /bot

RUN rm -rf /bot/.git

WORKDIR /bot

ENTRYPOINT ["bin/hubot","-a","slack"]
