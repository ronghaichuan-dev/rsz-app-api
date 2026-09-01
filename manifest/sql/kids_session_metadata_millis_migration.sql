-- 将既有会话时间列升级为毫秒精度；历史整秒数据保持原值，后续会话完整保留毫秒。
-- 请在部署新服务版本前，对对应 kids 数据库执行本迁移。

ALTER TABLE kids_identity_session
    MODIFY COLUMN expires_at DATETIME(3) NOT NULL COMMENT '访问令牌到期时间',
    MODIFY COLUMN refresh_expires_at DATETIME(3) NULL COMMENT '刷新令牌到期时间',
    MODIFY COLUMN revoked_at DATETIME(3) NULL COMMENT '撤销时间',
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间';
