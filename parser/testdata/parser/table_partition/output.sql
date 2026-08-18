ALTER TABLE `t1` TRUNCATE PARTITION `p0`
-- case
ALTER TABLE `t1` TRUNCATE PARTITION `p0`,`p1`
-- case
ALTER TABLE `t1` TRUNCATE PARTITION ALL
-- case
-- error: line 1 column 41 near "p0" 
-- case
-- error: line 1 column 41 near "ALL" 
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION `p0`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `p0`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `p0`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION `p0`,`p1`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `p0`,`p1`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `p0`,`p1`
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION ALL
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG ALL
-- case
ALTER TABLE `t1` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG ALL
-- case
-- error: line 1 column 41 near "p0" 
-- case
-- error: line 1 column 41 near "ALL" 
-- case
-- error: line 1 column 40 near "" 
-- case
ALTER TABLE `t_n` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `local`
-- case
ALTER TABLE `t_n` OPTIMIZE PARTITION NO_WRITE_TO_BINLOG `local`,`local`
-- case
ALTER TABLE `t1` REPAIR PARTITION `p0`
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG `p0`
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG `p0`
-- case
ALTER TABLE `t1` REPAIR PARTITION `p0`,`p1`
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG `p0`,`p1`
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG `p0`,`p1`
-- case
ALTER TABLE `t1` REPAIR PARTITION ALL
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG ALL
-- case
ALTER TABLE `t1` REPAIR PARTITION NO_WRITE_TO_BINLOG ALL
-- case
-- error: line 1 column 39 near "p0" 
-- case
-- error: line 1 column 39 near "ALL" 
-- case
-- error: line 1 column 38 near "" 
-- case
ALTER TABLE `t_n` REPAIR PARTITION NO_WRITE_TO_BINLOG `local`
-- case
ALTER TABLE `t_n` REPAIR PARTITION NO_WRITE_TO_BINLOG `local`,`local`
-- case
ALTER TABLE `t1` IMPORT PARTITION `p0` TABLESPACE
-- case
ALTER TABLE `t1` IMPORT PARTITION `p0`,`p1` TABLESPACE
-- case
ALTER TABLE `t1` IMPORT PARTITION ALL TABLESPACE
-- case
-- error: line 1 column 36 near ", p0 TABLESPACE" 
-- case
-- error: line 1 column 39 near "ALL TABLESPACE" 
-- case
ALTER TABLE `t1` DISCARD PARTITION `p0` TABLESPACE
-- case
ALTER TABLE `t1` DISCARD PARTITION `p0`,`p1` TABLESPACE
-- case
ALTER TABLE `t1` DISCARD PARTITION ALL TABLESPACE
-- case
-- error: line 1 column 37 near ", p0 TABLESPACE" 
-- case
-- error: line 1 column 40 near "ALL TABLESPACE" 
-- case
ALTER TABLE `t1` ADD PARTITION (PARTITION `p5` VALUES LESS THAN (2010) COMMENT = 'APSTART '' APEND')
-- case
ALTER TABLE `t1` ADD PARTITION (PARTITION `p5` VALUES LESS THAN (2010) COMMENT = 'xxx')
-- case
CREATE TABLE `t1` (`a` INT NOT NULL,`b` INT NOT NULL,`c` INT NOT NULL,PRIMARY KEY(`a`, `b`)) PARTITION BY RANGE (`a`) (PARTITION `x1` VALUES LESS THAN (5),PARTITION `x2` VALUES LESS THAN (10),PARTITION `x3` VALUES LESS THAN (MAXVALUE))
-- case
CREATE TABLE `t1` (`a` INT NOT NULL) PARTITION BY RANGE (`a`) (PARTITION `x1` VALUES LESS THAN (5) TABLESPACE = `ts1`)
-- case
CREATE TABLE `t` (`a` INT) PARTITION BY RANGE (`a`) (PARTITION `p0` VALUES LESS THAN (63340531200) ENGINE = MyISAM,PARTITION `p1` VALUES LESS THAN (63342604800) ENGINE = MyISAM)
-- case
CREATE TABLE `t` (`a` INT) PARTITION BY RANGE (`a`) (PARTITION `p0` VALUES LESS THAN (63340531200) ENGINE = MyISAM COMMENT = 'xxx',PARTITION `p1` VALUES LESS THAN (63342604800) ENGINE = MyISAM)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY RANGE (`a`) (PARTITION `p0` VALUES LESS THAN (63340531200) COMMENT = 'xxx' ENGINE = MyISAM,PARTITION `p1` VALUES LESS THAN (63342604800) ENGINE = MyISAM)
-- case
CREATE TABLE `t` (`id` INT) PARTITION BY RANGE (`id`) SUBPARTITION BY KEY (`id`) SUBPARTITIONS 2 (PARTITION `p0` VALUES LESS THAN (42))
-- case
CREATE TABLE `t` (`id` INT) PARTITION BY RANGE (`id`) SUBPARTITION BY HASH (`id`) (PARTITION `p0` VALUES LESS THAN (42))
-- case
CREATE TABLE `t1` (`a` VARCHAR(5),`b` INT,`c` VARCHAR(10),`d` DATETIME) PARTITION BY RANGE COLUMNS (`b`,`c`) SUBPARTITION BY HASH (TO_SECONDS(`d`)) (PARTITION `p0` VALUES LESS THAN (2, _UTF8MB4'b'),PARTITION `p1` VALUES LESS THAN (4, _UTF8MB4'd'),PARTITION `p2` VALUES LESS THAN (10, _UTF8MB4'za'))
-- case
CREATE TABLE `t1` (`a` INT,`b` TIMESTAMP DEFAULT _UTF8MB4'0000-00-00 00:00:00') ENGINE = INNODB PARTITION BY LINEAR HASH (`a`) PARTITIONS 1
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY HASH (`a`) (PARTITION `x`,PARTITION `y`)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY KEY (`a`) (PARTITION `x`,PARTITION `y`)
-- case
-- error: [ddl:1479]Syntax : RANGE PARTITIONING requires definition of VALUES LESS THAN for each partition
-- case
-- error: [ddl:1479]Syntax : LIST PARTITIONING requires definition of VALUES IN for each partition
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:1480]Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
-- case
-- error: [ddl:1480]Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY RANGE (`a`) (PARTITION `x` VALUES LESS THAN (MAXVALUE))
-- case
-- error: line 1 column 86 near "))" 
-- case
CREATE TABLE `t` (`a` VARCHAR(100),`b` INT) PARTITION BY LIST COLUMNS (`a`) (PARTITION `p1` VALUES IN (_UTF8MB4'a', _UTF8MB4'b', _UTF8MB4'DEFAULT'),PARTITION `pDef` DEFAULT)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY RANGE (`a`) (PARTITION `x` VALUES LESS THAN (10))
-- case
-- error: [ddl:1480]Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:1480]Only LIST PARTITIONING can use VALUES IN in partition definition
-- case
-- error: [ddl:1480]Only LIST PARTITIONING can use VALUES IN in partition definition
-- case
-- error: [ddl:1480]Only LIST PARTITIONING can use VALUES IN in partition definition
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY LIST (`a`) (PARTITION `x` VALUES IN (10))
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY LIST (`a`) (PARTITION `x` DEFAULT)
-- case
-- error: line 1 column 78 near "maxvalue))" 
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY LIST (`a`) (PARTITION `x` VALUES IN (DEFAULT, 10))
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:4113]Wrong partitioning type, expected type: `SYSTEM_TIME`
-- case
-- error: [ddl:4113]Wrong partitioning type, expected type: `SYSTEM_TIME`
-- case
-- error: [ddl:4113]Wrong partitioning type, expected type: `SYSTEM_TIME`
-- case
-- error: [ddl:4113]Wrong partitioning type, expected type: `SYSTEM_TIME`
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY SYSTEM_TIME (PARTITION `x` HISTORY,PARTITION `y` CURRENT)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY HASH (`a`) PARTITIONS 1
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY KEY (`a`) PARTITIONS 1
-- case
-- error: [ddl:1492]For RANGE partitions each partition must be defined
-- case
-- error: [ddl:1492]For LIST partitions each partition must be defined
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:4128]Wrong Partitions: must have at least one HISTORY and exactly one last CURRENT
-- case
-- error: [ddl:1657]Cannot have more than one value for this type of RANGE partitioning
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
CREATE TABLE `t1` (`a` INT,`b` INT) PARTITION BY RANGE COLUMNS (`a`,`b`) (PARTITION `x` VALUES LESS THAN (10, 20))
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: [ddl:1658]Row expressions in VALUES IN only allowed for multi-field column partitioning
-- case
CREATE TABLE `t1` (`a` INT,`b` INT) PARTITION BY LIST COLUMNS (`a`,`b`) (PARTITION `x` VALUES IN ((10, 20)))
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: [ddl:1653]Inconsistency in usage of column lists for partitioning
-- case
-- error: line 1 column 48 near ")" 
-- case
CREATE TABLE `t1` (`a` INT PRIMARY KEY) PARTITION BY KEY () PARTITIONS 1
-- case
-- error: line 1 column 53 near ") (partition x values less than maxvalue)" 
-- case
-- error: line 1 column 52 near ") (partition x default)" 
-- case
-- error: line 1 column 79 near "))" 
-- case
-- error: line 1 column 71 near "))" 
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY LIST (`a`) (PARTITION `x` DEFAULT)
-- case
-- error: line 1 column 75 near "range (b) (partition x values less than maxvalue)" 
-- case
-- error: [ddl:1484]Wrong number of partitions defined, mismatch with previous setting
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY HASH (`a`) (PARTITION `x`,PARTITION `y`)
-- case
-- error: [ddl:1485]Wrong number of subpartitions defined, mismatch with previous setting
-- case
CREATE TABLE `t1` (`a` INT,`b` INT) PARTITION BY RANGE (`a`) SUBPARTITION BY HASH (`b`) SUBPARTITIONS 2 (PARTITION `x` VALUES LESS THAN (MAXVALUE) (SUBPARTITION `y`,SUBPARTITION `z`))
-- case
CREATE TABLE `t1` (`a` INT,`b` INT) PARTITION BY RANGE (`a`) SUBPARTITION BY HASH (`b`) SUBPARTITIONS 2 (PARTITION `x` VALUES LESS THAN (10) (SUBPARTITION `y`,SUBPARTITION `z`),PARTITION `a` VALUES LESS THAN (20) (SUBPARTITION `b`,SUBPARTITION `c`))
-- case
-- error: [ddl:1485]Wrong number of subpartitions defined, mismatch with previous setting
-- case
-- error: [ddl:1500]It is only possible to mix RANGE/LIST partitioning with HASH/KEY partitioning for subpartitioning
-- case
-- error: [ddl:1504]Number of partitions = 0 is not an allowed value
-- case
-- error: [ddl:1504]Number of subpartitions = 0 is not an allowed value
-- case
-- error: line 1 column 69 near "limit 50000 (partition x history, partition y current)" 
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY SYSTEM_TIME INTERVAL 7 DAY (PARTITION `x` HISTORY,PARTITION `y` CURRENT)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY SYSTEM_TIME LIMIT 50000 (PARTITION `x` HISTORY,PARTITION `y` CURRENT)
-- case
CREATE TABLE `t1` (`a` INT) PARTITION BY HASH (`a`) (PARTITION `x` ENGINE = InnoDB COMMENT = 'xxxx' DATA DIRECTORY = '/var/data' INDEX DIRECTORY = '/var/index' MAX_ROWS = 70000 MIN_ROWS = 50 TABLESPACE = `innodb_file_per_table` NODEGROUP = 255)
-- case
CREATE TABLE `t1` (`a` INT,`b` INT) PARTITION BY RANGE (`a`) SUBPARTITION BY HASH (`b`) SUBPARTITIONS 1 (PARTITION `x` VALUES LESS THAN (MAXVALUE) (SUBPARTITION `y` ENGINE = InnoDB COMMENT = 'xxxx' DATA DIRECTORY = '/var/data' INDEX DIRECTORY = '/var/index' MAX_ROWS = 70000 MIN_ROWS = 50 TABLESPACE = `innodb_file_per_table` NODEGROUP = 255))
