-- 清理 kids 会话已被毫秒时间戳字段替代的旧日期字段。
-- 本迁移可在字段已删除的环境重复执行，不会影响新的毫秒时间字段。

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'expires_at'
);
SET @kids_sql := IF(@kids_column_exists = 1,
    'ALTER TABLE kids_identity_session DROP COLUMN expires_at',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'refresh_expires_at'
);
SET @kids_sql := IF(@kids_column_exists = 1,
    'ALTER TABLE kids_identity_session DROP COLUMN refresh_expires_at',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'revoked_at'
);
SET @kids_sql := IF(@kids_column_exists = 1,
    'ALTER TABLE kids_identity_session DROP COLUMN revoked_at',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'created_at'
);
SET @kids_sql := IF(@kids_column_exists = 1,
    'ALTER TABLE kids_identity_session DROP COLUMN created_at',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;

SET @kids_column_exists := (
    SELECT COUNT(*)
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_identity_session'
        AND column_name = 'updated_at'
);
SET @kids_sql := IF(@kids_column_exists = 1,
    'ALTER TABLE kids_identity_session DROP COLUMN updated_at',
    'SELECT 1');
PREPARE kids_migration_statement FROM @kids_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;
