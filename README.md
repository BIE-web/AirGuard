# AirGuard 项目

这是一个端到端的流量采集、特征提取与 VPN/穿透流识别与封禁演示平台，包含三个主要子模块：

- 数据采集节点 demo（`数据采集节点demo--autospider/`）：使用 Selenium 爬虫（spider）、v2ray 代理与 MySQL，通过 Docker Compose 快速部署示例流量生成与抓包环境。
- 特征提取工具（`特征提取工具--gotranal/`）：用 Go 实现的流（flow）特征提取器，支持在线抓包（libpcap）与离线 pcap 处理，输出 CSV 或通过 UDP 发送特征到识别服务。
- 网关识别与拦截（`网关识别拦截程序--VPN_Finder/`）：基于 Python 的接收/分析/告警与封禁模块，加载模型（`model/udp.pkl`）判定可疑 VPN 流，并通过 SSH 在路由器上下发 iptables 规则进行封禁，同时提供 Web 控制页面和日志展示。

## 总体工作流

1. 使用 `数据采集节点demo` 生成/代理流量并保存 pcap，或在采集节点上实时抓包。
2. 使用 `gotranal` 对 pcap 或网卡抓包数据进行六元组分流并计算 L3/L4 特征（如包长统计、到达间隔、TTL、TCP 窗口等），输出 CSV 或通过 UDP 发送至 `VPN_Finder`。
3. `VPN_Finder` 接收特征并用训练好的模型判断是否为 VPN/穿透流；若判定为可疑，会通过 SSH 在路由器上添加 iptables 规则进行阻断，并在 Web 界面展示日志与告警。

## 平台建议

- 建议在 Linux（或 Linux 虚拟机 / 服务器）上运行核心采集与抓包组件（尤其是 `gotranal` 的在线抓包模式），原因如下：
  - libpcap / 原始套接字和网卡访问在 Linux 上支持最好；容器可以通过 CAP_NET_RAW/CAP_NET_ADMIN 或 --net=host 获取抓包能力。
  - iptables 为 Linux 专有，封禁命令必须下发到支持 iptables 的设备（比如 Linux 路由器）。
  - Docker for Mac/Windows 在网络设备访问和特权方面有限制，不适合做真实的在线抓包或 iptables 操作。

## 快速开始（概览）

以下为从上到下的示例流程（假设在 Linux 主机）：

1) 启动数据采集 demo（生成流量 & pcap）

```bash
cd 数据采集节点demo--autospider
docker-compose up -d --build
```

2) 编译并运行 gotranal（在线模式，将特征发送到 VPN_Finder）

```bash
cd 特征提取工具--gotranal
go mod download
go build -o gotranal main.go
sudo ./gotranal on -I eth0 -b "tcp or udp" -s <VPN_FINDER_IP>:31115 -d 60
```

说明：在线模式通常需要 root 权限或 CAP_NET_RAW，或在容器中以特权模式运行。

3) 启动 VPN_Finder（接收并分析特征）

```bash
cd 网关识别拦截程序--VPN_Finder
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt  # 若不存在 requirements.txt 按 README 中列出的包安装: flask pandas joblib scikit-learn
python3 web.py
```

或单独运行分析进程（调试）：

```bash
python3 -c "import vpn_finder; vpn_finder.run('0.0.0.0', 31115)"
```

4) 抓取并发送离线 pcap（可在 macOS 本地导出 pcap 后上传到 Linux 处理）

```bash
./gotranal off -i /path/to/example.pcap -o out.csv -d 60
```

## 常见注意事项

- gotranal 在线抓包需要合适的抓包权限。若在容器中运行，请在 Linux 主机上使用 `--cap-add=NET_RAW --cap-add=NET_ADMIN --net=host` 或运行容器为特权模式。
- `VPN_Finder` 通过 SSH 在路由器上执行 iptables，请确保可以免交互 SSH（推荐使用密钥对）并在测试前确认回滚策略以免误封导致网络中断。
- 在开发阶段可以把 spider（Selenium 爬虫）在 macOS/本机上运行以便调试流量生成脚本，但 end-to-end 抓包/拦截建议在 Linux 环境验证。

## 目录与 README

- `数据采集节点demo--autospider/README.md`：包含 docker-compose、spider 内部调试与本地开发注意事项。
- `特征提取工具--gotranal/README.md`：包含编译、在线/离线运行示例与调试提示。
- `网关识别拦截程序--VPN_Finder/README.md`：包含配置、运行、ssh/iptable 注意事项与安全建议。

（每个子项目下已有更详细的 README，请按需查看。）

## 下一步建议

1. 如果你要在本机（macOS）做端到端测试：
  - 在本机上运行 `spider` 生成流量，并用 Wireshark/tshark 导出 pcap；将 pcap 上传到 Linux 主机并用 `gotranal off` 做离线特征提取与 `VPN_Finder` 测试。

2. 若能访问 Linux 主机：
  - 在 Linux 上以容器或二进制方式部署 `gotranal`（在线抓包），并把 `VPN_Finder` 放在同一网络中以减少网络配置问题。

3. 安全改进：
  - 为 `VPN_Finder` 的封禁操作增加幂等性和回滚接口，避免重复添加 iptables 规则。
  - 把关键配置（model 路径、阈值、ssh 凭证位置）移动到安全的配置或密钥管理中。

如果你愿意，我可以：
- 把 README 中的示例命令生成到一个 `scripts/` 目录，包含常用的 `docker-compose up`、gotranal build/run 与 vpn_finder 启动脚本；或
- 直接在你的 Linux 主机上（如果你授权并提供环境信息）帮你一步步部署并验证。请告诉我你下一步想做什么。
