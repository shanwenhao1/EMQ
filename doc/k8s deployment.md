# EMQ X 阿里云部署

- 使用`kubectl proxy --help`查看kubernetes监听的地址和端口, ![](../doc/picture/emq%20conf.png).
这里按照kubernetes集群内部监听地址192.168.1.143和监听端口8001
- 两种方式: 
    - 手动更改配置上传至自己的镜像库
        - 下载镜像并运行镜像
        - 更改镜像内的配置文件`etc/emqx.conf`
        ```bash
        cluster.discovery = k8s
        # kubernetes API server list, seperated by `,`
        cluster.k8s.apiserver = http://192.168.1.143:8001
        cluster.k8s.service_name = emqx
        # ip or dns
        cluster.k8s.address_type = ip
        cluster.k8s.app_name = emqx
        cluster.k8s.namespace = default
        ```
        - 保存镜像并上传至镜像库
    - 使用环境变量的方式修改`etc/emqx.conf`, [Github上关于EMQ X镜像的文档](https://github.com/emqx/emqx-docker/blob/master/README.md)
    ```bash
     # 指定如下环境变量
     - name: EMQX_CLUSTER__DISCOVERY
         value: k8s
     - name: EMQX_NAME
         value: emqx
     - name: EMQX_CLUSTER__K8S__APISERVER
         value: http://192.168.1.143:8001
     - name: EMQX_CLUSTER__K8S__NAMESPACE
         value: default
     - name: EMQX_CLUSTER__K8S__SERVICE_NAME
         value: emqx
     - name: EMQX_CLUSTER__K8S__ADDRESS_TYPE
         value: ip
     - name: EMQX_CLUSTER__K8S__APP_NAME
         value: emqx
    ```
- 编写k8s部署配置文件yaml, 创建deployment和pod. 部署两个EMQ X镜像的pod
```yaml
# 根据阿里云模板更改
apiVersion: apps/v1beta2 # for versions before 1.8.0 use apps/v1beta1
kind: Deployment
metadata:
  name: emqx
  labels:
    app: emqx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: emqx
  template:
    metadata:
      labels:
        app: emqx
    spec:
      containers:
      - name: emqx
        image: emqx/emqx:latest
        ports:
        - name: emqx-dashboard
          containerPort: 18083
        livenessProbe:
#          exec:
#            command:
#            - sh
#            - -c
#            - "mysqladmin ping -u root -p${MYSQL_ROOT_PASSWORD}"
#
#          tcpSocket:
#            port: 8080
          httpGet:
            path: /
            port: 18083
          initialDelaySeconds: 30
          timeoutSeconds: 5
          periodSeconds: 5
        readinessProbe:
#          exec:
#            command:
#            - sh
#            - -c
#            - "mysqladmin ping -u root -p${MYSQL_ROOT_PASSWORD}"
          httpGet:
            path: /
            port: 18083
          initialDelaySeconds: 5
          timeoutSeconds: 1
          periodSeconds: 5
        # specify user/password from existing secret
        env:
        - name: EMQX_CLUSTER__DISCOVERY
          value: k8s
        - name: EMQX_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__APISERVER
          value: http://192.168.1.143:8001
        - name: EMQX_CLUSTER__K8S__NAMESPACE
          value: default
        - name: EMQX_CLUSTER__K8S__SERVICE_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__ADDRESS_TYPE
          value: ip
        - name: EMQX_CLUSTER__K8S__APP_NAME
          value: emqx



# 参考官网资料编写
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: emqx
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: emqx
    spec:
      containers:
      - name: emqx
        image: emqx/emqx:latest
        ports:
        - name: emqx-dashboard
          containerPort: 18083
        env:
        - name: EMQX_CLUSTER__DISCOVERY
          value: k8s
        - name: EMQX_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__APISERVER
          value: http://192.168.1.143:8001
        - name: EMQX_CLUSTER__K8S__NAMESPACE
          value: default
        - name: EMQX_CLUSTER__K8S__SERVICE_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__ADDRESS_TYPE
          value: ip
        - name: EMQX_CLUSTER__K8S__APP_NAME
          value: emqx
      tty: true
```
- 创建Services, 使用NodePort的方式将emqx-dashboard的端口暴露出来
```yaml
apiVersion: v1
kind: Service
metadata:
  name: emqx
spec:
  ports:
  - port: 32333
    nodePort: 32333
    targetPort:  emqx-dashboard
    protocol: TCP
  selector:
    app: emqx
  type: NodePort
```
- 部署服务, 结合前两步编写`emqx.yml`并使用命令`kubectl create -f emqx.yml`部署EMQ X
```yaml
apiVersion: apps/v1
kind: Service
metadata:
  name: emqx
spec:
  ports:
  - port: 32333
    nodePort: 32333
    targetPort:  emqx-dashboard
    protocol: TCP
  selector:
    app: emqx
  type: NodePort

---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: emqx
  labels:
          app: emqx
spec:
  replicas: 2
  template:
    metadata:
      labels:
        app: emqx
    spec:
      containers:
      - name: emqx
        image: emqx/emqx:latest
        ports:
        - name: emqx-dashboard
          containerPort: 18083
        env:
        - name: EMQX_CLUSTER__DISCOVERY
          value: k8s
        - name: EMQX_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__APISERVER
          value: http://192.168.1.143:8001
        - name: EMQX_CLUSTER__K8S__NAMESPACE
          value: default
        - name: EMQX_CLUSTER__K8S__SERVICE_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__ADDRESS_TYPE
          value: ip
        - name: EMQX_CLUSTER__K8S__APP_NAME
          value: emqx
      tty: true
``` 


```yaml
# 阿里云方式创建
apiVersion: v1
kind: Service
metadata:
  name: emqx
spec:
  ports:
  - port: 32333
    nodePort: 32333
    targetPort:  emqx-dashboard
    protocol: TCP
  selector:
    app: emqx
  type: NodePort


---
apiVersion: apps/v1beta2 # for versions before 1.8.0 use apps/v1beta1
kind: Deployment
metadata:
  name: emqx
  labels:
    app: emqx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: emqx
  template:
    metadata:
      labels:
        app: emqx
    spec:
      containers:
      - name: emqx
        image: emqx/emqx:latest
        ports:
        - name: emqx-dashboard
          containerPort: 18083
        livenessProbe:
#          exec:
#            command:
#            - sh
#            - -c
#            - "mysqladmin ping -u root -p${MYSQL_ROOT_PASSWORD}"
#
#          tcpSocket:
#            port: 8080
          httpGet:
            path: /
            port: 18083
          initialDelaySeconds: 30
          timeoutSeconds: 5
          periodSeconds: 5
        readinessProbe:
#          exec:
#            command:
#            - sh
#            - -c
#            - "mysqladmin ping -u root -p${MYSQL_ROOT_PASSWORD}"
          httpGet:
            path: /
            port: 18083
          initialDelaySeconds: 5
          timeoutSeconds: 1
          periodSeconds: 5
        # specify user/password from existing secret
        env:
        - name: EMQX_CLUSTER__DISCOVERY
          value: k8s
        - name: EMQX_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__APISERVER
          value: http://192.168.1.143:8001
        - name: EMQX_CLUSTER__K8S__NAMESPACE
          value: default
        - name: EMQX_CLUSTER__K8S__SERVICE_NAME
          value: emqx
        - name: EMQX_CLUSTER__K8S__ADDRESS_TYPE
          value: ip
        - name: EMQX_CLUSTER__K8S__APP_NAME
          value: emqx
```

## 参考
- [EMQ k8s deployment and golang client](https://studygolang.com/articles/12858?fr=sidebar)

      
      
