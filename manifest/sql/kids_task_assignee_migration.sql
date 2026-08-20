-- 既有数据库升级：为任务分配表补充分儿童完成状态字段。
ALTER TABLE kids_task_assignee
    ADD COLUMN assignee_order INT NOT NULL DEFAULT 0 COMMENT '分配顺序，用于轮流模式' AFTER kid_id,
    ADD COLUMN completed TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '该儿童是否已完成' AFTER assignee_order,
    ADD COLUMN photo_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '该儿童照片凭证地址' AFTER completed,
    ADD COLUMN completed_at DATETIME NULL COMMENT '该儿童完成时间' AFTER photo_url,
    ADD COLUMN updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间' AFTER created_at,
    ADD KEY idx_task_completed (task_id, completed);
