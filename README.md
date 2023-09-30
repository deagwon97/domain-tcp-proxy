# subdomain 기반 다이나믹 리버스 프록시 서버 개발 및 성능 비교

![img](/asset/reverse-proxy-server.png)

 subdomain 기반으로 목적지 서버를 식별하여 사용자를 중계해주는 다이나믹 리버스 프록시 서버를 개발했습니다.  처음 개발할 때는 golang의 "net” 패키지와 gorutine을 활용했습니다.  하지만 생각보다 성능이 높지 않았고, nodejs 를 통해서 중계서버를 만든다면 성능상 어떤 차이가 있을 지 확인해보고 싶어 이 실험을 진행했습니다.

 결론적으로 Node.js 로 동작하는 중계 서버가 Golang으로 동작하는 중계서버보다 우수한 성능을 보였습니다.

 ![result](/asset/result.png)

 실험에 관한 설명은 아래 블로그 링크를 참고해 주시기 바랍니다.

 - https://deagwon.com/post/subdomain-기반-다이나믹-리버스-프록시-서버-개발-및-성능-비교



### how to run
```
# proxy-nodejs: node.js로 동작하는 중계 서버
cd ~
cd proxy-nodejs
yarn install
yarn build
yarn serve

# proxy-go:  go로 동작하는 중계 서버
cd ~
cd proxy-go
go mod tidy
go build .
./proxy-go

# proxy-test: 실험을 위한 코드
cd ~
cd proxy-test
go mod tidy
go build .
./run.sh
```