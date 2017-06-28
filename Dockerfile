FROM alpine:latest
RUN apk add --update ca-certificates

ADD ./githubql_exporter /usr/bin/githubql_exporter

EXPOSE 9212

ENTRYPOINT ["/usr/bin/githubql_exporter"]
