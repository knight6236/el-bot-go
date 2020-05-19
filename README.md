# el-bot-go

[el-bot](https://github.com/ElpsyCN/el-bot)的 go 版本。

一个基于 Mirai 的可快速配置的机器人模板。

# 文档

项目处于早期，暂无文档。

# 功能

[开发进度](https://github.com/ElpsyCN/el-bot-go/projects/1)

# 快速开始

## On Unix

1. 安装Golang和JRE
2. 下载依赖: `./install.sh`
3. 启动 mirai-console: ` ./start-console.sh`
4. 按照提示进行操作
5. 成功登录机器人后修改文件`plugins\MiraiAPIHTTP\setting.yml`为下列内容，并记住`authkey`
    ```yml
    authKey: qwertyuiop
    port: 8080
    enableWebsocket: false
    ```
6. 启动 el-bot-go: 在另一个shell运行脚本：`./start-el-bot.sh 机器人QQ号 AUTHKEY`

# 配置文件说明

<!-- config/custom/custom.yml -->

```yml
# 当接收到的群消息或好友消息为 hello 或「你好」时回复「Hello World!（你好 世界！）」
global:
  - when:
      message:
        - type: Plain
          text: hello
        - type: Plain
          text: 你好
    do:
      message:
        - type: Plain
          text: Hello World!
        - type: Plain
          text: （你好 世界！）

group:
# 当接收到的群消息为 say 时，调用「一言API」，原文发送接口返回的消息
  - when:
      message:
        - type: Plain
          text: say
    do:
      message:
        - type: Plain
          url: https://v1.hitokoto.cn?encode=text

# 当接收到的群消息为 jsay 时，调用「一言API」，解析返回后数据并拼接成文本消息发送
  - when:
      message:
        - type: Plain
          text: jsay
    do:
      message:
        - type: Plain
          url: https://v1.hitokoto.cn?encode=json&charset=utf-8
          jtext: '{hitokoto} ——— {from}'
```
<!-- 
+ `global`: 表示配置在接收到好友消息和群消息时都会生效。`friend`表示仅好友消息；`group`表仅群消息
+ `when`: 动作触发条件，满足任意一个即可触发。
+ `message`: 消息，写在`when`下表示接收到指定消息后触发，写在`do`下表示执行的动作
    + 写在`when`下: 表示任意一个消息即可触发，如上面的配置表示收到 hello 或 你好 时就触发动作。
    + 写在`do`下: 表示执行的动作，执行顺序为从上到下，如上面的配置表示动作为发送文本消息 「Hello World！（你好 世界！）」
+ `type`：消息类型
    + `Plain`：文本消息
+ `text`: 当`type`为 `Plain`时代表发送后面的原文。 -->


# 许可证

[GNU AFFERO GENERAL PUBLIC LICENSE version 3](https://github.com/ElpsyCN/el-bot-go/blob/master/LICENSE)

# 维护者

+ [ADD-SP](https://github.com/ADD-SP)
+ [YunYouJun](https://github.com/YunYouJun)

# 感谢

+ [mirai-api-http](https://github.com/mamoe/mirai-api-http)
+ [mirai-console](https://github.com/mamoe/mirai-console)
+ [mirai-console-wrapper](https://github.com/mamoe/mirai-console-wrapper)
+ [gomirai](https://github.com/Logiase/gomirai)