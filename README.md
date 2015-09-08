# Batchat

[WIP]

Batchat is a simple chat example of a go server with websocket features using [gorilla](https://github.com/gorilla/websocket). 

Also a simple html/jquery interface is used to demonstrate the functionality on a web page.

```
$ git clone git@github.com:ffhenkes/batchat.git 

$ go get github.com/gorilla/websocket

$ cd batchat && go build -o batchat && go install

$ batchat --addr="addr:PORT"

```

### Open two tabs in a browser..

```
http://addr:port