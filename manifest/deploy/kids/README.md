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

`rsz-deploy` 不应读取 `/etc/rsz/kids` 下的私有配置；发布脚本只切换版本并重启服务，配置文件由 `rsz-kids` 服务账号在启动时读取。

为部署用户添加最小 sudo 权限，执行 `sudo visudo -f /etc/sudoers.d/rsz-deploy` 并写入：

```sudoers
rsz-deploy ALL=(root) NOPASSWD: /usr/bin/systemctl restart rsz-kids-test, /usr/bin/systemctl status rsz-kids-test
```

Alibaba Cloud Linux 4 LTS 通常使用 `/usr/bin/systemctl`。如 `command -v systemctl` 返回不同路径，请同步调整 sudoers 和通用发布脚本的默认路径。

## GitHub Actions 配置

`deploy-kids-test.yml` 在推送 `test` 分支时调用可复用工作流 `_deploy-go-service.yml`；`deploy-kids-prod.yml` 在推送 `main` 分支时以 `prod` 环境部署。两个工作流分别使用 GitHub Actions 的 `test` 和 `production` Environment，避免测试和生产共用部署凭据。

在对应的 GitHub Environment 中配置以下 Actions Secrets：

| 名称 | 内容 |
| --- | --- |
| `RSZ_DEPLOY_HOST` | 当前环境服务器的 IP 或域名 |
| `RSZ_DEPLOY_USER` | `rsz-deploy` |
| `RSZ_DEPLOY_SSH_KEY` | 对应的 SSH 私钥 |
| `RSZ_DEPLOY_KNOWN_HOSTS` | `ssh-keyscan -H <服务器地址>` 得到的主机指纹 |

服务器为 ARM64 时，在对应 GitHub Environment 的 Variables 添加 `RSZ_DEPLOY_GOARCH=arm64`；x86_64 服务器不需要设置，默认使用 `amd64`。

`RSZ_DEPLOY_SSH_KEY` 必须是完整的私钥文本，首行为 `-----BEGIN OPENSSH PRIVATE KEY-----`，末行为 `-----END OPENSSH PRIVATE KEY-----`，不能填写 `.pub` 公钥、文件路径或 SSH 配置内容。工作流会在连接服务器前使用 `ssh-keygen` 校验私钥；校验失败时请重新复制本地私钥全文。

工作流固定执行以下流程：测试、构建 Linux 二进制、上传发布包、切换 `current`、重启 systemd 服务和健康检查。发布成功后只保留最新五个版本。

## 初始化 kids 生产服务器

生产服务器沿用测试环境的目录、用户与权限模型，但使用独立数据库、配置和 systemd 单元。部署前执行以下初始化，并将真实数据库、JWT 与 OAuth 参数写入服务器私有配置：

```bash
sudo install -o root -g rsz-kids -m 0640 \
  manifest/deploy/kids/config.prod.yaml.example \
  /etc/rsz/kids/config.prod.yaml

sudo install -o root -g root -m 0644 \
  manifest/deploy/kids/rsz-kids-prod.service \
  /etc/systemd/system/rsz-kids-prod.service
sudo systemctl daemon-reload
sudo systemctl enable rsz-kids-prod
```

为 `rsz-deploy` 增加仅针对生产服务的最小 sudo 权限：

```sudoers
rsz-deploy ALL=(root) NOPASSWD: /usr/bin/systemctl restart rsz-kids-prod, /usr/bin/systemctl status rsz-kids-prod
```

服务监听 `127.0.0.1:8002`；请为生产域名单独配置 Nginx 与 HTTPS。GitHub 中 `production` Environment 可按需配置审批规则，审批通过后每次推送到 `main` 都会自动发布。

## Nginx、HTTPS 与阿里云安全组

kids 固定监听 `127.0.0.1:18002`，不直接暴露公网端口。Nginx 是唯一的公网入口，负责 HTTP/HTTPS 终止和反向代理。

Nginx 配置额外保留 `/_nginx/` 路径用于验证默认静态首页，不经过 kids 服务：HTTP 初始化阶段访问 `http://apptest.rszai.com/_nginx/`；启用 HTTPS 后访问 `https://apptest.rszai.com/_nginx/`。其余根路径仍转发至 kids API。

