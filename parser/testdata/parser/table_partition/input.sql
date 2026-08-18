ALTER TABLE t1 TRUNCATE PARTITION p0
-- case
ALTER TABLE t1 TRUNCATE PARTITION p0, p1
-- case
ALTER TABLE t1 TRUNCATE PARTITION ALL
-- case
ALTER TABLE t1 TRUNCATE PARTITION ALL, p0
-- case
ALTER TABLE t1 TRUNCATE PARTITION p0, ALL
-- case
ALTER TABLE t1 OPTIMIZE PARTITION p0
-- case
ALTER TABLE t1 OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0
-- case
ALTER TABLE t1 OPTIMIZE PARTITION LOCAL p0
-- case
ALTER TABLE t1 OPTIMIZE PARTITION p0, p1
-- case
ALTER TABLE t1 OPTIMIZE PARTITION NO_WRITE_TO_BINLOG p0, p1
-- case
ALTER TABLE t1 OPTIMIZE PARTITION LOCAL p0, p1
-- case
ALTER TABLE t1 OPTIMIZE PARTITION ALL
-- case
ALTER TABLE t1 OPTIMIZE PARTITION NO_WRITE_TO_BINLOG ALL
-- case
ALTER TABLE t1 OPTIMIZE PARTITION LOCAL ALL
-- case
ALTER TABLE t1 OPTIMIZE PARTITION ALL, p0
-- case
ALTER TABLE t1 OPTIMIZE PARTITION p0, ALL
-- case
ALTER TABLE t_n OPTIMIZE PARTITION LOCAL
-- case
ALTER TABLE t_n OPTIMIZE PARTITION LOCAL local
-- case
ALTER TABLE t_n OPTIMIZE PARTITION LOCAL local, local
-- case
ALTER TABLE t1 REPAIR PARTITION p0
-- case
ALTER TABLE t1 REPAIR PARTITION NO_WRITE_TO_BINLOG p0
-- case
ALTER TABLE t1 REPAIR PARTITION LOCAL p0
-- case
ALTER TABLE t1 REPAIR PARTITION p0, p1
-- case
ALTER TABLE t1 REPAIR PARTITION NO_WRITE_TO_BINLOG p0, p1
-- case
ALTER TABLE t1 REPAIR PARTITION LOCAL p0, p1
-- case
ALTER TABLE t1 REPAIR PARTITION ALL
-- case
ALTER TABLE t1 REPAIR PARTITION NO_WRITE_TO_BINLOG ALL
-- case
ALTER TABLE t1 REPAIR PARTITION LOCAL ALL
-- case
ALTER TABLE t1 REPAIR PARTITION ALL, p0
-- case
ALTER TABLE t1 REPAIR PARTITION p0, ALL
-- case
ALTER TABLE t_n REPAIR PARTITION LOCAL
-- case
ALTER TABLE t_n REPAIR PARTITION LOCAL local
-- case
ALTER TABLE t_n REPAIR PARTITION LOCAL local, local
-- case
ALTER TABLE t1 IMPORT PARTITION p0 TABLESPACE
-- case
ALTER TABLE t1 IMPORT PARTITION p0, p1 TABLESPACE
-- case
ALTER TABLE t1 IMPORT PARTITION ALL TABLESPACE
-- case
ALTER TABLE t1 IMPORT PARTITION ALL, p0 TABLESPACE
-- case
ALTER TABLE t1 IMPORT PARTITION p0, ALL TABLESPACE
-- case
ALTER TABLE t1 DISCARD PARTITION p0 TABLESPACE
-- case
ALTER TABLE t1 DISCARD PARTITION p0, p1 TABLESPACE
-- case
ALTER TABLE t1 DISCARD PARTITION ALL TABLESPACE
-- case
ALTER TABLE t1 DISCARD PARTITION ALL, p0 TABLESPACE
-- case
ALTER TABLE t1 DISCARD PARTITION p0, ALL TABLESPACE
-- case
ALTER TABLE t1 ADD PARTITION (PARTITION `p5` VALUES LESS THAN (2010) COMMENT 'APSTART \' APEND')
-- case
ALTER TABLE t1 ADD PARTITION (PARTITION `p5` VALUES LESS THAN (2010) COMMENT = 'xxx')
-- case
CREATE TABLE t1 (a int not null,b int not null,c int not null,primary key(a,b))
		partition by range (a)
		partitions 3
		(partition x1 values less than (5),
		 partition x2 values less than (10),
		 partition x3 values less than maxvalue);
