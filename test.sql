-- ================================
-- test.sql
-- ================================


-- ================================
-- 用户测试数据（密码都是 123456）
-- ================================
INSERT INTO "user" (username, password, phone, email, is_logged_in, created_at, updated_at)
VALUES ('test1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13800000001', 'test1@test.com', false,
        NOW(), NOW()),
       ('test2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13800000002', 'test2@test.com', false,
        NOW(), NOW()),
       ('test3', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', '13800000003', 'test3@test.com', false,
        NOW(), NOW());

-- ================================
-- 好友关系（test1 和 test2 互为好友）
-- ================================
INSERT INTO friendships (user_id, friend_id, status, created_at, updated_at)
VALUES (1, 2, 'accepted', NOW(), NOW()),
       (2, 1, 'accepted', NOW(), NOW()),
       -- test3 申请加 test1，未处理
       (3, 1, 'pending', NOW(), NOW());

-- ================================
-- 私聊消息
-- ================================
INSERT INTO private_messages (from_id, to_id, content, type, read_at, created_at, updated_at)
VALUES (1, 2, 'test1: 你好', 'text', EXTRACT(EPOCH FROM NOW())::BIGINT, NOW() - INTERVAL '10 minutes', NOW()),
       (2, 1, 'test2: 你好啊', 'text', EXTRACT(EPOCH FROM NOW())::BIGINT, NOW() - INTERVAL '9 minutes', NOW()),
       (1, 2, 'test1: 今晚打游戏', 'text', EXTRACT(EPOCH FROM NOW())::BIGINT, NOW() - INTERVAL '8 minutes', NOW()),
       (2, 1, 'test2: 好啊几点', 'text', NULL, NOW() - INTERVAL '7 minutes', NOW()), -- 未读
       (2, 1, 'test2: 等我消息', 'text', NULL, NOW() - INTERVAL '6 minutes', NOW());
-- 未读

-- ================================
-- 群组
-- ================================
INSERT INTO groups (name, owner_id, notice, created_at, updated_at)
VALUES ('测试群', 1, '欢迎加入测试群', NOW(), NOW());

-- ================================
-- 群成员（test1 owner，test2 member）
-- ================================
INSERT INTO group_members (group_id, user_id, role, created_at, updated_at)
VALUES (1, 1, 'owner', NOW(), NOW()),
       (1, 2, 'member', NOW(), NOW());

-- ================================
-- 群消息
-- ================================
INSERT INTO group_messages (group_id, from_id, content, type, created_at, updated_at)
VALUES (1, 1, 'test1: 大家好', 'text', NOW() - INTERVAL '5 minutes', NOW()),
       (1, 2, 'test2: 大家好', 'text', NOW() - INTERVAL '4 minutes', NOW()),
       (1, 1, 'test1: 今晚开会', 'text', NOW() - INTERVAL '3 minutes', NOW());

-- ================================
-- 验证数据
-- ================================
SELECT 'user count' AS check_item, COUNT(*) AS cnt
FROM "user";
SELECT 'friendship' AS check_item, COUNT(*) AS cnt
FROM friendships;
SELECT 'private_msg' AS check_item, COUNT(*) AS cnt
FROM private_messages;
SELECT 'group' AS check_item, COUNT(*) AS cnt
FROM groups;
SELECT 'group_member' AS check_item, COUNT(*) AS cnt
FROM group_members;
SELECT 'group_msg' AS check_item, COUNT(*) AS cnt
FROM group_messages;