# 主要信息

+ 编号：0
+ 等级：Low
+ 发现人：ADD-SP
+ 修复人：ADD-SP
+ 发现日期：2020-05-06
+ 修复日期：2020-05-06
+ 缺陷所属模块：
    + Controller
    + ImageDoer
    + Message
+ 缺陷所属版本：v0.4.2
+ 缺陷状态：Closed

# 描述

## 使用的配置

```yml
  - when:
      sender:
        group:
          - 469327964
      message:
        - type: Plain
          regex: '.+'
        - type: Image
    do:
      receiver:
        group:
          - 1072803190
      message:
        - type: Plain
          text: '{el-message-text}'
        - type: Image
          url: '{el-message-image-url-0}'
        - type: Image
          url: '{el-message-image-url-1}'
        - type: Image
          url: '{el-message-image-url-2}'
        - type: Image
          url: '{el-message-image-url-3}'
```

## 症状

当消息只包含图片时不会发送任何消息。

当消息包含文本和图片时消息会重复一次。