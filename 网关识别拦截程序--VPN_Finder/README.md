# 网关识别与拦截 — VPN_Finder

此子项目负责接收来自 `gotranal` 的流特征，对流进行分类判定（是否为 VPN/穿透流），并在判定为可疑时通过 SSH 在路由器上下发 iptables 规则进行封禁，同时通过 Web 界面展示日志与告警。

主要组成
- `web.py`：Flask Web 控制界面，提供 `/start` `/stop` 接口，并启动 `websocketd` 将日志推送到前端。
- `vpn_finder.py`：核心分析逻辑，接收 UDP 特征流、用模型（`model/udp.pkl`）判断，并调用 SSH 执行 iptables 封禁。
- `model/`：保存训练好的模型文件（例如 `udp.pkl`）。
- `config.json`：运行配置（IP、端口、router_ip、iface 等）。

依赖（建议环境）
- Python 3.8+
- 建议创建虚拟环境并安装依赖（下面示例）：

```bash
cd 网关识别拦截程序--VPN_Finder
python3 -m venv .venv
source .venv/bin/activate
pip install flask pandas joblib scikit-learn
```

其他依赖：
- `websocketd` 二进制（用于把 `log` 文件通过 websocket 推送到前端）。如无需实时日志，可手动关闭或替换实现。

配置
- 编辑 `config.json`：设置 `ip`（本机/服务 IP）、`flask_port`、`websocketd_port`、`router_ip`（可 ssh 的路由器地址）、`iface`（gotranal 抓包的接口名）、`vpn_finder_port`（监听端口）。

启动

1) 以 Web 控制方式启动（`web.py` 会启动 websocketd 以及在本机运行 `vpn_finder`）：

```bash
cd 网关识别拦截程序--VPN_Finder
source .venv/bin/activate
python3 web.py
```

访问 `http://<this_host>:<flask_port>/` 查看日志页面（`config.json` 中 `flask_port` 指定）。

2) 直接运行 `vpn_finder`（不通过 Web）用于调试：

```bash
# 交互式运行（将函数直接调用为脚本）
python3 -c "import vpn_finder; vpn_finder.run('0.0.0.0', 31115)"
```

如何和 `gotranal` 联动
- 在 gotranal 的在线命令中使用 `-s <VPN_Finder_IP>:<vpn_finder_port>`，gotranal 会把特征通过 UDP 发送到 `VPN_Finder` 所指定的端口。

调试与建议
- SSH/iptables 权限：`vpn_finder` 使用 `ssh root@<router_ip> iptables ...` 下发规则，确保运行机器能通过密钥或密码免交互方式 ssh 到路由器，且路由器允许该操作。
- 幂等性：当前实现会不断添加 iptables 规则，建议在生产前修改为检查是否已存在规则再添加，或使用标签化的链并在 stop 流程中清理。
- 模型热加载：`udp.pkl` 使用 `joblib.load` 在启动时加载。如果要在运行时更新模型，建议添加接口或信号触发重新加载。
- 日志安全性：`log` 文件会被 websocketd 轮询读取，注意不要写入敏感信息。

常见问题
- 接收不到 gotranal 发来的数据：确认 `gotranal` 的 `-s` 参数地址和 `config.json` 中 `vpn_finder_port` 一致，且防火墙放通 UDP。
- ssh 命令无效：尝试手动在运行主机上执行 `ssh root@<router_ip> iptables -L` 确认可连通并有权限。

安全提示
- 自动下发 iptables 封禁规则会影响网络可达性，请在受控环境先做验证，并保留回滚策略（例如记录已下发规则并提供清除脚本）。
