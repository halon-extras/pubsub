FROM --platform=linux/amd64 ubuntu:24.04
LABEL org.opencontainers.image.authors="support@halon.io"

ARG HALON_REPO_USER
ENV HALON_REPO_USER=$HALON_REPO_USER
ARG HALON_REPO_PASS
ENV HALON_REPO_PASS=$HALON_REPO_PASS

RUN apt-get update

RUN apt-get install -y wget gnupg
RUN apt-get install -y apt-transport-https

RUN wget -qO - https://raw.githubusercontent.com/halon/changelog/master/pgp-keys/7F0A73B5.asc | apt-key add -
RUN echo "deb https://repo.halon.io/ noble stable" >> /etc/apt/sources.list.d/halon.list
RUN echo "machine repo.halon.io login ${HALON_REPO_USER} password ${HALON_REPO_PASS}" >> /etc/apt/auth.conf
RUN apt-get update && apt-get install -y halon=6.5.2

RUN /usr/bin/install -d /var/run/halon
ENV LD_LIBRARY_PATH=$LD_LIBRARY_PATH

RUN apt-get install -y build-essential cmake git curl

RUN echo -n "UBU2404" > /OSRELEASE.txt

RUN mkdir /tmp/halon
RUN cd /tmp/halon \
    && curl -O -L https://go.dev/dl/go1.24.0.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.24.0.linux-amd64.tar.gz

RUN mkdir /licenses
RUN cp /usr/local/go/LICENSE /licenses/LICENSE.go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/googleapis/google-cloud-go/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.google-cloud-go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/felixge/httpsnoop/refs/heads/master/LICENSE.txt \
    && cp LICENSE /licenses/LICENSE.httpsnoop
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/go-logr/logr/refs/heads/master/LICENSE \
    && cp LICENSE /licenses/LICENSE.go-logr
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/google/s2a-go/refs/heads/main/LICENSE.md \
    && cp LICENSE /licenses/LICENSE.s2a-go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/google/uuid/refs/heads/master/LICENSE \
    && cp LICENSE /licenses/LICENSE.uuid
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/googleapis/enterprise-certificate-proxy/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.enterprise-certificate-proxy
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/googleapis/gax-go/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.gax-go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/census-instrumentation/opencensus-go/refs/heads/master/LICENSE \
    && cp LICENSE /licenses/LICENSE.opencensus-go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/open-telemetry/opentelemetry-go-instrumentation/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.opentelemetry-go-instrumentation
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/googleapis/google-api-go-client/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.google-api-go-client
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/googleapis/go-genproto/refs/heads/main/LICENSE \
    && cp LICENSE /licenses/LICENSE.go-genproto
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/grpc/grpc-go/refs/heads/master/LICENSE \
    && cp LICENSE /licenses/LICENSE.grpc-go
RUN cd /tmp/halon \
    && curl -O -L https://raw.githubusercontent.com/protocolbuffers/protobuf-go/refs/heads/master/LICENSE \
    && cp LICENSE /licenses/LICENSE.protobuf-go

RUN rm -rf /tmp/halon

COPY .devcontainer/entrypoint.sh /entrypoint.sh
RUN chmod 755 /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

RUN apt-get install -y supervisor
COPY .devcontainer/supervisord.conf /etc/supervisor/conf.d/supervisord.conf
CMD ["/usr/bin/supervisord"]