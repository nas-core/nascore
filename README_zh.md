# nascore

[English](README.md) | [中文](README_zh.md)

一个轻量可以扩展的类似nextcloud filebrowser的项目

## 特性

- [x] 可以在低性能硬件上运行,可以在高性能硬件上拓展。
- [x] webdav功能 [rfc4918](https://datatracker.ietf.org/doc/html/rfc4918)
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
- [x] webui 在线解压 zip rar 7z tar tar.gz tar.xz tagz（部分依赖系统命令）
- [x] webui 在线压缩 tagz 歌手
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
- [ ] webui 在线盗版视频点播
- [ ] CalDAV
- [ ] CardDAV

## webdav 兼容性

```sh
nix-shell -p litmus
litmus -k webdav://yh:123@192.168.1.101:9000
```

- 关注 https://github.com/hacdias/webdav/issues/218 #RFC2518

## webdav fork

- webdav fork自 [hacdias/webdav/](https://github.com/hacdias/webdav/commit/04cee682fb42c7f684eba106c14aff3ba2fa20c0)
- ai给校验 到 https://github.com/hacdias/webdav/commit/5676a1c3823a3643a7025feeb85157f83cd91b0f
