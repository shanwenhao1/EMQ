# ACL鉴权

## ACL鉴权规则
官方提供如下方式来进行用户和ACL验证的存储: (对应的配置方式可以参考官网文档)
- LDAP
- HTTP: 开启HTTP插件后, 会终结ACL访问控制链, 采用http控制链
- MySQL
- Postgre
- Redis
- MongoDB

### 静态配置
官方提供了默认的鉴权, 在`etc/acl.conf`配置文件下. 
```bash
%% 允许'dashboard'用户订阅 '$SYS/#'
{allow, {user, "dashboard"}, subscribe, ["$SYS/#"]}.
%% 允许本机用户发布订阅全部主题
{allow, {ipaddr, "127.0.0.1"}, pubsub, ["$SYS/#", "#"]}.
%% 拒绝用户订阅'$SYS#'与'#'主题
{deny, all, subscribe, ["$SYS/#", {eq, "#"}]}.
%% 允许除以上规则之外的所有操作
{allow, all}.
```
规则如下
```bash
允许|拒绝  用户|IP地址|ClientID  发布|订阅  主题列表

## 访问控制规则采用 Erlang 元组格式，访问控制模块逐条匹配规则:
         ---------              ---------              ---------
Client -> | Rule1 | --nomatch--> | Rule2 | --nomatch--> | Rule3 | --> Default
          ---------              ---------              ---------
              |                      |                      |
            match                  match                  match
             \|/                    \|/                    \|/
        allow | deny           allow | deny           allow | deny
```

此种方式需要手动修改acl配置文件然后重启节点(所有节点都需要修改)

### MySQL实时配置ACL

首先需要关闭匿名认证(默认是开启的谁都能够登录)
```bash
vim /usr/local/emqttd/etc/emq.conf 
## Allow Anonymous authentication
mqtt.allow_anonymous = false
```
重启服务器之后不管是谁都会被链接拒绝, 我们需要准备用于检查用户和权限的mysql表:
```bash
CREATE TABLE `mqtt_user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(100) DEFAULT NULL,
  `password` varchar(100) DEFAULT NULL,
  `salt` varchar(20) DEFAULT NULL,
  `is_superuser` tinyint(1) DEFAULT 0,
  `created` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `mqtt_username` (`username`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;


CREATE TABLE `mqtt_acl` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `allow` int(1) DEFAULT NULL COMMENT '0: deny, 1: allow',
  `ipaddr` varchar(60) DEFAULT NULL COMMENT 'IpAddress',
  `username` varchar(100) DEFAULT NULL COMMENT 'Username',
  `clientid` varchar(100) DEFAULT NULL COMMENT 'ClientId',
  `access` int(2) NOT NULL COMMENT '1: subscribe, 2: publish, 3: pubsub',
  `topic` varchar(100) NOT NULL DEFAULT '' COMMENT 'Topic Filter',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- 建立ACL默认访问控制
INSERT INTO `mqtt_acl` (`id`, `allow`, `ipaddr`, `username`, `clientid`, `access`, `topic`)
VALUES
    (1,1,NULL,'$all',NULL,2,'#'),
    (2,0,NULL,'$all',NULL,1,'$SYS/#'),
    (3,0,NULL,'$all',NULL,1,'eq #'),
    (5,1,'127.0.0.1',NULL,NULL,2,'$SYS/#'),
    (6,1,'127.0.0.1',NULL,NULL,2,'#'),
    (7,1,NULL,'dashboard',NULL,1,'$SYS/#');
```
修改mysql配置文件`etc/plugins/emq_auth_mysql.conf`(或者使用dashboard修改配置)
```bash
auth.mysql.server = xxxxxxxxx:3306
auth.mysql.username = root
auth.mysql.password = xxxxxxxx
auth.mysql.database = emq
```
建立用户
```bash
# 用户名 server 密码 server 密码默认是sha256
insert `mqtt_user`(`username`,`password`) values('server','b3eacd33433b31b5252351032c9b3e7a2e7aa7738d5decdf0dd6c62680853c06');
# 用户名 cline 密码 cline
insert `mqtt_user`(`username`,`password`) values('cline','84829dbd815311888f0e3d85822e9b07d14be89a480a3c09ee67353f0e806e3b');
```
可以配置超级管理员(超级管理员会无视ACL规则对所有的topic都有订阅和推送的权限)
```bash
update `mqtt_user` set `is_superuser`=1 where `id`=1;
```


## 参考
- [EMQ百万级MQTT消息服务(ACL鉴权)](https://my.oschina.net/wenzhenxi/blog/1795748)