-- case
CREATE TABLE t1 (a int not null) partition by range (a) (partition x1 values less than (5) tablespace ts1)
-- case
create table t (a int) partition by range (a)
		  (PARTITION p0 VALUES LESS THAN (63340531200) ENGINE = MyISAM,
		   PARTITION p1 VALUES LESS THAN (63342604800) ENGINE MyISAM)
-- case
create table t (a int) partition by range (a)
		  (PARTITION p0 VALUES LESS THAN (63340531200) ENGINE = MyISAM COMMENT 'xxx',
		   PARTITION p1 VALUES LESS THAN (63342604800) ENGINE = MyISAM)
-- case
create table t1 (a int) partition by range (a)
		  (PARTITION p0 VALUES LESS THAN (63340531200) COMMENT 'xxx' ENGINE = MyISAM ,
		   PARTITION p1 VALUES LESS THAN (63342604800) ENGINE = MyISAM)
-- case
create table t (id int)
		    partition by range (id)
		    subpartition by key (id) subpartitions 2
		    (partition p0 values less than (42))
-- case
create table t (id int)
		    partition by range (id)
		    subpartition by hash (id)
		    (partition p0 values less than (42))
-- case
create table t1 (a varchar(5), b int signed, c varchar(10), d datetime)
		partition by range columns(b,c)
		subpartition by hash(to_seconds(d))
		( partition p0 values less than (2, 'b'),
		  partition p1 values less than (4, 'd'),
		  partition p2 values less than (10, 'za'));
