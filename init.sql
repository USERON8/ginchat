-- ================================
-- init.sql
-- ================================


-- ================================
-- 用户表
-- ================================
CREATE TABLE "user"
(
    id             BIGSERIAL PRIMARY KEY,
    created_at     TIMESTAMP,
    updated_at     TIMESTAMP,
    deleted_at     TIMESTAMP,
    username       VARCHAR(50)  NOT NULL,
    password       VARCHAR(255) NOT NULL,
    phone          VARCHAR(20),
    email          VARCHAR(100),
    identity       VARCHAR(255),
    client_ip      VARCHAR(50),
    client_port    VARCHAR(10),
    login_time     BIGINT,
    log_out_time   BIGINT,
    is_logged_in   BOOLEAN DEFAULT FALSE,
    is_admin       BOOLEAN DEFAULT FALSE,
    device_info    VARCHAR(255),
    heartbeat_time BIGINT
);

-- 用户表索引
CREATE UNIQUE INDEX idx_user_username ON "user" (username) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_deleted ON "user" (deleted_at);
CREATE INDEX idx_user_phone ON "user" (phone) WHERE deleted_at IS NULL;
CREATE INDEX idx_user_email ON "user" (email) WHERE deleted_at IS NULL;

-- ================================
-- 好友关系表
-- ================================
CREATE TABLE friendships
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    user_id    BIGINT NOT NULL,
    friend_id  BIGINT NOT NULL,
    status     VARCHAR(20) DEFAULT 'pending'
);

-- 好友表索引
CREATE INDEX idx_friendship_deleted ON friendships (deleted_at);
-- 查好友列表：WHERE user_id=? AND status=?
CREATE INDEX idx_friendship_user_status ON friendships (user_id, status) WHERE deleted_at IS NULL;
-- 查收到的申请：WHERE friend_id=? AND status=?
CREATE INDEX idx_friendship_friend_status ON friendships (friend_id, status) WHERE deleted_at IS NULL;
-- 防重复申请：WHERE user_id=? AND friend_id=?
CREATE UNIQUE INDEX idx_friendship_pair ON friendships (user_id, friend_id) WHERE deleted_at IS NULL;

-- ================================
-- 私聊消息表
-- ================================
CREATE TABLE private_messages
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    from_id    BIGINT NOT NULL,
    to_id      BIGINT NOT NULL,
    content    TEXT   NOT NULL,
    type       VARCHAR(20) DEFAULT 'text',
    read_at    BIGINT
);

-- 私聊消息索引
CREATE INDEX idx_private_msg_deleted ON private_messages (deleted_at);
-- 查聊天记录：WHERE (from_id=? AND to_id=?) OR (from_id=? AND to_id=?) ORDER BY id DESC
CREATE INDEX idx_private_msg_chat ON private_messages (from_id, to_id, id DESC) WHERE deleted_at IS NULL;
-- 查未读消息：WHERE to_id=? AND read_at IS NULL
CREATE INDEX idx_private_msg_unread ON private_messages (to_id, read_at) WHERE deleted_at IS NULL;

-- ================================
-- 群组表
-- ================================
CREATE TABLE groups
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    name       VARCHAR(50) NOT NULL,
    owner_id   BIGINT      NOT NULL,
    avatar     VARCHAR(255),
    notice     TEXT
);

-- 群组索引
CREATE INDEX idx_group_deleted ON groups (deleted_at);
-- 查某人创建的群
CREATE INDEX idx_group_owner ON groups (owner_id) WHERE deleted_at IS NULL;

-- ================================
-- 群成员表
-- ================================
CREATE TABLE group_members
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    group_id   BIGINT NOT NULL,
    user_id    BIGINT NOT NULL,
    role       VARCHAR(20) DEFAULT 'member'
);

-- 群成员索引
CREATE INDEX idx_group_member_deleted ON group_members (deleted_at);
-- 查群成员列表：WHERE group_id=?
CREATE INDEX idx_group_member_group ON group_members (group_id) WHERE deleted_at IS NULL;
-- 查我加入的群：WHERE user_id=?
CREATE INDEX idx_group_member_user ON group_members (user_id) WHERE deleted_at IS NULL;
-- 验证是否在群里：WHERE group_id=? AND user_id=?
CREATE UNIQUE INDEX idx_group_member_pair ON group_members (group_id, user_id) WHERE deleted_at IS NULL;

-- ================================
-- 群聊消息表
-- ================================
CREATE TABLE group_messages
(
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    group_id   BIGINT NOT NULL,
    from_id    BIGINT NOT NULL,
    content    TEXT   NOT NULL,
    type       VARCHAR(20) DEFAULT 'text'
);

-- 群消息索引
CREATE INDEX idx_group_msg_deleted ON group_messages (deleted_at);
-- 查群聊历史：WHERE group_id=? ORDER BY id DESC（游标分页）
CREATE INDEX idx_group_msg_chat ON group_messages (group_id, id DESC) WHERE deleted_at IS NULL;