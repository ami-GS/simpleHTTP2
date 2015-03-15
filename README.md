simpleHTTP2
===========

Simple implementation of HTTP2 in Golang

## Usage

### Install
```
$ go get github.com/ami-GS/GoHPACK
$ go get github.com/ami-GS/simpleHTTP2
```
or
```
$ go get github.com/mattn/gom
$ gom install
```

### Change directory and test
```
$ cd $GOPATH/src/github.com/ami-GS/simpleHTTP2/example
```

* Server
```
$ go run serv_Try.go <server port>
```

* Client
```
$ go run cli_Try.go <server ip> <server port>
```


## Function detail

```
Client                                         Server

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
```

## TODO
* Implement stream state
* Colorize output string

## Reference
* https://github.com/syucream/MinimumHTTP2
* https://speakerdeck.com/syucream/2-zui-su-shi-zhuang-v3

#### License
The MIT License (MIT) Copyright (c) 2015 ami-GS