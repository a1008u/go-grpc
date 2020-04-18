# grpc側  
## /command
```zsh
// .protoからpb.goは自動生成される。
protoc -I proto/ proto/greet.proto --go_out=plugins=grpc:proto
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```


# docker側
```docker

docker tag 3f3c23ef659b a1008u/webserver_grpc_go_webserver_client:v1.0.1
docker tag cc570ad3da71 a1008u/webserver_grpc_go_webserver_server:v1.0.1

docker push a1008u/webserver_grpc_go_webserver_client:v1.0.1
docker push a1008u/webserver_grpc_go_webserver_server:v1.0.1
```



# k8s側  
## /command
```k8s
cd k8s
# 作成
kubectl create --save-config -f greeter_client.yaml
kubectl create --save-config -f greeter_server.yaml

# 削除
kubectl delete -f greeter_server.yaml
kubectl delete -f greeter_client.yaml

kubectl apply -f greeter_client.yaml
kubectl apply -f greeter_server.yaml

# ip確認
minikube ip
```
