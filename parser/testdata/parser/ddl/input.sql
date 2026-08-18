CREATE
-- case
CREATE TABLE
-- case
CREATE TABLE foo (
-- case
CREATE TABLE foo ()
-- case
CREATE TABLE foo ();
-- case
CREATE TABLE foo.* (a varchar(50), b int);
-- case
CREATE TABLE foo (a varchar(50), b int);
-- case
CREATE TABLE foo (a TINYINT UNSIGNED);
-- case
CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED)
-- case
CREATE TABLE foo (a bigint unsigned, b bool);
-- case
CREATE TABLE foo (a TINYINT, b SMALLINT) CREATE TABLE bar (x INT, y int64)
-- case
CREATE TABLE foo (a int, b float); CREATE TABLE bar (x double, y float)
-- case
CREATE TABLE foo (a bytes)
-- case
CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED)
-- case
CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED) -- foo
-- case
CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED) // foo
-- case
CREATE TABLE foo (a SMALLINT UNSIGNED, b INT UNSIGNED) /* foo */
-- case
CREATE TABLE foo /* foo */ (a SMALLINT UNSIGNED, b INT UNSIGNED) /* foo */
-- case
CREATE TABLE foo (name CHAR(50) BINARY);
-- case
CREATE TABLE foo (name CHAR(50) COLLATE utf8_bin)
-- case
CREATE TABLE foo (id varchar(50) collate utf8_bin);
-- case
CREATE TABLE foo (name CHAR(50) CHARACTER SET UTF8)
-- case
CREATE TABLE foo (name CHAR(50) CHARACTER SET utf8 BINARY)
-- case
CREATE TABLE foo (name CHAR(50) CHARACTER SET utf8 BINARY CHARACTER set utf8)
-- case
CREATE TABLE foo (name CHAR(50) BINARY CHARACTER SET utf8 COLLATE utf8_bin)
-- case
CREATE TABLE foo (name CHAR(50) CHARACTER SET utf8 COLLATE utf8_bin COLLATE ascii_bin)
-- case
CREATE TABLE foo (name CHAR(50) COLLATE ascii_bin COLLATE latin1_bin)
-- case
CREATE TABLE foo (name CHAR(50) COLLATE ascii_bin PRIMARY KEY COLLATE latin1_bin)
-- case
CREATE TABLE foo (a.b, b);
-- case
CREATE TABLE foo (a, b.c);
-- case
CREATE TABLE (name CHAR(50) BINARY)
-- case
CREATE TABLE foo (name CHAR(50) COLLATE ascii_bin PRIMARY KEY COLLATE latin1_bin, INDEX (name ASC))
-- case
CREATE TABLE foo (name CHAR(50) COLLATE ascii_bin PRIMARY KEY COLLATE latin1_bin, INDEX (name DESC))
-- case
ALTER TABLE tmp CACHE
-- case
ALTER TABLE tmp NOCACHE
-- case
CREATE TEMPORARY TABLE t (a varchar(50), b int);
-- case
CREATE TEMPORARY TABLE t LIKE t1
-- case
DROP TEMPORARY TABLE t
-- case
create global temporary table t (a int, b varchar(255)) on commit delete rows
-- case
create temporary table t (a int, b varchar(255))
-- case
create global temporary table t (a int, b varchar(255))
-- case
create temporary table t (a int, b varchar(255)) on commit delete rows
-- case
create table t (a int, b varchar(255)) on commit delete rows
-- case
create global temporary table t (a int, b varchar(255)) partition by hash(a) partitions 10 on commit delete rows
-- case
create global temporary table t (a int, b varchar(255)) on commit preserve rows
-- case
drop global temporary table t
-- case
drop temporary table t
-- case
CREATE TABLE foo (node_id varchar(50), b int);
-- case
CREATE TABLE foo (node_state varchar(50), b int);
-- case
create table t (c int) avg_row_length = 3
-- case
create table t (c int) avg_row_length 3
-- case
create table t (c int) checksum = 0
-- case
create table t (c int) checksum 1
-- case
create table t (c int) table_checksum = 0
-- case
create table t (c int) table_checksum 1
-- case
create table t (c int) compression = 'NONE'
-- case
create table t (c int) compression 'lz4'
-- case
create table t (c int) connection = 'abc'
-- case
create table t (c int) connection 'abc'
-- case
create table t (c int) key_block_size = 1024
-- case
create table t (c int) key_block_size 1024
-- case
create table t (c int) max_rows = 1000
-- case
create table t (c int) max_rows 1000
-- case
create table t (c int) min_rows = 1000
-- case
create table t (c int) min_rows 1000
-- case
create table t (c int) password = 'abc'
-- case
create table t (c int) password 'abc'
-- case
create table t (c int) DELAY_KEY_WRITE=1
-- case
create table t (c int) DELAY_KEY_WRITE 1
-- case
create table t (c int) ROW_FORMAT = default
-- case
create table t (c int) ROW_FORMAT default
-- case
create table t (c int) ROW_FORMAT = fixed
-- case
create table t (c int) ROW_FORMAT = compressed
-- case
create table t (c int) ROW_FORMAT = compact
-- case
create table t (c int) ROW_FORMAT = redundant
-- case
create table t (c int) ROW_FORMAT = dynamic
-- case
create table t (c int) STATS_PERSISTENT = default
-- case
create table t (c int) STATS_PERSISTENT = 0
-- case
create table t (c int) STATS_PERSISTENT = 1
-- case
create table t (c int) STATS_SAMPLE_PAGES 0
-- case
create table t (c int) STATS_SAMPLE_PAGES 10
-- case
create table t (c int) STATS_SAMPLE_PAGES = 10
-- case
create table t (c int) STATS_SAMPLE_PAGES = default
-- case
create table t (c int) PACK_KEYS = 1
-- case
create table t (c int) PACK_KEYS = 0
-- case
create table t (c int) PACK_KEYS = DEFAULT
-- case
create table t (c int) STORAGE DISK
-- case
create table t (c int) STORAGE MEMORY
-- case
create table t (c int) SECONDARY_ENGINE null
-- case
create table t (c int) SECONDARY_ENGINE = innodb
-- case
create table t (c int) SECONDARY_ENGINE 'null'
-- case
create table testTableCompression (c VARCHAR(15000)) compression="ZLIB";
-- case
create table t1 (c1 int) compression="zlib";
-- case
create table t1 (c1 int) collate=binary;
-- case
create table t1 (c1 int) collate=utf8mb4_0900_as_cs;
-- case
create table t1 (c1 int) default charset=binary collate=binary;
-- case
create table t1 (c1 int) autoextend_size=4M
-- case
ALTER TABLE t_n UNION ( ), KEY_BLOCK_SIZE = 1
-- case
ALTER TABLE d_n.t_n UNION ( t_n ) REMOVE PARTITIONING
-- case
ALTER TABLE d_n.t_n LOCK DEFAULT , UNION = ( t_n , d_n.t_n ) REMOVE PARTITIONING
-- case
ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , MAX_ROWS 10, UNION ( d_n.t_n ) , ROW_FORMAT REDUNDANT, STATS_PERSISTENT = DEFAULT
-- case
create table t (b int) partition by range columns (b) (partition p0 values less than (not 3), partition p2 values less than (20));
-- case
create table t (b int) partition by range columns (b) (partition p0 values less than (1 or 3), partition p2 values less than (20));
-- case
create table t (b int) partition by range columns (b) (partition p0 values less than (3 is null), partition p2 values less than (20));
-- case
create table t (b int) partition by range (b is null) (partition p0 values less than (10));
-- case
create table t (b int) partition by list (not b) (partition p0 values in (10, 20));
-- case
create table t (b int) partition by hash ( not b );
-- case
create table t (b int) partition by range columns (b) (partition p0 values less than (3 in (3, 4, 5)), partition p2 values less than (20));
-- case
CREATE TABLE t (id int) ENGINE = INNDB PARTITION BY RANGE (id) (PARTITION p0 VALUES LESS THAN (10), PARTITION p1 VALUES LESS THAN (20));
-- case
create table t (c int) PARTITION BY HASH (c) PARTITIONS 32;
-- case
create table t (c int) PARTITION BY HASH (Year(VDate)) (PARTITION p1980 VALUES LESS THAN (1980) ENGINE = MyISAM, PARTITION p1990 VALUES LESS THAN (1990) ENGINE = MyISAM, PARTITION pothers VALUES LESS THAN MAXVALUE ENGINE = MyISAM)
-- case
create table t (c int) PARTITION BY RANGE (Year(VDate)) (PARTITION p1980 VALUES LESS THAN (1980) ENGINE = MyISAM, PARTITION p1990 VALUES LESS THAN (1990) ENGINE = MyISAM, PARTITION pothers VALUES LESS THAN MAXVALUE ENGINE = MyISAM)
-- case
create table t (c int, `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '') PARTITION BY RANGE (UNIX_TIMESTAMP(create_time)) (PARTITION p201610 VALUES LESS THAN(1477929600), PARTITION p201611 VALUES LESS THAN(1480521600),PARTITION p201612 VALUES LESS THAN(1483200000),PARTITION p201701 VALUES LESS THAN(1485878400),PARTITION p201702 VALUES LESS THAN(1488297600),PARTITION p201703 VALUES LESS THAN(1490976000))
-- case
CREATE TABLE `md_product_shop` (`shopCode` varchar(4) DEFAULT NULL COMMENT '地点') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 /*!50100 PARTITION BY KEY (shopCode) PARTITIONS 19 */;
-- case
CREATE TABLE `payinfo1` (`id` bigint(20) NOT NULL AUTO_INCREMENT, `oderTime` datetime NOT NULL) ENGINE=InnoDB AUTO_INCREMENT=641533032 DEFAULT CHARSET=utf8 ROW_FORMAT=COMPRESSED KEY_BLOCK_SIZE=8 /*!50500 PARTITION BY RANGE COLUMNS(oderTime) (PARTITION P2011 VALUES LESS THAN ('2012-01-01 00:00:00') ENGINE = InnoDB, PARTITION P1201 VALUES LESS THAN ('2012-02-01 00:00:00') ENGINE = InnoDB, PARTITION PMAX VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB)*/;
-- case
CREATE TABLE app_channel_daily_report (id bigint(20) NOT NULL AUTO_INCREMENT, app_version varchar(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT 'default', gmt_create datetime NOT NULL COMMENT '创建时间', PRIMARY KEY (id)) ENGINE=InnoDB AUTO_INCREMENT=33703438 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci
/*!50100 PARTITION BY RANGE (month(gmt_create)-1)
(PARTITION part0 VALUES LESS THAN (1) COMMENT = '1月份' ENGINE = InnoDB,
 PARTITION part1 VALUES LESS THAN (2) COMMENT = '2月份' ENGINE = InnoDB,
 PARTITION part2 VALUES LESS THAN (3) COMMENT = '3月份' ENGINE = InnoDB,
 PARTITION part3 VALUES LESS THAN (4) COMMENT = '4月份' ENGINE = InnoDB,
 PARTITION part4 VALUES LESS THAN (5) COMMENT = '5月份' ENGINE = InnoDB,
 PARTITION part5 VALUES LESS THAN (6) COMMENT = '6月份' ENGINE = InnoDB,
 PARTITION part6 VALUES LESS THAN (7) COMMENT = '7月份' ENGINE = InnoDB,
 PARTITION part7 VALUES LESS THAN (8) COMMENT = '8月份' ENGINE = InnoDB,
 PARTITION part8 VALUES LESS THAN (9) COMMENT = '9月份' ENGINE = InnoDB,
 PARTITION part9 VALUES LESS THAN (10) COMMENT = '10月份' ENGINE = InnoDB,
 PARTITION part10 VALUES LESS THAN (11) COMMENT = '11月份' ENGINE = InnoDB,
 PARTITION part11 VALUES LESS THAN (12) COMMENT = '12月份' ENGINE = InnoDB) */ ;
-- case
create table t (c int) primary_region="us";
-- case
create table t (c int) regions="us,3";
-- case
create table t (c int) followers="us,3";
-- case
create table t (c int) followers=3;
-- case
create table t (c int) followers=0;
-- case
create table t (c int) voters=3;
-- case
create table t (c int) learners="us,3";
-- case
create table t (c int) learners=3;
-- case
create table t (c int) schedule="even";
-- case
create table t (c int) constraints="ww";
-- case
create table t (c int) leader_constraints="ww";
-- case
create table t (c int) follower_constraints="ww";
-- case
create table t (c int) voter_constraints="ww";
-- case
create table t (c int) learner_constraints="ww";
-- case
create table t (c int) survival_preference="ww";
-- case
create table t (c int) /*T![placement] primary_region="us" */;
-- case
create table t (c int) /*T![placement] regions="us,3" */;
-- case
create table t (c int) /*T![placement] followers="us,3 */";
-- case
create table t (c int) /*T![placement] followers=3 */;
-- case
create table t (c int) /*T![placement] followers=0 */;
-- case
create table t (c int) /*T![placement] voters="us,3" */;
-- case
create table t (c int) /*T![placement] primary_region="us" regions="us,3"  */;
-- case
create table t (c int) placement policy="ww";
-- case
create table t (c int) /*T![placement] placement policy=`x` */;
-- case
create table t (c int) /*T![placement] placement policy="y" */;
-- case
alter table t primary_region="us";
-- case
alter table t regions="us,3";
-- case
alter table t followers=3;
-- case
alter table t followers=0;
-- case
alter table t voters=3;
-- case
alter table t learners=3;
-- case
alter table t schedule="even";
-- case
alter table t constraints="ww";
-- case
alter table t leader_constraints="ww";
-- case
alter table t follower_constraints="ww";
-- case
alter table t voter_constraints="ww";
-- case
alter table t learner_constraints="ww";
-- case
alter table t /*T![placement] primary_region="us" */;
-- case
alter table t placement policy="ww";
-- case
alter table t /*T![placement] placement policy="ww" */;
-- case
alter table t compact;
-- case
alter table t compact tiflash replica;
-- case
alter table t compact partition p1,p2;
-- case
alter table t compact partition p1,p2 tiflash replica;
-- case
create database t primary_region="us";
-- case
create database t regions="us,3";
-- case
create database t followers=3;
-- case
create database t followers=0;
-- case
create database t voters=3;
-- case
create database t learners=3;
-- case
create database t schedule="even";
-- case
create database t constraints="ww";
-- case
create database t leader_constraints="ww";
-- case
create database t follower_constraints="ww";
-- case
create database t voter_constraints="ww";
-- case
create database t learner_constraints="ww";
-- case
create database t /*T![placement] primary_region="us" */;
-- case
create database t placement policy="ww";
-- case
create database t default placement policy="ww";
-- case
create database t /*T![placement] placement policy="ww" */;
-- case
alter database t primary_region="us";
-- case
alter database t regions="us,3";
-- case
alter database t followers=3;
-- case
alter database t followers=0;
-- case
alter database t voters=3;
-- case
alter database t learners=3;
-- case
alter database t schedule="even";
-- case
alter database t constraints="ww";
-- case
alter database t leader_constraints="ww";
-- case
alter database t follower_constraints="ww";
-- case
alter database t voter_constraints="ww";
-- case
alter database t learner_constraints="ww";
-- case
alter database t /*T![placement] primary_region="us" */;
-- case
alter database t placement policy="ww";
-- case
alter database t default placement policy="ww";
-- case
alter database t PLACEMENT POLICY='DEFAULT';
-- case
alter database t PLACEMENT POLICY=DEFAULT;
-- case
alter database t PLACEMENT POLICY = DEFAULT;
-- case
alter database t PLACEMENT POLICY SET DEFAULT
-- case
alter database t PLACEMENT POLICY=`DEFAULT`;
-- case
alter database t /*T![placement] PLACEMENT POLICY='DEFAULT' */;
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) primary_region="us");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) regions="us,3");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) followers=3);
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) voters=3);
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) learners=3);
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) schedule="even");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) constraints="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) leader_constraints="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) follower_constraints="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) voter_constraints="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) learner_constraints="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) placement policy="ww");
-- case
create table m (c int) partition by range (c) (partition p1 values less than (200) /*T![placement] placement policy="ww" */);
-- case
alter table m partition t primary_region="us";
-- case
alter table m partition t regions="us,3";
-- case
alter table m partition t followers=3;
-- case
alter table m partition t primary_region="us" followers=3;
-- case
alter table m partition t voters=3;
-- case
alter table m partition t learners=3;
-- case
alter table m partition t schedule="even";
-- case
alter table m partition t constraints="ww";
-- case
alter table m partition t leader_constraints="ww";
-- case
alter table m partition t follower_constraints="ww";
-- case
alter table m partition t voter_constraints="ww";
-- case
alter table m partition t learner_constraints="ww";
-- case
alter table m partition t placement policy="ww";
-- case
alter table m partition t /*T![placement] placement policy="ww" */;
-- case
alter table m add partition (partition p1 values less than (200) primary_region="us");
-- case
alter table m add partition (partition p1 values less than (200) regions="us,3");
-- case
alter table m add partition (partition p1 values less than (200) followers=3);
-- case
alter table m add partition (partition p1 values less than (200) voters=3);
-- case
alter table m add partition (partition p1 values less than (200) learners=3);
-- case
alter table m add partition (partition p1 values less than (200) schedule="even");
-- case
alter table m add partition (partition p1 values less than (200) constraints="ww");
-- case
alter table m add partition (partition p1 values less than (200) leader_constraints="ww");
-- case
alter table m add partition (partition p1 values less than (200) follower_constraints="ww");
-- case
alter table m add partition (partition p1 values less than (200) voter_constraints="ww");
-- case
alter table m add partition (partition p1 values less than (200) learner_constraints="ww");
-- case
alter table m add partition (partition p1 values less than (200) placement policy="ww");
-- case
alter table m add partition (partition p1 values less than (200) /*T![placement] placement policy="ww" */);
-- case
alter table m add column a int, add partition (partition p1 values less than (200))
-- case
alter table m add partition (partition p1 values less than (200)), add column a int
-- case
create table t (c1 bool, c2 bool, check (c1 in (0, 1)) not enforced, check (c2 in (0, 1)))
-- case
CREATE TABLE Customer (SD integer CHECK (SD > 0), First_Name varchar(30));
-- case
CREATE TABLE Customer (SD integer CHECK (SD > 0) not enforced, SS varchar(30) check(ss='test') enforced);
-- case
CREATE TABLE Customer (SD integer CHECK (SD > 0) not null, First_Name varchar(30) comment 'string' not null);
-- case
CREATE TABLE Customer (SD integer comment 'string' CHECK (SD > 0) not null);
-- case
CREATE TABLE Customer (SD integer comment 'string' not enforced, First_Name varchar(30));
-- case
CREATE TABLE Customer (SD integer not enforced, First_Name varchar(30));
-- case
create database xxx
-- case
create database if exists xxx
-- case
create database if not exists xxx
-- case
create database xxx encryption = 'N'
-- case
create database xxx encryption 'N'
-- case
create database xxx default encryption = 'N'
-- case
create database xxx default encryption 'N'
-- case
create database xxx encryption = 'Y'
-- case
create database xxx encryption 'Y'
-- case
create database xxx default encryption = 'Y'
-- case
create database xxx default encryption 'Y'
-- case
create database xxx encryption = N
-- case
create schema xxx
-- case
create schema if exists xxx
-- case
create schema if not exists xxx
-- case
drop database xxx
-- case
drop database if exists xxx
-- case
drop database if not exists xxx
-- case
drop schema xxx
-- case
drop schema if exists xxx
-- case
drop schema if not exists xxx
-- case
drop table
-- case
drop table if exists t'xyz
-- case
drop table if exists t'
-- case
drop table if exists t`
-- case
drop table if exists t'
-- case
drop table if exists t"
-- case
drop table xxx
-- case
drop table xxx, yyy
-- case
drop tables xxx
-- case
drop tables xxx, yyy
-- case
drop table if exists xxx
-- case
drop table if exists xxx, yyy
-- case
drop table if not exists xxx
-- case
drop table xxx restrict
-- case
drop table xxx, yyy cascade
-- case
drop table if exists xxx restrict
-- case
drop view
-- case
drop view xxx
-- case
drop view xxx, yyy
-- case
drop view if exists xxx
-- case
drop view if exists xxx, yyy
-- case
drop stats t
-- case
drop stats t1, t2, t3
-- case
drop stats t global
-- case
drop stats t partition p0
-- case
drop stats t partition p0, p1, p2
-- case
CREATE TABLE address (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		create_at datetime NOT NULL,
		deleted tinyint(1) NOT NULL,
		update_at datetime NOT NULL,
		version bigint(20) DEFAULT NULL,
		address varchar(128) NOT NULL,
		address_detail varchar(128) NOT NULL,
		cellphone varchar(16) NOT NULL,
		latitude double NOT NULL,
		longitude double NOT NULL,
		name varchar(16) NOT NULL,
		sex tinyint(1) NOT NULL,
		user_id bigint(20) NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT FK_7rod8a71yep5vxasb0ms3osbg FOREIGN KEY (user_id) REFERENCES waimaiqa.user (id),
		INDEX FK_7rod8a71yep5vxasb0ms3osbg (user_id) comment ''
		) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARACTER SET UTF8 COLLATE UTF8_GENERAL_CI ROW_FORMAT=COMPACT COMMENT='' CHECKSUM=0 DELAY_KEY_WRITE=0;
-- case
CREATE TABLE test_data (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		create_at datetime NOT NULL,
		deleted tinyint(1) NOT NULL,
		update_at datetime NOT NULL,
		version bigint(20) DEFAULT NULL,
		address varchar(255) NOT NULL,
		amount decimal(19,2) DEFAULT NULL,
		charge_id varchar(32) DEFAULT NULL,
		paid_amount decimal(19,2) DEFAULT NULL,
		transaction_no varchar(64) DEFAULT NULL,
		wx_mp_app_id varchar(32) DEFAULT NULL,
		contacts varchar(50) DEFAULT NULL,
		deliver_fee decimal(19,2) DEFAULT NULL,
		deliver_info varchar(255) DEFAULT NULL,
		deliver_time varchar(255) DEFAULT NULL,
		description varchar(255) DEFAULT NULL,
		invoice varchar(255) DEFAULT NULL,
		order_from int(11) DEFAULT NULL,
		order_state int(11) NOT NULL,
		packing_fee decimal(19,2) DEFAULT NULL,
		payment_time datetime DEFAULT NULL,
		payment_type int(11) DEFAULT NULL,
		phone varchar(50) NOT NULL,
		store_employee_id bigint(20) DEFAULT NULL,
		store_id bigint(20) NOT NULL,
		user_id bigint(20) NOT NULL,
		payment_mode int(11) NOT NULL,
		current_latitude double NOT NULL,
		current_longitude double NOT NULL,
		address_latitude double NOT NULL,
		address_longitude double NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT food_order_ibfk_1 FOREIGN KEY (user_id) REFERENCES waimaiqa.user (id),
		CONSTRAINT food_order_ibfk_2 FOREIGN KEY (store_id) REFERENCES waimaiqa.store (id),
		CONSTRAINT food_order_ibfk_3 FOREIGN KEY (store_employee_id) REFERENCES waimaiqa.store_employee (id),
		UNIQUE FK_UNIQUE_charge_id USING BTREE (charge_id) comment '',
		INDEX FK_eqst2x1xisn3o0wbrlahnnqq8 USING BTREE (store_employee_id) comment '',
		INDEX FK_8jcmec4kb03f4dod0uqwm54o9 USING BTREE (store_id) comment '',
		INDEX FK_a3t0m9apja9jmrn60uab30pqd USING BTREE (user_id) comment ''
		) ENGINE=InnoDB AUTO_INCREMENT=95 DEFAULT CHARACTER SET utf8 COLLATE UTF8_GENERAL_CI ROW_FORMAT=COMPACT COMMENT='' CHECKSUM=0 DELAY_KEY_WRITE=0;
-- case
create table t (c int KEY);
-- case
CREATE TABLE address (
		id bigint(20) NOT NULL AUTO_INCREMENT,
		create_at datetime NOT NULL,
		deleted tinyint(1) NOT NULL,
		update_at datetime NOT NULL,
		version bigint(20) DEFAULT NULL,
		address varchar(128) NOT NULL,
		address_detail varchar(128) NOT NULL,
		cellphone varchar(16) NOT NULL,
		latitude double NOT NULL,
		longitude double NOT NULL,
		name varchar(16) NOT NULL,
		sex tinyint(1) NOT NULL,
		user_id bigint(20) NOT NULL,
		PRIMARY KEY (id),
		CONSTRAINT FK_7rod8a71yep5vxasb0ms3osbg FOREIGN KEY (user_id) REFERENCES waimaiqa.user (id) ON DELETE CASCADE ON UPDATE NO ACTION,
		INDEX FK_7rod8a71yep5vxasb0ms3osbg (user_id) comment ''
		) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARACTER SET utf8 COLLATE UTF8_GENERAL_CI ROW_FORMAT=COMPACT COMMENT='' CHECKSUM=0 DELAY_KEY_WRITE=0;
-- case
CREATE TABLE address (
id bigint(20) NOT NULL AUTO_INCREMENT,
create_at datetime NOT NULL,
deleted tinyint(1) NOT NULL,
update_at datetime NOT NULL,
version bigint(20) DEFAULT NULL,
address varchar(128) NOT NULL,
address_detail varchar(128) NOT NULL,
cellphone varchar(16) NOT NULL,
latitude double NOT NULL,
longitude double NOT NULL,
name varchar(16) NOT NULL,
sex tinyint(1) NOT NULL,
user_id bigint(20) NOT NULL,
PRIMARY KEY (id),
CONSTRAINT FK_7rod8a71yep5vxasb0ms3osbg FOREIGN KEY (user_id) REFERENCES waimaiqa.user (id) ON DELETE CASCADE ON UPDATE NO ACTION,
INDEX FK_7rod8a71yep5vxasb0ms3osbg (user_id) comment ''
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci ROW_FORMAT=COMPACT COMMENT='' CHECKSUM=0 DELAY_KEY_WRITE=0;
-- case
CREATE TABLE t1 (
		accout_id int(11) DEFAULT '0',
		summoner_id int(11) DEFAULT '0',
		union_name varbinary(52) NOT NULL,
		union_id int(11) DEFAULT '0',
		PRIMARY KEY (union_name)) ENGINE=MyISAM DEFAULT CHARSET=binary;
-- case
CREATE TABLE t (a DECIMAL(20,0), b DECIMAL(30), c FLOAT(25,0))
-- case
create table t (c int, index ci (c) USING BTREE COMMENT "123");
-- case
CREATE TABLE sbtest (id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT, k integer UNSIGNED DEFAULT '0' NOT NULL, c char(120) DEFAULT '' NOT NULL, pad char(60) DEFAULT '' NOT NULL, PRIMARY KEY  (id) )
-- case
create table test (create_date TIMESTAMP NOT NULL COMMENT '创建日期 create date' DEFAULT now());
-- case
create table ts (t int, v timestamp(3) default CURRENT_TIMESTAMP(3));
-- case
create table if not exists `t` (`id` int not null auto_increment comment '消息ID', primary key `pk_id` (`id`) );
-- case
create table a like b
-- case
create table a (id int REFERENCES a (id) ON delete NO ACTION )
-- case
create table a (id int REFERENCES a (id) ON update set default )
-- case
create table a (id int REFERENCES a (id) ON delete set null on update CASCADE)
-- case
create table a (id int REFERENCES a (id) ON update set default on delete RESTRICT)
-- case
create table a (id int REFERENCES a (id) MATCH FULL ON delete NO ACTION )
-- case
create table a (id int REFERENCES a (id) MATCH PARTIAL ON update NO ACTION )
-- case
create table a (id int REFERENCES a (id) MATCH SIMPLE ON update NO ACTION )
-- case
create table a (id int REFERENCES a (id) ON update set default )
-- case
create table a (id int REFERENCES a (id) ON update set default on update CURRENT_TIMESTAMP)
-- case
create table a (id int REFERENCES a (id) ON delete set default on update CURRENT_TIMESTAMP)
-- case
create table a (like b)
-- case
create table if not exists a like b
-- case
create table if not exists a (like b)
-- case
create table if not exists a like (b)
-- case
create table a (t int) like b
-- case
create table a (t int) like (b)
-- case
create table a select * from b
-- case
create table a as select * from b
-- case
create table a (m int, n datetime) as select * from b
-- case
create table a (unique(n)) as select n from b
-- case
create table a ignore as select n from b
-- case
create table a replace as select n from b
-- case
create table a (m int) replace as (select n as m from b union select n+1 as m from c group by 1 limit 2)
-- case
create table a
-- case
create table t (a timestamp default now)
-- case
create table t (a timestamp default now())
-- case
create table t (a timestamp default (((now()))))
-- case
create table t (a timestamp default now() on update now)
-- case
create table t (a timestamp default now() on update now())
-- case
create table t (a timestamp default now() on update (now()))
-- case
CREATE TABLE t (c TEXT) default CHARACTER SET utf8, default COLLATE utf8_general_ci;
-- case
CREATE TABLE t (c TEXT) shard_row_id_bits = 1;
-- case
CREATE TABLE t (c TEXT) shard_row_id_bits = 1, PRE_SPLIT_REGIONS = 1;
-- case
CREATE TABLE IF NOT EXISTS `general_log` (`event_time` timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),`user_host` mediumtext NOT NULL,`thread_id` bigint(20) unsigned NOT NULL,`server_id` int(10) unsigned NOT NULL,`command_type` varchar(64) NOT NULL,`argument` mediumblob NOT NULL) ENGINE=CSV DEFAULT CHARSET=utf8 COMMENT='General log'
-- case
CREATE TABLE followers ( f1 int NOT NULL REFERENCES user_profiles (uid) );
-- case
create table t (a int default rand())
-- case
create table t (a int default rand(1))
-- case
create table t (a int default (rand()))
-- case
create table t (a int default (rand(1)))
-- case
create table t (a int default (((rand()))))
-- case
create table t (a int default (((rand(1)))))
-- case
create table t (d date default current_date())
-- case
create table t (d date default current_date)
-- case
create table t (d date default (current_date()))
-- case
create table t (d date default (curdate()))
-- case
create table t (d date default curdate())
-- case
create table t (d date default current_date())
-- case
create table t (d date default date_format(now(),'%Y-%m'))
-- case
create table t (d date default (date_format(now(),'%Y-%m')))
-- case
create table t (d date default date_format(now(),'%Y-%m-%d'))
-- case
create table t (d date default date_format(now(),'%Y-%m-%d %H.%i.%s'))
-- case
create table t (d date default date_format(now(),'%Y-%m-%d %H:%i:%s'))
-- case
create table t (d date default date_format(now(),'%b %d %Y %h:%i %p'))
-- case
create table t (a varchar(32) default (replace(upper(uuid()), '-', '')))
-- case
create table t (a varchar(32) default replace(upper(uuid()), '-', ''))
-- case
create table t (a varchar(32) default (replace(convert(upper(uuid()) using utf8mb4), '-', '')))
-- case
create table t (a varchar(32) default replace(convert(upper(uuid()) using utf8mb4), '-', ''))
-- case
create table t (a int default upper(substring_index(user(),'@',1)))
-- case
create table t (a int default (upper(substring_index(user(),'@',1))))
-- case
create table t (a varchar(32) default (str_to_date('1980-01-01','%Y-%m-%d')))
-- case
create table t (a varchar(32) default str_to_date('1980-01-01','%Y-%m-%d'))
-- case
create table t (j json default (json_object()))
-- case
create table t (j json default (json_array()))
-- case
create table t (j json default (json_quote()))
-- case
create table t (j json default (json_object('foo', 5, 'bar', 'barfoo')))
-- case
create table t (j json default (json_array(1,2,3)))
-- case
create table t (j json default (json_quote('foobar')))
-- case
create table t (c char(33) default (nonexistingfunc('foobar')))
-- case
create table t (c char(33) default 'foobar')
-- case
create table t (c char(33) default ('foobar'))
-- case
create table t (i int default (0))
-- case
create table t (i int default (-1))
-- case
create table t (i int default (+1))
-- case
create table t (a int, b int, c char(33) default (b))
-- case
create table t (a int, b int, c char(33) default `b`)
-- case
create table t (a int) encryption = 'n';
-- case
create table t (a int) encryption 'n';
-- case
alter table t encryption = 'y';
-- case
alter table t encryption 'y';
-- case
create table t (a int key global)
-- case
create table t (a int key local)
-- case
create table t (a int primary key local)
-- case
create table t (a int primary key global)
-- case
create table t (a int UNIQUE local)
-- case
create table t (a int UNIQUE global)
-- case
create table t (a int UNIQUE key local)
-- case
create table t (a int UNIQUE key global)
-- case
alter table t add index (a)
-- case
alter table t add index (a) local
-- case
alter table t add index (a) global
-- case
alter table t add unique (a)
-- case
alter table t add unique (a) local
-- case
alter table t add unique (a) global
-- case
alter table t add unique key (a) global
-- case
alter table t add unique key (a)
-- case
alter table t add unique key (a) local
-- case
alter table t add primary key (a) global
-- case
alter table t add primary key (a)
-- case
alter table t add primary key (a) local
-- case
create index i on t (a)
-- case
create index i on t (a) local
-- case
create index i on t (a) global
-- case
ALTER DATABASE t CHARACTER SET = 'utf8'
-- case
ALTER DATABASE CHARACTER SET = 'utf8'
-- case
ALTER DATABASE t DEFAULT CHARACTER SET = 'utf8'
-- case
ALTER SCHEMA t DEFAULT CHARACTER SET = 'utf8'
-- case
ALTER SCHEMA DEFAULT CHARACTER SET = 'utf8'
-- case
ALTER SCHEMA t DEFAULT CHARSET = 'UTF8'
-- case
ALTER DATABASE t COLLATE = binary
-- case
ALTER DATABASE t CHARSET=binary COLLATE = binary
-- case
ALTER DATABASE t COLLATE = 'utf8_bin'
-- case
ALTER DATABASE COLLATE = 'utf8_bin'
-- case
ALTER DATABASE t DEFAULT COLLATE = 'utf8_bin'
-- case
ALTER SCHEMA t DEFAULT COLLATE = 'UTF8_BiN'
-- case
ALTER SCHEMA DEFAULT COLLATE = 'UTF8_BiN'
-- case
ALTER SCHEMA `` DEFAULT COLLATE = 'UTF8_BiN'
-- case
ALTER DATABASE t CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
-- case
ALTER DATABASE t DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
-- case
ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' COLLATE = 'utf8_bin'
-- case
ALTER DATABASE DEFAULT CHARSET = 'utf8mb4' DEFAULT COLLATE = 'utf8mb4_general_ci' CHARACTER SET = 'utf8' COLLATE = 'utf8mb4_bin'
-- case
ALTER TABLE t ADD COLUMN (a SMALLINT UNSIGNED)
-- case
ALTER TABLE t.* ADD COLUMN (a SMALLINT UNSIGNED)
-- case
ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED)
-- case
ALTER TABLE ADD COLUMN (a SMALLINT UNSIGNED)
-- case
ALTER TABLE t ADD COLUMN (a SMALLINT UNSIGNED, b varchar(255))
-- case
ALTER TABLE t ADD COLUMN IF NOT EXISTS (a SMALLINT UNSIGNED, b varchar(255))
-- case
ALTER TABLE t ADD COLUMN (a SMALLINT UNSIGNED FIRST)
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED FIRST
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED AFTER b
-- case
ALTER TABLE t ADD COLUMN IF NOT EXISTS a SMALLINT UNSIGNED AFTER b
-- case
ALTER TABLE employees ADD PARTITION
-- case
ALTER TABLE employees ADD PARTITION ( PARTITION P1 VALUES LESS THAN (2010))
-- case
ALTER TABLE employees ADD PARTITION ( PARTITION P2 VALUES LESS THAN MAXVALUE)
-- case
ALTER TABLE employees ADD PARTITION IF NOT EXISTS ( PARTITION P2 VALUES LESS THAN MAXVALUE)
-- case
ALTER TABLE employees ADD PARTITION IF NOT EXISTS PARTITIONS 5
-- case
ALTER TABLE employees ADD PARTITION (
				PARTITION P1 VALUES LESS THAN (2010),
				PARTITION P2 VALUES LESS THAN (2015),
				PARTITION P3 VALUES LESS THAN MAXVALUE)
