# el-bot-go

[![docs](https://github.com/ElpsyCN/el-bot-docs/workflows/docs/badge.svg)](https://docs.bot.elpsy.cn)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ElpsyCN/el-bot-go)
[![GitHub release (latest by date including pre-releases)](https://img.shields.io/github/v/release/ElpsyCN/el-bot-go?color=blue&include_prereleases)](https://github.com/ElpsyCN/el-bot-go/releases)
[![QQ Group](https://img.shields.io/badge/qq%20group-707408530-12B7F5)](https://shang.qq.com/wpa/qunwpa?idkey=5b0eef3e3256ce23981f3b0aa2457175c66ca9194efd266fd0e9a7dbe43ed653)
[![GitHub issues](https://img.shields.io/github/issues/ElpsyCN/el-bot-go)](https://github.com/ElpsyCN/el-bot-go/issues)
[![GitHub](https://img.shields.io/github/license/ElpsyCN/el-bot-go?color=%233eb370)](https://github.com/ElpsyCN/el-bot-go/blob/master/LICENSE)

[el-bot](https://github.com/ElpsyCN/el-bot)的 go 版本。

一个基于 Mirai 的可快速配置的 QQ机器人 模板。

# 文档

https://docs.bot.elpsy.cn

# 功能

只列出已经实现的功能，其它功能见[开发进度](https://github.com/ElpsyCN/el-bot-go/projects/1)。

+ 插件系统
+ 识别文本消息
  + 识别固定文本消息
  + 通过正则表达式识别文本消息
  + At & AtAll
+ 发送文本消息
  + 文本消息中支持嵌入一些预定义变量例如消息发送者昵称
  + 发送固定文本消息
  + 原文发送来自网络的文本
  + 发送解析后的来自网络的 JSON 文本
  + 发送通过`预定义变量`修饰的文本消息
  + At & AtAll
+ 识别事件
  + 新成员入群
  + 踢人 & 自己退群
  + 禁言 & 全员禁言
  + 全员禁言 & 解除全员禁言
+ 执行动作
  + 禁言/全员禁言
  + 解除禁言/解除全员禁言
+ 发送表情消息
  + 发送固定表情
+ 识别表情消息
  + 识别固定表情消息
+ 发送图片消息
  + 发送本地图片
  + 发送网络图片
+ 定时任务
  + 定时发送消息
+ Bot 控制
  + 暂停机器人
  + 恢复机器人
  + 退出机器人
  + 暂时屏蔽某些消息接收者
  + 取消屏蔽某些消息接收者
+ 消息自动转发
+ RSS 订阅
+ 配置触发次数统计
+ 反刷屏

# 反馈

有问题和建议欢迎提 Issue，谢谢！（在此之前，请确保您已仔细阅读文档。）

您也可以加入 QQ 群（707408530）进行反馈与讨论。

> 如果是通用的问题（如 BUG 反馈，新功能建议），最好在 Issue 中进行反馈，以便其他朋友参与讨论，减少重复。

# 已知问题



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