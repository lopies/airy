SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for role
-- ----------------------------
DROP TABLE IF EXISTS `role`;
CREATE TABLE `role` (
                            `id` bigint NOT NULL AUTO_INCREMENT,
                            `name` varchar(64) NOT NULL COMMENT '角色名称',
                            `desc` varchar(64) NOT NULL COMMENT '角色描述',
                            `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态：1正常(默认) 0停用',
                            `created_at` bigint NOT NULL,
                            `updated_at` bigint NOT NULL,
                            PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of role
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
                            `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                            `user_name` varchar(64) NOT NULL,
                            `nick_name` varchar(64) NOT NULL COMMENT '昵称',
                            `real_name` varchar(64) DEFAULT NULL COMMENT '真实姓名',
                            `phone` varchar(16) NOT NULL,
                            `avatar` varchar(128) DEFAULT NULL COMMENT '头像',
                            `password` varchar(64) NOT NULL COMMENT '密码',
                            `salt` varchar(32) NOT NULL COMMENT '密码',
                            `status` tinyint unsigned NOT NULL COMMENT '状态 1：正常 2：禁用',
                            `register_time` varchar(16) NOT NULL COMMENT '注册时间',
                            `register_ip` varchar(32) NOT NULL COMMENT '注册ip',
                            `lotime` varchar(16) NOT NULL COMMENT '登录时间',
                            `loip` varchar(32) NOT NULL COMMENT '登录ip',
                            `created_at` bigint unsigned NOT NULL,
                            `updated_at` bigint unsigned NOT NULL,
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user
-- ----------------------------
BEGIN;
INSERT INTO `user` (`id`, `user_name`, `nick_name`, `real_name`, `phone`, `avatar`, `password`, `salt`, `status`, `register_time`, `register_ip`, `lotime`, `loip`, `created_at`, `updated_at`) VALUES (1, 'hanlin', 'hanlin', 'hanlin', '13888888888', NULL, '123456', '123', 1, '1668084161', '127.0.0.1', '1668084161', '127.0.0.1', 1668084161, 1668084161);
COMMIT;

-- ----------------------------
-- Table structure for user_info
-- ----------------------------
DROP TABLE IF EXISTS `user_info`;
CREATE TABLE `user_info` (
                                 `id` bigint unsigned NOT NULL AUTO_INCREMENT,
                                 `user_id` bigint unsigned DEFAULT NULL COMMENT '用户ID',
                                 `role_ids` varchar(64) DEFAULT NULL COMMENT '角色ID 例如：1,2,3',
                                 `created_at` bigint unsigned DEFAULT NULL,
                                 PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of user_info
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
