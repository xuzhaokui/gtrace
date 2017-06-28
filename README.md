gtrace 使用指南
===

## QuickStart

**最少只需添加 3 行代码**即可完成请求的链路追踪，并默认会记录关键时间点、请求的基础信息（BodySize/Path/Host/Address/Hostname/Reqid 等）。


### 1. 初始化全局 tracer 配置

说明：若不对 trace 库做任何配置，默认不会记录任何 trace 相关的信息（但不妨碍代码正常调用）。

常见的配置过程如下（建议在服务启动后全局设置）：

__使用默认配置__：

```
gtrace.TracerEnable()
```

__指定服务名称__：

```
gtrace.TracerEnable(gtrace.SetService("<ServiceName>"))
```

其中：

- `ServiceName`: 用于告诉 trace 系统用什么名字指代该服务，不填则尝试取进程名为 serviceName



## 完整示例 demo

**TODO**