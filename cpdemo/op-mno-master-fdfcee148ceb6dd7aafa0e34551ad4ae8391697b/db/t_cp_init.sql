

CREATE TABLE `t_cp_sim` (
                            `id` varchar(50) COLLATE utf8mb4_general_ci NOT NULL COMMENT '主键',
                            `customer_id` bigint NOT NULL COMMENT '客户id',
                            `operator_id` bigint DEFAULT NULL COMMENT '运营商ID',
                            `client_id` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
                            `iccid` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'sim卡唯一识别码',
                            `imsi` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '网络标识码',
                            `msisdn` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '网络SDN',
                            `status` int NOT NULL COMMENT 'SIM状态 1：未生效；2：可测试；3：库存；4：激活；5：停卡',
                            `able_status` int DEFAULT '1' COMMENT '启用禁用状态1是启用  2是禁用  默认启用',
                            `is_delete` tinyint(1) NOT NULL COMMENT '是否删除 0：否；1：是',
                            `create_time` datetime NOT NULL COMMENT '创建时间',
                            `update_time` datetime DEFAULT NULL COMMENT '更新时间',
                            `push_status` int NOT NULL DEFAULT '1' COMMENT '推送状态；1：未推送，2：已推送，3：推送失败',
                            `push_time` datetime DEFAULT NULL COMMENT 'sim卡推送时间',
                            `simba_id` bigint DEFAULT NULL COMMENT 'xinbaID ',
                            `first_activate_time` datetime DEFAULT NULL COMMENT '首次激活时间',
                            `is_abnormal` tinyint DEFAULT '0' COMMENT '是否异常 1 否 2 是 ',
                            `contents` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '异常内容',
                            `effective_time` datetime DEFAULT NULL COMMENT '生效时间',
                            PRIMARY KEY (`id`),
                            UNIQUE KEY `iccid` (`iccid`) USING BTREE,
                            UNIQUE KEY `imsi` (`imsi`) USING BTREE,
                            UNIQUE KEY `msisdn` (`msisdn`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='客户SIM卡表';


CREATE TABLE `t_cp_running_bill_esim` (
                                          `id` bigint NOT NULL COMMENT '主键流水号',
                                          `iccid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                          `imsi` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                          `msisdn` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
                                          `json` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
                                          `create_time` datetime DEFAULT NULL COMMENT '创建时间',
                                          `status` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '0  失败 1成功',
                                          `effective_time` datetime DEFAULT NULL COMMENT '生效时间',
                                          PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci ROW_FORMAT=DYNAMIC COMMENT='变更流水表';





