# EMQ

EMQ X 是基于Erlang/OTP平台开发的开源物联网MQTT消息服务器. 是出色的软实时(Soft-Realtime)、低延时(Low-Latency)、
分布式(Distribute)的语言平台. MQTT是轻量的(Lightweight)、发布订阅模式(PubSub)的物联网消息协议.
[官方文档](https://developer.emqx.io/docs/broker/v3/cn/getstarted.html)



[EMQ X消息服务器要点](doc/EMQ%20X%20read%20note.md)

EMQ X 插件, [插件列表](http://docs.emqtt.cn/zh_CN/latest/plugins.html):
- [ACL鉴权](doc/plugins/acl.md)
    - 决定采用[EMQX_AUTH_HTTP认证插件](doc/plugins/emqx%20auth%20http.md)

- [emq web hook](doc/plugins/web%20hook.md)


[EMQ测试工具](https://www.jianshu.com/p/e5cf0c1fd55c)

[EMQ X在阿里云上部署步骤](doc/k8s%20deployment.md)

## 开发
使用[paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang)连接MQTT broker进行
发布、订阅消息等操作.
