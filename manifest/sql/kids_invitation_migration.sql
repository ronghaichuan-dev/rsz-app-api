-- Clearwave 接口邀请码事实表。
-- 邀请码仅保存安全摘要，原文只在 create 和 refresh 响应中出现一次。

CREATE TABLE IF NOT EXISTS kids_invitation (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    invite_id VARCHAR(128) NOT NULL COMMENT '接口邀请标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    target_role VARCHAR(32) NOT NULL COMMENT '目标角色',
    target_administrator_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '目标管理员标识',
    target_member_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '目标成员标识',
    permission_scope JSON NOT NULL COMMENT '邀请权限范围',
    code_hash CHAR(64) NOT NULL COMMENT '邀请码摘要',
    single_use TINYINT(1) NOT NULL DEFAULT 1 COMMENT '是否单次使用',
    generation INT UNSIGNED NOT NULL DEFAULT 1 COMMENT '邀请码代次',
    status VARCHAR(32) NOT NULL COMMENT '邀请状态',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    used_at DATETIME NULL COMMENT '使用时间',
    revoked_at DATETIME NULL COMMENT '撤销时间',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_by_account_id VARCHAR(128) NOT NULL COMMENT '创建账号标识',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_invite_id (invite_id),
    UNIQUE KEY uk_code_hash (code_hash),
    KEY idx_circle_status (circle_id, status),
    KEY idx_target_admin (target_administrator_id, status),
    KEY idx_target_member (target_member_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口邀请码';
