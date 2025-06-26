# nascore

English | [中文](README_zh.md)

A lightweight and extensible project similar to nextcloud and filebrowser.

## Features

- [x] Can run on low-performance hardware and scale on high-performance hardware.
- [x] WebDAV functionality
- [x] Multi-user management without database, using TOML format
- [x] WebUI implementation: login, upload, download, cut, copy, text editing, clipboard, right-click functionality
- [x] Automatic reverse proxy DDNS-go
- [x] Automatic reverse proxy AdGuard
- [x] Optional use of rclone as backend file system

## Roadmap

### Completed

- [x] WebUI visual multi-user management and visual configuration
- [x] WebUI multi-file upload
- [x] WebUI folder upload
- [x] WebUI online decompression of zip rar 7z tar tar.gz tar.xz tagz (partially depends on system commands and only available on Linux)
- [x] WebUI online compression to tagz format
- [x] WebUI image preview and switch
- [x] WebUI common text editing

- [x] Rclone mount command on startup

### File Related

- [ ] WebUI video thumbnail and switch
- [ ] WebUI image thumbnail compression
- [ ] Directory caching
- [ ] SFTP
- [ ] NFS
- [ ] SMB
- [ ] Online update
- [ ] Storage bucket / S3 direct connection

### Extensions

- [ ] WebUI reverse proxy
- [ ] WebUI websocket chat
- [ ] WebUI mail center
- [ ] WebUI alist
- [ ] WebUI frps frpc or ngrok
- [ ] WebUI m3u8 subscription-based video on demand
- [ ] CalDAV
- [ ] CardDAV
- [ ] WebUI one-click installation of rclone
- [ ] WebUI one-click installation of ddns-go
- [ ] WebUI one-click installation of ADGuard-home
- [ ] WebUI support for mosdns
- [ ] Package management priority: openwrt debian nixos archlinux winget

## WebDAV Compatibility

litmus test

```sh
litmus -k webdav://username:password@192.168.1.101:9000
```

- [rfc4918](https://datatracker.ietf.org/doc/html/rfc4918)
- WebDAV upstream [emersion/go-webdav](github.com/emersion/go-webdav)

## About Source Code

nascore is mainly developed in Golang, with some features using Rust and C. There is no significant technical complexity, for example, the WebDAV part relies entirely on the upstream library.

For many reasons, I am currently resistant to the act of open-sourcing this project. Many features are still under development, so a complete open source will not be available until the features are complete. But nascore promises

- The complete code of the main program will be open sourced after the core functions of the main program are completed.
- nascore tries to avoid reading sensitive configuration files, and even if it must be read, it will not collect relevant information.

You have the following methods to obtain the source code.

- This repository discloses some source code, including the complete front-end code and some back-end code.
- If you have the time and ability to participate in project collaboration and development, you can contact us to obtain Private repository permissions.

### Instructions for Merge Requests

- Because the webui is relatively simple, there are no plans to introduce npm/nodejs. May consider vue/react-based introduction based on js.
- webui requires as much compatibility as possible with lower version browsers.
- All API return HTTP status codes are 200, and the JSON code distinguishes error codes. In principle, only GET and POST methods are supported, and POST methods are the main method.
- Regarding the gin part of Golang, because only gin's routing and logging functions are used, the current performance test is good. Other functions are not planned to be introduced for the time being, but the handler needs to be compatible with the native library. In the future, it is planned to introduce HTTP libraries with better performance such as gnet/nbio.
- The nascore main program tries to reduce the system dependency requirements and the system core functions as much as possible. At the same time, nascore will also reduce the dependence on third-party libraries as much as possible in the future to ensure the stability and security of the code.
