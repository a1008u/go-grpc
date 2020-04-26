
## create_ca_certの成果物
 - ca.pem ... CA用の秘密鍵
 - ca.crt ... CAの自己署名証明書
 - ca.crl ... 失効リスト
> caディレクトリに成果物ができます。



```docker-compose
# ca
docker-compose up --build create_ca_cert

# client
docker-compose run create_client_cert go-grpc

# server
docker-compose run create_server_cert go-grpc.server.net

# nginxでのテスト
SERVER=go-grpc.server.net docker-compose up test_nginx
```
