
##### `CHANGELOG.md`
```markdown
# Changelog

## v0.1.0
- support bond config (v0.1.0)
- 从配置文件启动wireguard-go
- 自动启动wireguard-go
- 自动配置wireguard-go接口地址
- 自动配置allowedip为本地路由
- 导入启动配置
- 打印当前wireguard-go状态信息
```
[bond]
bondname = band0
bondmode = active-backup
bestslavepeer = 8sXko+IFTPN8OZua5XEFL1Ur798nwpBg+3Ikbm1oXiI=
slavepeer = Am5cA8bgtVLUQGuFCdzNlW+oIygGL5eSAw0q+EdGCic=
AllowIPs = 0.0.0.0/1,128.0.0.0/1,8.8.8.8/32
```
- Initial release
- Added basic functionality


