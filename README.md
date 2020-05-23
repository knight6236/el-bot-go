# el-bot-go

[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/ElpsyCN/el-bot-go?color=blue&include_prereleases)](https://github.com/ElpsyCN/el-bot-go/releases)
[![GitHub issues](https://img.shields.io/github/issues/ElpsyCN/el-bot-go)](https://github.com/ElpsyCN/el-bot-go/issues)
[![GitHub closed issues](https://img.shields.io/github/issues-closed/ElpsyCN/el-bot-go?color=success)](https://github.com/ElpsyCN/el-bot-go/issues)
[![GitHub](https://img.shields.io/github/license/ElpsyCN/el-bot-go?color=%233eb370)](https://github.com/ElpsyCN/el-bot-go/blob/master/LICENSE)

[el-bot](https://github.com/ElpsyCN/el-bot)的 go 版本。

一个基于 Mirai 的可快速配置的 QQ机器人 模板。

# 文档

[如何编写配置](docs/config-syntax.md)

# 功能

只列出已经实现的功能，其它功能见[开发进度](https://github.com/ElpsyCN/el-bot-go/projects/1)。

+ 识别文本消息
  + [x] 识别固定文本消息
  + [x] 通过正则表达式识别文本消息
+ 发送文本消息
  + [x] 文本消息中支持嵌入一些预定义变量例如消息发送者昵称
  + [x] 发送固定文本消息
  + [x] 原文发送来自网络的文本
  + [x] 发送解析后的来自网络的 JSON 文本
+ 识别事件
  + [x] 新成员入群
  + [x] 踢人/自己退群
  + [x] 禁言/全员禁言
  + [x] 全员禁言/解除全员禁言
+ 发送表情消息
  + [x] 发送固定表情
+ 识别表情消息
  + [x] 识别固定表情消息
+ 发送图片消息
  + [x] 发送本地图片
  + [x] 发送网络图片
+ 定时任务
  + [x] 定时发送消息
+ 配置触发次数统计

# 快速开始

## On Unix Like

1. 安装Golang和JRE
2. clone: `git clone git@github.com:ElpsyCN/el-bot-go.git`
3. 下载依赖: `./install.sh`
4. 由于安装 `golang` 的 `package` 比较麻烦，可以进入 [release](https://github.com/ElpsyCN/el-bot-go/releases) 下载对应的二进制文件到 `bin` 下，来使用编译好程序，跳过 `package `的安装的环节。
5. 创建文件`plugins/MiraiAPIHTTP/setting.yml`并填入下列内容
    ```yml
    authKey: qwertyuiop
    port: 8080
    enableWebsocket: false
    ```
6. 启动 mirai-console: ` ./start-console.sh`并按照提示进行操作
7. 启动 el-bot-go: 在另一个 shell 运行脚本：
    1. 选择合适的 shell 脚本 `start-el-bot-xxx-yyy.sh`
    2. `sh start-el-bot-xxx-yyy.sh 机器人QQ号`

## On Windows

1. 安装Golang和JRE
2. clone: `git clone git@github.com:ElpsyCN/el-bot-go.git`
3. 下载依赖：
    + 下载[mirai-console-wrapper-1.2.0-all](https://github.com/mamoe/mirai-console-wrapper/releases/download/1.2.0/mirai-console-wrapper-1.2.0-all.jar)到项目根目录
    + 下载[mirai-api-http-v1.7.0](https://github.com/mamoe/mirai-api-http/releases/download/v1.7.0/mirai-api-http-v1.7.0.jar)到`plugins/`
4. 由于安装 `golang` 的 `package` 比较麻烦，可以进入 [release](https://github.com/ElpsyCN/el-bot-go/releases) 下载对应的二进制文件到 `bin` 下，来使用编译好程序，跳过 `package `的安装的环节。
5. 创建文件`plugins/MiraiAPIHTTP/setting.yml`并填入下列内容
    ```yml
    authKey: qwertyuiop
    port: 8080
    enableWebsocket: false
    ```
6. 启动 mirai-console: `./start-console.bat`
7. 启动 el-bot-go: 在另一个 cmd 执行命令 `start-el-bot-xxxx.bat 机器人的QQ号 配置所在路径`
    1. 选择合适的 shell 脚本 `start-el-bot-xxx-yyy.sh`
    2. `start-el-bot-xxx.bat 机器人QQ号`

# 配置文件说明

<!-- config/custom/custom.yml -->

<details>
  <summary>点击查看</summary>

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
          text: '{el-url-text}'

  # 当接收到的群消息为 jsay 时，调用「一言API」，解析返回后数据并拼接成文本消息发送
  - when:
      message:
        - type: Plain
          text: jsay
    do:
      message:
        - type: Plain
          url: https://v1.hitokoto.cn?encode=json&charset=utf-8
          text: '{hitokoto} ——— {from}'
          json: true
  # 当某个成员被禁言时发送「「被禁言成员群昵称」喜提禁言套餐」
  - when:
      operation:
        - type: MemberMute
    do:
      message:
        - type: Plain
          text: 「{el-target-name}」喜提禁言套餐
  
  # 当某个成员被禁言时发送「恭喜「被禁言成员群昵称」出狱」
  - when:
      operation:
        - type: MemberUnmute
    do:
      message:
        - type: Plain
          text: '恭喜「{el-target-name}」出狱'

  # 当开启全体禁言时发送 「砸瓦鲁多！」
  - when:
      operation:
        - type: GroupMuteAll
    do:
      message:
        - type: Plain
          text: 砸瓦鲁多！
  
  #  当关闭全员禁言时发送「隐藏着黑暗力量的钥匙啊,在我面前显示你真正的力量！现在以你的主人，小樱之名命令你。封印解除！」
  - when:
      operation:
        - type: GroupUnMuteAll
    do:
      message:
        - type: Plain
          text: 隐藏着黑暗力量的钥匙啊,在我面前显示你真正的力量！现在以你的主人，小樱之名命令你。封印解除！
  
  # 当有新成员入群时发送「欢迎「新成员群昵称」进群」
  - when:
      operation:
        - type: MemberJoin
    do:
      message:
        - type: Plain
          text: 欢迎「{el-target-name}」进群

  # 当某成员被移除群聊时发送「管理员赠送「被移除的成员的群昵称」飞机票一张」
  - when:
      operation:
        - type: MemberLeaveByKick
    do:
      message:
        - type: Plain
          text: 管理员赠送「{el-target-name}」飞机票一张
  
  # 当某成员自行退出群聊是发送「有大佬走了，群地位+1。」
  - when:
      operation:
        - type: MemberLeaveByQuit
    do:
      message:
        - type: Plain
          text: 有大佬走了，群地位+1。
  
  # 当文本消息符合正则表达式时复读本次消息
  - when:
      message:
        - type: Plain
          regex: 复读
    do:
      message:
        - type: Plain
          text: '{el-message-text}'

  # 当收到的表情消息为 「撇嘴」时发送表情「微笑」
  - when:
      message:
        - type: Face
          name: piezui
    do:
      message:
        - type: Face
          name: weixiao


# el-message-text: 本次的文本消息
# el-sender-id: 发送消息的好友/群成员QQ号
# el-sender-name: 发送消息的好友/群成员的名称
# el-operator-id: 做出操作的好友/成员的QQ号
# el-operator-name: 做出操作的还有/群成员的名称
# el-target-id: 某些事件的目标成员的QQ号，如禁言，新成员进群，移除群成员等
# el-target-name: 某些事件的目标成员的名称，如禁言，新成员进群，移除群成员等
```
</details>


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
+ [一言开发者中心](https://developer.hitokoto.cn/)