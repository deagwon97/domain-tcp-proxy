FROM node:16 
COPY ./ /workdir
WORKDIR /workdir
RUN apt-get update -y  &&\
      apt-get install git -y &&\
      npm install -g npm@9.3.1 &&\
      npm install --global yarn --force 

# install go 1.20.7
RUN wget -c https://go.dev/dl/go1.20.7.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.7.linux-amd64.tar.gz
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc

RUN cat  << EOF >> /etc/hosts
0.0.0.0  7a4fe220e12bb0c312a47e3885990a10.service.com
0.0.0.0  7a4fe220e12bb0c3ba67abb5cef3a8c0.service.com
0.0.0.0  7a4fe220e12bb0c344cb29f16bbff67a.service.com
EOF