服务器 IP 的默认站点使用 `manifest/deploy/nginx/default.conf`。当前 Alibaba Cloud Linux 默认 `nginx.conf` 已在默认 server 中包含 `/etc/nginx/default.d/*.conf`，将该文件复制到 `/etc/nginx/default.d/default.conf` 后，可通过 `http://<服务器公网IP>/` 或未匹配域名访问 `/usr/share/nginx/html` 下的静态页面。

先在 DNS 中将测试域名 `apptest.rszai.com` 的 A 记录指向服务器公网 IP。将 Nginx 模板复制到服务器后安装：

```bash
sudo dnf install -y nginx

sudo install -o root -g root -m 0644 \
  /tmp/rsz-kids-test.nginx.conf \
  /etc/nginx/conf.d/rsz-kids-test.conf

sudo nginx -t
sudo systemctl enable --now nginx
```

模板文件为 `manifest/deploy/kids/rsz-kids-test.nginx.conf`，与配置模板、systemd 单元一样，应在初始化服务器前传到 `/tmp`：

```bash
scp manifest/deploy/kids/rsz-kids-test.nginx.conf \
  <管理用户>@<服务器IP>:/tmp/rsz-kids-test.nginx.conf

scp manifest/deploy/kids/rsz-kids-test.https.nginx.conf \
  <管理用户>@<服务器IP>:/tmp/rsz-kids-test.https.nginx.conf
```

Alibaba Cloud Linux 4 的默认仓库不提供 Certbot，因此本项目使用 `acme.sh` 通过 HTTP-01 签发 Let's Encrypt 证书。签发前，请先确认域名已解析、Nginx 已加载 HTTP 初始化配置且安全组已开放 `80`。

```bash
sudo install -d -o nginx -g nginx -m 0755 /var/www/acme
curl https://get.acme.sh | sudo sh -s email=<运维邮箱>
sudo /root/.acme.sh/acme.sh --set-default-ca --server letsencrypt
sudo /root/.acme.sh/acme.sh --issue \
  --server letsencrypt \
  -d apptest.rszai.com \
  -w /var/www/acme

sudo install -d -o root -g root -m 0700 /etc/nginx/ssl
sudo install -o root -g root -m 0600 /dev/null \
  /etc/nginx/ssl/apptest.rszai.com.key.pem
sudo install -o root -g root -m 0644 /dev/null \
  /etc/nginx/ssl/apptest.rszai.com.fullchain.pem

sudo /root/.acme.sh/acme.sh --install-cert -d apptest.rszai.com \
  --key-file /etc/nginx/ssl/apptest.rszai.com.key.pem \
  --fullchain-file /etc/nginx/ssl/apptest.rszai.com.fullchain.pem \
  --reloadcmd "/usr/bin/systemctl reload nginx"
```

签发完成后，将 HTTPS 模板替换为当前 Nginx 配置并重载：

```bash
sudo install -o root -g root -m 0644 \
  /tmp/rsz-kids-test.https.nginx.conf \
  /etc/nginx/conf.d/rsz-kids-test.conf
sudo nginx -t
sudo systemctl reload nginx
```

HTTPS 模板为 `manifest/deploy/kids/rsz-kids-test.https.nginx.conf`，初始化服务器时同样传到 `/tmp`。`acme.sh --install-cert` 会在续期后自动把新证书写入 Nginx 使用的路径并重载服务。

阿里云安全组仅需放行：

| 协议/端口 | 用途 | 建议来源 |
| --- | --- | --- |
| TCP 22 | SSH 与 GitHub Actions 部署 | 尽可能限制管理 IP；GitHub Actions 需要访问时使用密钥认证 |
| TCP 80 | Let's Encrypt 校验与 HTTP 自动跳转 | `0.0.0.0/0`、`::/0` |
| TCP 443 | 正式 HTTPS API | `0.0.0.0/0`、`::/0` |

不要在安全组或 firewalld 中放行 `18002`。若 firewalld 已启用，允许 Nginx 端口：

```bash
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --permanent --add-service=https
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