-- case
alter table t add partition (partition x values in ((3, 4), (5, 6)))
-- case
ALTER TABLE employees ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE employees ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10
-- case
ALTER TABLE employees ADD PARTITION LOCAL
-- case
ALTER TABLE employees ADD PARTITION LOCAL PARTITIONS 10
-- case
ALTER TABLE t_n REBUILD PARTITION ALL
-- case
ALTER TABLE d_n.t_n REBUILD PARTITION LOCAL ALL
-- case
ALTER TABLE t_n REBUILD PARTITION LOCAL ident
-- case
ALTER TABLE t_n REBUILD PARTITION NO_WRITE_TO_BINLOG ident , ident
-- case
ALTER TABLE t_n REBUILD PARTITION LOCAL
-- case
ALTER TABLE t_n REBUILD PARTITION LOCAL local
-- case
ALTER TABLE t_n REBUILD PARTITION LOCAL local, local
-- case
alter table t drop partition p1;
-- case
alter table t drop partition p2;
-- case
alter table t drop partition if exists p2;
-- case
alter table t drop partition p1, p2;
-- case
alter table t drop partition if exists p1, p2;
-- case
alter table t check partition all;
-- case
alter table t check partition p;
-- case
alter table t check partition p1, p2;
-- case
alter table employees add partition partitions 1;
-- case
alter table employees add partition partitions 2;
-- case
alter table clients coalesce partition 3;
-- case
alter table clients coalesce partition 4;
-- case
alter table clients coalesce partition no_write_to_binlog 4;
-- case
alter table clients coalesce partition local 4;
-- case
ALTER TABLE t DISABLE KEYS
-- case
ALTER TABLE t ENABLE KEYS
-- case
ALTER TABLE t MODIFY COLUMN a varchar(255)
-- case
ALTER TABLE t MODIFY COLUMN IF EXISTS a varchar(255)
-- case
ALTER TABLE t CHANGE COLUMN a b varchar(255)
-- case
ALTER TABLE t CHANGE COLUMN IF EXISTS a b varchar(255)
-- case
ALTER TABLE t CHANGE COLUMN a b varchar(255) CHARACTER SET UTF8 BINARY
-- case
ALTER TABLE t CHANGE COLUMN a b varchar(255) FIRST
-- case
ALTER TABLE db.t RENAME to db1.t1
-- case
ALTER TABLE db.t RENAME db1.t1
-- case
ALTER TABLE db.t RENAME = db1.t1
-- case
ALTER TABLE db.t RENAME as db1.t1
-- case
ALTER TABLE t RENAME to t1
-- case
ALTER TABLE t RENAME t1
-- case
ALTER TABLE t RENAME = t1
-- case
ALTER TABLE t RENAME as t1
-- case
ALTER TABLE t_n ORDER BY ident
-- case
ALTER TABLE t_n ORDER BY ident ASC
-- case
ALTER TABLE t_n ORDER BY ident DESC
-- case
ALTER TABLE t_n ORDER BY ident1, ident2
-- case
ALTER TABLE t_n ORDER BY ident1 ASC, ident2
-- case
ALTER TABLE t_n ORDER BY ident1 ASC, ident2 ASC
-- case
ALTER TABLE t_n ORDER BY ident1 ASC, ident2 DESC
-- case
ALTER TABLE t_n ORDER BY ident1 DESC, ident2
-- case
ALTER TABLE t_n ORDER BY ident1 DESC, ident2 ASC
-- case
ALTER TABLE t_n ORDER BY ident1 DESC, ident2 DESC
-- case
ALTER TABLE t_n ORDER BY ident1, ident2, ident3
-- case
ALTER TABLE t_n ORDER BY ident1, ident2, ident3 ASC
-- case
ALTER TABLE t_n ORDER BY ident1, ident2, ident3 DESC
-- case
ALTER TABLE t_n ORDER BY ident1 ASC, ident2 ASC, ident3 ASC
-- case
ALTER TABLE t_n ORDER BY ident1 DESC, ident2 DESC, ident3 DESC
-- case
ALTER TABLE t RENAME COLUMN a TO b
-- case
ALTER TABLE t RENAME COLUMN t.a TO t.b
-- case
ALTER TABLE t RENAME COLUMN a TO t.b
-- case
ALTER TABLE t RENAME COLUMN t.a TO b
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT 1
-- case
ALTER TABLE t ALTER a SET DEFAULT 1
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT CURRENT_TIMESTAMP
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT NOW()
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT 1+1
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT (CURRENT_TIMESTAMP())
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT (NOW())
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT (1+1)
-- case
ALTER TABLE t ALTER COLUMN a SET DEFAULT (1)
-- case
ALTER TABLE t ALTER COLUMN a DROP DEFAULT
-- case
ALTER TABLE t ALTER a DROP DEFAULT
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock=none
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock=default
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock=shared
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock=exclusive
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock none
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock default
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock shared
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, lock exclusive
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, LOCK=NONE
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, LOCK=DEFAULT
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, LOCK=SHARED
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, LOCK=EXCLUSIVE
-- case
ALTER TABLE t ADD FULLTEXT KEY `FullText` (`name` ASC)
-- case
ALTER TABLE t ADD FULLTEXT `FullText` (`name` ASC)
-- case
ALTER TABLE t ADD FULLTEXT INDEX `FullText` (`name` ASC)
-- case
ALTER TABLE t ADD INDEX (a) USING BTREE COMMENT 'a'
-- case
ALTER TABLE t ADD INDEX IF NOT EXISTS (a) USING BTREE COMMENT 'a'
-- case
ALTER TABLE t ADD INDEX (a) USING RTREE COMMENT 'a'
-- case
ALTER TABLE t ADD INDEX (a) PRE_SPLIT_REGIONS = 4
-- case
ALTER TABLE t ADD INDEX (a) PRE_SPLIT_REGIONS 4
-- case
ALTER TABLE t ADD INDEX (a) PRE_SPLIT_REGIONS = 'a'
-- case
ALTER TABLE t ADD PRIMARY KEY (a) CLUSTERED PRE_SPLIT_REGIONS = 4
-- case
ALTER TABLE t ADD PRIMARY KEY (a) PRE_SPLIT_REGIONS = 4 NONCLUSTERED
-- case
ALTER TABLE t ADD INDEX (a) PRE_SPLIT_REGIONS = (between (1, 'a') and (2, 'b') regions 4);
-- case
ALTER TABLE t ADD INDEX (a) PRE_SPLIT_REGIONS = (by (1, 'a'), (2, 'b'), (3, 'c'));
-- case
ALTER TABLE t ADD INDEX (a) comment 'a' PRE_SPLIT_REGIONS = (between (1, 'a') and (2, 'b') regions 4);
-- case
CREATE INDEX idx ON t (a, b) pre_split_regions = 100
-- case
CREATE INDEX idx ON t (a, b) PRE_SPLIT_REGIONS = (between (1, 'a') and (2, 'b') regions 4);
-- case
ALTER TABLE t ADD INDEX idx(a) pre_split_regions = 100, ADD INDEX idx2(b) pre_split_regions = (by(1),(2),(3))
-- case
ALTER TABLE t ADD KEY (a) USING HASH COMMENT 'a'
-- case
ALTER TABLE t ADD INDEX (a) USING BTREE /*T![global_index] GLOBAL */ COMMENT 'a'
-- case
ALTER TABLE t ADD UNIQUE INDEX (a) /*T![global_index] GLOBAL */
-- case
ALTER TABLE t ADD UNIQUE INDEX (a) LOCAL
-- case
ALTER TABLE t ADD KEY IF NOT EXISTS (a) USING HASH COMMENT 'a'
-- case
ALTER TABLE t ADD PRIMARY KEY ident USING RTREE ( a DESC , b   )
-- case
ALTER TABLE t ADD KEY USING RTREE   ( a ) 
-- case
ALTER TABLE t ADD KEY USING RTREE ( ident ASC , ident ( 123 ) )
-- case
ALTER TABLE t ADD PRIMARY KEY (a) COMMENT 'a'
-- case
ALTER TABLE t ADD UNIQUE (a) COMMENT 'a'
-- case
ALTER TABLE t ADD UNIQUE KEY (a) COMMENT 'a'
-- case
ALTER TABLE t ADD UNIQUE INDEX (a) COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR (a) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR (VEC_COSINE_DISTANCE(a)) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR ((VEC_COSINE_DISTANCE(a))) USING HASH COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR KEY ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((lower(a))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a), a)) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX (a, (VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) USING HYPO COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX ((VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS ((VEC_COSINE_DISTANCE(a))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD VECTOR INDEX IF NOT EXISTS ((VEC_COSINE_DISTANCE(a))) ADD_COLUMNAR_REPLICA_ON_DEMAND USING HNSW COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR (a) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR ((a - 1)) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR (a) USING HASH COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR (a, b) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR KEY (a, b) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR INDEX (a, b) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR INDEX (a) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR INDEX (a) USING HYPO COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR INDEX ((a - 1)) USING HYPO COMMENT 'a'
-- case
ALTER TABLE t ADD COLUMNAR INDEX IF NOT EXISTS (a) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE t ADD CONSTRAINT fk_t2_id FOREIGN KEY (t2_id) REFERENCES t(id)
-- case
ALTER TABLE t ADD CONSTRAINT fk_t2_id FOREIGN KEY IF NOT EXISTS (t2_id) REFERENCES t(id)
-- case
ALTER TABLE t ADD CONSTRAINT c_1 CHECK (1+1) NOT ENFORCED, ADD UNIQUE (a)
-- case
ALTER TABLE t ADD CONSTRAINT c_1 CHECK (1+1) ENFORCED, ADD UNIQUE (a)
-- case
ALTER TABLE t ADD CONSTRAINT c_1 CHECK (1+1), ADD UNIQUE (a)
-- case
ALTER TABLE t ENGINE ''
-- case
ALTER TABLE t ENGINE = ''
-- case
ALTER TABLE t ENGINE = 'innodb'
-- case
ALTER TABLE t ENGINE = innodb
-- case
ALTER TABLE `db`.`t` ENGINE = ``
-- case
ALTER TABLE t INSERT_METHOD = FIRST
-- case
ALTER TABLE t INSERT_METHOD LAST
-- case
ALTER TABLE t ADD COLUMN a SMALLINT UNSIGNED, ADD COLUMN a SMALLINT
-- case
ALTER TABLE t ADD COLUMN a SMALLINT, ENGINE = '', default COLLATE = UTF8_GENERAL_CI
-- case
ALTER TABLE t ENGINE = '', COMMENT='', default COLLATE = UTF8_GENERAL_CI
-- case
ALTER TABLE t ENGINE = '', ADD COLUMN a SMALLINT
-- case
ALTER TABLE t default COLLATE = UTF8_GENERAL_CI, ENGINE = '', ADD COLUMN a SMALLINT
-- case
ALTER TABLE t shard_row_id_bits = 1
-- case
ALTER TABLE t AUTO_INCREMENT 3
-- case
ALTER TABLE t AUTO_INCREMENT = 3
-- case
ALTER TABLE t FORCE AUTO_INCREMENT 3
-- case
ALTER TABLE t FORCE AUTO_INCREMENT = 3
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` mediumtext CHARACTER SET UTF8MB4 COLLATE UTF8MB4_UNICODE_CI NOT NULL , ALGORITHM = DEFAULT;
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` mediumtext CHARACTER SET UTF8MB4 COLLATE UTF8MB4_UNICODE_CI NOT NULL , ALGORITHM = INPLACE;
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` mediumtext CHARACTER SET UTF8MB4 COLLATE UTF8MB4_UNICODE_CI NOT NULL , ALGORITHM = COPY;
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` MEDIUMTEXT CHARACTER SET UTF8MB4 COLLATE UTF8MB4_UNICODE_CI NOT NULL, ALGORITHM = INSTANT;
-- case
ALTER TABLE t CONVERT TO CHARACTER SET UTF8;
-- case
ALTER TABLE t CONVERT TO CHARSET UTF8;
-- case
ALTER TABLE t CONVERT TO CHARACTER SET UTF8 COLLATE UTF8_BIN;
-- case
ALTER TABLE t CONVERT TO CHARSET UTF8 COLLATE UTF8_BIN;
-- case
alter table d_n.t_n convert to character set default
-- case
alter table d_n.t_n convert to charset default
-- case
alter table d_n.t_n convert to char set default
-- case
alter table d_n.t_n convert to character set default collate utf8mb4_0900_ai_ci
-- case
ALTER TABLE t FORCE
-- case
ALTER TABLE t DROP INDEX;
-- case
ALTER TABLE t DROP INDEX a
-- case
ALTER TABLE t DROP INDEX IF EXISTS a
-- case
ALTER TABLE t ALTER INDEX a INVISIBLE
-- case
ALTER TABLE t ALTER INDEX a VISIBLE
-- case
ALTER TABLE t DROP FOREIGN KEY a
-- case
ALTER TABLE t DROP COLUMN a CASCADE
-- case
ALTER TABLE t DROP COLUMN IF EXISTS a CASCADE
-- case
ALTER TABLE testTableCompression COMPRESSION="LZ4";
-- case
ALTER TABLE t1 COMPRESSION="zlib";
-- case
ALTER TABLE t1
-- case
ALTER TABLE t1 ,
-- case
ALTER TABLE t RENAME KEY a TO b;
-- case
ALTER TABLE t RENAME INDEX a TO b;
-- case
ALTER TABLE d_n.t_n DROP CHECK ident;
-- case
ALTER TABLE t_n LOCK = DEFAULT , DROP CHECK ident;
-- case
ALTER TABLE t_n ALTER CHECK ident ENFORCED;
-- case
ALTER TABLE t_n ALTER CHECK ident NOT ENFORCED;
-- case
ALTER TABLE t_n DROP CONSTRAINT ident
-- case
ALTER TABLE t_n DROP CHECK ident
-- case
ALTER TABLE t_n ALTER CONSTRAINT ident
-- case
ALTER TABLE t_n ALTER CONSTRAINT ident enforced
-- case
ALTER TABLE t_n ALTER CHECK ident not enforced
-- case
alter table t analyze partition a
-- case
alter table t analyze partition a with 4 buckets
-- case
alter table t analyze partition a index b
-- case
alter table t analyze partition a index b with 4 buckets
-- case
alter table t partition by hash(a)
-- case
alter table t add column a int partition by hash(a)
-- case
alter table t add column a int partition by hash(a) update indexes (idx_a global)
-- case
alter table t add column a int partition by hash(a) update indexes (idx_a global, idx_b local)
-- case
alter table t add column a int partition by hash(a) update indexes (idx_a normal)
-- case
alter table t add column a int partition by hash(a) update indexes (global)
-- case
alter table t partition by range(a)
-- case
alter table t partition by range(a) update indexes (a local)
-- case
alter table t partition by range(a) (partition x values less than (75))
-- case
alter table t add column a int, partition by range(a) (partition x values less than (75))
-- case
alter table t comment 'cmt' partition by hash(a)
-- case
alter table t enable keys, comment = 'cmt' partition by hash(a)
-- case
alter table t enable keys, comment = 'cmt', partition by hash(a)
-- case
alter table t partition by hash(a) enable keys
-- case
alter table t partition by hash(a), enable keys
-- case
alter table t partition by range FIELDS(a) (partition x values less than maxvalue)
-- case
alter table t partition by list FIELDS(a) (PARTITION p0 VALUES IN (5, 10, 15))
-- case
alter table t partition by range FIELDS(a,b,c) (partition p1 values less than (1,1,1));
-- case
alter table t partition by list FIELDS(a,b,c) (PARTITION p0 VALUES IN ((5, 10, 15)))
-- case
alter table t with validation, add column b int as (a + 1)
-- case
alter table t without validation, add column b int as (a + 1)
-- case
alter table t without validation, with validation, add column b int as (a + 1)
-- case
alter table t with validation, modify column b int as (a + 2) 
-- case
alter table t with validation, change column b c int as (a + 2)
-- case
ALTER TABLE d_n.t_n ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE d_n.t_n ADD PARTITION LOCAL
-- case
alter table t with validation, exchange partition p with table nt without validation;
-- case
alter table t exchange partition p with table nt with validation;
-- case
alter table t reorganize partition;
-- case
alter table t reorganize partition local;
-- case
alter table t reorganize partition no_write_to_binlog;
-- case
ALTER TABLE members REORGANIZE PARTITION n0 INTO (PARTITION s0 VALUES LESS THAN (1960), PARTITION s1 VALUES LESS THAN (1970));
-- case
ALTER TABLE members REORGANIZE PARTITION LOCAL n0 INTO (PARTITION s0 VALUES LESS THAN (1960), PARTITION s1 VALUES LESS THAN (1970));
-- case
ALTER TABLE members REORGANIZE PARTITION p1,p2,p3 INTO ( PARTITION s0 VALUES LESS THAN (1960), PARTITION s1 VALUES LESS THAN (1970));
-- case
alter table t reorganize partition remove partition;
-- case
alter table t reorganize partition no_write_to_binlog remove into (partition p0 VALUES LESS THAN (1991));
-- case
ALTER TABLE t ATTRIBUTES='str'
-- case
ALTER TABLE t ATTRIBUTES='str1,str2'
-- case
ALTER TABLE t ATTRIBUTES="str1,str2"
-- case
ALTER TABLE t ATTRIBUTES 'str1,str2'
-- case
ALTER TABLE t ATTRIBUTES "str1,str2"
-- case
ALTER TABLE t ATTRIBUTES=DEFAULT
-- case
ALTER TABLE t ATTRIBUTES=default
-- case
ALTER TABLE t ATTRIBUTES=DeFaUlT
-- case
ALTER TABLE t ATTRIBUTES
-- case
ALTER TABLE t PARTITION p ATTRIBUTES='str'
-- case
ALTER TABLE t PARTITION p ATTRIBUTES='str1,str2'
-- case
ALTER TABLE t PARTITION p ATTRIBUTES="str1,str2"
-- case
ALTER TABLE t PARTITION p ATTRIBUTES 'str1,str2'
-- case
ALTER TABLE t PARTITION p ATTRIBUTES "str1,str2"
-- case
ALTER TABLE t PARTITION p ATTRIBUTES=DEFAULT
-- case
ALTER TABLE t PARTITION p ATTRIBUTES=default
-- case
ALTER TABLE t PARTITION p ATTRIBUTES=DeFaUlT
-- case
ALTER TABLE t PARTITION p ATTRIBUTES
-- case
CREATE TABLE t1 (attributes int);
-- case
CREATE INDEX idx ON t (a)
-- case
CREATE INDEX IF NOT EXISTS idx ON t (a)
-- case
CREATE UNIQUE INDEX idx ON t (a)
-- case
CREATE UNIQUE INDEX IF NOT EXISTS idx ON t (a)
-- case
CREATE UNIQUE INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE BTREE
-- case
CREATE UNIQUE INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE HASH
-- case
CREATE UNIQUE INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE RTREE
-- case
CREATE UNIQUE INDEX ident TYPE BTREE ON d_n.t_n ( ident , ident ASC )
-- case
CREATE UNIQUE INDEX ident USING BTREE ON d_n.t_n ( ident , ident ASC )
-- case
CREATE SPATIAL INDEX idx ON t (a)
-- case
CREATE SPATIAL INDEX IF NOT EXISTS idx ON t (a)
-- case
CREATE FULLTEXT INDEX idx ON t (a)
-- case
CREATE FULLTEXT INDEX IF NOT EXISTS idx ON t (a)
-- case
CREATE FULLTEXT INDEX idx ON t (a) WITH PARSER ident
-- case
CREATE FULLTEXT INDEX idx ON t (a) WITH PARSER ident comment 'string'
-- case
CREATE FULLTEXT INDEX idx ON t (a) comment 'string' with parser ident
-- case
CREATE FULLTEXT INDEX idx ON t (a) WITH PARSER ident comment 'string' lock default
-- case
CREATE INDEX idx ON t (a) USING HASH
-- case
CREATE INDEX idx ON t (a) COMMENT 'foo'
-- case
CREATE INDEX idx ON t (a) USING HASH COMMENT 'foo'
-- case
CREATE INDEX idx ON t (a) LOCK=NONE
-- case
CREATE INDEX idx USING BTREE ON t (a) USING HASH COMMENT 'foo'
-- case
CREATE INDEX idx USING BTREE ON t (a)
-- case
CREATE INDEX idx ON t ( a ) VISIBLE
-- case
CREATE INDEX idx ON t ( a ) INVISIBLE
-- case
CREATE INDEX idx ON t ( a ) INVISIBLE VISIBLE
-- case
CREATE INDEX idx ON t ( a ) VISIBLE INVISIBLE
-- case
CREATE INDEX idx ON t ( a ) USING HASH VISIBLE
-- case
CREATE INDEX idx ON t ( a ) USING HASH INVISIBLE
-- case
CREATE VECTOR INDEX idx ON t (a) USING HNSW 
-- case
CREATE VECTOR INDEX idx ON t (a, b) USING HNSW 
-- case
CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)))
-- case
CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE BTREE
-- case
CREATE VECTOR INDEX idx ON t USING HNSW ((VEC_COSINE_DISTANCE(a)))
-- case
CREATE VECTOR idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
-- case
CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a)), a) USING HNSW
-- case
CREATE VECTOR INDEX idx ON t (a, (VEC_COSINE_DISTANCE(a))) USING HNSW
-- case
CREATE VECTOR KEY idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
-- case
CREATE VECTOR INDEX idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
-- case
CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW
-- case
CREATE VECTOR INDEX IF NOT EXISTS idx ON t ((VEC_COSINE_DISTANCE(a))) TYPE HNSW
-- case
CREATE VECTOR INDEX ident TYPE HNSW ON d_n.t_n ((VEC_COSINE_DISTANCE(a)))
-- case
CREATE VECTOR INDEX idx USING HNSW ON t ((VEC_COSINE_DISTANCE(a)))
-- case
CREATE VECTOR INDEX ident ON d_n.t_n ( ident , ident ASC ) TYPE HNSW
-- case
CREATE UNIQUE INDEX ident USING HNSW ON d_n.t_n ( ident , ident ASC )
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = DEFAULT
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM DEFAULT
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = INPLACE
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM INPLACE
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = COPY
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM COPY
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = DEFAULT LOCK = DEFAULT
-- case
CREATE INDEX idx ON t ( a ) LOCK = DEFAULT ALGORITHM = DEFAULT
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = INPLACE LOCK = EXCLUSIVE
-- case
CREATE INDEX idx ON t ( a ) LOCK = EXCLUSIVE ALGORITHM = INPLACE
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM = ident
-- case
CREATE INDEX idx ON t ( a ) ALGORITHM ident
-- case
drop index a on t
-- case
drop index a on db.t
-- case
drop index a on db.`tb-ttb`
-- case
drop index if exists a on t
-- case
drop index if exists a on db.t
-- case
drop index if exists a on db.`tb-ttb`
-- case
drop index idx on t algorithm = default
-- case
drop index idx on t algorithm default
-- case
drop index idx on t algorithm = inplace
-- case
drop index idx on t algorithm inplace
-- case
drop index idx on t lock = default
-- case
drop index idx on t lock default
-- case
drop index idx on t lock = shared
-- case
drop index idx on t lock shared
-- case
drop index idx on t algorithm = default lock = default
-- case
drop index idx on t lock = default algorithm = default
-- case
drop index idx on t algorithm = inplace lock = exclusive
-- case
drop index idx on t lock = exclusive algorithm = inplace
-- case
drop index idx on t algorithm = algorithm_type
-- case
drop index idx on t algorithm algorithm_type
-- case
drop index idx on t lock = lock_type
-- case
drop index idx on t lock lock_type
-- case
RENAME TABLE t TO t1
-- case
RENAME TABLE t t1
-- case
RENAME TABLE d.t TO d1.t1
-- case
RENAME TABLE t1 TO t2, t3 TO t4
-- case
TRUNCATE TABLE t1
-- case
TRUNCATE t1
-- case
ALTER TABLE t ADD INDEX () 
-- case
ALTER TABLE t ADD UNIQUE ()
-- case
ALTER TABLE t ADD UNIQUE INDEX ()
-- case
ALTER TABLE t ADD UNIQUE KEY ()
-- case
ALTER TABLE d_n.t_n SECONDARY_LOAD
-- case
ALTER TABLE d_n.t_n SECONDARY_UNLOAD
-- case
ALTER TABLE t_n LOCK = DEFAULT , SECONDARY_LOAD
-- case
ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , SECONDARY_LOAD
-- case
ALTER TABLE d_n.t_n ALGORITHM = DEFAULT , SECONDARY_UNLOAD
-- case
create table a (process double)
-- case
create table t (a int1, b int2, c int3, d int4, e int8)
-- case
create table t (lv long varchar null)
-- case
CREATE TABLE cdp_test.`test2-1` (id int(11) DEFAULT NULL,key(id));
-- case
CREATE TABLE miantiao (`扁豆焖面`       INT(11));
-- case
CREATE TABLE bar (m INT)  SELECT n FROM foo;
-- case
CREATE TABLE bar (m INT) IGNORE SELECT n FROM foo;
-- case
CREATE TABLE bar (m INT) REPLACE SELECT n FROM foo;
-- case
create table t (a timestamp, b timestamp as (a) not null on update current_timestamp);
-- case
create table t (a bigint, b bigint as (a) primary key auto_increment);
-- case
create table t (a bigint, b bigint as (a) not null default 10);
-- case
create table t (a bigint, b bigint as (a+1) not null);
-- case
create table t (a bigint, b bigint as (a+1) not null);
-- case
create table t (a bigint, b bigint as (a+1) not null comment 'ttt');
-- case
create table t(a int, index idx((cast(a as binary(1)))));
-- case
alter table t add column (f timestamp as (a+1) default '2019-01-01 11:11:11');
-- case
alter table t modify column f int as (a+1) default 55;
-- case
create table t (a int column_format fixed)
-- case
create table t (a int column_format default)
-- case
create table t (a int column_format dynamic)
-- case
alter table t modify column a bigint column_format default
-- case
recover table by job 11
-- case
recover table by job 11,12,13
-- case
recover table by job
-- case
recover table by job 0
-- case
recover table t1
-- case
recover table t1,t2
-- case
recover table 
-- case
recover table t1 100
-- case
recover table t1 abc
-- case
flashback table t
-- case
flashback table t TO t1
-- case
flashback table t TO timestamp
-- case
flashback database db1
-- case
flashback schema db1
-- case
flashback database db1 to db2
-- case
flashback schema db1 to db2
-- case
flashback cluster to timestamp '2021-05-26 16:45:26'
-- case
flashback table t to timestamp '2021-05-26 16:45:26'
-- case
flashback table t,t1 to timestamp '2021-05-26 16:45:26'
-- case
flashback database test to timestamp '2021-05-26 16:45:26'
-- case
flashback schema test to timestamp '2021-05-26 16:45:26'
-- case
flashback cluster to timestamp TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())
-- case
flashback cluster to timestamp DATE_SUB(NOW(), INTERVAL 3 SECOND)
-- case
flashback table to timestamp '2021-05-26 16:45:26'
-- case
flashback database to timestamp '2021-05-26 16:45:26'
-- case
flashback cluster to tso 445494955052105721
-- case
flashback table t to tso 445494955052105722
-- case
flashback table t,t1 to tso 445494955052105723
-- case
flashback database test to tso 445494955052105724
-- case
flashback schema test to tso 445494955052105725
-- case
flashback table to tso 445494955052105726
-- case
flashback database to tso 445494955052105727
-- case
flashback schema test to tso 0
-- case
flashback schema test to tso -100
-- case
alter table t remove partitioning
-- case
alter table db.ident remove partitioning
-- case
alter table t lock = default remove partitioning
-- case
alter table t add column a int remove partitioning
-- case
alter table t add column a int, add index (c) remove partitioning
-- case
alter table t add column a int, remove partitioning
-- case
alter table t add column a int, add index (c), remove partitioning
-- case
alter table t remove partitioning add column a int
-- case
alter table t remove partitioning, add column a int
-- case
alter table t add column a double (4,2) zerofill references b match full on update set null first
-- case
alter table d_n.t_n add constraint foreign key ident (ident(1)) references d_n.t_n match full on delete set null
-- case
alter table t_n add constraint ident foreign key (ident,ident(1)) references t_n match full on update set null on delete restrict
-- case
alter table d_n.t_n add foreign key ident (ident, ident(1) asc) references t_n match partial on delete cascade remove partitioning
-- case
alter table d_n.t_n add constraint foreign key (ident asc) references d_n.t_n match simple on update cascade on delete cascade
-- case
create table t (a character varying(1));
-- case
create table t (a character varying(255));
-- case
create table t (a char varying(50));
-- case
create table t (a varcharacter(1));
-- case
create table t (a varcharacter(50));
-- case
create table t (a varcharacter(1), b varcharacter(255));
-- case
create table t (a char);
-- case
create table t (a character);
-- case
create table t (a character varying(50), b int);
-- case
create table t (a character, b int);
-- case
create table t (a national character varying(50));
-- case
create table t (a national char varying(50));
-- case
create table t (a national char);
-- case
create table t (a national character);
-- case
create table t (a nchar);
-- case
create table t (a nchar varchar(50));
-- case
create table t (a nchar varcharacter(50));
-- case
create table t (a national varchar);
-- case
create table t (a national varchar(50));
-- case
create table t (a national varcharacter(50));
-- case
create table t (a nchar varying(50));
-- case
create table t (a nvarchar(50));
-- case
create table nchar (a int);
-- case
create table nchar (a int, b nchar);
-- case
create table nchar (a int, b nchar(50));
-- case
alter table t_n storage disk , modify ident national varcharacter(12) column_format fixed first;
-- case
create table t (a serial);
-- case
create table t (a serial null);
-- case
create table t (b int, a serial);
-- case
create table t (a int serial default value);
-- case
create table t (a int serial default value null);
-- case
create table t (a bigint serial default value);
-- case
create table t (a smallint serial default value);
-- case
create table t (a long);
-- case
create table t (a long varchar);
-- case
create table t (a long varcharacter);
-- case
create table t (a long char varying);
-- case
create table t (a long character varying);
-- case
create table t (a mediumtext, b long varchar, c long, d long varcharacter, e long char varying, f long character varying, g long);
-- case
create table t (a long varbinary);
-- case
create table t (a long char varying, b long varbinary);
-- case
create table t (a long char set utf8);
-- case
create table t (a long char varying char set utf8);
-- case
create table t (a long character set utf8);
-- case
create table t (a long character varying character set utf8);
-- case
alter table d_n.t_n modify column ident long after ident remove partitioning
-- case
alter table d_n.t_n modify column ident long char varying after ident remove partitioning
-- case
alter table d_n.t_n modify column ident long character varying after ident remove partitioning
-- case
alter table d_n.t_n modify column ident long varchar after ident remove partitioning
-- case
alter table d_n.t_n modify column ident long varcharacter after ident remove partitioning
-- case
alter table t_n change column ident ident long char varying binary charset utf8 first , tablespace ident
-- case
alter table t_n change column ident ident long character varying binary charset utf8 first , tablespace ident
-- case
create table t (a int) stats_auto_recalc 2;
-- case
create table t (a int) stats_auto_recalc = 10;
-- case
create table t (a int) stats_auto_recalc 0;
-- case
create table t (a int) stats_auto_recalc default;
-- case
create table t (a int) stats_auto_recalc = 0;
-- case
create table t (a int) stats_auto_recalc = 1;
-- case
create table t (a int) stats_auto_recalc=default;
-- case
create table t (a int) stats_persistent = 1, stats_auto_recalc = 1;
-- case
create table t (a int) stats_auto_recalc = 1, stats_sample_pages = 25;
-- case
alter table t modify a bigint, ENGINE=InnoDB, stats_auto_recalc = 0
-- case
create table stats_auto_recalc (a int);
-- case
create table stats_auto_recalc (a int) stats_auto_recalc=1;
-- case
create table t (a int, primary key type type btree (a));
-- case
create table t (a int, primary key type btree (a));
-- case
create table t (a int, primary key using btree (a));
-- case
create table t (a int, primary key (a) type btree);
-- case
create table t (a int, primary key (a) using btree);
-- case
create table t (a int, unique index type type btree (a));
-- case
create table t (a int, unique index type using btree (a));
-- case
create table t (a int, unique index type btree (a));
-- case
create table t (a int, unique index using btree (a));
-- case
create table t (a int, unique index (a) using btree);
-- case
create table t (a int, unique key (a) using btree);
-- case
create table t (a int, index type type btree (a));
-- case
create table t (a int, index type btree (a));
-- case
create table t (a int, index type using btree (a));
-- case
create table t (a int, index using btree (a));
-- case
ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES LESS THAN ( MAXVALUE ) STORAGE ENGINE text_string MAX_ROWS 12 )
-- case
ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION ident VALUES LESS THAN MAXVALUE STORAGE ENGINE = text_string, PARTITION ident VALUES LESS THAN ( MAXVALUE ) (SUBPARTITION text_string MIN_ROWS 11))
-- case
ALTER TABLE d_n.t_n WITHOUT VALIDATION , ADD PARTITION ( PARTITION ident VALUES IN ( DEFAULT ) STORAGE ENGINE text_string MAX_ROWS 12 )
-- case
ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG ( PARTITION ident VALUES IN ( DEFAULT ) STORAGE ENGINE text_string MAX_ROWS 12 )
-- case
ALTER TABLE d_n.t_n ADD PARTITION ( PARTITION ident VALUES IN ( DEFAULT ), partition ptext values in ('default') )
-- case
ALTER TABLE d_n.t_n WITH VALIDATION , ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION ident VALUES LESS THAN MAXVALUE STORAGE ENGINE = text_string, PARTITION ident VALUES IN ( DEFAULT ) (SUBPARTITION text_string MIN_ROWS 11))
-- case
ALTER TABLE d_n.t_n ADD PARTITION (PARTITION ident VALUES IN ( DEFAULT ))
-- case
ALTER TABLE d_n.t_n ADD PARTITION (PARTITION ident VALUES IN (1, default ))
-- case
ALTER TABLE t IMPORT TABLESPACE;
-- case
ALTER TABLE t DISCARD TABLESPACE;
-- case
ALTER TABLE db.t IMPORT TABLESPACE;
-- case
ALTER TABLE db.t DISCARD TABLESPACE;
-- case
ALTER TABLE t ADD ( CHECK ( true ) )
-- case
ALTER TABLE t ADD ( CONSTRAINT CHECK ( true ) )
-- case
ALTER TABLE t ADD COLUMN ( CONSTRAINT ident CHECK ( 1>2 ) NOT ENFORCED )
-- case
alter table t add column (b int, constraint c unique key (b))
-- case
ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) )
-- case
ALTER TABLE t ADD COLUMN ( CONSTRAINT CHECK ( true ) ENFORCED , CHECK ( true ) )
-- case
ALTER TABLE t ADD COLUMN (a1 int, CONSTRAINT b1 CHECK (a1>0))
-- case
ALTER TABLE t ADD COLUMN (a1 int, a2 int, CONSTRAINT b1 CHECK (a1>0), CONSTRAINT b2 CHECK (a2<10))
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY (`a1`))
-- case
ALTER TABLE t ADD (a1 int, CONSTRAINT PRIMARY KEY (a1))
-- case
ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), UNIQUE (a2))
-- case
ALTER TABLE t ADD (a1 int, a2 int, PRIMARY KEY (a1), CONSTRAINT b2 UNIQUE (a2))
-- case
ALTER TABLE ident ADD ( CONSTRAINT FOREIGN KEY ident ( EXECUTE ( 123 ) ) REFERENCES t ( a ) MATCH SIMPLE ON DELETE CASCADE ON UPDATE SET NULL )
-- case
ALTER TABLE t ADD COLUMN a DATE CHECK ( a > 0 ) FIRST
-- case
ALTER TABLE t ADD a1 int CONSTRAINT ident CHECK ( a1 > 1 ) REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE;
-- case
ALTER TABLE t ADD COLUMN a DATE CONSTRAINT CHECK ( a > 0 ) FIRST
-- case
ALTER TABLE t ADD a TINYBLOB CONSTRAINT ident CHECK ( 1>2 ) REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE
-- case
ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) ENFORCED
-- case
ALTER TABLE t ADD a2 int CONSTRAINT ident CHECK (a2 > 1) NOT ENFORCED
-- case
ALTER TABLE t ADD a2 int CONSTRAINT ident primary key REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE;
-- case
ALTER TABLE t ADD a2 int CONSTRAINT ident primary key (a2))
-- case
ALTER TABLE t ADD a2 int CONSTRAINT ident unique key (a2))
-- case
ALTER TABLE t SET TIFLASH REPLICA 2 LOCATION LABELS 'a','b'
-- case
ALTER TABLE t SET TIFLASH REPLICA 0
-- case
ALTER DATABASE t SET TIFLASH REPLICA 2 LOCATION LABELS 'a','b'
-- case
ALTER DATABASE t SET TIFLASH REPLICA 0
-- case
ALTER DATABASE t SET TIFLASH REPLICA 1 SET TIFLASH REPLICA 2 LOCATION LABELS 'a','b'
-- case
ALTER DATABASE t SET TIFLASH REPLICA 1 SET TIFLASH REPLICA 2
-- case
ALTER DATABASE t SET TIFLASH REPLICA 1 LOCATION LABELS 'a','b' SET TIFLASH REPLICA 2
-- case
ALTER DATABASE t SET TIFLASH REPLICA 1 LOCATION LABELS 'a','b' SET TIFLASH REPLICA 2 LOCATION LABELS 'a', 'b'
-- case
CREATE TABLE IF NOT EXISTS table_ident (a SQL_TSI_YEAR(4), b SQL_TSI_YEAR);
-- case
CREATE TABLE IF NOT EXISTS table_ident (ident1 BOOL COMMENT "text_string" unique, ident2 SQL_TSI_YEAR(4) ZEROFILL);
-- case
create table t (y sql_tsi_year(4), y1 sql_tsi_year)
-- case
create table t (y sql_tsi_year(4) unsigned zerofill zerofill, y1 sql_tsi_year signed unsigned zerofill)
-- case
insert into t set a = default
-- case
insert into t set a := default
-- case
replace t set a = default
-- case
update t set a = default
-- case
insert into t set a = default on duplicate key update a = default
-- case
create table t (a text byte ascii)
-- case
create table t (a text byte charset latin1)
-- case
create table t (a longtext ascii)
-- case
create table t (a mediumtext ascii)
-- case
create table t (a tinytext ascii)
-- case
create table t (a text byte)
-- case
create table t (a long byte, b text ascii)
-- case
create table t (a text ascii, b mediumtext ascii, c int)
-- case
create table t (a int, b text ascii, c mediumtext ascii)
-- case
create table t (a long ascii, b long ascii)
-- case
create table t (a long character set utf8mb4, b long charset utf8mb4, c long char set utf8mb4)
-- case
create table t (a int STORAGE MEMORY, b varchar(255) STORAGE MEMORY)
-- case
create table t (a int storage DISK, b varchar(255) STORAGE DEFAULT)
-- case
create table t (a int STORAGE DEFAULT, b varchar(255) STORAGE DISK)
-- case
create table t (a fixed(6, 3), b fixed key)
-- case
create table t (a numeric, b fixed(6))
-- case
create table t (a fixed(65, 30) zerofill, b numeric, c fixed(65) unsigned zerofill)
-- case
create table a(a int, key(lower(a)));
-- case
create table a(a int, key(a+1));
-- case
create table a(a int, key(a, a+1));
-- case
create table a(a int, b int, key((a+1), (b+1)));
-- case
create table a(a int, b int, key(a, (b+1)));
-- case
create table a(a int, b int, key((a+1), b));
-- case
create table a(a int, b int, key((a + 1) desc));
-- case
create sequence sequence
-- case
create sequence seq
-- case
create sequence if not exists seq
-- case
create sequence seq
-- case
create sequence seq
-- case
create sequence if not exists seq
-- case
create sequence if not exists seq
-- case
create sequence if not exists seq increment
-- case
create sequence if not exists seq increment 1
-- case
create sequence if not exists seq increment = 1
-- case
create sequence if not exists seq increment by 1
-- case
create sequence if not exists seq minvalue
-- case
create sequence if not exists seq minvalue 1
-- case
create sequence if not exists seq minvalue = 1
-- case
create sequence if not exists seq no
-- case
create sequence if not exists seq nominvalue
-- case
create sequence if not exists seq no minvalue
-- case
create sequence if not exists seq maxvalue
-- case
create sequence if not exists seq maxvalue 1
-- case
create sequence if not exists seq maxvalue = 1
-- case
create sequence if not exists seq no
-- case
create sequence if not exists seq nomaxvalue
-- case
create sequence if not exists seq no maxvalue
-- case
create sequence if not exists seq start
-- case
create sequence if not exists seq start with
-- case
create sequence if not exists seq start =
-- case
create sequence if not exists seq start with
-- case
create sequence if not exists seq start 1
-- case
create sequence if not exists seq start = 1
-- case
create sequence if not exists seq start with 1
-- case
create sequence if not exists seq cache
-- case
create sequence if not exists seq cache 1
-- case
create sequence if not exists seq cache = 1
-- case
create sequence if not exists seq nocache
-- case
create sequence if not exists seq no cache
-- case
create sequence if not exists seq cycle
-- case
create sequence if not exists seq nocycle
-- case
create sequence if not exists seq no cycle
-- case
create sequence seq increment 1 start with 0 minvalue 0 maxvalue 1000
-- case
create sequence seq increment 1 start with 0 minvalue 0 maxvalue 1000
-- case
create sequence seq increment 10 start with 0 minvalue 0 maxvalue 1000
-- case
create sequence if not exists seq cache 1 increment 1 start with -1 minvalue 0 maxvalue 1000
-- case
create sequence sEq start with 0 minvalue 0 maxvalue 1000
-- case
create sequence if not exists seq increment 1 start with 0 minvalue -2 maxvalue 1000
-- case
create sequence seq increment -1 start with -1 minvalue -1 maxvalue -1000 cache = 10 nocycle
-- case
create table sequence (a int)
-- case
create table t (sequence int)
-- case
drop sequence
-- case
drop sequence seq
-- case
drop sequence if exists seq
-- case
drop sequence seq
-- case
drop sequence if exists seq
-- case
drop sequence if exists seq, seq2, seq3
-- case
drop sequence seq seq2
-- case
drop sequence seq, seq2
-- case
create table t (a bigint auto_random(3) primary key, b varchar(255))
-- case
create table t (a bigint auto_random primary key, b varchar(255))
-- case
create table t (a bigint primary key auto_random(4), b varchar(255))
-- case
create table t (a bigint primary key auto_random(3) primary key unique, b varchar(255))
-- case
create table t (a bigint auto_random(5, 53) primary key, b varchar(255))
-- case
create table t (a bigint auto_random(15, 32) primary key, b varchar(255))
-- case
create table t (a int) auto_id_cache=1
-- case
create table t (a int auto_increment key) auto_id_cache 10
-- case
create table t (a bigint, b varchar(255)) auto_id_cache 50
-- case
create table t (a bigint auto_random(3) primary key) auto_random_base = 10
-- case
create table t (a bigint primary key auto_random(4), b varchar(100)) auto_random_base 200
-- case
alter table t auto_random_base = 50
-- case
alter table t auto_increment 30, auto_random_base 40
-- case
alter table t force auto_random_base = 50
-- case
alter table t auto_increment 30, force auto_random_base 40
-- case
alter sequence seq
-- case
alter sequence seq comment="haha"
-- case
alter sequence seq start = 1
-- case
alter sequence seq start with 1 increment by 1
-- case
alter sequence seq start with 1 increment by 2 minvalue 0 maxvalue 100
-- case
alter sequence seq increment -1 start with -1 minvalue -1 maxvalue -1000 cache = 10 nocycle
-- case
alter sequence if exists seq2 increment = 2
-- case
alter sequence seq restart
-- case
alter sequence seq start with 3 restart with 5
-- case
alter sequence seq restart = 5
-- case
create sequence seq restart = 5
-- case
create table t (a int, index ``(a))
-- case
create table t (a int, b varchar(255), primary key(b, a) clustered)
-- case
create table t (a int, b varchar(255), primary key(b, a) nonclustered)
-- case
create table t (a int primary key nonclustered, b varchar(255))
-- case
create table t (a int, b varchar(255) primary key clustered)
-- case
create table t (a int, b varchar(255) default 'a' primary key clustered)
-- case
create table t (a int, b varchar(255) primary key nonclustered, primary key(b, a) nonclustered)
-- case
create table t (a int, b varchar(255), primary key(b, a) using RTREE nonclustered)
-- case
create table t (a int, b varchar(255), primary key(b, a) using RTREE clustered nonclustered)
-- case
create table t (a int, b varchar(255), primary key(b, a) using RTREE nonclustered clustered)
-- case
create table t (a int, b varchar(255) clustered primary key)
-- case
create table t (a int, b varchar(255) primary key nonclustered clustered)
-- case
alter table t add primary key (`a`, `b`) clustered
-- case
alter table t add primary key (`a`, `b`) nonclustered
-- case
create table t(a int, b vector(3), vector index(b) USING HNSW);
-- case
create table t(a int, b vector(3), vector index(a, b) USING HNSW);
-- case
create table t(a int, b vector(3), vector index((VEC_COSINE_DISTANCE(b))));
-- case
create table t(a int, b vector(3), vector index((VEC_COSINE_DISTANCE(b))) USING HASH);
-- case
create table t(a int, b vector(3), vector index(a, (VEC_COSINE_DISTANCE(b))) USING HNSW);
-- case
create table t(a int, b vector(3), vector index((VEC_COSINE_DISTANCE(b)), a) USING HNSW);
-- case
create table t(a int, b vector(3), vector index(VEC_COSINE_DISTANCE(b)) USING HNSW);
-- case
create table t(a int, b vector(3), vector key((VEC_COSINE_DISTANCE(b))) TYPE HNSW);
-- case
create table t(a int, b vector(3), vector index((b+1)) USING HNSW);
-- case
create table t(a int, b vector(3), vector index((VEC_COSINE_DISTANCE(a, b))) USING HNSW);
-- case
create table t(a int, b vector(3), vector index((VEC_COSINE_DISTANCE(b))) USING HNSW);
-- case
drop placement policy x
-- case
drop placement policy x, y
-- case
drop placement policy if exists x
-- case
drop placement policy if exists x, y
-- case
show create placement policy x
-- case
show create placement policy if exists x
-- case
show create placement policy x, y
-- case
show create placement policy `placement`
-- case
create placement policy x primary_region='us'
-- case
create placement policy x region='us, 3'
-- case
create placement policy x followers=3
-- case
create placement policy x followers=0
-- case
create placement policy x voters=3
-- case
create placement policy x learners=3
-- case
create placement policy x schedule='even'
-- case
create placement policy x constraints='ww'
-- case
create placement policy x leader_constraints='ww'
-- case
create placement policy x follower_constraints='ww'
-- case
create placement policy x voter_constraints='ww'
-- case
create placement policy x learner_constraints='ww'
-- case
create placement policy x primary_region='cn' regions='us' schedule='even'
-- case
create placement policy x primary_region='cn', leader_constraints='ww', leader_constraints='yy'
-- case
create placement policy if not exists x regions = 'us', follower_constraints='yy'
-- case
create or replace placement policy x regions='us'
-- case
create placement policy x placement policy y
-- case
create table masking (id int)
-- case
select masking from t
-- case
alter table t add masking int
-- case
alter table t drop masking
-- case
alter table t modify masking int
-- case
create masking policy p on t(c) as c
-- case
create masking policy p on t(c) as c enable
-- case
create masking policy if not exists p on t(c) as c disable
-- case
create or replace masking policy p on t(c) as c
-- case
create masking policy p on t(c) as c restrict on none
-- case
create or replace masking policy if not exists p on t(c) as c
-- case
create masking policy p on t(c) as case when current_user() = 'root' then c else 'xxx' end enable
-- case
create masking policy p on t(c) as c restrict on (insert_into_select, delete_select) enable
-- case
create masking policy p on t(c) as case when current_user() not in ('root@%', 'u@%') then 'x' else c end
-- case
alter table t add masking policy p on (c) as c
-- case
alter table t add masking policy p on (c) as c restrict on none
-- case
alter table t add masking policy p on (c) as c disable
-- case
alter table t add masking policy p on (c) as c restrict on (update_select, ctas) disable
-- case
alter table t modify masking policy p set expression = case when current_role() in ('r1', 'r2') then c else 'x' end
-- case
alter table t modify masking policy p set expression case when current_role() in ('r1', 'r2') then c else 'x' end
-- case
alter table t modify masking policy p set restrict on (insert_into_select, update_select, delete_select, ctas)
-- case
alter table t modify masking policy p set restrict on none
-- case
alter table t enable masking policy p
-- case
alter table t disable masking policy p
-- case
alter table t drop masking policy p
-- case
show masking policies for t
-- case
show masking policies for t where column_name = 'c'
-- case
alter placement policy x primary_region='us'
-- case
alter placement policy x region='us, 3'
-- case
alter placement policy x followers=3
-- case
alter placement policy x voters=3
-- case
alter placement policy x learners=3
-- case
alter placement policy x schedule='even'
-- case
alter placement policy x constraints='ww'
-- case
alter placement policy x leader_constraints='ww'
-- case
alter placement policy x follower_constraints='ww'
-- case
alter placement policy x voter_constraints='ww'
-- case
alter placement policy x learner_constraints='ww'
-- case
alter placement policy x primary_region='cn' regions='us' schedule='even'
-- case
alter placement policy x primary_region='cn', leader_constraints='ww', leader_constraints='yy'
-- case
alter placement policy if exists x regions = 'us', follower_constraints='yy'
-- case
alter placement policy x placement policy y
-- case
create resource group x cpu ='8c'
-- case
create resource group x region ='us, 3'
-- case
create resource group x cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'
-- case
create resource group x ru_per_sec=2000
-- case
create resource group x ru_per_sec=200000
-- case
create resource group x ru_per_sec=UNLIMITED
-- case
create resource group x ru_per_sec=unlimited
-- case
create resource group x ru_per_sec='check'
-- case
create resource group x followers=0
-- case
create resource group x burstable=true
-- case
create resource group x burstable=false
-- case
create resource group x burstable=disable
-- case
create resource group x ru_per_sec=1000, burstable
-- case
create resource group x burstable, ru_per_sec=2000
-- case
create resource group x ru_per_sec=3000 burstable
-- case
create resource group x burstable ru_per_sec=4000
-- case
create resource group x BURSTABLE = UNLIMITED ru_per_sec=4000
-- case
create resource group x BURSTABLE = MODERATED ru_per_sec=4000
-- case
create resource group x BURSTABLE = OFF ru_per_sec=4000
-- case
create resource group x ru_per_sec=20, priority=LOW, burstable
-- case
create resource group default ru_per_sec=20, priority=LOW, burstable
-- case
create resource group default ru_per_sec=UNLIMITED, priority=LOW, burstable
-- case
create resource group x ru_per_sec=1000
-- case
create resource group x ru_per_sec=1000 burstable=unlimited
-- case
create resource group x ru_per_sec=1000 burstable=off
-- case
create resource group x ru_per_sec=1000 burstable=moderated
-- case
create resource group x burstable=unlimited, ru_per_sec=2000
-- case
create resource group x burstable=off, ru_per_sec=2000
-- case
create resource group x burstable=moderated, ru_per_sec=2000
-- case
create resource group x ru_per_sec=1000 ,burstable=unlimited
-- case
create resource group x ru_per_sec=1000 ,burstable=off
-- case
create resource group x ru_per_sec=1000 ,burstable=moderated
-- case
create resource group x ru_per_sec=1000 , priority=LOW,burstable=unlimited
-- case
create resource group x ru_per_sec=1000 , priority=LOW,burstable=off
-- case
create resource group x ru_per_sec=1000 , priority=LOW,burstable=moderated
-- case
create resource group x ru_per_sec=UNLIMITED , priority=LOW,burstable=unlimited
-- case
create resource group x ru_per_sec=UNLIMITED , priority=LOW,burstable=off
-- case
create resource group x ru_per_sec=UNLIMITED , priority=LOW,burstable=moderated
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION DRYRUN)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10m' ACTION COOLDOWN)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=(ACTION KILL EXEC_ELAPSED='10m')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' WATCH=SIMILAR DURATION '10m' ACTION COOLDOWN)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED "10s" ACTION COOLDOWN WATCH EXACT DURATION='10m')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED '9s' ACTION COOLDOWN WATCH EXACT)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED '8s' ACTION COOLDOWN WATCH EXACT DURATION = UNLIMITED)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED '7s' ACTION COOLDOWN WATCH EXACT DURATION UNLIMITED)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED '7s' ACTION COOLDOWN WATCH EXACT DURATION 'UNLIMITED')
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' RU 100 ACTION DRYRUN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' PROCESSED_KEYS 100 ACTION DRYRUN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' WATCH SIMILAR DURATION '10m' ACTION COOLDOWN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION COOLDOWN WATCH EXACT DURATION '10m')
-- case
create resource group x ru_per_sec=1000 background = (task_types='')
-- case
create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT=50)
-- case
create resource group x ru_per_sec=1000 background = (UTILIZATION_LIMIT="NAN")
-- case
create resource group x ru_per_sec=1000 background (task_types='br,lightning')
-- case
create resource group x ru_per_sec=1000 background (task_types='br,lightning',utilization_limit=50)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED "10s" ACTION COOLDOWN WATCH EXACT DURATION='10m')  background (task_types 'br,lightning')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED "10s" ACTION COOLDOWN WATCH PLAN DURATION='10m')  background (task_types 'br,lightning')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT (EXEC_ELAPSED "10s" ACTION COOLDOWN WATCH PLAN DURATION='10m')  background (task_types 'br,lightning', utilization_limit 10)
-- case
create resource group x ru_per_sec=UNLIMITED QUERY_LIMIT (EXEC_ELAPSED "10s" ACTION COOLDOWN WATCH PLAN DURATION='10m')  background (task_types 'br,lightning')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s')
-- case
create resource group x ru_per_sec=1000 QUERY=(EXEC_ELAPSED '10s')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=EXEC_ELAPSED '10s'
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s'
-- case
create resource group x ru_per_sec=1000 LIMIT=(EXEC_ELAPSED '10s')
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN ACTION KILL)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS=100)
-- case
create resource group x ru_per_sec=1000 QUERY=(PROCESSED_KEYS 100)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=PROCESSED_KEYS 100
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 100
-- case
create resource group x ru_per_sec=1000 LIMIT=(PROCESSED_KEYS 100)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (PROCESSED_KEYS 100 ACTION DRYRUN ACTION KILL)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU=100)
-- case
create resource group x ru_per_sec=1000 QUERY=(RU 100)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT=RU 100
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 100
-- case
create resource group x ru_per_sec=1000 LIMIT=(RU 100)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (RU 100 ACTION DRYRUN ACTION KILL)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s' PROCESSED_KEYS=100)
-- case
create resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED='10s', PROCESSED_KEYS=100, RU=100)
-- case
alter resource group x cpu ='8c'
-- case
alter resource group x region ='us, 3'
-- case
alter resource group x burstable=true
-- case
alter resource group x burstable=false
-- case
alter resource group x burstable=disable
-- case
alter resource group default priority = high
-- case
alter resource group x cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'
-- case
alter resource group x ru_per_sec=1000
-- case
alter resource group x ru_per_sec=2000, BURSTABLE
-- case
alter resource group x ru_per_sec=UNLIMITED
-- case
alter resource group x ru_per_sec=UNLIMITED, BURSTABLE
-- case
alter resource group x ru_per_sec=unlimited, BURSTABLE
-- case
alter resource group x ru_per_sec='check', BURSTABLE
-- case
alter resource group x BURSTABLE, ru_per_sec=3000
-- case
alter resource group x BURSTABLE ru_per_sec=4000
-- case
alter resource group x ru_per_sec=2000, BURSTABLE=unlimited
-- case
alter resource group x ru_per_sec=2000, BURSTABLE=moderated
-- case
alter resource group x ru_per_sec=2000, BURSTABLE=off
-- case
alter resource group x ru_per_sec=UNLIMITED
-- case
alter resource group x ru_per_sec=UNLIMITED, BURSTABLE=unlimited
-- case
alter resource group x ru_per_sec=UNLIMITED, BURSTABLE=moderated
-- case
alter resource group x ru_per_sec=unlimited, BURSTABLE=off
-- case
alter resource group x ru_per_sec='check', BURSTABLE
-- case
alter resource group x ru_per_sec=2000, BURSTABLE=yes
-- case
alter resource group x BURSTABLE=unlimited, ru_per_sec=3000
-- case
alter resource group x BURSTABLE=moderated, ru_per_sec=3000
-- case
alter resource group x BURSTABLE=off, ru_per_sec=3000
-- case
alter resource group x BURSTABLE=unlimited ru_per_sec=4000
-- case
alter resource group x BURSTABLE=moderated ru_per_sec=4000
-- case
alter resource group x BURSTABLE=off ru_per_sec=4000
-- case
alter resource group x BURSTABLE
-- case
alter resource group x ru_per_sec=200000 BURSTABLE
-- case
alter resource group x followers=0
-- case
alter resource group x ru_per_sec=20 priority=MID BURSTABLE
-- case
alter resource group x BURSTABLE=unlimited
-- case
alter resource group x BURSTABLE=moderated
-- case
alter resource group x BURSTABLE=off
-- case
alter resource group x ru_per_sec=200000 BURSTABLE=unlimited
-- case
alter resource group x ru_per_sec=200000 BURSTABLE=moderated
-- case
alter resource group x ru_per_sec=200000 BURSTABLE=off
-- case
alter resource group x followers=0
-- case
alter resource group x ru_per_sec=20 priority=MID
-- case
alter resource group x ru_per_sec=20 priority=HIGH BURSTABLE=unlimited
-- case
alter resource group x ru_per_sec=20 priority=HIGH BURSTABLE=moderated
-- case
alter resource group x ru_per_sec=20 priority=HIGH BURSTABLE=off
-- case
alter resource group x QUERY_LIMIT=NULL
-- case
alter resource group x QUERY_LIMIT=()
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION DRYRUN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=()
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10m' ACTION COOLDOWN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=( ACTION KILL EXEC_ELAPSED '10m')
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' WATCH SIMILAR DURATION '10m' ACTION COOLDOWN)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION COOLDOWN WATCH EXACT DURATION '10m')
-- case
alter resource group x ru_per_sec=UNLIMITED QUERY_LIMIT=(EXEC_ELAPSED '10s' ACTION COOLDOWN WATCH EXACT DURATION '10m')
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s')
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT EXEC_ELAPSED '10s'
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN ACTION KILL)
-- case
alter resource group x ru_per_sec=1000 QUERY_LIMIT = (EXEC_ELAPSED '10s' ACTION DRYRUN WATCH SIMILAR DURATION '10m' ACTION COOLDOWN)
-- case
alter resource group x background=()
-- case
alter resource group x background NULL
-- case
alter resource group default priority=low background = ( task_types "ttl" )
-- case
alter resource group default burstable=unlimited background ( task_types = 'a,b,c' )
-- case
alter resource group default burstable=moderated background ( task_types = 'a,b,c' )
-- case
alter resource group default burstable=off background ( task_types = 'a,b,c' )
-- case
alter resource group default burstable=unlimited background ( utilization_limit = 20 )
-- case
alter resource group default burstable=moderated background ( utilization_limit = 20 )
-- case
alter resource group default burstable=off background ( utilization_limit = 20 )
-- case
alter resource group default burstable=unlimited background ( task_types = 'a,b,c', utilization_limit = 20 )
-- case
alter resource group default burstable=moderated background ( task_types = 'a,b,c', utilization_limit = 20 )
-- case
alter resource group default burstable=off background ( task_types = 'a,b,c', utilization_limit = 20 )
-- case
alter resource group default burstable=unlimited background ( utilization_limit = 'abc' )
-- case
alter resource group default burstable=moderated background ( utilization_limit = 'abc' )
-- case
alter resource group default burstable=off background ( utilization_limit = 'abc' )
-- case
drop resource group x;
-- case
drop resource group DEFAULT;
-- case
drop resource group if exists x;
-- case
drop resource group x,y
-- case
drop resource group if exists x,y
-- case
set resource group x;
-- case
set resource group ``;
-- case
set resource group `default`;
-- case
set resource group default;
-- case
set resource group x y
-- case
CREATE ROLE `RESOURCE`
-- case
CREATE ROLE RESOURCE
-- case
CREATE TABLE t (a int) STATS_BUCKETS=1
-- case
CREATE TABLE t (a int) STATS_BUCKETS='abc'
-- case
CREATE TABLE t (a int) STATS_BUCKETS=
-- case
CREATE TABLE t (a int) STATS_TOPN=1
-- case
CREATE TABLE t (a int) STATS_TOPN='abc'
-- case
CREATE TABLE t (a int) STATS_AUTO_RECALC=1
-- case
CREATE TABLE t (a int) STATS_AUTO_RECALC='abc'
-- case
CREATE TABLE t(a int) STATS_SAMPLE_RATE=0.1
-- case
CREATE TABLE t (a int) STATS_SAMPLE_RATE='abc'
-- case
CREATE TABLE t (a int) STATS_COL_CHOICE='all'
-- case
CREATE TABLE t (a int) STATS_COL_CHOICE='list'
-- case
CREATE TABLE t (a int) STATS_COL_CHOICE=1
-- case
CREATE TABLE t (a int, b int) STATS_COL_LIST='a,b'
-- case
CREATE TABLE t (a int, b int) STATS_COL_LIST=1
-- case
CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1
-- case
CREATE TABLE t (a int) STATS_BUCKETS=1,STATS_TOPN=1 PARTITION BY RANGE (a) (PARTITION p1 VALUES LESS THAN (200))
-- case
ALTER TABLE t STATS_OPTIONS='str'
-- case
ALTER TABLE t STATS_OPTIONS='str1,str2'
-- case
ALTER TABLE t STATS_OPTIONS="str1,str2"
-- case
ALTER TABLE t STATS_OPTIONS 'str1,str2'
-- case
ALTER TABLE t STATS_OPTIONS "str1,str2"
-- case
ALTER TABLE t STATS_OPTIONS=DEFAULT
-- case
ALTER TABLE t STATS_OPTIONS=default
-- case
ALTER TABLE t STATS_OPTIONS=DeFaUlT
-- case
ALTER TABLE t STATS_OPTIONS
-- case
CREATE TABLE t (a int) INSERT_METHOD=FIRST
