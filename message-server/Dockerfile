FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang git make && \
    go env -w GOPROXY=https://goproxy.cn,direct

RUN git clone https://github.com/git-lfs/git-lfs.git -b v3.4.0 && \
    cd git-lfs && \
    make
# build binary
COPY . /go/src/github.com/opensourceways/software-package-server
RUN cd /go/src/github.com/opensourceways/software-package-server/message-server && GO111MODULE=on CGO_ENABLED=0 go build

WORKDIR /go/src/github.com/opensourceways/software-package-server/download
RUN cd /go/src/github.com/opensourceways/software-package-server/download && GO111MODULE=on CGO_ENABLED=0 go build -a -o download .

# copy binary config and utils
FROM openeuler/openeuler:22.03
WORKDIR /opt/app/

COPY  --from=BUILDER /go/src/github.com/opensourceways/software-package-server/message-server/message-server /opt/app/message-server
COPY  --from=BUILDER /go/src/github.com/opensourceways/software-package-server/download/download /opt/app/download
COPY  --chown=root --chmod=755 --from=BUILDER /git-lfs/bin/git-lfs /usr/local/bin/git-lfs

COPY softwarepkg/infrastructure/pkgciimpl/clone_repo.sh /opt/app/clone_repo.sh
COPY softwarepkg/infrastructure/pkgciimpl/pull_request.sh /opt/app/pull_request.sh

RUN chmod +x /opt/app/*.sh  \
    && dnf update -y  \
    && dnf in -y curl git

ENTRYPOINT ["/opt/app/message-server"]
