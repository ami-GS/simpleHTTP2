simpleHTTP2
===========

Simple implementation of HTTP2 in Golang

## Usage

* Install and change Directory
```
$ go get github.com/ami-GS/simpleHTTP2
$ cd $GOPATH/src/github.com/ami-GS/simpleHTTP2/test
```

* Server
```
$ go run serv_Try.go
```

* Client
```
$ go run cli_Try.go
```


## Function detail

Client                                                         Server

      ----------------------------------------->

            connection preface

      ----------------------------------------->

            Settings frame

      <-----------------------------------------

            Settings frame (Flag ACK)

      ----------------------------------------->

            Headers frame (Flag END_HEADERS)

      <-----------------------------------------

            Data frame

      ----------------------------------------->

            GoAway frame


## Reference
* https://github.com/syucream/MinimumHTTP2
* https://speakerdeck.com/syucream/2-zui-su-shi-zhuang-v3