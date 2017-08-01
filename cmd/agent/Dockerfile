FROM alpine

MAINTAINER Mike Roetgers <mike@roetgers.org>

ENV K8S_VERSION="v1.7.0"

RUN apk update \
 && apk add ca-certificates curl \
 && curl -L https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && chmod +x /usr/local/bin/kubectl

ADD ./agent /usr/local/bin/deploi-agent

ENTRYPOINT ["deploi-agent"]
CMD ["-h"]
