# Web Hook
webhook的概念很简单,就是在各类hook点给预设的API发送数据,用以动态监控指定动作.
这跟CPU上的调试用的埋点非常类似,GDB就是用这种思想监控程序的执行过程.

emq的webhook插件使用很简单,只需要设置API的URL,再勾选要监控的动作即可,它硬编码了要发送的格式,
即在Body中发送json格式的文本,json中的内容也根据action硬编码,无法配置,目前我们可以适应该接口,如果有需求,
可以将其改成可配置的软接口,也可以直接更改硬编码的接口.

```bash
web.hook.api.url = http://192.168.1.94:8081/mqtt/webhook

## Encode message payload field
##
## Value: base64 | base62
##
## Default: undefined
## web.hook.encode_payload = base64

web.hook.rule.client.connected.1     = {"action": "on_client_connected"}
web.hook.rule.client.disconnected.1  = {"action": "on_client_disconnected"}
web.hook.rule.client.subscribe.1     = {"action": "on_client_subscribe"}
web.hook.rule.client.unsubscribe.1   = {"action": "on_client_unsubscribe"}
web.hook.rule.session.created.1      = {"action": "on_session_created"}
web.hook.rule.session.subscribed.1   = {"action": "on_session_subscribed"}
web.hook.rule.session.unsubscribed.1 = {"action": "on_session_unsubscribed"}
web.hook.rule.session.terminated.1   = {"action": "on_session_terminated"}
web.hook.rule.message.publish.1      = {"action": "on_message_publish"}
web.hook.rule.message.deliver.1    = {"action": "on_message_deliver"}
web.hook.rule.message.acked.1        = {"action": "on_message_acked"}
```


## 参考
- [EMQ插件组合实现](https://www.cnblogs.com/bforever/p/10518122.html)
- [EMQ github地址, 包含硬编码API json 格式](https://github.com/emqx/emqx-web-hook)