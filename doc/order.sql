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
   ID                   int(11) not null auto_increment comment '���ﳵID',
   USER_ID              int(11) comment '�û�ID',
   GOODS_ID             int(11) not null comment '��ƷID',
   GOODS_PRICE          double(13,2) not null comment '��Ʒ�۸�',
   GOODS_NUM            int(11) not null comment '��Ʒ����',
   GOODS_IMG            varchar(255) not null comment '��ƷͼƬURL',
   GOODS_NAME           varchar(255) not null comment '��Ʒ����',
   GOODS_URL            varchar(255) not null comment '��ƷURL',
   SHOP_ID              int(11) comment '��ID',
   SHOP_NAME            varchar(255) comment '����',
   SHOP_URL             varchar(255) comment '��URL',
   SELLER_ID            int(11) not null comment '����ID',
   SELLER_NAME          varchar(255) not null comment '��������',
   FROM_URL             varchar(255) not null comment '��ϵͳURL',
   STATUS               int(11) default 0 comment '����״̬��0δ����10�ѹ���20ֱ�ӹ���,30�ӹ��ﳵ�Ƴ���',
   TIME                 timestamp not null comment '����ʱ��',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_cart comment '������ϵͳ��Ʒ����ͬһ���ﳵ��';

/*==============================================================*/
/* Table: o_orders                                              */
/*==============================================================*/
create table o_orders
(
   ID                   int(11) not null auto_increment comment '����ID',
   ORD_NO               varchar(50) not null comment '������',
   PUB_ORD_NO           varchar(50) not null comment '�ϲ�������',
   USR_ID               int(11) not null comment '�û�ID',
   SHOP_ID              int(11) comment '��ID',
   SHOP_NAME            varchar(255) comment '����',
   SHOP_URL             varchar(255) comment '��URL',
   SELLER_ID            int(11) not null comment '����ID',
   SELLER_NAME          varchar(255) not null comment '��������',
   ORD_TYPE             int(11) default 0 comment '��������(0WEB������10�ֻ�����)',
   STATUS               int(11) default 0 comment '����״̬(0δ���10�Ѹ��20������ɣ�30����ȡ����40����ʧ�ܣ�50�����˿�,60���˿�,70��Ѷ���,80ɾ���Ķ���)',
   FROM_URL             varchar(255) not null comment '��ϵͳURL',
   PRICE                double(13,2) comment '֧����֧���ܼ�',
   CREATE_TIME          datetime not null comment '����ʱ��',
   PAY_TIME             datetime comment '֧��ʱ��',
   TIME                 timestamp not null comment '����ʱ��',
   SYNC                 int(11) default 0 comment '�Ƿ�ͬ����ɣ�0δͬ����1��ͬ����',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_orders comment '�������Ÿ�����ϵͳ�����ж���';

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
   ID                   int(11) not null auto_increment comment '��������ID',
   ORD_ID               int(11) not null comment '����ID',
   GOODS_PRICE          double(13,2) not null comment '��Ʒ�۸�',
   GOODS_ID             int(11) not null comment '��ƷID',
   GOODS_NUM            int(11) not null comment '��Ʒ����',
   GOODS_IMG            varchar(255) not null comment '��ƷͼƬURL',
   GOODS_NAME           varchar(255) not null comment '��Ʒ����',
   GOODS_URL            varchar(255) not null comment '��ƷURL',
   STATUS               int(11) default 0 comment '��Ʒ״̬(0����,10�����˿�,20���˿�)',
   TIME                 timestamp not null comment '����ʱ��',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_orders_item comment '�����µ���ϸ����';

/*==============================================================*/
/* Table: o_usr                                                 */
/*==============================================================*/
create table o_usr
(
   ID                   int(11) not null auto_increment,
   USR_ID               int(11) not null comment '�û�ID',
   GOODS_ID             int(11) not null comment '��ƷID',
   STATUS               int(11) default 0 comment '״̬��0δ����10�ѹ���11��ȷ���ջ�,20���˿',
   TOKEN                varchar(255) not null comment '����TOKEN(Ψһ��ʶ)',
   TIME                 timestamp not null comment '����ʱ��',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);

alter table o_usr comment '��Ʒ���û��󶨱��������ϵͳ���Զ�����';

