## EMQ部分知识

EMQ X默认端口
- 1883: MQTT协议端口
- 8883: MQTT/SSL端口
- 8083: MQTT/WebSocket端口
- 8080: HTTP API端口
- 18083: Dashboard管理控制台端口

EMQ X节点间桥接有两种协议
- RPC桥接: RPC桥接只能在EMQ X Broker间使用, 且不支持订阅远程节点的主题去同步数据
- MQTT桥接: MQTT桥接同时支持转发和通过订阅主题来实现数据同步两种方式


## 部署

### Docker部署启动
```bash
# 拉取最新的镜像
docker pull emqx/emqx:v3.1.2
# 启动docker容器
docker run -d --name emqx31 -p 1883:1883 -p 8083:8083 -p 8883:8883 -p 8084:8084 -p 18083:18083 emqx/emqx:v3.1.2

docker run -d --name emqx31 --net=host -p 1883:1883 -p 8083:8083 -p 8883:8883 -p 8084:8084 -p 18083:18083 emqx/emqx:v3.1.2
      # 可使用以下命令进入docker容器查看
      ker exec -it ContainerId sh
```

### EMQ X配置文件

| 配置文件 | 说明 |
| :---: | :---: |
| etc/emqx.conf | EMQ X 消息服务器配置文件 |
| etc/acl.conf | EMQ X 默认ACL规则配置文件 |
| etc/plugins/*.conf | EMQ X 各类插件配置文件 |



EMQ X集群设置
- etc/emqx.conf
    - 默认设置 
        ```bash
        # 集群名称
        cluster.name = emqxcl
        # 集群发现策略
        cluster.discovery = manual
        # 启用集群自愈
        cluster.autoheal = on
        # 宕机节点自动清除周期
        cluster.autoclean = 5m

        # 允许客户端以匿名身份通过验证, 生产环境建议改为false
        allow_anonymous = true
        # 设置所有ACL规则都不能匹配时是否允许访问
        acl_nomatch = allow
        # 设置存储ACL规则的默认文件
        acl_file = etc/acl.conf
        # 设置是否允许ACL缓存
        enable_acl_cache = on
        # 设置每个客户端 ACL 最大缓存数量
        acl_cache_max_size = 32
        # 设置 ACL 缓存的有效时间
        acl_cache_ttl = 1m
        ```
        - 集群发现策略:
            - manual: 手动命令创建集群
            - static: 静态节点列表自动集群
            - mcast: UDP 组播方式自动集群
            - dns: DNS A 记录自动集群
            - etcd: 通过 etcd 自动集群
            - k8s: Kubernetes 服务自动集群
        - 基于kubernetes自动集群
        ```bash
        cluster.discovery = k8s
        # Kubernetes API服务器列表, 以 , 进行隔离
        cluster.k8s.apiserver = http://10.110.111.204:8080cluster.k8s.service_name = emqx
        # 帮助查找集群中的EMQ X节点的服务名称
        cluster.k8s.service_name = emqx
        # 用于从k8s服务中提取host的地址类型
        cluster.k8s.address_type = ip
        # EMQ X的节点名称
        cluster.k8s.app_name = emqx
        # kubernetes的命名空间
        cluster.k8s.namespace = default
        ```
        - acl访问控制规则定义. EMQ X消息服务器接收到MQTT客户端发布的或订阅请求时, 会逐条匹配ACL规则,
        匹配成功返回allow或deny.
        ```bash
        # 格式
        允许|拒绝  用户|IP地址|ClientID  发布|订阅  主题列表
        # 样例1: 允许dashboard用户订阅 $SYS/#
        {allow, {user, "dashboard"}, subscribe, ["$SYS/#"]}
        # 样例2: 允许本机用户发布订阅全部主题
        {allow, {ipaddr, "127.0.0.1"}, pubsub, ["$SYS/#", "#"]}
        # 样例3: 拒绝本机用户以外的其他用户订阅 $SYS/# 与 # 主题
        {deny, all, subscribe, ["$SYS/#", {eq, "#"}]}
        # 允许上述规则以外的任何情形
        {allow, all}
        ```
    
## 插件

`data/loaded_plugins`配置中为需要在系统启动时默认启动的插件    

命令行控制插件
```bash
## 显示所有可用的插件列表
./bin/emqx_ctl plugins list
## 加载某插件
./bin/emqx_ctl plugins load emqx_auth_username
## 卸载某插件
./bin/emqx_ctl plugins unload emqx_auth_username
## 重新加载某插件
./bin/emqx_ctl plugins reload emqx_auth_username
```

## EMQ 架构设计

EMQ X 消息服务器在设计上, 分离了前端协议(FrontEnd)与后端集成(Backend)，其次分离了消息路由平面(Flow Plane)与
监控管理平面(Monitor/Control Plane)

- 支持百万TCP连接
- 全异步架构: 基于Erlang/OTP平台的全异步的架构
    - 异步TCP连接处理
    - 异步主题(Topic)订阅、异步消息发布
    - 部分资源负载限制部分采用同步设计, 比如TCP连接创建和Mnesia数据库事务执行
- 消息持久化: EMQ在设计上分离消息路由与消息存储职责后, 数据复制容灾备份甚至应用集成, 可在数据层面灵活实现.
    - 消息路由基于内存
    - 消息存储基于磁盘
    - EMQ X企业版产品中, 可以通过规则引擎或插件的方式, 持久化消息到Redis、MongoDB、MySQL、PostgreSQL等数据库中,
    或Kafka等消息队列中.
- 设计原则:
    - 核心解决的问题: 处理海量的并发MQTT连接与路由消息
    - 充分利用Erlang/OTP平台软实时、低延时、高并发、分布容错的优势
    - 连接(Connection)、会话(Session)、路由(Router)、集群分层(Cluster)
    - 消息路由平面与控制管理平面(Control Plane)分离
    - 支持后端数据库或NoSQL实现数据持久化、容灾备份与应用集成
- 系统分层:
    - 连接层(Connection Layer): 负责 TCP 连接处理、 MQTT 协议编解码
    - 会话层(Session Layer): 处理 MQTT 协议发布订阅消息交互流程
    - 路由层(Router Layer): 节点内路由派发 MQTT 消息. 路由层维护订阅者与订阅关系表, 
    并在本节点发布订阅模式派发(Dispatch)消息
    - 分布层(Distributed Layer): 分布层通过匹配主题树(Topic Trie)和查找路由表(Route Table),在集群的节点间转发路由
    MQTT 消息. 分布层维护全局主题树(Topic Trie)与路由表(Route Table).主题树由通配主题构成, 路由表映射主题到节点.
    - 认证与访问控制(ACL): 连接层支持可扩展的认证与访问控制模块
    - 钩子(Hooks)与插件(Plugins): 系统每层提供可扩展的钩子，支持插件方式扩展服务器


## 分布集群

- Erlang/OTP分布式编程: 由分布互联的Erlang运行系统组成.
- 节点之间通过TCP互联
- Erlang节点由唯一的节点名称标识, 节点间通过名称进行通信寻址.
- epmd： Erlang端口映射服务程序
- 安全: Erlang节点间通过一个相同的cookie进行互连认证
```bash
$HOME/.erlang.cookie 文件
erl -setcookie <Cookie>
```
- EMQ X消息服务器每个集群节点, 都保存一份主题树(Topic Trie)和路由表.
- 客户端的主题订阅关系, 只保存在客户端所在节点, 用于本节点内派发消息到客户端.

![](../doc/picture/msg%20push%20in%20cluster.png)集群消息在节点间的路由与派发流程
- 节点加入集群、退出集群命令形式: 
```bash
# 启动两台节点, 在emqx@s2.emqx.io上执行
./bin/emqx_ctl cluster join emqx@s1.emqx.io
# 任意节点上查询集群状态
./bin/emqx_ctl cluster status
# 节点退出集群
./bin/emqx_ctl cluster leave
# 在某个节点上强制删除集群中的其他某个节点
./bin/emqx_ctl cluster force-leave emqx@s2.emqx.io
```
- 节点发现与自动集群
- 跨节点会话: EMQ X集群模式下, MQTT连接的持久会话跨节点. 比如一MQTT客户端在node1节点上创建持久会话, 
客户端断线重连至node2时, MQTT的连接在node2节点, 但持久会话仍在node1节点,.


## 部署架构
EMQ X可作为物联网接入服务(IOT Hub)部署在阿里云等公有云平台上. 

典型部署结构![](../doc/picture/Deployment.png)
- LB(负载均衡): 分发设备的MQTT连接与消息到EMQ X集群
- 推荐在LB终结SSL连接, 设备与LB之间TLS安全连接, LB与EMQ X之间普通TCP连接. 可轻松支持100万设备.

## 协议
- MQTT协议: 一个轻量的发布订阅模式消息传输协议, 专门针对低带宽和不稳定网络环境的物联网应用设计.
    - 开放消息协议, 简单易实现
    - 发布订阅模式, 一对多消息发布
    - 基于TCP/IP网络连接
    - 1字节固定报头, 2字节心跳报文, 报文结构紧凑
    - 消息Qos支持, 可靠传输保证: Qos保证不是端到端的, 是客户端到服务器之间的
- MQTT-SN协议: 与MQTT不同的是, 使用UDP进行通信.
    - 不同点:
        - Topic使用TopicId(16-byte的数字)代替
        - MQTT-SN可随时更改will的内容, 甚至取消, 而MQTT只允许Connect时设定will的内容且不允许修改.
        - MQTT-SN的网络中有网关这种设备, 负责把MQTT-SN转换成MQTT, 和云端的MQTT Broker通信. 支持自动发现网关的功能
        - MQTT-SN还支持设备的睡眠功能, 如果设备进入睡眠状态, 无法接收UDP数据, 网关将把下行的PUBLISH消息缓存起来, 
        直到设备苏醒后再传送.


## 参考
- [EMQ管理控制台使用](https://www.jianshu.com/p/ae76ac570f51)