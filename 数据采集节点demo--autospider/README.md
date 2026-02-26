# 数据采集节点 Demo — autospider

本子项目用于构建一个容器化的爬虫 + 代理 + 数据持久化演示环境，用来生成测试流量并将抓包文件保存供后续特征提取使用。

目录结构要点：
- `docker-compose.yml`：一键启动的服务编排（spider、v2ray、datbase、cron）。
- `spider/`：爬虫镜像构建上下文，包含 `code/spider.py`（Selenium 爬虫）和 `chrome_soft/chromedriver_v100`。
- `v2ray/`：v2ray 代理容器及示例 pcap、工具脚本。
- `mysql/init/`：数据库初始化 SQL 文件（容器启动时会导入）。

先决条件
- Docker 与 docker-compose 已安装并能运行（Linux/macOS）。
- 若在本机调试 `spider.py`，需安装 Chrome 浏览器与匹配版本的 Chromedriver。容器运行时 chromedriver 已在镜像内。

启动（推荐在项目根目录下执行）：

```bash
cd 数据采集节点demo--autospider
docker-compose up -d --build
```

主要服务说明
- spider：运行基于 Selenium 的爬虫并通过 v2ray 代理访问目标站点，用于生成网络流量样本。镜像会挂载 `./spider/code` 到容器内 `/code`，可在运行时修改脚本并即时生效。
- v2ray：代理与抓包端点，可生成/转发流量，并保存 pcap 到 `./v2ray/pcap`。
- datbase：MySQL（用于存储爬虫相关数据）；初始化 SQL 在 `mysql/init/` 中。
- cron：调度服务（gocron），用于定时任务或触发爬虫。

在容器中手动运行/调试 spider
- 进入 spider 容器：

```bash
docker-compose exec spider /bin/bash
cd /code
# 以 headless 模式在容器里执行（示例：运行 youtube 流量）
python3 spider.py youtube
```

本地调试（非容器）注意事项
- 请把 `driver_path` 或 chromedriver 放到本机可访问路径，并在 `spider.py` 中修改 `Service("/chrome_soft/chromedriver_v100")` 为本机 chromedriver 路径。
- 运行前安装依赖：

```bash
python3 -m venv .venv
source .venv/bin/activate
pip install selenium pytz
```

常见问题
- 如果 Selenium 无法连接到 chromedriver，确认 chromedriver 版本与 Chrome 版本匹配，并且可执行权限正确。
- 代理连接失败时，检查 `v2ray` 服务是否在网络中可达（docker 网络中 IP 在 `docker-compose.yml` 中被固定）。

安全提示
- spider 脚本会访问外部站点，请在可控网络环境下运行，避免对目标站点造成不当访问。

更多
- 若需将生成的 pcap 交给特征提取模块（`gotranal`）处理，请查看 `特征提取工具--gotranal/README.md`。
