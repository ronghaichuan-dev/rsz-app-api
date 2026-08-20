#!/usr/bin/env bash

set -Eeuo pipefail

# DEPLOY_ROOT 定义所有微服务发布目录的公共根路径。
readonly DEPLOY_ROOT="${RSZ_DEPLOY_ROOT:-/opt/rsz}"
# SYSTEMCTL_BIN 指定 Alibaba Cloud Linux 上 systemctl 的绝对路径，便于与 sudoers 精确匹配。
readonly SYSTEMCTL_BIN="${RSZ_SYSTEMCTL_BIN:-/usr/bin/systemctl}"

# check_arguments 校验服务、环境、版本、发布包和健康检查参数。
check_arguments() {
  if [[ $# -ne 5 ]]; then
    echo "用法: $0 <服务名> <环境> <版本号> <发布包路径> <健康检查地址>" >&2
    exit 2
  fi
  if [[ ! $1 =~ ^[a-z][a-z0-9-]*$ ]]; then
    echo "服务名只能包含小写字母、数字和连字符。" >&2
    exit 2
  fi
  if [[ ! $2 =~ ^(dev|test|prod)$ ]]; then
    echo "环境只能是 dev、test 或 prod。" >&2
    exit 2
  fi
  if [[ ! $3 =~ ^[0-9a-f]{40}$ ]]; then
    echo "版本号必须是 40 位 Git Commit SHA。" >&2
    exit 2
  fi
  if [[ ! -f $4 ]]; then
    echo "未找到发布包: $4" >&2
    exit 2
  fi
  if [[ ! $5 =~ ^https?:// ]]; then
    echo "健康检查地址必须以 http:// 或 https:// 开头。" >&2
    exit 2
  fi
}

# verify_health 等待服务启动，并通过健康检查确认新版本可以提供 HTTP 服务。
verify_health() {
  local service_name="$1"
  local health_url="$2"
  local attempt
  for attempt in {1..15}; do
    if curl --fail --silent --show-error --max-time 3 "$health_url" >/dev/null; then
      echo "健康检查通过: $health_url"
      return 0
    fi
    sleep 2
  done
  echo "健康检查失败: $health_url" >&2
  sudo "$SYSTEMCTL_BIN" --no-pager --full status "$service_name" || true
  return 1
}

# cleanup_old_releases 仅保留最近五个发布目录，便于故障时快速切回先前版本。
cleanup_old_releases() {
  local releases_dir="$1"
  local release
  while IFS= read -r release; do
    [[ -z $release ]] && continue
    rm -rf -- "$release"
  done < <(
    find "$releases_dir" -mindepth 1 -maxdepth 1 -type d -printf '%T@ %p\n' \
      | sort -rn \
      | awk 'NR > 5 { $1=""; sub(/^ /, ""); print }'
  )
}

# main 解压发布包、切换当前版本、重启服务并验证健康状态。
main() {
  local app_name="$1"
  local app_env="$2"
  local version="$3"
  local archive="$4"
  local health_url="$5"
  local service_name="rsz-${app_name}-${app_env}"
  local base_dir="${DEPLOY_ROOT}/${app_name}"
  local releases_dir="${base_dir}/releases"
  local release_dir="${releases_dir}/${version}"

  if [[ ! -x $SYSTEMCTL_BIN ]]; then
    echo "未找到 systemctl: $SYSTEMCTL_BIN" >&2
    exit 1
  fi

  mkdir -p "$release_dir"
  tar -xzf "$archive" -C "$release_dir"
  rm -f -- "$archive"

  if [[ ! -x "${release_dir}/${app_name}" || ! -d "${release_dir}/manifest/i18n" ]]; then
    echo "发布包缺少 ${app_name} 二进制或多语言文件。" >&2
    exit 1
  fi

  ln -sfn "$release_dir" "${releases_dir}/current"
  sudo "$SYSTEMCTL_BIN" restart "$service_name"
  verify_health "$service_name" "$health_url"
  cleanup_old_releases "$releases_dir"
  echo "${app_name} 已发布: ${version}"
}

check_arguments "$@"
main "$@"
