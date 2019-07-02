# EMQ X AUTH HTTP

## 插件详解
emqx_auth_http将每个终端的接入事件、ACL事件抛给用户自己的WebServer以实现接入认证和ACL鉴权的功能.

实际上, emqx_auth_http相对于用户的Web Services来讲只是个简单的、无状态的HTTP Client. 它只是将EMQX内部的登录认证、
和ACL控制的请求转发到用户的Web Services, 并做一定逻辑处理而已.

架构逻辑如下:
![](../../doc/picture/emqx%20auth%20http.png)
- 认证: 每当终端一个CONNECT请求上来时, 将其携带的ClientId、Username、Password等参数, 向用户自己
配置Web Services发起一个认证请求. 成功则允许该终端连接
- ACL: 每当终端执行PUBLISH和SUBSCRIBE操作时, 将ClientId和Topic等参数, 向用户自己配置Web Services
发起一个ACL的请求. 成功则允许此次PUBLISH/SUBSCRIBE


配置文件: 配置了三个HTTP Request的参数
- 终端接入认证(auth_req)
```bash
# 配置auth_req请求所需要访问的URL路径地址
auth.http.auth_req = http://192.168.1.94:8081/mqtt/auth                  
# 配置auth_req请求所使用的HTTP Method, 仅支持GET/POST/PUT                                           
auth.http.auth_req.method = post                                      
# 配置auth_req请求所携带的参数列表, params参数支持占位符%u %c %a %P                                                    
auth.http.auth_req.params = clientid=%c,username=%u,password=%P
```
- 判断是否为超级用户(super_req)
```bash
## Value: URL                                                         
auth.http.super_req = http://192.168.1.94:8081/mqtt/superuser            
## Value: post | get | put                                            
auth.http.super_req.method = post                                     
## Value: Params                                                      
auth.http.super_req.params = clientid=%c,username=%u  
```
- ACL请求(acl_req)
```bash
## Value: URL                                                         
auth.http.acl_req = http://192.168.1.94:8081/mqtt/acl                    
## Value: post | get | put                                            
auth.http.acl_req.method = get                                        
## Value: Params                                                      
auth.http.acl_req.params = access=%A,username=%u,clientid=%c,ipaddr=%a,topic=%t
```
 

## 参考
- [EMQX_AUTH_HTTP 认证插件使用指南](https://www.jianshu.com/p/7918974c026d)
- [http服务鉴权](https://www.cnblogs.com/shihuc/p/10679800.html)