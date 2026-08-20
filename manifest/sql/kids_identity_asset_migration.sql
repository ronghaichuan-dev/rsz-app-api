-- Clearwave Backend Protocol v1 持久化基础表。
-- 本迁移只属于 kids 微服务自己的数据库，不与其他应用共用。

CREATE TABLE IF NOT EXISTS kids_identity_session (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    session_id VARCHAR(128) NOT NULL COMMENT '接口会话标识',
    account_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '接口账号标识',
    principal_kind VARCHAR(32) NOT NULL COMMENT '主体类型',
    access_token_hash CHAR(64) NOT NULL COMMENT '访问令牌摘要',
    refresh_token_hash CHAR(64) NOT NULL DEFAULT '' COMMENT '刷新令牌摘要',
    guest_upgrade_grant_hash CHAR(64) NOT NULL DEFAULT '' COMMENT '游客升级凭据摘要',
    expires_at DATETIME NOT NULL COMMENT '访问令牌到期时间',
    refresh_expires_at DATETIME NULL COMMENT '刷新令牌到期时间',
    revoked_at DATETIME NULL COMMENT '撤销时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_session_id (session_id),
    UNIQUE KEY uk_access_token_hash (access_token_hash),
    UNIQUE KEY uk_refresh_token_hash (refresh_token_hash),
    KEY idx_account_id (account_id),
    KEY idx_principal_expires (principal_kind, expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口会话';

CREATE TABLE IF NOT EXISTS kids_request_deduplication (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    principal_scope VARCHAR(256) NOT NULL COMMENT '主体和成员作用域',
    idempotency_key VARCHAR(256) NOT NULL COMMENT '幂等键',
    operation_id VARCHAR(128) NOT NULL COMMENT '接口操作标识',
    route_fingerprint CHAR(64) NOT NULL COMMENT '路由摘要',
    body_fingerprint CHAR(64) NOT NULL COMMENT '请求体摘要',
    response_status SMALLINT UNSIGNED NOT NULL COMMENT '首次响应状态码',
    response_body JSON NOT NULL COMMENT '首次规范响应',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_scope_key (principal_scope, idempotency_key),
    KEY idx_operation_id (operation_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口幂等记录';

CREATE TABLE IF NOT EXISTS kids_sync_commit (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    commit_id VARCHAR(128) NOT NULL COMMENT '接口提交标识',
    circle_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '圈子标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '单调提交序列',
    change_payload JSON NOT NULL COMMENT '完整变更集合',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_commit_id (commit_id),
    UNIQUE KEY uk_commit_sequence (commit_sequence),
    KEY idx_circle_sequence (circle_id, commit_sequence)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口同步提交';

CREATE TABLE IF NOT EXISTS kids_mutation_receipt (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    receipt_id VARCHAR(128) NOT NULL COMMENT '接口回执标识',
    commit_id VARCHAR(128) NOT NULL COMMENT '提交标识',
    operation_id VARCHAR(128) NOT NULL COMMENT '接口操作标识',
    result_kind VARCHAR(32) NOT NULL COMMENT '结果类型',
    committed_at DATETIME NOT NULL COMMENT '提交时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_receipt_id (receipt_id),
    KEY idx_commit_id (commit_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口写入回执';

CREATE TABLE IF NOT EXISTS kids_asset_upload (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    upload_id VARCHAR(128) NOT NULL COMMENT '上传标识',
    account_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '账号标识',
    circle_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '圈子标识',
    purpose VARCHAR(64) NOT NULL COMMENT '上传用途',
    content_type VARCHAR(128) NOT NULL COMMENT '内容类型',
    byte_size BIGINT UNSIGNED NOT NULL COMMENT '字节大小',
    sha256 CHAR(64) NOT NULL COMMENT '内容摘要',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    status VARCHAR(32) NOT NULL COMMENT '上传状态',
    expires_at DATETIME NOT NULL COMMENT '到期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_upload_id (upload_id),
    KEY idx_account_purpose (account_id, purpose)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口资产上传';

CREATE TABLE IF NOT EXISTS kids_asset (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    asset_id VARCHAR(128) NOT NULL COMMENT '资产标识',
    upload_id VARCHAR(128) NOT NULL COMMENT '上传标识',
    circle_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '圈子标识',
    purpose VARCHAR(64) NOT NULL COMMENT '资产用途',
    content_type VARCHAR(128) NOT NULL COMMENT '内容类型',
    byte_size BIGINT UNSIGNED NOT NULL COMMENT '字节大小',
    sha256 CHAR(64) NOT NULL COMMENT '内容摘要',
    state VARCHAR(32) NOT NULL COMMENT '资产状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    committed_at DATETIME NOT NULL COMMENT '提交时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_asset_id (asset_id),
    UNIQUE KEY uk_upload_id (upload_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口已提交资产';

CREATE TABLE IF NOT EXISTS kids_feedback (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    feedback_id VARCHAR(128) NOT NULL COMMENT '反馈标识',
    account_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '账号标识',
    category VARCHAR(64) NOT NULL COMMENT '反馈分类',
    content TEXT NOT NULL COMMENT '反馈内容',
    contact_type VARCHAR(64) NOT NULL DEFAULT '' COMMENT '联系类型',
    contact VARCHAR(512) NOT NULL DEFAULT '' COMMENT '联系方式',
    privacy_consent_version VARCHAR(128) NOT NULL COMMENT '隐私同意版本',
    attachment_asset_ids JSON NOT NULL COMMENT '附件资产标识',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_feedback_id (feedback_id),
    KEY idx_account_created (account_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口反馈';

CREATE TABLE IF NOT EXISTS kids_entitlement (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    entitlement_id VARCHAR(128) NOT NULL COMMENT '权益标识',
    account_id VARCHAR(128) NOT NULL COMMENT '账号标识',
    product_id VARCHAR(256) NOT NULL DEFAULT '' COMMENT '商品标识',
    status VARCHAR(32) NOT NULL COMMENT '权益状态',
    valid_until_at DATETIME NULL COMMENT '有效截止时间',
    verified_at DATETIME NULL COMMENT '验证时间',
    revoked_at DATETIME NULL COMMENT '撤销时间',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_entitlement_id (entitlement_id),
    UNIQUE KEY uk_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口权益快照';
