-- kids家庭成员去除角色字段：成员不再区分 adult/kid。
-- 执行前请确认当前环境的 kids_family_member 表仍存在 idx_role 索引和 role 字段。
ALTER TABLE kids_family_member
    DROP INDEX idx_role,
    DROP COLUMN role;
