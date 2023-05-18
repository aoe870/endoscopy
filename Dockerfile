FROM docker.das-security.cn/golang as build

ENV GO111MODULE=on
ENV GOPROXY=https://ci.das-security.cn/repository/go_cn/

WORKDIR /go/release

ADD . .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o endoscopy cmd/main.go

FROM docker.das-security.cn/alpine as prod

COPY --from=build /go/release/endoscopy /

ENTRYPOINT ["/endoscopy", "server"]