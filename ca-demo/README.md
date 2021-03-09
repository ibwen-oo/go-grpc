# 生成证书

## 下载cfssl工具


## 使用cfssl工具生成证书
```
wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -O /usr/bin/cfssl
wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -O /usr/bin/cfssl-json
wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64 -O /usr/bin/cfssl-certinfo
chmod a+x /usr/bin/cfssl* 
```

## 生成ca证书
```
[root@test ~]# mkdir -p /data/cert
[root@test ~]# cd /data/cert

# 创建CA证书信息文件
[root@test cert]# cat ca-csr.json
{
    "CN": "lbingwen",
    "hosts": [
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "ST": "beijing",
            "L": "beijing",
            "O": "bing",
            "OU": "ops"
        }
    ],
    "ca": {
        "expiry": "175200h"
    }
}

"""
CN: Common Name,浏览器使用该字段验证网站是否合法,一般写的是域名,非常重要.
C: Country,国家.
ST: State,州,省.
L: Locality,地区,城市.
O: Organization Name,组织名称,公司名称.
OU: Organization Unit Name,组织单位名称,公司部门.
"""

# 生成CA证书
[root@test cert]# cfssl gencert -initca ca-csr.json | cfssl-json -bare ca
[root@test cert]# ll
total 16
-rw-r--r-- 1 root root  997 Jan 30 17:12 ca.csr  # 证书签名请求
-rw-r--r-- 1 root root  329 Jan 30 17:12 ca-csr.json
-rw------- 1 root root 1679 Jan 30 17:12 ca-key.pem  # 私钥
-rw-r--r-- 1 root root 1346 Jan 30 17:12 ca.pem  # 证书
```

## 创建基于跟证书的config配置文件
```
[root@test cert]# cat ca-config.json
{
    "signing": {
        "default": {
            "expiry": "175200h"
        },
        "profiles": {
            "server": {
                "expiry": "175200h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth"
                ]
            },
            "client": {
                "expiry": "175200h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "peer": {
                "expiry": "175200h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth",
                    "client auth"
                ]
            }
        }
    }
}
```

## 生成server端证书
```
# 生成grpc server 端证书,用于grpc server对外提供服务
# 创建生成自签证书签名请求(csr)的json配置文件
# hosts 中为grpc server可能使用的域名,
[root@test apiserver]# cert]# cat server-csr.json
{
    "CN": "grpc-server",
    "hosts": [
        "127.0.0.1"
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "ST": "beijing",
            "L": "beijing",
            "O": "bing",
            "OU": "ops"
        }
    ]
}

# 生成证书
[root@test cert]# cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=server server-csr.json | cfssl-json -bare server
[root@dns apiserver]# ll
-rw-r--r-- 1 root root  836 Jan 30 17:19 ca-config.json
-rw-r--r-- 1 root root  997 Jan 30 17:12 ca.csr
-rw-r--r-- 1 root root  329 Jan 30 17:12 ca-csr.json
-rw------- 1 root root 1679 Jan 30 17:12 ca-key.pem
-rw-r--r-- 1 root root 1346 Jan 30 17:12 ca.pem
-rw-r--r-- 1 root root 1045 Jan 30 18:41 server.csr
-rw-r--r-- 1 root root  305 Jan 30 18:30 server-csr.json
-rw------- 1 root root 1675 Jan 30 18:41 server-key.pem
-rw-r--r-- 1 root root 1399 Jan 30 18:41 server.pem
```

## 生成client端证书
```
# 生成client和server通信的证书
创建生成自签证书签名请求(csr)的json配置文件
[root@test cert]# cat client-csr.json
{
    "CN": "grpc-client",
    "hosts": [
    ],
    "key": {
        "algo": "rsa",
        "size": 2048
    },
    "names": [
        {
            "C": "CN",
            "ST": "beijing",
            "L": "beijing",
            "O": "bing",
            "OU": "ops"
        }
    ]
}

# 生成证书
[root@dns apiserver]# cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=ca-config.json -profile=client client-csr.json | cfssl-json -bare client
[root@dns apiserver]# ll
-rw-r--r-- 1 root root  836 Jan 30 17:19 ca-config.json
-rw-r--r-- 1 root root  997 Jan 30 17:12 ca.csr
-rw-r--r-- 1 root root  329 Jan 30 17:12 ca-csr.json
-rw------- 1 root root 1679 Jan 30 17:12 ca-key.pem
-rw-r--r-- 1 root root 1346 Jan 30 17:12 ca.pem
-rw-r--r-- 1 root root 1001 Jan 30 17:19 client.csr
-rw-r--r-- 1 root root  285 Jan 30 17:14 client-csr.json
-rw------- 1 root root 1679 Jan 30 17:19 client-key.pem
-rw-r--r-- 1 root root 1371 Jan 30 17:19 client.pem
-rw-r--r-- 1 root root 1045 Jan 30 18:41 server.csr
-rw-r--r-- 1 root root  305 Jan 30 18:30 server-csr.json
-rw------- 1 root root 1675 Jan 30 18:41 server-key.pem
-rw-r--r-- 1 root root 1399 Jan 30 18:41 server.pem
```

# 目录结构
```
# 服务端
go-grpc-example
├── cert
│   ├── ca.pem
│   ├── server.pem
│   ├── server-key.pem
├── proto
│   ├── prod.proto
└── server
    ├── prod.pb.go
    ├── server.go

# 客户端
grpc-client
├── cert
│   ├── ca.pem
│   ├── client.pem
│   ├── client-key.pem
├── proto
│   ├── prod.pb.go
├── client.go 
```