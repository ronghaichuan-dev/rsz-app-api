-- Clearwave 接口幂等首响外层元数据。
-- 重放必须保留首次响应的 change_cursor 与 etag，不能只缓存 data 子对象。

ALTER TABLE kids_request_deduplication
    ADD COLUMN response_change_cursor VARCHAR(263) NULL COMMENT '首次响应变更游标' AFTER response_body,
    ADD COLUMN response_etag VARCHAR(256) NULL COMMENT '首次响应实体标签' AFTER response_change_cursor;
