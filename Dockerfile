FROM node:16 
COPY ./ /workdir
WORKDIR /workdir
RUN apt-get update -y  &&\
      apt-get install git -y &&\
      npm install -g npm@9.3.1 &&\
      npm install --global yarn --force 
<<<<<<< HEAD
RUN apt install net-tools -y
# install go 1.20.6
RUN wget -c https://go.dev/dl/go1.20.6.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
    
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
=======

# install go 1.20.7
RUN wget -c https://go.dev/dl/go1.20.7.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.7.linux-amd64.tar.gz


ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOROOT=/usr/local/go \
    GOPATH=/workdir/go \
    PATH=$PATH:/usr/local/go/bin

RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v github.com/cweill/gotests/gotests@v1.6.0
RUN go install -v github.com/fatih/gomodifytags@v1.16.0
RUN go install -v github.com/josharian/impl@v1.1.0
RUN go install -v github.com/haya14busa/goplay/cmd/goplay@v1.0.0
>>>>>>> e4ea4cf (update)
