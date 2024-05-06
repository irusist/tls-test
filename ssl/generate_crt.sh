#!/bin/bash

#
#-newkey rsa:2048 生成一个新的2048位RSA密钥。
#-nodes 表示不对生成的密钥进行加密。
#-keyout 指定输出私钥文件的名称。
#-x509 表示生成一个自签名证书而非证书请求。
#-days 365 指定证书的有效期。
#-out 指定输出证书文件的名称。
#-subj 后面跟的字符串指定了证书主题字段的各个部分：
#/C=CN 表示国家代码（例如中国）。
#/ST=Shanghai 表示州或省份。
#/L=Shanghai 表示地点（城市）。
#/O=YourCompany 表示组织（公司）名称。
#/OU=YourDepartment 表示组织单位（部门）。
#/CN=localhost （服务端）或/CN=client （客户端）用于指定常用名称（通常是服务器的域名或客户端的标识）


# 生成CA证书
openssl req -newkey rsa:2048 -nodes -x509 -days 3650 -keyout ca.key -out ca.pem  \
-subj "/C=CN/ST=Shanghai/L=Shanghai/O=YourCACompany/OU=CertificationAuthority/CN=YourCA"


# 为服务端创建密钥和证书
#openssl req -newkey rsa:2048 -nodes -keyout server.key -x509 -days 365 -out server.pem \
#-subj "/C=CN/ST=Shanghai/L=Shanghai/O=YourCompany/OU=YourDepartment/CN=localhost"

openssl req -newkey rsa:2048 -nodes -keyout server.key -out server.csr \
-subj "/C=CN/ST=Shanghai/L=Shanghai/O=YourCompany/OU=YourServerDepartment/CN=localhost" \
&& openssl x509 -req -days 365 -in server.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out server.pem

# 为客户端创建密钥和证书
#openssl req -newkey rsa:2048 -nodes -keyout client.key -x509 -days 365 -out client.pem \
#-subj "/C=CN/ST=Shanghai/L=Shanghai/O=YourCompany/OU=YourDepartment/CN=client"

openssl req -newkey rsa:2048 -nodes -keyout client.key -out client.csr \
-subj "/C=CN/ST=Shanghai/L=Shanghai/O=YourCompany/OU=YourClientDepartment/CN=client" \
&& openssl x509 -req -days 365 -in client.csr -CA ca.pem -CAkey ca.key -CAcreateserial -out client.pem

#cat client.key client.pem > combined.pem
cat client.key client.pem ca.pem > combined.pem

# 客户端访问
# curl -k -s --cert-type PEM --cert combined.pem https://localhost:9443

# 创建证书 secret
# kubectl create secret tls proxy-ssl-secret --key combined.pem --cert combined.pem -o yaml --dry-run > secret.yaml

kubectl create secret generic proxy-ssl-secret --from-file=tls.crt=combined.pem  --from-file=tls.key=combined.pem --namespace=default -o yaml  --dry-run=client > secret.yaml

kubectl create secret generic proxy-ssl-secret --from-file=tls.crt=combined.pem  --from-file=tls.key=combined.pem --from-file=ca.crt=combined.pem --namespace=default -o yaml  --dry-run=client > secret.yaml