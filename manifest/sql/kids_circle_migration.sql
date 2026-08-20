-- 既有数据库升级：补充圈子、邀请码和家庭成员形象字段。
CREATE TABLE IF NOT EXISTS kids_circle (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '圈子ID',
    name VARCHAR(64) NOT NULL COMMENT '圈子名称',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '圈子图标标识',
    owner_user_id BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_owner_user_id (owner_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids圈子';

CREATE TABLE IF NOT EXISTS kids_circle_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '圈子用户关系ID',
    circle_id BIGINT UNSIGNED NOT NULL COMMENT '圈子ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role VARCHAR(16) NOT NULL COMMENT '圈子角色：admin/member',
    member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '绑定的家庭成员ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_circle_user (circle_id, user_id),
    KEY idx_user_id (user_id),
    KEY idx_member_id (member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids圈子用户关系';

CREATE TABLE IF NOT EXISTS kids_invite_code (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '邀请码ID',
    code CHAR(6) NOT NULL COMMENT '六位邀请码',
    circle_id BIGINT UNSIGNED NOT NULL COMMENT '圈子ID',
    inviter_user_id BIGINT UNSIGNED NOT NULL COMMENT '邀请人用户ID',
    invite_role VARCHAR(16) NOT NULL COMMENT '邀请加入角色：admin/member',
    target_member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标家庭成员ID，0表示不指定',
    expired_at DATETIME NOT NULL COMMENT '过期时间',
    used_at DATETIME NULL COMMENT '使用时间',
    used_by_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '使用者用户ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code),
    KEY idx_circle_role (circle_id, invite_role),
    KEY idx_target_member_id (target_member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids邀请码';

ALTER TABLE kids_family_member
    ADD COLUMN circle_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属圈子ID' AFTER id,
    ADD COLUMN gender VARCHAR(16) NOT NULL DEFAULT '' COMMENT '性别：male/female' AFTER name,
    ADD COLUMN avatar_style VARCHAR(64) NOT NULL DEFAULT '' COMMENT '虚拟形象风格标识' AFTER avatar,
    ADD KEY idx_circle_id (circle_id);
