# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.154.0/containers/debian/.devcontainer/base.Dockerfile

# [Choice] Debian version: buster, stretch
ARG VARIANT="buster"
FROM mcr.microsoft.com/vscode/devcontainers/base:0-${VARIANT}

# ** [Optional] Uncomment this section to install additional packages. **
# RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
#     && apt-get -y install --no-install-recommends <your-package-list-here>


# replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

# update the repository sources list
# and install dependencies
RUN apt-get update \
    && apt-get install -y curl \
    && apt-get -y autoclean

# nvm environment variables
ENV NVM_DIR /usr/local/nvm
ENV NODE_VERSION 12.18.2

# install nvm
# https://github.com/creationix/nvm#install-script
RUN curl --silent -o- https://raw.githubusercontent.com/creationix/nvm/v0.31.2/install.sh | bash

# install node and npm
RUN source $NVM_DIR/nvm.sh \
    && nvm install $NODE_VERSION \
    && nvm alias default $NODE_VERSION \
    && nvm use default

# add node and npm to path so the commands are available
ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
ENV PATH $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

# confirm installation
RUN node -v
RUN npm -v
RUN curl -L --silent -O https://golang.org/dl/go1.13.4.linux-amd64.tar.gz
RUN tar -xzvf ./go1.13.4.linux-amd64.tar.gz -C /usr/local
# RUN tar -C /usr/local -xzf go1.13.4.linux-amd64.tar.gzk
ENV PATH /usr/local/go/bin:$PATH
# RUN go version

RUN curl -L --silent -O https://gethstore.blob.core.windows.net/builds/geth-alltools-linux-amd64-1.9.11-6a62fe39.tar.gz 
RUN tar -xzvf ./geth-alltools-linux-amd64-1.9.11-6a62fe39.tar.gz
RUN mv geth-alltools-linux-amd64-1.9.11-6a62fe39/* /usr/local/bin

RUN apt install -y build-essential
RUN apt install -y python python3 python-dev libffi-dev libc-dev make

# RUN curl -L "https://github.com/docker/compose/releases/download/1.26.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
# RUN chmod +x /usr/local/bin/docker-compose

# RUN curl -fsSL https://get.docker.com | bash

# Install Docker CE CLI
# RUN apt-get update \
#     && apt-get install -y apt-transport-https ca-certificates curl gnupg2 lsb-release \
#     && curl -fsSL https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]')/gpg | apt-key add - 2>/dev/null \
#     && echo "deb [arch=amd64] https://download.docker.com/linux/$(lsb_release -is | tr '[:upper:]' '[:lower:]') $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list \
#     && apt-get update \
#     && apt-get install -y docker-ce-cli

# # Install Docker Compose
# RUN LATEST_COMPOSE_VERSION=$(curl -sSL "https://api.github.com/repos/docker/compose/releases/latest" | grep -o -P '(?<="tag_name": ").+(?=")') \
#     && curl -sSL "https://github.com/docker/compose/releases/download/${LATEST_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose \
#     && chmod +x /usr/local/bin/docker-compose


COPY library-scripts/docker-debian.sh /tmp/library-scripts/
RUN apt-get update && bash /tmp/library-scripts/docker-debian.sh
ENTRYPOINT ["/usr/local/share/docker-init.sh"]
CMD ["sleep", "infinity"]