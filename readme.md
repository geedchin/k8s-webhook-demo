
# demo简介
> 该demo是k8s webhook的一个简单实践。

# 通过k8s签发证书

## 创建k8s csr(CertificateSigningRequest)资源

> 首先在空白目录下生成自己的csr文件及私钥，根据实际情况修改hosts和其他内容
```shell
cat <<EOF | cfssl genkey - | cfssljson -bare server
{
  "hosts": [
    "10.1.1.50"
  ],
  "CN": "kubernetes",
  "key": {
    "algo": "ecdsa",
    "size": 256
  },
  "names": [
    {
      "C": "CN",
      "ST": "SZ",
      "L": "SZ",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF
```
> 此时生成两个文件：server.csr, server-key.pem，server.csr用于向ca申请公钥，server-key.pem为服务的私钥。
> 
> 使用以下命令生成k8s中的csr资源对象：
```shell
cat <<EOF | kubectl apply -f -
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: my-csr
spec:
  request: $(cat server.csr | base64 | tr -d '\n')
  usages:
  - digital signature
  - key encipherment
  - server auth
EOF
```
> 执行以下命令通过申请： `kubectl certificate approve my-csr` <br>
> 执行命令获取签发后的文件： `kubectl get csr my-csr -o jsonpath='{.status.certificate}' | base64 --decode > server.crt` <br>
> server.crt文件即为签发后的公钥

# 运行demo
> 将生成的 server-key.pem server.crt 文件放至 file目录下，然后找到main运行即可