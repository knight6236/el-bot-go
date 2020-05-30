# 配置语法说明

# 通用配置

## 生效范围

```yml
global:
friend:
group:
```

| 关键字 | 必要 | 说明                                             |
| ------ | ---- | ------------------------------------------------ |
| global | 否   | 此关键字下的配置在接收到好友消息和群消息时均生效 |
| friend | 否   | 此关键字下的配置在接收到好友消息时生效           |
| group  | 否   | 此关键字下的配置群消息时生效                     |



## 基本结构

```yml
when:
  countID:
  operation:
  message:
do:
  isCount:
  operation:
  message:
```

| 关键字    | 必要 | 说明                                                |
| --------- | ---- | --------------------------------------------------- |
| when      | 是   | 表示此配置触发的条件。                              |
| do        | 是   | 表示此配置触发后执行的动作。                        |
| countID   | 否   | 当 `count: true` 是必须填写，用于记录配置触发日志。 |
| operation | 否   | 表示一些动作或事件，如新成员入群，禁言等。          |
| message   | 否   | 表示消息，包括文本、图片、表情和 XML。              |

## message

```yml
message:
  sender:
  receiver:
  detail:
```

| 关键字   | 必要 | 接受消息时的作用 | 发送消息时的作用 |
| -------- | ---- | ---------------- | ---------------- |
| sender   | 否   | 消息的发送者     | 无               |
| receiver | 否   | 无               | 消息的接收者     |
| detail   | 是   | 消息详情         | 消息详情         |

## detail

```yml
detail:
  - type:
    other:
```

| 关键字 | 必要 | 接受消息时的作用                        | 发送消息时的作用 |
| ------ | ---- | --------------------------------------- | ---------------- |
| type   | 是   | 接收到的消息的类型                      | 发送的消息的类型 |
| other  | 否   | 根据 `type`的不同会有不同的名字和作用。 |                  |

### 消息类型

| 关键字 | 接受消息时的作用      | 发送消息时的作用      | 附属关键字 | 接受消息时附属关键字的作用         | 发送消息时附属关键字的作用                                   |
| ------ | --------------------- | --------------------- | ---------- | ---------------------------------- | ------------------------------------------------------------ |
| Plain  | 表示接收到的文本消息  | 表示发送的文本消息    | text       | 表示要接受到的文本                 | 表示要发送的文本                                             |
|        |                       |                       | regex      | 使用对应正则表达式匹配接收到的文本 | 无                                                           |
|        |                       |                       | url        | 无                                 | 向`url`发送`GET`请求获取返回的文本                           |
|        |                       |                       | json       | 无                                 | 表示`url`返回的文本是否为`json`                              |
| Image  | 表示接收到的图片      | 表示发送的图片        | url        | 无                                 | 要发送的图片的 URL                                           |
|        |                       |                       | path       | 无                                 | 要发送的图片的路径（相对于`plugins/MiraiAPIHTTP/images`）    |
|        |                       |                       | reDirect   | 无                                 | 如果要发送的图片的 URL 会重定向到其它  URL 则设置为`true`，反之则忽略。 |
| Face   | 表示接收到的表情      | 表示发送的表情        | faceName   | 接收到的表情的名称                 | 要发送的表情的名称（详见`config/face-map.yml`）              |
| Xml    | 表示接收到的 XML 文本 | 表示要发送的 XML 文本 | text       | 无                                 | 表示要发送的 XML 文本内容                                    |


## sender

```yml
sender:
  group:
    - 群号
    ...
  user:
    - QQ号
    ...
```

| 关键字 | 必要 | 说明                                                         |
| ------ | ---- | ------------------------------------------------------------ |
| group  | 否   | 表示消息来源的群号，可包括若干个群号。                       |
| user   | 否   | 表示消息来源的「群成员」或好友的「QQ号」，可以包括若干个「QQ号」。 |


## receiver

```yml
receiver:
  group:
    - 群号
    ...
  user:
    - QQ号
    ...
```

| 关键字 | 必要 | 说明                                                         |
| ------ | ---- | ------------------------------------------------------------ |
| group  | 否   | 表示接受消息的群的群号，可包括若干个群号。                   |
| user   | 否   | 表示接受消息的「群成员」或好友的「QQ号」，可以包括若干个「QQ号」。 |

## operation

```yml
operation:
  - type: 
    other:
```

| 关键字 | 必要 | 接受消息时的作用                        | 发送消息时的作用                        |
| ------ | ---- | --------------------------------------- | --------------------------------------- |
| type   | 是   | 表示事件/操作的类型。                   | 表示操作的类型。                        |
| other  | 否   | 根据 `type`的不同会有不同的名字和作用。 | 根据 `type`的不同会有不同的名字和作用。 |

### 事件/操作类型

