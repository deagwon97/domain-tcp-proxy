FROM node:16 
COPY ./ /workdir
WORKDIR /workdir
RUN apt-get update -y  &&\
      apt-get install git -y &&\
      npm install -g npm@9.3.1 &&\
      npm install --global yarn --force 
RUN apt install net-tools -y
# install go 1.20.6
RUN wget -c https://go.dev/dl/go1.20.6.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
    
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc