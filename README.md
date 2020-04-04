```zsh
// .protoからpb.goは自動生成される。
protoc -I proto/ proto/greet.proto --go_out=plugins=grpc:proto
```
protoc -I work/ work/helloworld.proto --go_out=plugins=grpc:work
