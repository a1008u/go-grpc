```zsh
// .protoからpb.goは自動生成される。
protoc -I proto/ proto/greet.proto --go_out=plugins=grpc:proto
```
protoc -I work/ work/helloworld.proto --go_out=plugins=grpc:work



# k8s側  
## /command
```k8s
cd k8s
# 作成
kubectl create -f greeter_client.yaml
kubectl create -f greeter_server.yaml

# 削除
kubectl delete -f greeter_server.yaml
kubectl delete -f greeter_client.yaml

# ip確認
minikube ip
```
