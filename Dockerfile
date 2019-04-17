FROM docker:18.06.1-ce as static-docker-source

FROM debian:stretch

ENV GOPATH="/go"

ENV PATH="${PATH}:/usr/local/go/bin:${GOPATH}/bin:/opt/google-cloud-sdk/bin"

RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        jq \
        bash \
        dumb-init \
        procps \
        unzip \
        python \
        git && \
    apt-get autoremove -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN export CLOUDSDK_INSTALL_DIR=/opt && \
    curl https://dl.google.com/dl/cloudsdk/release/install_google_cloud_sdk.bash | bash && \
    gcloud components install kubectl && \
    curl -fsSl -o go.tar.gz https://dl.google.com/go/go1.10.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm -f go.tar.gz

RUN mkdir -p ${GOPATH}/src/github.com/AlexsJones/ ${GOPATH}/bin && \
    curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
    cd ${GOPATH}/src/github.com/AlexsJones && \
    git clone https://github.com/AlexsJones/vortex.git && \
    cd vortex && \
    dep ensure -v && \
    CGO_ENABLED=0 GOOS=linux go build --ldflags="-s -w" -o /usr/bin/vortex && \
    cd /

WORKDIR /go/src/github.com/AlexsJones/

COPY . /go/src/github.com/AlexsJones/gravitywell

RUN cd gravitywell && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X 'main.version=$(cat VERSION)' -X 'main.revision=$(git rev-parse --short HEAD)' -X 'main.buildtime=$(date -u +%Y-%m-%d.%H:%M:%S)'" -o /gravitywell && \
    rm -rf ${GOPATH}

ENTRYPOINT ["dumb-init", "/bin/bash"]
