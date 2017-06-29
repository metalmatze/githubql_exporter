FROM alpine:latest
RUN apk add --update ca-certificates

ADD ./githubql_exporter /usr/bin/githubql_exporter

EXPOSE 9276

ENTRYPOINT ["/usr/bin/githubql_exporter"]
