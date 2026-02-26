# 特征提取工具 — gotranal

这是一个用 Go 实现的流（flow）特征提取工具。它可以从离线 pcap 文件或实时网络接口抓包，按六元组（srcIP:srcPort, dstIP:dstPort, proto）聚合流并计算多种 L3/L4 特征，输出为 CSV 或通过 UDP 发送给下游服务（例如 `VPN_Finder`）。

主要功能
- 在线模式（`on`）：直接从网卡抓包并按时间窗口导出或发送特征。
- 离线模式（`off`）：读取 pcap 文件并对其中的流进行特征计算，输出 CSV。

依赖
- Go 1.16+（建议最新稳定版）。
- libpcap（系统库）。在 macOS 可通过 `brew install libpcap`，在 Debian/Ubuntu 上是 `libpcap-dev`。
- Go module 依赖在仓库根目录通过 `go.mod` 管理。

编译

```bash
cd 特征提取工具--gotranal
go mod download
go build -o gotranal main.go
```

使用示例

- 离线处理 pcap（生成 CSV）：

```bash
./gotranal off -i /path/to/sample.pcap -o out_features.csv -d 60
```

- 在线抓包并把结果发送给远端识别服务（例如 `VPN_Finder`）：

```bash
# 在带有抓包权限的主机上运行（或通过 ssh 在路由器/采集节点上运行）
./gotranal on -I any -b "tcp or udp" -s 192.168.0.2:31115 -d 60
```

参数说明（常用）
- `on` / `off`：在线或离线子命令。
- `-I`：抓包网卡（如 eth0、any 等）。
- `-b`：BPF 过滤器，默认 `tcp or udp`。
- `-s`：发送特征至的服务器地址（host:port）。仅在线模式有效。
- `-i`：离线输入 pcap 路径（离线模式必需）。
- `-o`：输出文件名或前缀。
- `-d`：时长（秒），在线模式表示每个时间窗口的 duration。

调试与注意事项
- 抓包通常需要 root 权限，在线模式请以 root 或 sudo 权限运行。
- 若在容器中运行，请保证容器有 CAP_NET_RAW 权限或以特权模式运行以允许抓包。
- 当输出特征发送到 `VPN_Finder` 时，需确保目标 IP/端口和防火墙策略允许 UDP 流量。

文件与代码位置
- `main.go`：命令行解析、在线/离线模式调度。
- `flowDuration/`：流特征计算的实现（`flow.go` 为核心逻辑）。

扩展建议
- 增加日志和速率限制，避免短时间大量 UDP 发送造成网络拥塞。
- 将输出格式文档化（字段顺序与含义），以便下游解析器更健壮地兼容。
