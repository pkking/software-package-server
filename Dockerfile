FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang git make && \
    go env -w GOPROXY=https://goproxy.cn,direct

# build binary
COPY . /go/src/github.com/opensourceways/software-package-server
RUN cd /go/src/github.com/opensourceways/software-package-server && GO111MODULE=on CGO_ENABLED=0 go build

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow curl && \
    groupadd -g 1000 software-package-server && \
    useradd -u 1000 -g software-package-server -s /bin/bash -m software-package-server

USER software-package-server

WORKDIR /opt/app/

COPY --chown=software-package-server --from=BUILDER /go/src/github.com/opensourceways/software-package-server/software-package-server /opt/app

ENTRYPOINT ["/opt/app/software-package-server"]
