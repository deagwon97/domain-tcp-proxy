# subdomain 기반 다이나믹 리버스 프록시 서버 개발 및 성능 비교

![img](/asset/reverse-proxy-server.png)

 subdomain 기반으로 목적지 서버를 식별하여 사용자를 중계해주는 다이나믹 리버스 프록시 서버를 개발했습니다.  처음 개발할 때는 nodejs  활용했습니다.  하지만 생각보다 성능이 높지 않았고,  golang를 통해서 중계서버를 만든다면 성능상 어떤 차이가 있을 지 확인해보고 싶어 이 실험을 진행했습니다.

 결론적으로 Golang으로 동작하는 중계 서버가 Nodejs로 동작하는 중계서버보다 안정성 측면에서 더 우수하다는 것을 알게되었습니다. 실험에 관한 설명은 아래 블로그 링크를 참고해 주시기 바랍니다.
 - <a href="https://deagwon97.github.io/%EB%84%A4%ED%8A%B8%EC%9B%8C%ED%81%AC/2023/09/30/subdomain-%EA%B8%B0%EB%B0%98-%EB%8B%A4%EC%9D%B4%EB%82%98%EB%AF%B9-%EB%A6%AC%EB%B2%84%EC%8A%A4-%ED%94%84%EB%A1%9D%EC%8B%9C-%EC%84%9C%EB%B2%84-%EA%B0%9C%EB%B0%9C-%EB%B0%8F-%EC%84%B1%EB%8A%A5-%EB%B9%84%EA%B5%90.html">subdomain-기반-다이나믹-리버스-프록시-서버-개발-및-성능-비교</a>



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
