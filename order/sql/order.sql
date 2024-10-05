create database if not exists `dt_order`;
use `dt_order`;

drop table if exists `order`;
create table `order` (
    `id` bigint not null auto_increment,
    `user_id` bigint not null default '0' comment '用户ID',
    `goods_id` bigint not null default '0' comment '商品ID',
    `num` int not null default '0' comment '购买数量',
    `row_state` tinyint not null default '0' comment '-1:下单回滚 0:待支付 1:已支付',
    PRIMARY KEY (`id`),
) Engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;