# 在 Linux 虚拟机（VMware）上运行 AirGuard 实验的完整指南

本指南假设你在一台 Linux 虚拟机（推荐 Ubuntu 22.04 / Debian 11）中执行全部实验步骤（不使用 WSL）。Python 环境使用 Conda 管理。文档覆盖：系统准备、Docker、Go、Conda 安装、构建与运行 `gotranal`、启动 `autospider`（docker-compose）、启动 `VPN_Finder`（conda 环境）、端到端测试与排障建议。

重要前提
- 本指南以 Ubuntu/Debian 为主；若使用 RedHat/CentOS/AlmaLinux，请用对应包管理器（yum/dnf）替代 apt 的安装命令。
- 建议 VM 分配足够资源（至少 2 CPU、4GB 内存，测试环境推荐更多）。
- 需要可 SSH 到路由器或另一台可执行 iptables 的 Linux 主机，用于 `VPN_Finder` 下发封禁命令。

—— 目录结构说明（你仓库中已有）
- `数据采集节点demo--autospider/`：docker-compose 演示（spider、v2ray、mysql、cron）。
- `特征提取工具--gotranal/`：Go 程序，负责抓包与流特征提取。
- `网关识别拦截程序--VPN_Finder/`：Python 服务，接收特征并判定/封禁。

一、系统准备（Ubuntu/Debian）

以普通用户（后续需要 sudo 权限）打开终端：

```bash
# 更新系统
sudo apt update && sudo apt upgrade -y

# 基本构建工具与网络包
sudo apt install -y build-essential git wget curl ca-certificates unzip pkg-config \
    libpcap-dev libpcap0.8-dev openssh-client net-tools iproute2 socat
```

二、安装 Docker（用于 autospider）

按官方方式安装 Docker CE：

```bash
# 添加 Docker 官方 GPG
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo \"deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable\" \
  | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 允许当前用户使用 docker（重新登录或 newgrp docker 生效）
sudo usermod -aG docker $USER
newgrp docker || true

# 验证
docker version
```

说明与建议：如果你的 VM 与宿主宿主的网络或防火墙有特殊限制，请先确保 Docker pull、网络访问正常。

三、安装 Go（用于构建 gotranal）

```bash
# 下载并安装 Go（以 1.20.x 为例，按需替换版本）
wget https://go.dev/dl/go1.20.10.linux-amd64.tar.gz -O /tmp/go.tar.gz
sudo tar -C /usr/local -xzf /tmp/go.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
source ~/.profile
go version
```

四、安装 Miniconda（用于 Python/Conda 管理）

```bash
# 下载并安装 Miniconda（无交互安装）
wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh -O ~/miniconda.sh
bash ~/miniconda.sh -b -p $HOME/miniconda
eval "$($HOME/miniconda/bin/conda shell.bash hook)"
conda init
source ~/.bashrc
conda --version
```

五、把代码放到虚拟机（克隆或复制）

建议把仓库放在 VM 的本地文件系统（性能比挂载目录好）：

```bash
cd ~
git clone <your-repo-url> AirGuard
cd AirGuard
```

如果代码已经在宿主机上，通过 scp/共享文件夹或直接把压缩包拷贝进 VM。

六、构建 gotranal

```bash
cd ~/AirGuard/特征提取工具--gotranal
go mod download
go build -o gotranal main.go
# 确认二进制存在
ls -lh ./gotranal
```

七、准备并启动 autospider（docker-compose）

```bash
cd ~/AirGuard/数据采集节点demo--autospider
docker-compose up -d --build
docker-compose ps
```

注意：`docker-compose.yml` 为服务分配了静态私有 IP（172.19.0.x）。如果你的 Docker 网络或实验要求不同，可编辑 `docker-compose.yml` 删除 `ipv4_address:` 字段以使用 Docker 的动态分配。

八、为 VPN_Finder 创建 Conda 环境并安装依赖

```bash
cd ~/AirGuard/网关识别拦截程序--VPN_Finder
conda create -n vpnfinder python=3.10 -y
conda activate vpnfinder
pip install flask pandas joblib scikit-learn
```

（可选）把依赖写入文件：

```bash
pip freeze | grep -E "flask|pandas|joblib|scikit-learn" > requirements_vpnfinder.txt
```

九、生成 SSH key 并配置路由器（用于下发 iptables）

如果你要在路由器上实际添加 iptables 规则，需要免交互 SSH（推荐密钥方式）：

```bash
ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""
# 将公钥复制到路由器（需在路由器上允许 root 登录并支持 ssh）
ssh-copy-id root@<router_ip>
# 测试
ssh root@<router_ip> iptables -L -n --line-numbers
```

如果你没有可以下发 iptables 的路由器，可以准备一台同网段的 Linux 机器做为“受控主机”进行测试（行为与路由器类似）。

十、配置 `VPN_Finder`（修改 `config.json`）

在 `VPN_Finder/config.json` 中设置以下关键项：
- `ip`：`VPN_Finder` 绑定 IP（通常 `0.0.0.0` 或当前 VM 的局域网 IP，若 gotranal 与 VPN_Finder 在同一台 VM 则可用 `127.0.0.1`）。
- `flask_port`：Flask Web 界面端口（例如 `5000`）。
- `websocketd_port`：websocketd 的端口（例如 `8080`）。
- `router_ip`：要下发 iptables 的路由器/控制主机 IP（用于 ssh）。
- `iface`：gotranal 抓包时的网络接口名（如 `eth0`、`ens33`、`any`）。
- `vpn_finder_port`：gotranal 发送特征的 UDP 端口（例如 `31115`）。

