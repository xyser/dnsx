CREATE TABLE `records` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `name` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '记录名称',
  `type` varchar(32) COLLATE utf8_unicode_ci NOT NULL COMMENT '资源类型',
  `value` varchar(255) COLLATE utf8_unicode_ci NOT NULL COMMENT '记录值',
  `ttl` int(10) unsigned DEFAULT '600' COMMENT '缓存生存时间',
  `priority` int(10) unsigned DEFAULT '10' COMMENT '优先级(越小优先级越高)',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '记录软删除时间，默认NULL',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
