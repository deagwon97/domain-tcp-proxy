FROM node:16 
COPY ./ /workdir
WORKDIR /workdir
RUN apt-get update -y 
RUN apt-get install git -y
# RUN apt install nginx -y 
RUN apt install luajit -y
RUN apt install libluajit-5.1-dev -y
RUN apt install lua5.1 -y
RUN apt install liblua5.1-0-dev -y
RUN apt install -y luarocks
RUN apt install -y nettle-dev
RUN apt install -y lua-check
## https://github.com/openresty/lua-nginx-module#installation


# RUN mkdir -p /usr/local/lib/lua/5.1/ && cp /workdir/proxy-nginx/lua-blowfish/build/blowfish.so /usr/local/lib/lua/5.1/blowfish.so
WORKDIR /workdir/proxy-nginx/nginx-1.19.3
# RUN ./install
WORKDIR /workdir



RUN npm install -g npm@9.3.1
RUN npm install --global yarn --force

# install go 1.20.6
RUN wget -c https://go.dev/dl/go1.20.6.linux-amd64.tar.gz && \
    rm -rf /usr/local/go && tar -C /usr/local -xzf go1.20.6.linux-amd64.tar.gz
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux
    
RUN echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc