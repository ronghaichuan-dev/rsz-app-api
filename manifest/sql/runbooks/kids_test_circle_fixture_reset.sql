-- kids 测试圈子 fixture 受控重置手册。
--
-- 本文件不是数据库 migration，禁止纳入任何自动执行的 migration 队列，禁止在生产库执行。
-- 它只删除指定测试圈子的可重建 v1 领域事实和同步投影，保留账号、会话、圈子、管理员与成员身份，
-- 以便操作人员随后使用 API 重建全新的测试奖励 fixture。它绝不更新既有 ledger_id、exchange_id
-- 或同步提交载荷中的任何字段，避免同一标识出现不同载荷。
--
-- 使用前必须：
-- 1. 停止该测试圈子的写入和部署 smoke；
-- 2. 确认目标为专用测试账号及圈子，并记录工单号；
-- 3. 确认所有测试客户端会清除该圈子的本地可重建投影后重新登录；
-- 4. 从受控终端连接测试数据库后，先执行本文件创建存储过程，再使用文末 CALL 执行。

CREATE TABLE IF NOT EXISTS kids_test_fixture_reset_audit (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    reset_id CHAR(36) NOT NULL COMMENT '重置标识',
    ticket VARCHAR(128) NOT NULL COMMENT '变更工单标识',
    environment VARCHAR(16) NOT NULL COMMENT '执行环境',
    database_name VARCHAR(64) NOT NULL COMMENT '执行数据库名称',
    circle_id VARCHAR(128) NOT NULL COMMENT '重置圈子标识',
    circle_name VARCHAR(80) NOT NULL COMMENT '重置前圈子名称',
    operator_name VARCHAR(128) NOT NULL COMMENT '执行人员标识',
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '执行时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_reset_id (reset_id),
    KEY idx_circle_created (circle_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids测试fixture重置审计';

DELIMITER //

DROP PROCEDURE IF EXISTS kids_reset_test_circle_fixture //

CREATE PROCEDURE kids_reset_test_circle_fixture(
    IN p_environment VARCHAR(16),
    IN p_expected_database VARCHAR(64),
    IN p_circle_id VARCHAR(128),
    IN p_expected_circle_name VARCHAR(80),
    IN p_ticket VARCHAR(128),
    IN p_operator_name VARCHAR(128),
    IN p_confirmation VARCHAR(64)
)
main: BEGIN
    DECLARE v_circle_name VARCHAR(80) DEFAULT NULL;
    DECLARE v_database_name VARCHAR(64) DEFAULT NULL;
    DECLARE v_reset_id CHAR(36);

    -- 环境、目标、人工确认与工单均为必填，避免空条件或非测试库误删。
    IF p_environment IS NULL OR p_environment <> 'test' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '仅允许在 test 环境重置测试 fixture';
    END IF;
    IF p_expected_database IS NULL OR CHAR_LENGTH(TRIM(p_expected_database)) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '必须提供当前测试数据库名称';
    END IF;
    IF p_circle_id IS NULL OR CHAR_LENGTH(TRIM(p_circle_id)) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '必须提供测试圈子标识';
    END IF;
    IF p_expected_circle_name IS NULL OR CHAR_LENGTH(TRIM(p_expected_circle_name)) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '必须提供重置前圈子精确名称';
    END IF;
    IF p_ticket IS NULL OR CHAR_LENGTH(TRIM(p_ticket)) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '必须提供变更工单标识';
    END IF;
    IF p_operator_name IS NULL OR CHAR_LENGTH(TRIM(p_operator_name)) = 0 THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '必须提供执行人员标识';
    END IF;
    IF p_confirmation IS NULL OR p_confirmation <> 'RESET_ISOLATED_TEST_FIXTURE' THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '确认字符串不匹配，拒绝重置';
    END IF;

    SET v_database_name = DATABASE();
    IF v_database_name IS NULL OR v_database_name <> p_expected_database THEN
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '当前数据库与确认的测试数据库不一致';
    END IF;

    START TRANSACTION;

    -- 锁定并核对目标圈子名称，防止错误圈子标识被误用于重置。
    SELECT name
    INTO v_circle_name
    FROM kids_circle_info
    WHERE circle_id = p_circle_id
    FOR UPDATE;

    IF v_circle_name IS NULL OR v_circle_name <> p_expected_circle_name THEN
        ROLLBACK;
        SIGNAL SQLSTATE '45000' SET MESSAGE_TEXT = '测试圈子不存在或名称与预期不一致';
    END IF;

    SET v_reset_id = UUID();

    -- 先固定待删除标识，所有后续删除均只能命中该圈子当前的可重建 fixture。
    CREATE TEMPORARY TABLE kids_reset_target_commits (
        commit_id VARCHAR(128) NOT NULL,
        PRIMARY KEY (commit_id)
    ) ENGINE=Memory;
    INSERT INTO kids_reset_target_commits (commit_id)
    SELECT commit_id
    FROM kids_sync_commit
    WHERE circle_id = p_circle_id;

    CREATE TEMPORARY TABLE kids_reset_target_tasks (
        task_id VARCHAR(128) NOT NULL,
        PRIMARY KEY (task_id)
    ) ENGINE=Memory;
    INSERT INTO kids_reset_target_tasks (task_id)
    SELECT task_id
    FROM kids_task_definition
    WHERE circle_id = p_circle_id;

    CREATE TEMPORARY TABLE kids_reset_target_rewards (
        reward_id VARCHAR(128) NOT NULL,
        PRIMARY KEY (reward_id)
    ) ENGINE=Memory;
    INSERT INTO kids_reset_target_rewards (reward_id)
    SELECT reward_id
    FROM kids_reward_definition
    WHERE circle_id = p_circle_id;

    CREATE TEMPORARY TABLE kids_reset_target_completions (
        completion_id VARCHAR(128) NOT NULL,
        PRIMARY KEY (completion_id)
    ) ENGINE=Memory;
    INSERT INTO kids_reset_target_completions (completion_id)
    SELECT completion_id
    FROM kids_task_completion
    WHERE circle_id = p_circle_id;

    -- 先删除依赖审计事实，再删除其投影和定义；不对任何历史行执行 UPDATE。
    DELETE cancellation
    FROM kids_task_cancellation AS cancellation
    INNER JOIN kids_reset_target_completions AS completion
        ON completion.completion_id = cancellation.completion_id;

    DELETE FROM kids_task_completion WHERE circle_id = p_circle_id;
    DELETE FROM kids_task_occurrence WHERE circle_id = p_circle_id;

    DELETE assignment_record
    FROM kids_task_assignment AS assignment_record
    INNER JOIN kids_reset_target_tasks AS task_record
        ON task_record.task_id = assignment_record.task_id;
    DELETE FROM kids_task_definition WHERE circle_id = p_circle_id;
    DELETE FROM kids_task_tag_definition WHERE circle_id = p_circle_id;

    DELETE FROM kids_notification_outbox WHERE circle_id = p_circle_id;
    DELETE FROM kids_reward_exchange WHERE circle_id = p_circle_id;
    DELETE FROM kids_star_ledger WHERE circle_id = p_circle_id;
    DELETE FROM kids_star_balance WHERE circle_id = p_circle_id;

    DELETE assignment_record
    FROM kids_reward_assignment AS assignment_record
    INNER JOIN kids_reset_target_rewards AS reward_record
        ON reward_record.reward_id = assignment_record.reward_id;
    DELETE cooldown
    FROM kids_reward_cooldown AS cooldown
    INNER JOIN kids_reset_target_rewards AS reward_record
        ON reward_record.reward_id = cooldown.reward_id;
    DELETE FROM kids_reward_definition WHERE circle_id = p_circle_id;

    DELETE FROM kids_invitation WHERE circle_id = p_circle_id;
    DELETE FROM kids_asset WHERE circle_id = p_circle_id;
    DELETE FROM kids_asset_upload WHERE circle_id = p_circle_id;

    -- 幂等快照和回执同样属于旧 fixture，必须随同步提交一起删除，不能改写其返回载荷。
    DELETE deduplication
    FROM kids_request_deduplication AS deduplication
    INNER JOIN kids_reset_target_commits AS target_commit
        ON JSON_UNQUOTE(JSON_EXTRACT(deduplication.response_body, '$.receipt.commit_id')) = target_commit.commit_id;
    DELETE receipt
    FROM kids_mutation_receipt AS receipt
    INNER JOIN kids_reset_target_commits AS target_commit
        ON target_commit.commit_id = receipt.commit_id;
    DELETE FROM kids_sync_commit WHERE circle_id = p_circle_id;

    INSERT INTO kids_test_fixture_reset_audit (
        reset_id, ticket, environment, database_name, circle_id, circle_name, operator_name
    ) VALUES (
        v_reset_id, p_ticket, p_environment, v_database_name, p_circle_id, v_circle_name, p_operator_name
    );

    COMMIT;

    DROP TEMPORARY TABLE IF EXISTS kids_reset_target_completions;
    DROP TEMPORARY TABLE IF EXISTS kids_reset_target_rewards;
    DROP TEMPORARY TABLE IF EXISTS kids_reset_target_tasks;
    DROP TEMPORARY TABLE IF EXISTS kids_reset_target_commits;

    SELECT v_reset_id AS reset_id, p_circle_id AS circle_id, v_circle_name AS circle_name;
END //

DELIMITER ;

-- 示例：仅在确认连接的是测试数据库、目标圈子与名称均正确后，替换尖括号内参数再执行。
-- CALL kids_reset_test_circle_fixture(
--     'test',
--     '<测试数据库名称>',
--     '<测试圈子ID>',
--     '<测试圈子精确名称>',
--     '<变更工单号>',
--     '<执行人员>',
--     'RESET_ISOLATED_TEST_FIXTURE'
-- );

-- 重建完成后，使用下列查询确认旧审计标识已不可见，再让客户端清除本地测试数据并重新登录。
-- SELECT COUNT(*) AS 剩余流水 FROM kids_star_ledger WHERE circle_id = '<测试圈子ID>';
-- SELECT COUNT(*) AS 剩余兑换 FROM kids_reward_exchange WHERE circle_id = '<测试圈子ID>';
-- SELECT COUNT(*) AS 剩余同步提交 FROM kids_sync_commit WHERE circle_id = '<测试圈子ID>';
