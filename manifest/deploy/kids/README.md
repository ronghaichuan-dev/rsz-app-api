# rsz 微服务部署约定

本项目使用一套通用发布流程部署所有 Go 微服务。可替换发布包位于 `/opt/rsz/<服务名>`，服务器私有配置位于 `/etc/rsz/<服务名>`，二者分离，避免密码和密钥进入 Git 或发布包。

```text
/opt/rsz/kids/
├── releases/
│   ├── <commit-sha>/
│   └── current -> <commit-sha>
└── data/                         # 可选：上传文件等持久化数据

/etc/rsz/kids/
├── config.test.yaml
└── config.prod.yaml
```

`APP_CONFIG_DIR` 会让服务优先从 `/etc/rsz/<服务名>/config.<环境>.yaml` 读取配置；本地开发未设置该环境变量时，仍使用仓库内的 `config/<服务名>/config.<环境>.yaml`。

## 初始化 kids 测试服务器

以下命令适用于 Alibaba Cloud Linux 4 LTS。默认发布目标为 `linux/amd64`，可执行 `uname -m` 确认输出为 `x86_64`。

```bash
sudo dnf install -y curl tar sudo

sudo useradd --system --create-home --shell /usr/sbin/nologin rsz-kids
sudo useradd --create-home --shell /bin/bash rsz-deploy

sudo install -d -o rsz-deploy -g rsz-deploy -m 0755 /opt/rsz/kids/releases
sudo install -d -o root -g rsz-kids -m 0750 /etc/rsz/kids
sudo install -o root -g rsz-kids -m 0640 \
  manifest/deploy/kids/config.test.yaml.example \
  /etc/rsz/kids/config.test.yaml

sudo install -o root -g root -m 0644 \
  manifest/deploy/kids/rsz-kids-test.service \
  /etc/systemd/system/rsz-kids-test.service
sudo systemctl daemon-reload
sudo systemctl enable rsz-kids-test
```

编辑 `/etc/rsz/kids/config.test.yaml`，填入测试数据库连接、独立 JWT 密钥与 OAuth 参数。该文件不能提交到 Git，也不能由 GitHub Actions 上传覆盖。

为部署用户添加最小 sudo 权限，执行 `sudo visudo -f /etc/sudoers.d/rsz-deploy` 并写入：

```sudoers
rsz-deploy ALL=(root) NOPASSWD: /usr/bin/systemctl restart rsz-kids-test, /usr/bin/systemctl status rsz-kids-test
```

Alibaba Cloud Linux 4 LTS 通常使用 `/usr/bin/systemctl`。如 `command -v systemctl` 返回不同路径，请同步调整 sudoers 和通用发布脚本的默认路径。

## GitHub Actions 配置

`deploy-kids-test.yml` 在推送 `test` 分支时调用可复用工作流 `_deploy-go-service.yml`。为仓库配置以下 Actions Secrets：

| 名称 | 内容 |
| --- | --- |
| `RSZ_DEPLOY_HOST` | 测试服务器的 IP 或域名 |
| `RSZ_DEPLOY_USER` | `rsz-deploy` |
| `RSZ_DEPLOY_SSH_KEY` | 对应的 SSH 私钥 |
| `RSZ_DEPLOY_KNOWN_HOSTS` | `ssh-keyscan -H <服务器地址>` 得到的主机指纹 |

服务器为 ARM64 时，在仓库 Actions Variables 添加 `RSZ_DEPLOY_GOARCH=arm64`；x86_64 服务器不需要设置，默认使用 `amd64`。

工作流固定执行以下流程：测试、构建 Linux 二进制、上传发布包、切换 `current`、重启 systemd 服务和健康检查。发布成功后只保留最新五个版本。

## 阿里云网络放行

建议由 Nginx 对外提供 `80`/`443`，kids 保持监听本机 `18002`，无需在阿里云安全组公开 `18002`。若临时直接访问 kids，需要在安全组中增加 TCP `18002` 入站规则；若实例启用了 firewalld，还要执行：

```bash
sudo firewall-cmd --permanent --add-port=18002/tcp
sudo firewall-cmd --reload
```

## 新增微服务

以 `admin` 测试环境为例，新增时只需：

1. 创建 `/etc/rsz/admin/config.test.yaml` 与 `rsz-admin-test.service`，服务单元设置 `APP_CONFIG_DIR=/etc/rsz/admin`。
2. 创建 `/opt/rsz/admin/releases`，所有者设置为 `rsz-deploy`；为进程创建独立 `rsz-admin` 系统用户。
3. 在 sudoers 中仅允许 `rsz-deploy` 重启和查看 `rsz-admin-test`。
4. 新建薄工作流并复用基础工作流：

```yaml
name: 部署 admin 到测试服务器

on:
  push:
    branches: [test]

jobs:
  deploy:
    uses: ./.github/workflows/_deploy-go-service.yml
    with:
      service_name: admin
      app_env: test
      health_url: http://127.0.0.1:18001/v1/health
    secrets: inherit
```

通用脚本会基于服务名自动选择 `/opt/rsz/<服务名>`、`/etc/rsz/<服务名>` 和 `rsz-<服务名>-<环境>`，不需要复制业务相关的发布脚本。

## 回滚

找出旧版本目录后执行：

```bash
sudo -u rsz-deploy ln -sfn /opt/rsz/kids/releases/<旧提交SHA> /opt/rsz/kids/releases/current
sudo systemctl restart rsz-kids-test
curl --fail http://127.0.0.1:18002/v1/health
```

数据库 migration 必须先在测试库单独执行和验证。当前 SQL migration 未引入版本记录机制，自动部署流程不会执行 SQL 文件。