示例：

```json
{
  "ip": "127.0.0.1",
  "flask_port": "5000",
  "websocketd_port": "8080",
  "router_ip": "192.168.1.1",
  "iface": "eth0",
  "vpn_finder_port": "31115"
}
```

十一、启动 `VPN_Finder` 服务

```bash
cd ~/AirGuard/网关识别拦截程序--VPN_Finder
conda activate vpnfinder
python3 web.py
```

`web.py` 会启动 `websocketd`（如果本目录下有 `websocketd` 二进制）并在 `flask_port` 上运行 Web 界面；也会通过 `/start` 接口尝试在 `router_ip` 或远端启动 gotranal（如果配置为 ssh 启动），但通常我们在本机/VM 上直接运行 `gotranal`。

十二、运行 gotranal（在线模式），将特征发送到 `VPN_Finder`

注意：抓包需要 root 权限或相应 CAP 权限。直接在 VM 上以 sudo 运行：

```bash
cd ~/AirGuard/特征提取工具--gotranal
sudo ./gotranal on -I <iface> -b "tcp or udp" -s 127.0.0.1:31115 -d 60
```

说明：
- 将 `<iface>` 替换为实际接口（例如 `eth0`、`ens33`），可用 `ip link` / `ip addr` 查看。
- `-s 127.0.0.1:31115` 表示把特征通过 UDP 发到本机的 `VPN_Finder`（端口需与 `config.json` 中一致）。

可选：在容器中运行 gotranal（如果你想把它容器化）：

```bash
docker run --rm --cap-add=NET_RAW --cap-add=NET_ADMIN --net=host \
  -v $(pwd):/work -w /work myrepo/gotranal:latest \
  ./gotranal on -I <iface> -b "tcp or udp" -s 127.0.0.1:31115 -d 60
```

十三、端到端快速测试

1) 本地模拟发送一条特征到 `VPN_Finder`（用于验证监听与队列）

```bash
echo "dummy_flow_line_example" | nc -u -w1 127.0.0.1 31115
# 或发送多行：
printf "flow1\nflow2\n" | nc -u -w1 127.0.0.1 31115
```

2) 检查 `VPN_Finder` 的日志文件 `./log` 或 Web 界面，看是否记录收到数据与分析结果。

3) 运行真实端到端：先启动 `autospider`（docker-compose），再用 `gotranal on` 在线抓取并发送；观察 `VPN_Finder` 是否有分析/封禁记录。

十四、常见问题与排查建议

- gotranal 无法打开接口或报错权限不足：确保以 root 运行（sudo）；或在容器中添加 `CAP_NET_RAW`。
- 找不到接口：用 `ip link` / `ip addr` 列出接口，别用 `any` 作为所有环境的默认值，某些系统需要具体接口名。
- `VPN_Finder` 无法接收 UDP：确认 `gotranal -s` 目标 IP/端口与 `config.json` 中 `vpn_finder_port` 一致，并用 `ss -uln` / `netstat -uln` 检查端口监听。
- ssh 相关问题：先在命令行手工 `ssh root@router_ip` 测试连通和权限，确保 `ssh-copy-id` 已正确部署公钥。
- Docker 网络冲突：若 `docker-compose.yml` 中静态 IP 与 VM 网络冲突，建议删除 `ipv4_address` 字段或调整子网配置。

十五、回滚与清理（重要）

如果 `VPN_Finder` 在路由器上下发了 iptables 规则导致误封，登录路由器并查看规则：

```bash
ssh root@<router_ip>
iptables -L FORWARD -n --line-numbers
# 删除第 N 条规则（根据实际 line-numbers）
iptables -D FORWARD <N>
```

建议：在 `VPN_Finder` 中实现“已下发规则清单”并提供清理脚本，便于回滚。

十六、可选增强（建议）

- 把 `VPN_Finder` 的封禁命令替换为先记录、后由运维人工确认再下发的流程，降低误封风险。
- 为 `gotranal` 添加日志与输出样例，或把输出同时保存为 CSV 以便离线审计。
- 在 VM 上把关键服务设置为 systemd 服务（如 `vpn_finder.service`、`gotranal.service`）以便开机自启与管理。

附：常用命令速查

```bash
# 列出网络接口
ip addr

# 查看 UDP 端口监听
ss -uln

# 查看 docker 容器
docker ps -a

# 查看 docker 网络
docker network ls

# 构建并运行 gotranal
cd ~/AirGuard/特征提取工具--gotranal
go build -o gotranal main.go
sudo ./gotranal on -I eth0 -b "tcp or udp" -s 127.0.0.1:31115 -d 60

# 启动 autospider
cd ~/AirGuard/数据采集节点demo--autospider
docker-compose up -d --build

# 启动 VPN_Finder
cd ~/AirGuard/网关识别拦截程序--VPN_Finder
conda activate vpnfinder
python3 web.py
```

如果你希望，我可以继续：
- 把本文件内容自动写入仓库的 `LINUX_VM_GUIDE.md`（我已经执行），并生成一个 `vpnfinder/requirements.txt` 或 `scripts/bootstrap_vm.sh` 的一键安装脚本（可直接在 VM 上运行）。

祝实验顺利，遇到错误把日志和具体报错贴过来我会继续帮你排查。
