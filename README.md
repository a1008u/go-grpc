# grpc側  
## /command
```zsh
// .protoからpb.goは自動生成される。
// protoc -I protoが存在するpathを指定　コンパイルするprotoファイルpath　作成ファイルの置き場を指定
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

# configMapの設定
kubectl create -f config.yaml 

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

---

# Helm側
## command
```helm
helm package go-grpc

heml create go-grpc

# templatesの確認 -> helmでビルド + upgrade
helm install . --name-template go-grpc --debug --dry-run
helm install --name-template go-grpc .
helm upgrade go-grpc .
helm list 

# helmのrepoでstableを確認する
helm search repo stable

```

## 参考
[hellmの公式](https://helm.sh/)  
[事実上の標準ツールとなっているKubernetes向けデプロイツール「Helm」入門](https://knowledge.sakura.ad.jp/23603/)  
[Helm VS Kustomize](https://qiita.com/ttr_tkmkb/items/638ad7acbc3b6fa537df)
