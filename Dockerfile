FROM golang:1.14.3 AS build
LABEL maintainer="Gary Ellis <gary.luis.ellis@gmail.com>"

WORKDIR /github.com/garyellis/secrets-ctl

COPY . /github.com/garyellis/secrets-ctl
RUN package=github.com/garyellis/secrets-ctl/pkg/cmd VERSION=v0.1.0 && \
    BUILD_DATE="-X '${package}.BuildDate=$(date)'" && \
    GIT_COMMIT="-X ${package}.GitCommit=$(git rev-list -1 HEAD)" && \
    _VERSION="-X ${package}.Version=$VERSION" && \
    FLAGS="$GIT_COMMIT $_VERSION $BUILD_DATE" && \
    export GOOS=linux GOARCH=amd64 && go build -o /release/secrets-ctl-${VERSION}_${GOOS}-${GOARCH} -ldflags "${FLAGS}" && \
    export GOOS=darwin GOARCH=amd64 && go build -o /release/secrets-ctl-${VERSION}_${GOOS}-${GOARCH} -ldflags "${FLAGS}"

FROM ubuntu
COPY --from=build /release/ /release/