| 关键字            | 接受消息时的作用     | 发送消息时的作用              | 附属关键字 | 接受消息时附属关键字的作用   | 发送消息时附属关键字的作用                 |
| ----------------- | -------------------- | ----------------------------- | ---------- | ---------------------------- | ------------------------------------------ |
| At                | 表示某成员被 At      | At 某成员                     | groupID    | 无                           | 被 At 的成员所在群的「群号」               |
|                   |                      |                               | userID     | 被 At 的成员所在群的「群号」 | 被 At 的成员的「QQ号」                     |
| AtAll             | 管理员 At 了全体成员 | At 全体成员                   | groupID    | 无                           | 接收「@全体成员」消息的群的群号            |
| MemeberMute       | 某成员被禁言         | 「禁言 」某个群的某个成员     | groupID    | 无                           | 被 「禁言 」的成员所在群的「群号」         |
|                   |                      |                               | userID     | 无                           | 被 「禁言 」的成员的「QQ号」               |
| MemberUnMute      | 某成员被解除禁言     | 「解除禁言 」某个群的某个成员 | groupID    | 无                           | 被 「解除禁言 」的成员所在群的「群号」     |
|                   |                      |                               | userID     | 无                           | 被 「解除禁言 」的成员的「QQ号」           |
| GroupMuteAll      | 管理员开启群员禁言   | 开启某个群的「全员禁言 」     | groupID    | 无                           | 被 「全员禁言 」的成员所在群的「群号」     |
| GroupUnMuteAll    | 管理员关闭群员禁言   | 解除某个群的「全员禁言 」     | groupID    | 无                           | 被 「解除全员禁言 」的成员所在群的「群号」 |
| MemberJoin        | 新成员入群           | 无                            | 无         | 无                           | 无                                         |
| MemberLeaveByKick | 成员被管理员移除     | 无                            | 无         | 无                           | 无                                         |
| MemberLeaveByQuit | 成员主动退群         | 无                            | 无         | 无                           | 无                                         |

## 触发规则

配置的触发规则为 `sender && (message || operation)` ，即 `sender` 必须返回 `true`，在此条件下 `message` 和 `operation` 任意一个返回 `true` 即可。

通常情况下`message`下会有多个消息类型，满足任意一个 `message` 就会返回`true`。

通常情况下`sender`下也会有多个接受消息的群号和 QQ号，满足任意一个群号或 QQ号`sender`就会返回`true`。

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
  # 当文本消息符合正则表达式时复读第一个子表达式的内容并 At 消息发送者
  - when:
      message:
          detail:
          - type: Plain
            regex: 复读\s(.+)
    do:
      operation:
        - type: At
          groupID: '{el-sender-group-id}'
          userID: '{el-sender-user-id}'
      message:
        detail:
          - type: Plain
            text: '{el-regex-0}'
```

### 发送网络图片

```yml
global:
  # 当文本消息为 img 时发送 URL 指向的图片
  - when:
      message:
        detail:
          - type: Plain
            text: img
    do:
      message:
        detail:
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
        detail:
          - type: Plain
            text: 「{el-target-user-name}」喜提禁言套餐
```

更多例子见 [config/default.yml](../config/default.yml)


# 预定义变量

有时我们可能需要一些额外信息，如本次接收到的消息，新进群的成员的 QQ 号，被禁言的成员的群名片等。

[预定义变量](pre-def-var.md)

# 定时任务

## 基本结构

```yml
crontab:
  - cron: '* * * * * *'
    do:
      operation:
      message:
```

| 关键字  | 说明                               |      |
| ------- | ---------------------------------- | ---- |
| crontab | 此关键字下的配置会被解析为定时任务 |      |
| cron    | cron 表达式，精确到秒。            |      |

区别于「通用配置」，定时任务中的 `do.message.receiver` 是必须的，不允许省略，即至少要指定一个消息的接收者。

## 配置举例

```yml
# 每隔一分钟就发送消息「一分钟过去了」
crontab:
  - cron: '0 * * * * *'
    do:
      message:
      	receiver:
          group:
            - 接收消息的群号
        detail:
          - type: Plain
            text: 一分钟过去了
```

# 消息转发

## 基本结构

```yml
transfer:
  - listen:
      group:
      user:
    target:
      grpup:
      user:
  - listen:
      group:
      user:
    target:
      group:
      user:
```

| 关键字   | 是否必要 | 说明                           |      |
| -------- | -------- | ------------------------------ | ---- |
| transfer | 是       | 该关键字下的配置为消息转发配置 |      |
| listen   | 是       | 表示监听哪些群和好友的消息     |      |
| group    | 否       | 表示被监听/接收消息的若干个群  |      |
| user     | 否       | 表示被监听/接收消息的若干好友  |      |
| target   | 是       | 消息的接收者                   |      |



## 配置举例

```yml
# 当接收到的指定的群或指定的好友的消息时，自动将消息转发给指定的群和好友
transfer:
  - listen:
      group:
        - 群号
        - 群号
        ...
      user:
        - QQ号
        - QQ号
        ...
    target:
      group:
        - 群号
        - 群号
        ...
      user:
        - QQ号
        - QQ号
        ...
```

# RSS 订阅

## 基本结构

```yml
rss:
  - url: 
    do:
      operation:
      message:
        receiver:
```

| 关键字 | 必要 | 说明                           |
| ------ | ---- | ------------------------------ |
| rss    | 是   | `rss`下的配置均为`rss`订阅配置 |
| url    | 是   | `rss`URL                       |

区别于「通用配置」，定时任务中的 `do.message.receiver` 是必须的，不允许省略，即至少要指定一个消息的接收者。

## 配置举例

```yml
rss:
  - url: https://xxxx/atom.xml
    do:
      message:
        receiver:
          group:
            - 群号
        detail:
          - type: Plain
            text: '标题：{el-rss-title}{\n}'
          - type: Plain
            text: '作者：{el-rss-author}{\n}'
          - type: Plain
            text: '链接：{el-rss-link}{\n}'
          - type: Plain
            text: '时间：{el-rss-year}-{el-rss-month}-{el-rss-day} {el-rss-hour}:{el-rss-minute}:{el-rss-second}'
```

关于上面配置中所用到的 {el-rss-link} 等预定义变量的说明请参考[预定义变量](pre-def-var.md)