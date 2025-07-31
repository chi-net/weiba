# weiba

## 不仅仅是一个简单的Telegram频道/群组入群/频验证解决方案

## 小功能

1. GIC群组/频道验证以及白名单功能
2. 检测群组/频道的加入退出情况并通知给管理员
3. 频道附属讨论群组自动unpin频道发出的自动pin消息
4. 类似于Livegram的聊天机器人(可以让对方在不直接PM的情况下通过bot和你聊天)
5. 贴贴功能 发送 "/贴" 即可体验！

## 如何使用

注意：你需要登录Telegram桌面版APP才能导出数据

*(目前1,2,4步骤皆为可选步骤，如果你需要进行GIC验证才需要做，不需要的话，只要把`config.yml` `features`里面的`gic_auth`设置为`false`即可。)*
1. 使用桌面客户端导出数据，数据只勾选Private Groups和Private Channels(当然你想验公开群也可以一并勾选，建议把仅导出我的消息去掉)，导出格式选择机器可读JSON
2. 进入导出文件夹，把`result.json`放置与`convert.py`同目录，运行`convert.py`(此时需要关注是不是全部都是您想要验证的频道和群组内容，如果里面存在其他内容，可以编辑`convert.py`的skip列表自行跳过) 获得`result_converted.json`
3. 编辑`config.sample.yml` 设置机器人token和管理员的UID 修改为`config.yml`并与编译出来的Golang APP放置于同一目录
4. 将`result_converted.json`更名为`data.json`与编译出来的Golang APP(可在releases下载)放置于同一目录，运行即可

如果你不需要GIC验证功能的话，直接获取`config.sample.yml`进行编辑并重命名为`config.yml`，放置在与可执行文件同一目录即可。

你也可以使用容器服务：

```shell
podman run -dit --name weibabot -v config.yml:/app/config.yml -v data.json:/app/data.json ghcr.io/chi-net/weiba
```

## Group in Common的验证原理

weiba依赖你杜叔叔帐号所加入的私有群组和频道消息进行验证，由于杜叔叔中文圈子是一个一个小群组的交集与并集，因此有很多人在看到一个陌生人的时候都会不自主的去翻他的Group In Common以对他进行验证。

验证的简单原理就是把私有群组和频道的消息导出 生成私有群组和频道的链接 如果那个验证者并没有加入那个群组/频道的话，他是没有办法访问里面的内容的，反之则可以获取消息内容，再将消息内容与数据库中的进行比对即可验证是否存在共同私群。

隐私声明：本应用在`convert.py`阶段下就将原文信息进行了sha1哈希加密 后续都进行的是sha1校验，仅收集了消息的id和哈希 不会存在任何原文消息泄漏的情况

## 开源协议 & 免责声明

本应用的开源协议是 GPL3

免责声明：您不得将本应用用于任何违反任何地区法律法规，由您本人使用此应用造成的所有后果，开发者概不承担责任。