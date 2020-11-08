### 1.　Install protobuf

 **Method 1:**  (Always timeout and fails in China mainland)

```shell
go get -u github.com/golang/protobuf/proto
go get -u github.com/golang/protobuf/protoc-gen-go
```

**Method 2:**  (Recommend)

```shell
git clone https://github.com/golang/protobuf.git $GOPATH/src/github.com/golang/protobuf
cd $GOPATH/src/github.com/golang/protobuf
go install ./proto
go install ./protoc-gen-go
```

### 2.  Install GRPC

~~~shell
git clone https://github.com/golang/net.git $GOPATH/src/golang.org/x/net
git clone https://github.com/golang/text.git $GOPATH/src/golang.org/x/text
git clone https://github.com/google/go-genproto.git $GOPATH/src/google.golang.org/genproto
git clone https://github.com/grpc/grpc-go.git $GOPATH/src/google.golang.org/grpc
cd $GOPATH/src/
go install google.golang.org/grpc
~~~

If some package can't find in the process of install, you should download it from github. 

For example, if I miss a protobuf-go, I should run the blow command to download the package to my local 

environment.

> git clone https://github.com/protocolbuffers/protobuf-go.git $GOPATH/src/google.golang.org

