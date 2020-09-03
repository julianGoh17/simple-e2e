FROM golang:alpine

ENV GOROOT="/usr/local/go" \
    GOPATH="/go" \
    SIMPLE_E2E_PATH="${GOPATH}/src/github.com/julianGoh17/simple-e2e"

ARG USER="e2e"

# Install dependencies
RUN set -ex && \
    apk update && \
    apk add --no-cache sudo bash git openssh gcc go git mercurial

# Add non-root user
RUN adduser -D $USER \
    && echo "$USER ALL=(ALL) NOPASSWD: ALL" > /etc/sudoers.d/$USER \
    && chmod 0440 /etc/sudoers.d/$USER
USER $USER
WORKDIR /home/${USER}

# Set up volume mount
RUN mkdir -p ${SIMPLE_E2E_PATH} && \
    mkdir /home/${USER}/tests && \
    mkdir /home/${USER}/Dockerfiles

COPY ./framework ${SIMPLE_E2E_PATH}/framework
# Need to copy these for tests as Github action doesn't use volumes
COPY ./tests /home/${USER}/tests
COPY ./Dockerfiles /home/${USER}/Dockerfiles

VOLUME [ ${SIMPLE_E2E_PATH}/framework ]
VOLUME [ /home/${USER}/tests ]
VOLUME [ /home/${USER}/Dockerfiles ]

COPY "entrypoint.sh" "/entrypoint.sh"

ENTRYPOINT [ "/entrypoint.sh" ]