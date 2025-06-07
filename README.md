# Nezha-Paper
一个可能具有魔法的哪吒面板

## 上游

[https://github.com/nezhahq/nezha](https://github.com/nezhahq/nezha)

## 魔法改动

- 支持白名单跨域域名, 允许前端全静态，完全前后端分离, 配合我的前端分支使用
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

### 服务器国家码强制指定

在编辑服务器信息时，可以通过设置 `force_country_code` 参数来强制指定服务器显示的国家码，覆盖自动检测的地理位置信息。

#### API变动

服务器编辑API (`PATCH /server/{id}`) 新增字段：

```json
{
  "force_country_code": "cn" // 可选，强制指定国家码，使用小写的ISO国家代码
}
```

- 当此参数为空时，系统将使用GeoIP自动检测的国家码
- 当此参数不为空时，系统将优先使用指定的国家码，忽略GeoIP检测结果
- 国家码应使用小写的ISO国家代码，如：cn（中国）、us（美国）、jp（日本）等

此功能可以让您手动纠正错误的地理位置检测结果，或者为特定服务器指定自定义的显示位置。
