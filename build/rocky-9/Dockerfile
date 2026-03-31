FROM rockylinux:9
LABEL org.opencontainers.image.authors="support@halon.io"

ARG TARGETARCH
ARG TARGETOS

ARG HALON_REPO_USER
ARG HALON_REPO_PASS

RUN yum update -y && yum install -y epel-release

RUN echo "[halon]" >> /etc/yum.repos.d/Halon.repo
RUN echo "name=Halon repository" >> /etc/yum.repos.d/Halon.repo
RUN echo "baseurl=https://${HALON_REPO_USER}:${HALON_REPO_PASS}@repo.halon.io/dists/rhel9/${TARGETARCH}/stable" >> /etc/yum.repos.d/Halon.repo
RUN echo "gpgcheck=1" >> /etc/yum.repos.d/Halon.repo
RUN echo "gpgkey=https://raw.githubusercontent.com/halon/changelog/master/pgp-keys/7F0A73B5.asc" >> /etc/yum.repos.d/Halon.repo
RUN echo "enabled=1" >> /etc/yum.repos.d/Halon.repo
RUN yum install -y halon-6.0.0

RUN yum install -y cmake3 make gcc-c++ rpm-build

RUN echo -n "ROCKY9" > /OSRELEASE.txt

RUN mkdir /tmp/halon
RUN <<EOF
set -eu
cd /tmp/halon
case "$TARGETARCH" in
  amd64)
    SHA256=031f088e5d955bab8657ede27ad4e3bc5b7c1ba281f05f245bcc304f327c987a
    ;;
  arm64)
    SHA256=a290581cfe4fe28ddd737dde3095f3dbeb7f2e4065cab4eae44dfc53b760c2f7
    ;;
esac
curl -O -L https://go.dev/dl/go1.26.1.${TARGETOS}-${TARGETARCH}.tar.gz
echo "${SHA256} go1.26.1.${TARGETOS}-${TARGETARCH}.tar.gz" | sha256sum -c
tar -C /usr/local -xzf go1.26.1.${TARGETOS}-${TARGETARCH}.tar.gz
EOF
RUN /usr/local/go/bin/go install github.com/CycloneDX/cyclonedx-gomod/cmd/cyclonedx-gomod@latest
ENV PATH="$PATH:/usr/local/go/bin:/root/go/bin"

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

COPY build.sh /build.sh
CMD ["/build.sh"]
