# Nezha-Paper
一个可能具有魔法的哪吒面板

## 上游

[https://github.com/nezhahq/nezha](https://github.com/nezhahq/nezha)

## 魔法改动

- 支持白名单跨域域名, 允许前端全静态, 配合我的前端分支使用
  [https://github.com/Sunakier/nezha-paper](https://github.com/Sunakier/nezha-paper)

## 配置说明

### WebSocket跨域配置

在配置文件中添加 `ws_allow_origins` 参数可以控制WebSocket允许的跨域来源。

```yaml
# 允许的WebSocket跨域来源，多个域名用逗号分隔，为空则只允许同源
ws_allow_origins: "example.com,api.example.org"
```

- 当此参数为空时，WebSocket连接只允许同源请求
- 当此参数不为空时，WebSocket连接将允许来自指定域名的跨域请求
- 可以使用多个域名，用逗号分隔
- 域名格式应为 `domain.com` 或 `subdomain.domain.com`，不需要包含协议（http/https）

此功能可以让您将前端部署在不同的域名下，同时保持与后端API的WebSocket通信。
