package orders
//table
const (
	o_usr = `
/*==============================================================*/
/* Table: o_usr                                                 */
/*==============================================================*/
create table o_usr
(
   ID                   int(11) not null auto_increment,
   USR_ID               int(11) not null comment '用户ID',
   GOODS_ID             int(11) not null comment '商品ID',
   GOODS_DETAIL         varchar(255) comment '商品详情(试卷，答疑)（格式为"，id:name，id:name"）',
   STATUS               int(11) not null default 0 comment '状态（0已购买，10未购买）',
   TIME                 timestamp not null comment '操作时间',
   ADD1                 varchar(255),
   ADD2                 varchar(255),
   primary key (ID)
);
alter table o_usr comment '商品与用户绑定表，存放在子系统，自动生成';

`
//select
	chk_o_usr = ` SELECT * FROM o_usr WHERE  ID<? `

   ins_o_usr = `
   INSERT INTO o_usr (USR_ID, GOODS_ID, STATUS, TOKEN, TIME) VALUES
   (?,?,?,?,?)
   `
   update_o_usr = `
   UPDATE o_usr set STATUS=? WHERE USR_ID=? and GOODS_ID=? and TOKEN=?
   `
   //电商
   selectGoodsInf = `
   SELECT a.Tid,b.PRICE,a.IMGS,a.NAME,a.STORE_ID,c.name as storeName,c.u_id,c.u_name from rcp_resource_new a,rcp_resource_detail b,rcp_store c
   WHERE a.TID = b.R_ID AND a.STORE_ID=c.tid AND b.TID = ?`
)