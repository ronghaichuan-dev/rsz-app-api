-- 将 kids 接口会话迁移为 Unix epoch milliseconds 的唯一事实来源。
-- 请先执行既有会话毫秒精度迁移，再停止旧版服务并执行本文件。

ALTER TABLE kids_identity_session
    ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '会话状态' AFTER principal_kind,
    ADD COLUMN issued_at_ms BIGINT UNSIGNED NULL COMMENT '会话签发毫秒时间戳' AFTER guest_upgrade_grant_hash,
    ADD COLUMN access_expires_at_ms BIGINT UNSIGNED NULL COMMENT '访问令牌到期毫秒时间戳' AFTER issued_at_ms,
    ADD COLUMN refresh_expires_at_ms BIGINT UNSIGNED NULL COMMENT '刷新令牌到期毫秒时间戳' AFTER access_expires_at_ms,
    ADD COLUMN revoked_at_ms BIGINT UNSIGNED NULL COMMENT '撤销毫秒时间戳' AFTER refresh_expires_at_ms;

UPDATE kids_identity_session
SET
    issued_at_ms = CAST(UNIX_TIMESTAMP(created_at) * 1000 AS UNSIGNED),
    access_expires_at_ms = CAST(UNIX_TIMESTAMP(expires_at) * 1000 AS UNSIGNED),
    refresh_expires_at_ms = CASE WHEN refresh_expires_at IS NULL THEN NULL ELSE CAST(UNIX_TIMESTAMP(refresh_expires_at) * 1000 AS UNSIGNED) END,
    revoked_at_ms = CASE WHEN revoked_at IS NULL THEN NULL ELSE CAST(UNIX_TIMESTAMP(revoked_at) * 1000 AS UNSIGNED) END,
    status = CASE WHEN revoked_at IS NULL THEN 'active' ELSE 'revoked' END;

ALTER TABLE kids_identity_session
    MODIFY COLUMN issued_at_ms BIGINT UNSIGNED NOT NULL COMMENT '会话签发毫秒时间戳',
    MODIFY COLUMN access_expires_at_ms BIGINT UNSIGNED NOT NULL COMMENT '访问令牌到期毫秒时间戳',
    DROP COLUMN expires_at,
    DROP COLUMN refresh_expires_at,
    DROP COLUMN revoked_at,
    DROP COLUMN created_at,
    DROP COLUMN updated_at,
    ADD KEY idx_status_expires (status, access_expires_at_ms);
