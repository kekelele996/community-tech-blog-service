-- ============================================================
-- techblog_db 初始化脚本（MySQL 容器首次启动时执行）
-- 说明：
--   1. 全部业务表由后端 GORM AutoMigrate 在启动时自动创建/校验
--      （见 backend/internal/database/database.go）。
--   2. 管理员账号由后端 SeedAdmin 幂等初始化
--      （见 backend/internal/database/seed.go）。
--   3. 本脚本仅确保数据库字符集正确，避免重复建表导致漂移。
-- ============================================================
CREATE DATABASE IF NOT EXISTS techblog_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE techblog_db;
