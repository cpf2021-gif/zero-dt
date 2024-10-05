create database if not exists `dt_stock`;
use `dt_stock`;

drop table if exists `stock`;
create table `stock` (
    `id` bigint not null auto_increment,
    `goods_id` bigint not null default '0' comment '商品ID',
    `num` int not null default '0' comment '库存数量',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_goods_id` (`goods_id`)
) Engine=InnoDB default charset=utf8mb4 collate=utf8mb4_general_ci;

begin;
insert into `stock` (`id`, `goods_id`, `num`) values (1, 1, 100);
commit;