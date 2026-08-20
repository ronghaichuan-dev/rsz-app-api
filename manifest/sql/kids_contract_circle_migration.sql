-- Clearwave 合同账号与圈子领域事实表。
-- 所有标识均使用合同 opaque ID，避免旧 kids 自增标识成为外部授权或同步事实。

CREATE TABLE IF NOT EXISTS kids_contract_account (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    account_id VARCHAR(128) NOT NULL COMMENT '合同账号标识',
    status VARCHAR(32) NOT NULL COMMENT '账号状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_account_id (account_id),
    KEY idx_status_updated (status, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同账号';

CREATE TABLE IF NOT EXISTS kids_contract_account_binding (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    binding_id VARCHAR(128) NOT NULL COMMENT '合同绑定标识',
    account_id VARCHAR(128) NOT NULL COMMENT '合同账号标识',
    environment VARCHAR(16) NOT NULL COMMENT '绑定环境',
    migration_policy VARCHAR(32) NOT NULL COMMENT '迁移策略',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    issued_at DATETIME NOT NULL COMMENT '签发时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_binding_id (binding_id),
    UNIQUE KEY uk_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同账号绑定';

CREATE TABLE IF NOT EXISTS kids_contract_circle (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    name VARCHAR(80) NOT NULL COMMENT '圈子名称',
    icon JSON NOT NULL COMMENT '圈子视觉引用',
    owner_administrator_id VARCHAR(128) NOT NULL COMMENT '所有者管理员标识',
    status VARCHAR(32) NOT NULL COMMENT '圈子状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_circle_id (circle_id),
    KEY idx_status_updated (status, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同圈子';

CREATE TABLE IF NOT EXISTS kids_contract_administrator (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    administrator_id VARCHAR(128) NOT NULL COMMENT '合同管理员标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    bound_account_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '绑定账号标识',
    display_name VARCHAR(120) NOT NULL COMMENT '显示名称',
    avatar JSON NOT NULL COMMENT '头像视觉引用',
    role VARCHAR(32) NOT NULL COMMENT '管理员角色',
    permissions JSON NOT NULL COMMENT '权限集合',
    status VARCHAR(32) NOT NULL COMMENT '管理员状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_administrator_id (administrator_id),
    KEY idx_circle_status (circle_id, status),
    KEY idx_account_status (bound_account_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同管理员';

CREATE TABLE IF NOT EXISTS kids_contract_member (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    member_id VARCHAR(128) NOT NULL COMMENT '合同成员标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    bound_account_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '绑定账号标识',
    display_name VARCHAR(120) NOT NULL COMMENT '显示名称',
    gender VARCHAR(16) NOT NULL COMMENT '性别',
    avatar JSON NOT NULL COMMENT '头像视觉引用',
    status VARCHAR(32) NOT NULL COMMENT '成员状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_member_id (member_id),
    KEY idx_circle_status (circle_id, status),
    KEY idx_account_status (bound_account_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同成员';

CREATE TABLE IF NOT EXISTS kids_contract_membership (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    membership_id VARCHAR(128) NOT NULL COMMENT '合同成员身份标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    account_id VARCHAR(128) NOT NULL COMMENT '合同账号标识',
    actor_type VARCHAR(32) NOT NULL COMMENT '授权主体类型',
    actor_id VARCHAR(128) NOT NULL COMMENT '授权主体标识',
    role VARCHAR(32) NOT NULL COMMENT '成员角色',
    permissions JSON NOT NULL COMMENT '权限集合',
    status VARCHAR(32) NOT NULL COMMENT '成员身份状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_membership_id (membership_id),
    UNIQUE KEY uk_circle_account (circle_id, account_id),
    KEY idx_account_status (account_id, status),
    KEY idx_actor_id (actor_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同圈子成员身份';

CREATE TABLE IF NOT EXISTS kids_contract_circle_selection (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    selection_id VARCHAR(128) NOT NULL COMMENT '合同选择标识',
    account_id VARCHAR(128) NOT NULL COMMENT '合同账号标识',
    current_circle_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '当前圈子标识',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_selection_id (selection_id),
    UNIQUE KEY uk_account_id (account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同圈子选择';