-- case
CREATE TABLE t1 (a INT, b TIMESTAMP DEFAULT '0000-00-00 00:00:00')
ENGINE=INNODB PARTITION BY LINEAR HASH (a) PARTITIONS 1;
-- case
create table t1 (a int) partition by hash (a) (partition x, partition y)
-- case
create table t1 (a int) partition by key (a) (partition x, partition y)
-- case
create table t1 (a int) partition by range (a) (partition x, partition y)
-- case
create table t1 (a int) partition by list (a) (partition x, partition y)
-- case
create table t1 (a int) partition by system_time (partition x, partition y)
-- case
create table t1 (a int) partition by hash (a) (partition x values less than (10))
-- case
create table t1 (a int) partition by key (a) (partition x values less than (10))
-- case
create table t1 (a int) partition by range (a) (partition x values less than (maxvalue))
-- case
create table t1 (a int) partition by range (a) (partition x values less than (default))
-- case
create table t (a varchar(100), b int) partition by list columns (a) (partition p1 values in ('a','b','DEFAULT'), partition pDef values in (default))
-- case
create table t1 (a int) partition by range (a) (partition x values less than (10))
-- case
create table t1 (a int) partition by list (a) (partition x values less than (10))
-- case
create table t1 (a int) partition by system_time (partition x values less than (10))
-- case
create table t1 (a int) partition by hash (a) (partition x values in (10))
-- case
create table t1 (a int) partition by key (a) (partition x values in (10))
-- case
create table t1 (a int) partition by range (a) (partition x values in (10))
-- case
create table t1 (a int) partition by list (a) (partition x values in (10))
-- case
create table t1 (a int) partition by list (a) (partition x values in (default))
-- case
create table t1 (a int) partition by list (a) (partition x values in (maxvalue))
-- case
create table t1 (a int) partition by list (a) (partition x values in (default, 10))
-- case
create table t1 (a int) partition by system_time (partition x values in (10))
-- case
create table t1 (a int) partition by hash (a) (partition x history, partition y current)
-- case
create table t1 (a int) partition by key (a) (partition x history, partition y current)
-- case
create table t1 (a int) partition by range (a) (partition x history, partition y current)
-- case
create table t1 (a int) partition by list (a) (partition x history, partition y current)
-- case
create table t1 (a int) partition by system_time (partition x history, partition y current)
-- case
create table t1 (a int) partition by hash (a)
-- case
create table t1 (a int) partition by key (a)
-- case
create table t1 (a int) partition by range (a)
-- case
create table t1 (a int) partition by list (a)
-- case
create table t1 (a int) partition by system_time
-- case
create table t1 (a int) partition by system_time (partition x history)
-- case
create table t1 (a int) partition by system_time (partition x current)
-- case
create table t1 (a int, b int) partition by range (a) (partition x values less than (10, 20))
-- case
create table t (id int) partition by range columns (id) (partition p0 values less than (1, 2))
-- case
create table t1 (a int, b int) partition by range columns (a, b) (partition x values less than (10, 20))
-- case
create table t1 (a int, b int) partition by range columns (a, b) (partition x values less than (10))
-- case
create table t1 (a int, b int) partition by range columns (a, b) (partition x values less than maxvalue)
-- case
create table t1 (a int, b int) partition by list (a) (partition x values in ((10, 20)))
-- case
create table t1 (a int, b int) partition by list columns (a, b) (partition x values in ((10, 20)))
-- case
create table t1 (a int, b int) partition by list columns (a, b) (partition x values in (10, 20))
-- case
create table t1 (a int, b int) partition by list columns (a, b) (partition x values in (10, (20, 30)))
-- case
create table t1 (a int, b int) partition by list columns (a, b) (partition x values in ((10, 20), 30))
-- case
create table t1 (a int, b int) partition by list columns (a, b) (partition x values in ((10, 20), (30, 40, 50)))
-- case
create table t1 (a int) partition by hash (a) ()
-- case
create table t1 (a int primary key) partition by key ()
-- case
create table t1 (a int) partition by range columns () (partition x values less than maxvalue)
-- case
create table t1 (a int) partition by list columns () (partition x default)
-- case
create table t1 (a int) partition by range (a) (partition x values less than ())
-- case
create table t1 (a int) partition by list (a) (partition x values in ())
-- case
create table t1 (a int) partition by list (a) (partition x default)
-- case
create table t1 (a int, b int) partition by range (a) subpartition by range (b) (partition x values less than maxvalue)
-- case
create table t1 (a int) partition by hash (a) partitions 2 (partition x)
-- case
create table t1 (a int) partition by hash (a) partitions 2 (partition x, partition y)
-- case
create table t1 (a int, b int) partition by range (a) subpartition by hash (b) subpartitions 2 (partition x values less than maxvalue (subpartition y))
-- case
create table t1 (a int, b int) partition by range (a) subpartition by hash (b) subpartitions 2 (partition x values less than maxvalue (subpartition y, subpartition z))
-- case
create table t1 (a int, b int) partition by range (a) subpartition by hash (b) (partition x values less than (10) (subpartition y,subpartition z),partition a values less than (20) (subpartition b,subpartition c))
-- case
create table t1 (a int, b int) partition by range (a) subpartition by hash (b) (partition x values less than (10) (subpartition y),partition a values less than (20) (subpartition b,subpartition c))
-- case
create table t1 (a int, b int) partition by range (a) (partition x values less than (10) (subpartition y))
-- case
create table t1 (a int) partition by hash (a) partitions 0
-- case
create table t1 (a int, b int) partition by range (a) subpartition by hash (b) subpartitions 0 (partition x values less than (10))
-- case
create table t1 (a int) partition by system_time interval 7 day limit 50000 (partition x history, partition y current)
-- case
create table t1 (a int) partition by system_time interval 7 day (partition x history, partition y current)
-- case
create table t1 (a int) partition by system_time limit 50000 (partition x history, partition y current)
-- case
create table t1 (a int) partition by hash(a) (partition x engine InnoDB comment 'xxxx' data directory '/var/data' index directory '/var/index' max_rows 70000 min_rows 50 tablespace `innodb_file_per_table` nodegroup 255)
-- case
create table t1 (a int, b int) partition by range(a) subpartition by hash(b) (partition x values less than maxvalue (subpartition y engine InnoDB comment 'xxxx' data directory '/var/data' index directory '/var/index' max_rows 70000 min_rows 50 tablespace `innodb_file_per_table` nodegroup 255))
