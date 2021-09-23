ARG GOVERSION="1.16.2"

FROM golang:${GOVERSION}-alpine AS buildenv

ARG GOOS="linux"

COPY . $GOPATH/src/
WORKDIR $GOPATH/src

RUN	apk add --quiet --no-cache \
		build-base \
		make \
		git && \
	make clean build STATIC=true

FROM scratch
ARG VERSION="0.7.0"
LABEL org.opencontainers.image.title="tea - CLI for Gitea - git with a cup of tea"
LABEL org.opencontainers.image.description="A command line tool to interact with Gitea servers"
LABEL org.opencontainers.image.version="${VERSION}"
LABEL org.opencontainers.image.authors="Tamás Gérczei <tamas@gerczei.eu>"
LABEL org.opencontainers.image.vendor="The Gitea Authors"
COPY --from=buildenv /go/src/tea /
ENV HOME="/app"
ENTRYPOINT ["/tea"]
