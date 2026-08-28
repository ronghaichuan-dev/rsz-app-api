-- 将 kids 接口会话迁移为 Unix epoch milliseconds 的唯一事实来源。
-- 本迁移兼容 status 已存在但其他列尚未迁移的中断状态，可安全再次执行。
-- 请先执行既有会话毫秒精度迁移，再停止旧版服务并执行本文件。

-- 按列分别补齐，避免多列 ALTER 因单个已存在列整体失败。
SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'status'
);
SET @kids_sql := IF(@kids_column_exists = 0,
    'ALTER TABLE kids_identity_session ADD COLUMN status VARCHAR(32) NOT NULL DEFAULT ''active'' COMMENT ''会话状态'' AFTER principal_kind',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'issued_at_ms'
);
SET @kids_sql := IF(@kids_column_exists = 0,
    'ALTER TABLE kids_identity_session ADD COLUMN issued_at_ms BIGINT UNSIGNED NULL COMMENT ''会话签发毫秒时间戳'' AFTER guest_upgrade_grant_hash',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'access_expires_at_ms'
);
SET @kids_sql := IF(@kids_column_exists = 0,
    'ALTER TABLE kids_identity_session ADD COLUMN access_expires_at_ms BIGINT UNSIGNED NULL COMMENT ''访问令牌到期毫秒时间戳'' AFTER issued_at_ms',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'refresh_expires_at_ms'
);
SET @kids_sql := IF(@kids_column_exists = 0,
    'ALTER TABLE kids_identity_session ADD COLUMN refresh_expires_at_ms BIGINT UNSIGNED NULL COMMENT ''刷新令牌到期毫秒时间戳'' AFTER access_expires_at_ms',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'revoked_at_ms'
);
SET @kids_sql := IF(@kids_column_exists = 0,
    'ALTER TABLE kids_identity_session ADD COLUMN revoked_at_ms BIGINT UNSIGNED NULL COMMENT ''撤销毫秒时间戳'' AFTER refresh_expires_at_ms',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

-- 仅在旧日期字段仍完整存在时回填；已完成迁移的环境不会重复计算时间。
SET @kids_legacy_columns := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name IN ('created_at', 'expires_at', 'refresh_expires_at', 'revoked_at')
);
SET @kids_sql := IF(@kids_legacy_columns = 4,
    'UPDATE kids_identity_session SET issued_at_ms = CAST(UNIX_TIMESTAMP(created_at) * 1000 AS UNSIGNED), access_expires_at_ms = CAST(UNIX_TIMESTAMP(expires_at) * 1000 AS UNSIGNED), refresh_expires_at_ms = CASE WHEN refresh_expires_at IS NULL THEN NULL ELSE CAST(UNIX_TIMESTAMP(refresh_expires_at) * 1000 AS UNSIGNED) END, revoked_at_ms = CASE WHEN revoked_at IS NULL THEN NULL ELSE CAST(UNIX_TIMESTAMP(revoked_at) * 1000 AS UNSIGNED) END, status = CASE WHEN revoked_at IS NULL THEN ''active'' ELSE ''revoked'' END WHERE issued_at_ms IS NULL OR access_expires_at_ms IS NULL',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

-- 回填完成后收紧新字段约束。
ALTER TABLE kids_identity_session
    MODIFY COLUMN status VARCHAR(32) NOT NULL DEFAULT 'active' COMMENT '会话状态',
    MODIFY COLUMN issued_at_ms BIGINT UNSIGNED NOT NULL COMMENT '会话签发毫秒时间戳',
    MODIFY COLUMN access_expires_at_ms BIGINT UNSIGNED NOT NULL COMMENT '访问令牌到期毫秒时间戳';

-- 逐列清理旧日期字段，兼容此前已执行完成的环境。
SET @kids_column_exists := (
    SELECT COUNT(*) FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'kids_identity_session' AND column_name = 'expires_at'
);
SET @kids_sql := IF(@kids_column_exists = 1, 'ALTER TABLE kids_identity_session DROP COLUMN expires_at', 'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*) FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'kids_identity_session' AND column_name = 'refresh_expires_at'
);
SET @kids_sql := IF(@kids_column_exists = 1, 'ALTER TABLE kids_identity_session DROP COLUMN refresh_expires_at', 'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*) FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'kids_identity_session' AND column_name = 'revoked_at'
);
SET @kids_sql := IF(@kids_column_exists = 1, 'ALTER TABLE kids_identity_session DROP COLUMN revoked_at', 'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*) FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'kids_identity_session' AND column_name = 'created_at'
);
SET @kids_sql := IF(@kids_column_exists = 1, 'ALTER TABLE kids_identity_session DROP COLUMN created_at', 'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*) FROM information_schema.columns
    WHERE table_schema = DATABASE() AND table_name = 'kids_identity_session' AND column_name = 'updated_at'
);
SET @kids_sql := IF(@kids_column_exists = 1, 'ALTER TABLE kids_identity_session DROP COLUMN updated_at', 'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

-- 仅在索引尚未存在时创建，保证中断后重跑不会失败。
SET @kids_index_exists := (
    SELECT COUNT(*)
    FROM information_schema.statistics
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND index_name = 'idx_status_expires'
);
SET @kids_sql := IF(@kids_index_exists = 0,
    'ALTER TABLE kids_identity_session ADD KEY idx_status_expires (status, access_expires_at_ms)',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;
