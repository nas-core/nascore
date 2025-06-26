# nascore

[English](README.md) | [中文](README_zh.md)

一个轻量可以扩展的类似nextcloud filebrowser的项目

## 特性

- [x] 可以在低性能硬件上运行,可以在高性能硬件上拓展。
- [x] webdav功能
- [x] 多用户管理不依赖数据库,使用toml格式
- [x] webui的实现 登陆 上传 下载 剪切 复制 文本编辑 粘贴板 右键功能
- [x] 自动反向代理 DDNS-go
- [x] 自动反向代理 AdGuard
- [x] 可选使用rclone作为后端文件系统

## 计划书

### 已经完成

- [x] webui 可视化多用户管理 可视化配置
- [x] webui 多文件上传
- [x] webui 文件夹上传
- [x] webui 在线解压 zip rar 7z tar tar.gz tar.xz tagz（部分依赖系统命令 且只有linux可用）
- [x] webui 在线压缩 tagz 格式
- [x] webui 图片预览和开关
- [x] webui 常见文本各种在线编辑
- [x] Rclone 跟随启动挂载命令

### 文件相关

- [ ] webui 视频缩略图和开关
- [ ] webui 图片缩略图压缩
- [ ] 目录缓存功能
- [ ] sftp
- [ ] nfs
- [ ] smb
- [ ] 在线更新
- [ ] 储存桶/s3直连

### 拓展

- [ ] webui 反向代理
- [ ] webui websocket chat
- [ ] webui mail center
- [ ] webui alist
- [ ] webui frps frpc 或者 ngrok
- [ ] webui 基于m3u8订阅的视频点播
- [ ] CalDAV
- [ ] CardDAV
- [ ] webui 一键安装 rclone
- [ ] webui 一键安装 ddns-go
- [ ] webui 一键安装 ADGuard-home
- [ ] webui 对 mosdns 的支持
- [ ] 包管理 优先级 openwrt debian nixos archlinux winget

## webdav 兼容性

litmus 测试

```sh
litmus -k webdav://username:password@192.168.1.101:9000
```

- [rfc4918](https://datatracker.ietf.org/doc/html/rfc4918)
- webdav 上游 [emersion/go-webdav](github.com/emersion/go-webdav)

## 关于源码

nascore主要采用golang开发，部分功能使用rust和c。没有太大技术含量，比如webdav部分完全依赖上游库。

因为很多原因，让我对这个项目目前就开源这个行为较为抵触，目前很多功能还在开发中，所以在功能完善之前不会完整开源。但nascore承诺

- 主程序的核心功能完善后 会开源主程序的完整的代码。
- nascore 尽量避免读取敏感配置文件，即便是必须读取的也不会收集相关信息。

你有以下方法可以获得源码。

- 本仓库公开了部分源码，包括完整的前端代码，和部分后端代码。
- 有时间有能力参与项目协作开发的，可以联系获取Private仓库权限。

### 合并请求的说明

- 因为webui较为简单，所以不打算引入npm/nodejs。可能会考虑基于js引入的vue/react方式。
- webui 要求尽可能的兼容低版本浏览器。
- api部分所有返回http状态码均为200，json的code区分错误码，原则上只支持get和post方式，且以post方式为主。
- 关于golang部分的gin,因为只用到了gin的路由和日志功能，目前性能测试良好。其他功能暂时不打算引入，但是handler需要一直兼容原生库存。后期计划引入gnet/nbio等性能更好的http库。
- nascore 主程序 尽量降低对系统依赖要求，并尽可能的系统核心功能。同时，nascore后续也会尽可能的减少对第三方库的依赖，以保证代码的稳定性和安全性。
