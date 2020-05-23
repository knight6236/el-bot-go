# 配置语法说明

# 通用配置

## 结构说明

通用配置通常具有下面这种结构

```yml
area:
    - when:
        sender: 
            group:
                - 群号
            user:
                - 好友/群成员QQ号
        message:
            - type: Plain
              text: 文本消息
              regex: 正则表达式
            - type: Image
              path: 图片路径（相对于 plugins\MiraiAPIHTTP\setting.yml）
              url: 图片 URL
              direct: true | false
            - type: Face
              path: piezui
            - type: Xml
              text: '{el-message-xml}'
        operation:
            - type: MemberMute | MemberUnmute | GroupMuteAll | GroupUnMuteAll
      do:
        receiver:
          group:
                - 群号
          user:
                - 好友/群成员QQ号
        message:
          - type: Plain
            text: 文本消息
            regex: 正则表达式
          - type: Image
            path: 图片路径（相对于 plugins\MiraiAPIHTTP\setting.yml）
            url: 图片 URL
            direct: true | false
          - type: Face
            path: piezui
          - type: Xml
            text: '{el-message-xml}'
```

+ area: 表示配置的生效范围，取值只能为下面所列 
    + global: 接收到「好友消息」和「群消息」时均生效
    + friend：仅对接收到的「好友消息」生效
    + group: 仅对接收到「群消息」生效
+ when: 所有 when 下的配置均作为配置的触发条件，当且仅当 `sender`、`message`和`operation`均返回 true 是触发配置，执行动作。
+ sender: 表示消息的发送者，其下包含
    + group: 可包含一个或多个群号，表示接收到的消息所在的群号。满足任意一个 sender 就返回 true
    + user: 可包含一个或多个 QQ 号，表示接收到的消息的发送者的 QQ 号。满足任意一个 sender 就返回 true
+ message: 表示消息。当 `when`下的`message`中存在多个类型的消息时，符合任意一个均返回 true
    + Plain: 文本消息
        + text：文本消息的内容
        + regex：尝试使用给定的正则表达式匹配文本消息
        + url: 文本来自 GET 请求
        + json: GET 返回的消息为 JSON 格式时为 true，反之可以忽略。
    + Image：图片消息
        + url：图片来自 GET 请求
        + path: 图片来自本地，图片路径相对于 plugins\MiraiAPIHTTP\setting.yml
        + direct: 如果 URL 会重定向到另外一个 URL则为 ture，反之可以忽略。
    + Face: 表情消息
        + name: 表情名称
    + Xml: XML 消息
        + text: XML 文本
+ operation: 表示操作/事件，符合任意一个就返回 true，若不写 operation 则默认为 true。
    + MemberMute:
    + MemberUnmute:
    + GroupMuteAll:
    + GroupUnMuteAll: 
+ do: 所有 do 下的配置均位要执行的动作/发送的消息
+ receiver: 表示消息的接收者，其下包含
    + group: 可包含一个或多个群号，表示将消息发送给的群的群号
    + user: 可包含一个或多个 QQ 号，表示将消息发送给的还有/群成员的 QQ 号

## 预定义变量

有时我们可能需要一些额外信息，如本次接收到的消息，新进群的成员的 QQ 号，被禁言的成员的群名片等。

[预定义变量](pre-def-var.md)

## 配置举例

### 一问一答型

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
```

### 复读机型

```yml
global:
  # 当文本消息符合正则表达式时复读本次消息
  - when:
      message:
        - type: Plain
          regex: 复读
    do:
      message:
        - type: Plain
          text: '{el-message-text}'
```

### 发送网络图片

```yml
global:
  # 当文本消息为 img 时发送 URL 指向的图片
  - when:
      message:
        - type: Plain
          text: img
    do:
      message:
        - type: Image
          url: https://xxxx
```

### 成员被禁言时发送消息

```yml
group:
  # 当某个成员被禁言时发送「「被禁言成员群昵称」喜提禁言套餐」
  - when:
      operation:
        - type: MemberMute
    do:
      message:
        - type: Plain
          text: 「{el-target-name}」喜提禁言套餐
```

更多例子见 `config/default.yml`

# 定时任务

定时任务举例

```yml
crontab:
  - cron: '0 * * * * *'
    do:
      receiver:
        group:
          - 群号
        friend:
          - QQ号
      message:
        - type: Plain
          url: 一分钟过去了
```

# 配置触发统计

配置触发统计举例

```yml
  - countID: 打招呼
    when:
      message:
        - type: Plain
          text: hello
        - type: Plain
          text: 你好
    do:
      count: true
      message:
        - type: Plain
          text: Hello World!
        - type: Plain
          text: （你好 世界！）

   - when:
      message:
        - type: Plain
          text: 统计结果
    do:
      count: true
      message:
        - type: Plain
          text: '{el-bot-overall}'
```
