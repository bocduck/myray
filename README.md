this is AI work , use with caution

以下代码是一个最简go vless httpupgrade transport 服务端实现，其以反代形式工作在安全的tls web服务器背后，且在反代层以path做了安全认证，与现有客户端兼容。

审计有无安全问题、性能问题。
