# grpc側  
## /command
```zsh
// .protoからpb.goは自動生成される。
protoc -I proto/ proto/greet.proto --go_out=plugins=grpc:proto
protoc -I helloworld/ helloworld/helloworld.proto --go_out=plugins=grpc:helloworld
```



# docker側
```docker
# イメージIDの確認
docker images webserver_grpc_go_webserver_client
docker images webserver_grpc_go_webserver_server

# tagの作成
docker tag イメージid a1008u/webserver_grpc_go_webserver_client:v1.0.2
docker tag イメージid a1008u/webserver_grpc_go_webserver_server:v1.0.2

# docker hubにpush
docker push a1008u/webserver_grpc_go_webserver_client:v1.0.2
docker push a1008u/webserver_grpc_go_webserver_server:v1.0.2
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

# apply
kubectl apply -f greeter_client.yaml
kubectl apply -f greeter_server.yaml

# yamlの更新を確認
kubectl apply -f greeter_client.yaml --server-dry-run
kubectl diff -f greeter_client.yaml

kubectl apply -f greeter_server.yaml --server-dry-run
kubectl diff -f greeter_server.yaml
```


# minikube側  
## /command
```minikube
# dashboard
minikube dashboard --url

# ip確認
minikube ip
```
