/*==============================================================*/
/* DBMS name:      MySQL 5.0                                    */
/* Created on:     2015/8/3 18:12:35                            */
/*==============================================================*/


drop table if exists o_cart;

drop index Index_usr on o_orders;

drop index Index_url on o_orders;

drop table if exists o_orders;

drop table if exists o_orders_item;

drop table if exists o_usr;

/*==============================================================*/
/* Table: o_cart                                                */
/*==============================================================*/
create table o_cart
(
   ID                   int(11) not null auto_increment comment '购物车ID',
   USER_ID              int(11) comment '用户ID',
   GOODS_ID             int(11) not null comment '商品ID',
   GOODS_PRICE          double(13,2) not null comment '商品价格',
   GOODS_NUM            int(11) not null comment '商品数量',
   GOODS_IMG            varchar(255) not null comment '商品图片URL',
   GOODS_NAME           varchar(255) not null comment '商品名称',
   GOODS_URL            varchar(255) not null comment '商品URL',
   SHOP_ID              int(11) comment '店ID',
   SHOP_NAME            varchar(255) comment '店名',
   SHOP_URL             varchar(255) comment '店URL',
   SELLER_ID            int(11) not null comment '卖家ID',
   SELLER_NAME          varchar(255) not null comment '卖家姓名',
   FROM_URL             varchar(255) not null comment '子系统URL',
   STATUS               int(11) default 0 comment '购买状态（0未购买，10已购买，20直接购买,30从购物车移除）',
   TIME                 timestamp not null comment '操作时间',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_cart comment '所有子系统物品放在同一购物车。';

/*==============================================================*/
/* Table: o_orders                                              */
/*==============================================================*/
create table o_orders
(
   ID                   int(11) not null auto_increment comment '订单ID',
   ORD_NO               varchar(50) not null comment '订单号',
   PUB_ORD_NO           varchar(50) not null comment '合并订单号',
   USR_ID               int(11) not null comment '用户ID',
   SHOP_ID              int(11) comment '店ID',
   SHOP_NAME            varchar(255) comment '店名',
   SHOP_URL             varchar(255) comment '店URL',
   SELLER_ID            int(11) not null comment '卖家ID',
   SELLER_NAME          varchar(255) not null comment '卖家姓名',
   ORD_TYPE             int(11) default 0 comment '订单类型(0WEB订单，10手机订单)',
   STATUS               int(11) default 0 comment '订单状态(0未付款，10已付款，20交易完成，30订单取消，40交易失败，50包含退款,60已退款,70免费订单,80删除的订单)',
   FROM_URL             varchar(255) not null comment '子系统URL',
   PRICE                double(13,2) comment '支付宝支付总价',
   CREATE_TIME          datetime not null comment '创建时间',
   PAY_TIME             datetime comment '支付时间',
   TIME                 timestamp not null comment '操作时间',
   SYNC                 int(11) default 0 comment '是否同步完成（0未同步，1已同步）',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_orders comment '订单表存放各个子系统的所有订单';

/*==============================================================*/
/* Index: Index_url                                             */
/*==============================================================*/
create index Index_url on o_orders
(
   FROM_URL
);

/*==============================================================*/
/* Index: Index_usr                                             */
/*==============================================================*/
create index Index_usr on o_orders
(
   USR_ID
);

/*==============================================================*/
/* Table: o_orders_item                                         */
/*==============================================================*/
create table o_orders_item
(
   ID                   int(11) not null auto_increment comment '订单详情ID',
   ORD_ID               int(11) not null comment '订单ID',
   GOODS_PRICE          double(13,2) not null comment '商品价格',
   GOODS_ID             int(11) not null comment '商品ID',
   GOODS_NUM            int(11) not null comment '商品数量',
   GOODS_IMG            varchar(255) not null comment '商品图片URL',
   GOODS_NAME           varchar(255) not null comment '商品名称',
   GOODS_URL            varchar(255) not null comment '商品URL',
   STATUS               int(11) default 0 comment '商品状态(0正常,10申请退款,20已退款)',
   TIME                 timestamp not null comment '操作时间',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_orders_item comment '订单下的详细数据';

/*==============================================================*/
/* Table: o_usr                                                 */
/*==============================================================*/
create table o_usr
(
   ID                   int(11) not null auto_increment,
   USR_ID               int(11) not null comment '用户ID',
   GOODS_ID             int(11) not null comment '商品ID',
   STATUS               int(11) default 0 comment '状态（0未购买，10已购买，11已确认收货,20已退款）',
   TOKEN                varchar(255) not null comment '操作TOKEN(唯一标识)',
   TIME                 timestamp not null comment '操作时间',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_usr comment '商品与用户绑定表，存放在子系统，自动生成';

