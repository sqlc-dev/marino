-- error: line 1 column 6 near "" 
-- case
-- error: line 1 column 12 near "" 
-- case
-- error: line 1 column 18 near "" 
-- case
-- error: line 1 column 19 near ")" 
-- case
-- error: line 1 column 19 near ");" 
-- case
-- error: line 1 column 18 near "* (a varchar(50), b int);" 
-- case
CREATE TABLE `foo` (`a` VARCHAR(50),`b` INT)
-- case
CREATE TABLE `foo` (`a` TINYINT UNSIGNED)
-- case
CREATE TABLE `foo` (`a` SMALLINT UNSIGNED,`b` INT UNSIGNED)
-- case
CREATE TABLE `foo` (`a` BIGINT UNSIGNED,`b` TINYINT(1))
-- case
-- error: line 1 column 47 near "CREATE TABLE bar (x INT, y int64)" 
-- case
CREATE TABLE `foo` (`a` INT,`b` FLOAT); CREATE TABLE `bar` (`x` DOUBLE,`y` FLOAT)
-- case
-- error: line 1 column 25 near "bytes)" 
-- case
CREATE TABLE `foo` (`a` SMALLINT UNSIGNED,`b` INT UNSIGNED)
-- case
CREATE TABLE `foo` (`a` SMALLINT UNSIGNED,`b` INT UNSIGNED)
-- case
-- error: line 1 column 56 near "// foo" 
-- case
CREATE TABLE `foo` (`a` SMALLINT UNSIGNED,`b` INT UNSIGNED)
-- case
CREATE TABLE `foo` (`a` SMALLINT UNSIGNED,`b` INT UNSIGNED)
-- case
CREATE TABLE `foo` (`name` CHAR(50) BINARY)
-- case
CREATE TABLE `foo` (`name` CHAR(50) COLLATE utf8_bin)
-- case
CREATE TABLE `foo` (`id` VARCHAR(50) COLLATE utf8_bin)
-- case
CREATE TABLE `foo` (`name` CHAR(50) CHARACTER SET UTF8)
-- case
CREATE TABLE `foo` (`name` CHAR(50) BINARY CHARACTER SET UTF8)
-- case
-- error: line 1 column 67 near "CHARACTER set utf8)" 
-- case
CREATE TABLE `foo` (`name` CHAR(50) BINARY CHARACTER SET UTF8 COLLATE utf8_bin)
-- case
-- error: [parser:1064]Multiple COLLATE clauses line 1 column 86 near ")" 
-- case
-- error: [parser:1064]Multiple COLLATE clauses line 1 column 69 near ")" 
-- case
-- error: [parser:1064]Multiple COLLATE clauses line 1 column 81 near ")" 
-- case
-- error: line 1 column 22 near ", b);" 
-- case
-- error: line 1 column 20 near ", b.c);" 
-- case
-- error: line 1 column 14 near "(name CHAR(50) BINARY)" 
-- case
-- error: [parser:1064]Multiple COLLATE clauses line 1 column 81 near ", INDEX (name ASC))" 
-- case
-- error: [parser:1064]Multiple COLLATE clauses line 1 column 81 near ", INDEX (name DESC))" 
-- case
ALTER TABLE `tmp` CACHE
-- case
ALTER TABLE `tmp` NOCACHE
-- case
CREATE TEMPORARY TABLE `t` (`a` VARCHAR(50),`b` INT)
-- case
CREATE TEMPORARY TABLE `t` LIKE `t1`
-- case
DROP TEMPORARY TABLE `t`
-- case
CREATE GLOBAL TEMPORARY TABLE `t` (`a` INT,`b` VARCHAR(255)) ON COMMIT DELETE ROWS
-- case
CREATE TEMPORARY TABLE `t` (`a` INT,`b` VARCHAR(255))
-- case
-- error: line 1 column 55 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
-- case
-- error: line 1 column 70 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
-- case
-- error: line 1 column 60 near ""GLOBAL TEMPORARY and ON COMMIT DELETE ROWS must appear together 
-- case
CREATE GLOBAL TEMPORARY TABLE `t` (`a` INT,`b` VARCHAR(255)) PARTITION BY HASH (`a`) PARTITIONS 10 ON COMMIT DELETE ROWS
-- case
CREATE GLOBAL TEMPORARY TABLE `t` (`a` INT,`b` VARCHAR(255)) ON COMMIT PRESERVE ROWS
-- case
DROP GLOBAL TEMPORARY TABLE `t`
-- case
DROP TEMPORARY TABLE `t`
-- case
CREATE TABLE `foo` (`node_id` VARCHAR(50),`b` INT)
-- case
CREATE TABLE `foo` (`node_state` VARCHAR(50),`b` INT)
-- case
CREATE TABLE `t` (`c` INT) AVG_ROW_LENGTH = 3
-- case
CREATE TABLE `t` (`c` INT) AVG_ROW_LENGTH = 3
-- case
CREATE TABLE `t` (`c` INT) CHECKSUM = 0
-- case
CREATE TABLE `t` (`c` INT) CHECKSUM = 1
-- case
CREATE TABLE `t` (`c` INT) TABLE_CHECKSUM = 0
-- case
CREATE TABLE `t` (`c` INT) TABLE_CHECKSUM = 1
-- case
CREATE TABLE `t` (`c` INT) COMPRESSION = 'NONE'
-- case
CREATE TABLE `t` (`c` INT) COMPRESSION = 'lz4'
-- case
CREATE TABLE `t` (`c` INT) CONNECTION = 'abc'
-- case
CREATE TABLE `t` (`c` INT) CONNECTION = 'abc'
-- case
CREATE TABLE `t` (`c` INT) KEY_BLOCK_SIZE = 1024
-- case
CREATE TABLE `t` (`c` INT) KEY_BLOCK_SIZE = 1024
-- case
CREATE TABLE `t` (`c` INT) MAX_ROWS = 1000
-- case
CREATE TABLE `t` (`c` INT) MAX_ROWS = 1000
-- case
CREATE TABLE `t` (`c` INT) MIN_ROWS = 1000
-- case
CREATE TABLE `t` (`c` INT) MIN_ROWS = 1000
-- case
CREATE TABLE `t` (`c` INT) PASSWORD = 'abc'
-- case
CREATE TABLE `t` (`c` INT) PASSWORD = 'abc'
-- case
CREATE TABLE `t` (`c` INT) DELAY_KEY_WRITE = 1
-- case
CREATE TABLE `t` (`c` INT) DELAY_KEY_WRITE = 1
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = DEFAULT
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = DEFAULT
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = FIXED
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = COMPRESSED
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = COMPACT
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = REDUNDANT
-- case
CREATE TABLE `t` (`c` INT) ROW_FORMAT = DYNAMIC
-- case
CREATE TABLE `t` (`c` INT) STATS_PERSISTENT = DEFAULT /* TableOptionStatsPersistent is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) STATS_PERSISTENT = DEFAULT /* TableOptionStatsPersistent is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) STATS_PERSISTENT = DEFAULT /* TableOptionStatsPersistent is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) STATS_SAMPLE_PAGES = 0
-- case
CREATE TABLE `t` (`c` INT) STATS_SAMPLE_PAGES = 10
-- case
CREATE TABLE `t` (`c` INT) STATS_SAMPLE_PAGES = 10
-- case
CREATE TABLE `t` (`c` INT) STATS_SAMPLE_PAGES = DEFAULT
-- case
CREATE TABLE `t` (`c` INT) PACK_KEYS = DEFAULT /* TableOptionPackKeys is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) PACK_KEYS = DEFAULT /* TableOptionPackKeys is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) PACK_KEYS = DEFAULT /* TableOptionPackKeys is not supported */ 
-- case
CREATE TABLE `t` (`c` INT) STORAGE DISK
-- case
CREATE TABLE `t` (`c` INT) STORAGE MEMORY
-- case
CREATE TABLE `t` (`c` INT) SECONDARY_ENGINE = NULL
-- case
CREATE TABLE `t` (`c` INT) SECONDARY_ENGINE = 'innodb'
-- case
CREATE TABLE `t` (`c` INT) SECONDARY_ENGINE = 'null'
-- case
CREATE TABLE `testTableCompression` (`c` VARCHAR(15000)) COMPRESSION = 'ZLIB'
-- case
CREATE TABLE `t1` (`c1` INT) COMPRESSION = 'zlib'
-- case
CREATE TABLE `t1` (`c1` INT) DEFAULT COLLATE = BINARY
-- case
CREATE TABLE `t1` (`c1` INT) DEFAULT COLLATE = UTF8MB4_0900_AS_CS
-- case
CREATE TABLE `t1` (`c1` INT) DEFAULT CHARACTER SET = BINARY DEFAULT COLLATE = BINARY
-- case
CREATE TABLE `t1` (`c1` INT) AUTOEXTEND_SIZE = 4M
-- case
ALTER TABLE `t_n` UNION = (), KEY_BLOCK_SIZE = 1
-- case
ALTER TABLE `d_n`.`t_n` UNION = (`t_n`) REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` LOCK = DEFAULT, UNION = (`t_n`,`d_n`.`t_n`) REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` ALGORITHM = DEFAULT, MAX_ROWS = 10, UNION = (`d_n`.`t_n`), ROW_FORMAT = REDUNDANT, STATS_PERSISTENT = DEFAULT /* TableOptionStatsPersistent is not supported */ 
-- case
-- error: line 1 column 89 near "not 3), partition p2 values less than (20));" 
-- case
-- error: line 1 column 90 near "or 3), partition p2 values less than (20));" 
-- case
-- error: line 1 column 90 near "is null), partition p2 values less than (20));" 
-- case
-- error: line 1 column 47 near "is null) (partition p0 values less than (10));" 
-- case
-- error: line 1 column 45 near "not b) (partition p0 values in (10, 20));" 
-- case
-- error: line 1 column 46 near "not b );" 
-- case
-- error: line 1 column 90 near "in (3, 4, 5)), partition p2 values less than (20));" 
-- case
CREATE TABLE `t` (`id` INT) ENGINE = INNDB PARTITION BY RANGE (`id`) (PARTITION `p0` VALUES LESS THAN (10),PARTITION `p1` VALUES LESS THAN (20))
-- case
CREATE TABLE `t` (`c` INT) PARTITION BY HASH (`c`) PARTITIONS 32
-- case
-- error: [ddl:1480]Only RANGE PARTITIONING can use VALUES LESS THAN in partition definition
-- case
CREATE TABLE `t` (`c` INT) PARTITION BY RANGE (YEAR(`VDate`)) (PARTITION `p1980` VALUES LESS THAN (1980) ENGINE = MyISAM,PARTITION `p1990` VALUES LESS THAN (1990) ENGINE = MyISAM,PARTITION `pothers` VALUES LESS THAN (MAXVALUE) ENGINE = MyISAM)
-- case
CREATE TABLE `t` (`c` INT,`create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP() COMMENT '') PARTITION BY RANGE (UNIX_TIMESTAMP(`create_time`)) (PARTITION `p201610` VALUES LESS THAN (1477929600),PARTITION `p201611` VALUES LESS THAN (1480521600),PARTITION `p201612` VALUES LESS THAN (1483200000),PARTITION `p201701` VALUES LESS THAN (1485878400),PARTITION `p201702` VALUES LESS THAN (1488297600),PARTITION `p201703` VALUES LESS THAN (1490976000))
-- case
CREATE TABLE `md_product_shop` (`shopCode` VARCHAR(4) DEFAULT NULL COMMENT '地点') ENGINE = InnoDB DEFAULT CHARACTER SET = UTF8MB4 PARTITION BY KEY (`shopCode`) PARTITIONS 19
-- case
CREATE TABLE `payinfo1` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`oderTime` DATETIME NOT NULL) ENGINE = InnoDB AUTO_INCREMENT = 641533032 DEFAULT CHARACTER SET = UTF8 ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8 PARTITION BY RANGE COLUMNS (`oderTime`) (PARTITION `P2011` VALUES LESS THAN (_UTF8MB4'2012-01-01 00:00:00') ENGINE = InnoDB,PARTITION `P1201` VALUES LESS THAN (_UTF8MB4'2012-02-01 00:00:00') ENGINE = InnoDB,PARTITION `PMAX` VALUES LESS THAN (MAXVALUE) ENGINE = InnoDB)
-- case
CREATE TABLE `app_channel_daily_report` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`app_version` VARCHAR(32) COLLATE utf8_unicode_ci NOT NULL DEFAULT _UTF8MB4'default',`gmt_create` DATETIME NOT NULL COMMENT '创建时间',PRIMARY KEY(`id`)) ENGINE = InnoDB AUTO_INCREMENT = 33703438 DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_UNICODE_CI PARTITION BY RANGE (MONTH(`gmt_create`)-1) (PARTITION `part0` VALUES LESS THAN (1) COMMENT = '1月份' ENGINE = InnoDB,PARTITION `part1` VALUES LESS THAN (2) COMMENT = '2月份' ENGINE = InnoDB,PARTITION `part2` VALUES LESS THAN (3) COMMENT = '3月份' ENGINE = InnoDB,PARTITION `part3` VALUES LESS THAN (4) COMMENT = '4月份' ENGINE = InnoDB,PARTITION `part4` VALUES LESS THAN (5) COMMENT = '5月份' ENGINE = InnoDB,PARTITION `part5` VALUES LESS THAN (6) COMMENT = '6月份' ENGINE = InnoDB,PARTITION `part6` VALUES LESS THAN (7) COMMENT = '7月份' ENGINE = InnoDB,PARTITION `part7` VALUES LESS THAN (8) COMMENT = '8月份' ENGINE = InnoDB,PARTITION `part8` VALUES LESS THAN (9) COMMENT = '9月份' ENGINE = InnoDB,PARTITION `part9` VALUES LESS THAN (10) COMMENT = '10月份' ENGINE = InnoDB,PARTITION `part10` VALUES LESS THAN (11) COMMENT = '11月份' ENGINE = InnoDB,PARTITION `part11` VALUES LESS THAN (12) COMMENT = '12月份' ENGINE = InnoDB)
-- case
-- error: line 1 column 37 near "primary_region="us";" 
-- case
-- error: line 1 column 30 near "regions="us,3";" 
-- case
-- error: line 1 column 32 near "followers="us,3";" 
-- case
-- error: line 1 column 32 near "followers=3;" 
-- case
-- error: line 1 column 32 near "followers=0;" 
-- case
-- error: line 1 column 29 near "voters=3;" 
-- case
-- error: line 1 column 31 near "learners="us,3";" 
-- case
-- error: line 1 column 31 near "learners=3;" 
-- case
-- error: line 1 column 31 near "schedule="even";" 
-- case
-- error: line 1 column 34 near "constraints="ww";" 
-- case
-- error: line 1 column 41 near "leader_constraints="ww";" 
-- case
-- error: line 1 column 43 near "follower_constraints="ww";" 
-- case
-- error: line 1 column 40 near "voter_constraints="ww";" 
-- case
-- error: line 1 column 42 near "learner_constraints="ww";" 
-- case
-- error: line 1 column 42 near "survival_preference="ww";" 
-- case
-- error: line 1 column 53 near "primary_region="us" */;" 
-- case
-- error: line 1 column 46 near "regions="us,3" */;" 
-- case
-- error: line 1 column 48 near "followers="us,3 */";" 
-- case
-- error: line 1 column 48 near "followers=3 */;" 
-- case
-- error: line 1 column 48 near "followers=0 */;" 
-- case
-- error: line 1 column 45 near "voters="us,3" */;" 
-- case
-- error: line 1 column 53 near "primary_region="us" regions="us,3"  */;" 
-- case
CREATE TABLE `t` (`c` INT) PLACEMENT POLICY = `ww`
-- case
CREATE TABLE `t` (`c` INT) PLACEMENT POLICY = `x`
-- case
CREATE TABLE `t` (`c` INT) PLACEMENT POLICY = `y`
-- case
-- error: line 1 column 28 near "primary_region="us";" 
-- case
-- error: line 1 column 21 near "regions="us,3";" 
-- case
-- error: line 1 column 23 near "followers=3;" 
-- case
-- error: line 1 column 23 near "followers=0;" 
-- case
-- error: line 1 column 20 near "voters=3;" 
-- case
-- error: line 1 column 22 near "learners=3;" 
-- case
-- error: line 1 column 22 near "schedule="even";" 
-- case
-- error: line 1 column 25 near "constraints="ww";" 
-- case
-- error: line 1 column 32 near "leader_constraints="ww";" 
-- case
-- error: line 1 column 34 near "follower_constraints="ww";" 
-- case
-- error: line 1 column 31 near "voter_constraints="ww";" 
-- case
-- error: line 1 column 33 near "learner_constraints="ww";" 
-- case
-- error: line 1 column 44 near "primary_region="us" */;" 
-- case
ALTER TABLE `t` PLACEMENT POLICY = `ww`
-- case
ALTER TABLE `t` PLACEMENT POLICY = `ww`
-- case
ALTER TABLE `t` COMPACT
-- case
ALTER TABLE `t` COMPACT TIFLASH REPLICA
-- case
ALTER TABLE `t` COMPACT PARTITION `p1`,`p2`
-- case
ALTER TABLE `t` COMPACT PARTITION `p1`,`p2` TIFLASH REPLICA
-- case
-- error: line 1 column 32 near "primary_region="us";" 
-- case
-- error: line 1 column 25 near "regions="us,3";" 
-- case
-- error: line 1 column 27 near "followers=3;" 
-- case
-- error: line 1 column 27 near "followers=0;" 
-- case
-- error: line 1 column 24 near "voters=3;" 
-- case
-- error: line 1 column 26 near "learners=3;" 
-- case
-- error: line 1 column 26 near "schedule="even";" 
-- case
-- error: line 1 column 29 near "constraints="ww";" 
-- case
-- error: line 1 column 36 near "leader_constraints="ww";" 
-- case
-- error: line 1 column 38 near "follower_constraints="ww";" 
-- case
-- error: line 1 column 35 near "voter_constraints="ww";" 
-- case
-- error: line 1 column 37 near "learner_constraints="ww";" 
-- case
-- error: line 1 column 48 near "primary_region="us" */;" 
-- case
CREATE DATABASE `t` PLACEMENT POLICY = `ww`
-- case
CREATE DATABASE `t` PLACEMENT POLICY = `ww`
-- case
CREATE DATABASE `t` PLACEMENT POLICY = `ww`
-- case
-- error: line 1 column 31 near "primary_region="us";" 
-- case
-- error: line 1 column 24 near "regions="us,3";" 
-- case
-- error: line 1 column 26 near "followers=3;" 
-- case
-- error: line 1 column 26 near "followers=0;" 
-- case
-- error: line 1 column 23 near "voters=3;" 
-- case
-- error: line 1 column 25 near "learners=3;" 
-- case
-- error: line 1 column 25 near "schedule="even";" 
-- case
-- error: line 1 column 28 near "constraints="ww";" 
-- case
-- error: line 1 column 35 near "leader_constraints="ww";" 
-- case
-- error: line 1 column 37 near "follower_constraints="ww";" 
-- case
-- error: line 1 column 34 near "voter_constraints="ww";" 
-- case
-- error: line 1 column 36 near "learner_constraints="ww";" 
-- case
-- error: line 1 column 47 near "primary_region="us" */;" 
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `ww`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `ww`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
ALTER DATABASE `t` PLACEMENT POLICY = `DEFAULT`
-- case
-- error: line 1 column 97 near "primary_region="us");" 
-- case
-- error: line 1 column 90 near "regions="us,3");" 
-- case
-- error: line 1 column 92 near "followers=3);" 
-- case
-- error: line 1 column 89 near "voters=3);" 
-- case
-- error: line 1 column 91 near "learners=3);" 
-- case
-- error: line 1 column 91 near "schedule="even");" 
-- case
-- error: line 1 column 94 near "constraints="ww");" 
-- case
-- error: line 1 column 101 near "leader_constraints="ww");" 
-- case
-- error: line 1 column 103 near "follower_constraints="ww");" 
-- case
-- error: line 1 column 100 near "voter_constraints="ww");" 
-- case
-- error: line 1 column 102 near "learner_constraints="ww");" 
-- case
CREATE TABLE `m` (`c` INT) PARTITION BY RANGE (`c`) (PARTITION `p1` VALUES LESS THAN (200) PLACEMENT POLICY = `ww`)
-- case
CREATE TABLE `m` (`c` INT) PARTITION BY RANGE (`c`) (PARTITION `p1` VALUES LESS THAN (200) PLACEMENT POLICY = `ww`)
-- case
-- error: line 1 column 40 near "primary_region="us";" 
-- case
-- error: line 1 column 33 near "regions="us,3";" 
-- case
-- error: line 1 column 35 near "followers=3;" 
-- case
-- error: line 1 column 40 near "primary_region="us" followers=3;" 
-- case
-- error: line 1 column 32 near "voters=3;" 
-- case
-- error: line 1 column 34 near "learners=3;" 
-- case
-- error: line 1 column 34 near "schedule="even";" 
-- case
-- error: line 1 column 37 near "constraints="ww";" 
-- case
-- error: line 1 column 44 near "leader_constraints="ww";" 
-- case
-- error: line 1 column 46 near "follower_constraints="ww";" 
-- case
-- error: line 1 column 43 near "voter_constraints="ww";" 
-- case
-- error: line 1 column 45 near "learner_constraints="ww";" 
-- case
ALTER TABLE `m` PARTITION `t` PLACEMENT POLICY = `ww`
-- case
ALTER TABLE `m` PARTITION `t` PLACEMENT POLICY = `ww`
-- case
-- error: line 1 column 79 near "primary_region="us");" 
-- case
-- error: line 1 column 72 near "regions="us,3");" 
-- case
-- error: line 1 column 74 near "followers=3);" 
-- case
-- error: line 1 column 71 near "voters=3);" 
-- case
-- error: line 1 column 73 near "learners=3);" 
-- case
-- error: line 1 column 73 near "schedule="even");" 
-- case
-- error: line 1 column 76 near "constraints="ww");" 
-- case
-- error: line 1 column 83 near "leader_constraints="ww");" 
-- case
-- error: line 1 column 85 near "follower_constraints="ww");" 
-- case
-- error: line 1 column 82 near "voter_constraints="ww");" 
-- case
-- error: line 1 column 84 near "learner_constraints="ww");" 
-- case
ALTER TABLE `m` ADD PARTITION (PARTITION `p1` VALUES LESS THAN (200) PLACEMENT POLICY = `ww`)
-- case
ALTER TABLE `m` ADD PARTITION (PARTITION `p1` VALUES LESS THAN (200) PLACEMENT POLICY = `ww`)
-- case
ALTER TABLE `m` ADD COLUMN `a` INT, ADD PARTITION (PARTITION `p1` VALUES LESS THAN (200))
-- case
ALTER TABLE `m` ADD PARTITION (PARTITION `p1` VALUES LESS THAN (200)), ADD COLUMN `a` INT
-- case
CREATE TABLE `t` (`c1` TINYINT(1),`c2` TINYINT(1),CHECK(`c1` IN (0,1)) NOT ENFORCED,CHECK(`c2` IN (0,1)) ENFORCED)
-- case
CREATE TABLE `Customer` (`SD` INT CHECK(`SD`>0) ENFORCED,`First_Name` VARCHAR(30))
-- case
CREATE TABLE `Customer` (`SD` INT CHECK(`SD`>0) NOT ENFORCED,`SS` VARCHAR(30) CHECK(`ss`=_UTF8MB4'test') ENFORCED)
-- case
CREATE TABLE `Customer` (`SD` INT CHECK(`SD`>0) ENFORCED NOT NULL,`First_Name` VARCHAR(30) COMMENT 'string' NOT NULL)
-- case
CREATE TABLE `Customer` (`SD` INT COMMENT 'string' CHECK(`SD`>0) ENFORCED NOT NULL)
-- case
-- error: line 1 column 63 near "enforced, First_Name varchar(30));" 
-- case
-- error: line 1 column 46 near "enforced, First_Name varchar(30));" 
-- case
CREATE DATABASE `xxx`
-- case
-- error: line 1 column 25 near "exists xxx" 
-- case
CREATE DATABASE IF NOT EXISTS `xxx`
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'N'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'N'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'N'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'N'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'Y'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'Y'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'Y'
-- case
CREATE DATABASE `xxx` ENCRYPTION = 'Y'
-- case
-- error: line 1 column 34 near "N" 
-- case
CREATE DATABASE `xxx`
-- case
-- error: line 1 column 23 near "exists xxx" 
-- case
CREATE DATABASE IF NOT EXISTS `xxx`
-- case
DROP DATABASE `xxx`
-- case
DROP DATABASE IF EXISTS `xxx`
-- case
-- error: line 1 column 20 near "not exists xxx" 
-- case
DROP DATABASE `xxx`
-- case
DROP DATABASE IF EXISTS `xxx`
-- case
-- error: line 1 column 18 near "not exists xxx" 
-- case
-- error: line 1 column 10 near "" 
-- case
-- error: line 1 column 26 near "'xyz" 
-- case
-- error: line 1 column 23 near "'" 
-- case
-- error: line 1 column 23 near "`" 
-- case
-- error: line 1 column 23 near "'" 
-- case
-- error: line 1 column 23 near """ 
-- case
DROP TABLE `xxx`
-- case
DROP TABLE `xxx`, `yyy`
-- case
DROP TABLE `xxx`
-- case
DROP TABLE `xxx`, `yyy`
-- case
DROP TABLE IF EXISTS `xxx`
-- case
DROP TABLE IF EXISTS `xxx`, `yyy`
-- case
-- error: line 1 column 17 near "not exists xxx" 
-- case
DROP TABLE `xxx`
-- case
DROP TABLE `xxx`, `yyy`
-- case
DROP TABLE IF EXISTS `xxx`
-- case
-- error: line 1 column 9 near "" 
-- case
DROP VIEW `xxx`
-- case
DROP VIEW `xxx`, `yyy`
-- case
DROP VIEW IF EXISTS `xxx`
-- case
DROP VIEW IF EXISTS `xxx`, `yyy`
-- case
DROP STATS `t`
-- case
DROP STATS `t1`, `t2`, `t3`
-- case
DROP STATS `t` GLOBAL
-- case
DROP STATS `t` PARTITION `p0`
-- case
DROP STATS `t` PARTITION `p0`,`p1`,`p2`
-- case
CREATE TABLE `address` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`create_at` DATETIME NOT NULL,`deleted` TINYINT(1) NOT NULL,`update_at` DATETIME NOT NULL,`version` BIGINT(20) DEFAULT NULL,`address` VARCHAR(128) NOT NULL,`address_detail` VARCHAR(128) NOT NULL,`cellphone` VARCHAR(16) NOT NULL,`latitude` DOUBLE NOT NULL,`longitude` DOUBLE NOT NULL,`name` VARCHAR(16) NOT NULL,`sex` TINYINT(1) NOT NULL,`user_id` BIGINT(20) NOT NULL,PRIMARY KEY(`id`),CONSTRAINT `FK_7rod8a71yep5vxasb0ms3osbg` FOREIGN KEY (`user_id`) REFERENCES `waimaiqa`.`user`(`id`),INDEX `FK_7rod8a71yep5vxasb0ms3osbg`(`user_id`)) ENGINE = InnoDB AUTO_INCREMENT = 30 DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_GENERAL_CI ROW_FORMAT = COMPACT COMMENT = '' CHECKSUM = 0 DELAY_KEY_WRITE = 0
-- case
CREATE TABLE `test_data` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`create_at` DATETIME NOT NULL,`deleted` TINYINT(1) NOT NULL,`update_at` DATETIME NOT NULL,`version` BIGINT(20) DEFAULT NULL,`address` VARCHAR(255) NOT NULL,`amount` DECIMAL(19,2) DEFAULT NULL,`charge_id` VARCHAR(32) DEFAULT NULL,`paid_amount` DECIMAL(19,2) DEFAULT NULL,`transaction_no` VARCHAR(64) DEFAULT NULL,`wx_mp_app_id` VARCHAR(32) DEFAULT NULL,`contacts` VARCHAR(50) DEFAULT NULL,`deliver_fee` DECIMAL(19,2) DEFAULT NULL,`deliver_info` VARCHAR(255) DEFAULT NULL,`deliver_time` VARCHAR(255) DEFAULT NULL,`description` VARCHAR(255) DEFAULT NULL,`invoice` VARCHAR(255) DEFAULT NULL,`order_from` INT(11) DEFAULT NULL,`order_state` INT(11) NOT NULL,`packing_fee` DECIMAL(19,2) DEFAULT NULL,`payment_time` DATETIME DEFAULT NULL,`payment_type` INT(11) DEFAULT NULL,`phone` VARCHAR(50) NOT NULL,`store_employee_id` BIGINT(20) DEFAULT NULL,`store_id` BIGINT(20) NOT NULL,`user_id` BIGINT(20) NOT NULL,`payment_mode` INT(11) NOT NULL,`current_latitude` DOUBLE NOT NULL,`current_longitude` DOUBLE NOT NULL,`address_latitude` DOUBLE NOT NULL,`address_longitude` DOUBLE NOT NULL,PRIMARY KEY(`id`),CONSTRAINT `food_order_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `waimaiqa`.`user`(`id`),CONSTRAINT `food_order_ibfk_2` FOREIGN KEY (`store_id`) REFERENCES `waimaiqa`.`store`(`id`),CONSTRAINT `food_order_ibfk_3` FOREIGN KEY (`store_employee_id`) REFERENCES `waimaiqa`.`store_employee`(`id`),UNIQUE `FK_UNIQUE_charge_id`(`charge_id`) USING BTREE,INDEX `FK_eqst2x1xisn3o0wbrlahnnqq8`(`store_employee_id`) USING BTREE,INDEX `FK_8jcmec4kb03f4dod0uqwm54o9`(`store_id`) USING BTREE,INDEX `FK_a3t0m9apja9jmrn60uab30pqd`(`user_id`) USING BTREE) ENGINE = InnoDB AUTO_INCREMENT = 95 DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_GENERAL_CI ROW_FORMAT = COMPACT COMMENT = '' CHECKSUM = 0 DELAY_KEY_WRITE = 0
-- case
CREATE TABLE `t` (`c` INT PRIMARY KEY)
-- case
CREATE TABLE `address` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`create_at` DATETIME NOT NULL,`deleted` TINYINT(1) NOT NULL,`update_at` DATETIME NOT NULL,`version` BIGINT(20) DEFAULT NULL,`address` VARCHAR(128) NOT NULL,`address_detail` VARCHAR(128) NOT NULL,`cellphone` VARCHAR(16) NOT NULL,`latitude` DOUBLE NOT NULL,`longitude` DOUBLE NOT NULL,`name` VARCHAR(16) NOT NULL,`sex` TINYINT(1) NOT NULL,`user_id` BIGINT(20) NOT NULL,PRIMARY KEY(`id`),CONSTRAINT `FK_7rod8a71yep5vxasb0ms3osbg` FOREIGN KEY (`user_id`) REFERENCES `waimaiqa`.`user`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,INDEX `FK_7rod8a71yep5vxasb0ms3osbg`(`user_id`)) ENGINE = InnoDB AUTO_INCREMENT = 30 DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_GENERAL_CI ROW_FORMAT = COMPACT COMMENT = '' CHECKSUM = 0 DELAY_KEY_WRITE = 0
-- case
CREATE TABLE `address` (`id` BIGINT(20) NOT NULL AUTO_INCREMENT,`create_at` DATETIME NOT NULL,`deleted` TINYINT(1) NOT NULL,`update_at` DATETIME NOT NULL,`version` BIGINT(20) DEFAULT NULL,`address` VARCHAR(128) NOT NULL,`address_detail` VARCHAR(128) NOT NULL,`cellphone` VARCHAR(16) NOT NULL,`latitude` DOUBLE NOT NULL,`longitude` DOUBLE NOT NULL,`name` VARCHAR(16) NOT NULL,`sex` TINYINT(1) NOT NULL,`user_id` BIGINT(20) NOT NULL,PRIMARY KEY(`id`),CONSTRAINT `FK_7rod8a71yep5vxasb0ms3osbg` FOREIGN KEY (`user_id`) REFERENCES `waimaiqa`.`user`(`id`) ON DELETE CASCADE ON UPDATE NO ACTION,INDEX `FK_7rod8a71yep5vxasb0ms3osbg`(`user_id`)) ENGINE = InnoDB AUTO_INCREMENT = 30 DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_GENERAL_CI ROW_FORMAT = COMPACT COMMENT = '' CHECKSUM = 0 DELAY_KEY_WRITE = 0
-- case
CREATE TABLE `t1` (`accout_id` INT(11) DEFAULT _UTF8MB4'0',`summoner_id` INT(11) DEFAULT _UTF8MB4'0',`union_name` VARBINARY(52) NOT NULL,`union_id` INT(11) DEFAULT _UTF8MB4'0',PRIMARY KEY(`union_name`)) ENGINE = MyISAM DEFAULT CHARACTER SET = BINARY
-- case
CREATE TABLE `t` (`a` DECIMAL(20,0),`b` DECIMAL(30),`c` FLOAT(25,0))
-- case
CREATE TABLE `t` (`c` INT,INDEX `ci`(`c`) USING BTREE COMMENT '123')
-- case
CREATE TABLE `sbtest` (`id` INT UNSIGNED NOT NULL AUTO_INCREMENT,`k` INT UNSIGNED DEFAULT _UTF8MB4'0' NOT NULL,`c` CHAR(120) DEFAULT _UTF8MB4'' NOT NULL,`pad` CHAR(60) DEFAULT _UTF8MB4'' NOT NULL,PRIMARY KEY(`id`))
-- case
CREATE TABLE `test` (`create_date` TIMESTAMP NOT NULL COMMENT '创建日期 create date' DEFAULT CURRENT_TIMESTAMP())
-- case
CREATE TABLE `ts` (`t` INT,`v` TIMESTAMP(3) DEFAULT CURRENT_TIMESTAMP(3))
-- case
CREATE TABLE IF NOT EXISTS `t` (`id` INT NOT NULL AUTO_INCREMENT COMMENT '消息ID',PRIMARY KEY `pk_id`(`id`))
-- case
CREATE TABLE `a` LIKE `b`
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) ON DELETE NO ACTION)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) ON UPDATE SET DEFAULT)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) ON DELETE SET NULL ON UPDATE CASCADE)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) ON DELETE RESTRICT ON UPDATE SET DEFAULT)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) MATCH FULL ON DELETE NO ACTION)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) MATCH PARTIAL ON UPDATE NO ACTION)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) MATCH SIMPLE ON UPDATE NO ACTION)
-- case
CREATE TABLE `a` (`id` INT REFERENCES `a`(`id`) ON UPDATE SET DEFAULT)
-- case
-- error: line 1 column 72 near "update CURRENT_TIMESTAMP)" 
-- case
-- error: line 1 column 90 near "CURRENT_TIMESTAMP)" 
-- case
CREATE TABLE `a` LIKE `b`
-- case
CREATE TABLE IF NOT EXISTS `a` LIKE `b`
-- case
CREATE TABLE IF NOT EXISTS `a` LIKE `b`
-- case
-- error: line 1 column 35 near "(b)" 
-- case
-- error: line 1 column 27 near "like b" 
-- case
-- error: line 1 column 27 near "like (b)" 
-- case
CREATE TABLE `a` AS SELECT * FROM `b`
-- case
CREATE TABLE `a` AS SELECT * FROM `b`
-- case
CREATE TABLE `a` (`m` INT,`n` DATETIME) AS SELECT * FROM `b`
-- case
CREATE TABLE `a` (UNIQUE(`n`)) AS SELECT `n` FROM `b`
-- case
CREATE TABLE `a` IGNORE AS SELECT `n` FROM `b`
-- case
CREATE TABLE `a` REPLACE AS SELECT `n` FROM `b`
-- case
CREATE TABLE `a` (`m` INT) REPLACE AS (SELECT `n` AS `m` FROM `b` UNION SELECT `n`+1 AS `m` FROM `c` GROUP BY 1 LIMIT 2)
-- case
CREATE TABLE `a`
-- case
-- error: line 1 column 39 near "now)" 
-- case
CREATE TABLE `t` (`a` TIMESTAMP DEFAULT CURRENT_TIMESTAMP())
-- case
CREATE TABLE `t` (`a` TIMESTAMP DEFAULT CURRENT_TIMESTAMP())
-- case
-- error: line 1 column 55 near "now)" 
-- case
CREATE TABLE `t` (`a` TIMESTAMP DEFAULT CURRENT_TIMESTAMP() ON UPDATE CURRENT_TIMESTAMP())
-- case
-- error: line 1 column 53 near "(now()))" 
-- case
CREATE TABLE `t` (`c` TEXT) DEFAULT CHARACTER SET = UTF8 DEFAULT COLLATE = UTF8_GENERAL_CI
-- case
CREATE TABLE `t` (`c` TEXT) SHARD_ROW_ID_BITS = 1
-- case
CREATE TABLE `t` (`c` TEXT) SHARD_ROW_ID_BITS = 1 PRE_SPLIT_REGIONS = 1
-- case
CREATE TABLE IF NOT EXISTS `general_log` (`event_time` TIMESTAMP(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),`user_host` MEDIUMTEXT NOT NULL,`thread_id` BIGINT(20) UNSIGNED NOT NULL,`server_id` INT(10) UNSIGNED NOT NULL,`command_type` VARCHAR(64) NOT NULL,`argument` MEDIUMBLOB NOT NULL) ENGINE = CSV DEFAULT CHARACTER SET = UTF8 COMMENT = 'General log'
-- case
CREATE TABLE `followers` (`f1` INT NOT NULL REFERENCES `user_profiles`(`uid`))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND()))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND(1)))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND()))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND(1)))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND()))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (RAND(1)))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (CURRENT_DATE()))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%Y-%m')))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%Y-%m')))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%Y-%m-%d')))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%Y-%m-%d %H.%i.%s')))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%Y-%m-%d %H:%i:%s')))
-- case
CREATE TABLE `t` (`d` DATE DEFAULT (DATE_FORMAT(NOW(), _UTF8MB4'%b %d %Y %h:%i %p')))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (REPLACE(UPPER(UUID()), _UTF8MB4'-', _UTF8MB4'')))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (REPLACE(UPPER(UUID()), _UTF8MB4'-', _UTF8MB4'')))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (REPLACE(CONVERT(UPPER(UUID()) USING 'utf8mb4'), _UTF8MB4'-', _UTF8MB4'')))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (REPLACE(CONVERT(UPPER(UUID()) USING 'utf8mb4'), _UTF8MB4'-', _UTF8MB4'')))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (UPPER(SUBSTRING_INDEX(USER(), _UTF8MB4'@', 1))))
-- case
CREATE TABLE `t` (`a` INT DEFAULT (UPPER(SUBSTRING_INDEX(USER(), _UTF8MB4'@', 1))))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (STR_TO_DATE(_UTF8MB4'1980-01-01', _UTF8MB4'%Y-%m-%d')))
-- case
CREATE TABLE `t` (`a` VARCHAR(32) DEFAULT (STR_TO_DATE(_UTF8MB4'1980-01-01', _UTF8MB4'%Y-%m-%d')))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_OBJECT()))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_ARRAY()))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_QUOTE()))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_OBJECT(_UTF8MB4'foo', 5, _UTF8MB4'bar', _UTF8MB4'barfoo')))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_ARRAY(1, 2, 3)))
-- case
CREATE TABLE `t` (`j` JSON DEFAULT (JSON_QUOTE(_UTF8MB4'foobar')))
-- case
CREATE TABLE `t` (`c` CHAR(33) DEFAULT (NONEXISTINGFUNC(_UTF8MB4'foobar')))
-- case
CREATE TABLE `t` (`c` CHAR(33) DEFAULT _UTF8MB4'foobar')
-- case
CREATE TABLE `t` (`c` CHAR(33) DEFAULT _UTF8MB4'foobar')
-- case
CREATE TABLE `t` (`i` INT DEFAULT 0)
-- case
CREATE TABLE `t` (`i` INT DEFAULT -1)
-- case
CREATE TABLE `t` (`i` INT DEFAULT +1)
-- case
CREATE TABLE `t` (`a` INT,`b` INT,`c` CHAR(33) DEFAULT (`b`))
-- case
-- error: line 1 column 53 near ")" 
-- case
CREATE TABLE `t` (`a` INT) ENCRYPTION = 'n'
-- case
CREATE TABLE `t` (`a` INT) ENCRYPTION = 'n'
-- case
ALTER TABLE `t` ENCRYPTION = 'y'
-- case
ALTER TABLE `t` ENCRYPTION = 'y'
-- case
CREATE TABLE `t` (`a` INT PRIMARY KEY GLOBAL)
-- case
CREATE TABLE `t` (`a` INT PRIMARY KEY)
-- case
CREATE TABLE `t` (`a` INT PRIMARY KEY)
-- case
CREATE TABLE `t` (`a` INT PRIMARY KEY GLOBAL)
-- case
CREATE TABLE `t` (`a` INT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` INT UNIQUE KEY GLOBAL)
-- case
CREATE TABLE `t` (`a` INT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` INT UNIQUE KEY GLOBAL)
-- case
ALTER TABLE `t` ADD INDEX(`a`)
-- case
ALTER TABLE `t` ADD INDEX(`a`)
-- case
ALTER TABLE `t` ADD INDEX(`a`) GLOBAL
-- case
ALTER TABLE `t` ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) GLOBAL
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) GLOBAL
-- case
ALTER TABLE `t` ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`) GLOBAL
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`)
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`)
-- case
CREATE INDEX `i` ON `t` (`a`)
-- case
CREATE INDEX `i` ON `t` (`a`)
-- case
CREATE INDEX `i` ON `t` (`a`) GLOBAL
-- case
ALTER DATABASE `t` CHARACTER SET = utf8
-- case
ALTER DATABASE CHARACTER SET = utf8
-- case
ALTER DATABASE `t` CHARACTER SET = utf8
-- case
ALTER DATABASE `t` CHARACTER SET = utf8
-- case
ALTER DATABASE CHARACTER SET = utf8
-- case
ALTER DATABASE `t` CHARACTER SET = utf8
-- case
ALTER DATABASE `t` COLLATE = binary
-- case
ALTER DATABASE `t` CHARACTER SET = binary COLLATE = binary
-- case
ALTER DATABASE `t` COLLATE = utf8_bin
-- case
ALTER DATABASE COLLATE = utf8_bin
-- case
ALTER DATABASE `t` COLLATE = utf8_bin
-- case
ALTER DATABASE `t` COLLATE = utf8_bin
-- case
ALTER DATABASE COLLATE = utf8_bin
-- case
ALTER DATABASE `` COLLATE = utf8_bin
-- case
ALTER DATABASE `t` CHARACTER SET = utf8mb4 COLLATE = utf8_bin
-- case
ALTER DATABASE `t` CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci CHARACTER SET = utf8 COLLATE = utf8mb4_bin
-- case
ALTER DATABASE CHARACTER SET = utf8mb4 COLLATE = utf8_bin
-- case
ALTER DATABASE CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci CHARACTER SET = utf8 COLLATE = utf8mb4_bin
-- case
ALTER TABLE `t` ADD COLUMN (`a` SMALLINT UNSIGNED)
-- case
-- error: line 1 column 15 near "* ADD COLUMN (a SMALLINT UNSIGNED)" 
-- case
ALTER TABLE `t` ADD COLUMN IF NOT EXISTS (`a` SMALLINT UNSIGNED)
-- case
-- error: line 1 column 15 near "ADD COLUMN (a SMALLINT UNSIGNED)" 
-- case
ALTER TABLE `t` ADD COLUMN (`a` SMALLINT UNSIGNED, `b` VARCHAR(255))
-- case
ALTER TABLE `t` ADD COLUMN IF NOT EXISTS (`a` SMALLINT UNSIGNED, `b` VARCHAR(255))
-- case
-- error: line 1 column 51 near "FIRST)" 
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED FIRST
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED AFTER `b`
-- case
ALTER TABLE `t` ADD COLUMN IF NOT EXISTS `a` SMALLINT UNSIGNED AFTER `b`
-- case
ALTER TABLE `employees` ADD PARTITION
-- case
ALTER TABLE `employees` ADD PARTITION (PARTITION `P1` VALUES LESS THAN (2010))
-- case
ALTER TABLE `employees` ADD PARTITION (PARTITION `P2` VALUES LESS THAN (MAXVALUE))
-- case
ALTER TABLE `employees` ADD PARTITION IF NOT EXISTS (PARTITION `P2` VALUES LESS THAN (MAXVALUE))
-- case
ALTER TABLE `employees` ADD PARTITION IF NOT EXISTS PARTITIONS 5
-- case
ALTER TABLE `employees` ADD PARTITION (PARTITION `P1` VALUES LESS THAN (2010), PARTITION `P2` VALUES LESS THAN (2015), PARTITION `P3` VALUES LESS THAN (MAXVALUE))
-- case
ALTER TABLE `t` ADD PARTITION (PARTITION `x` VALUES IN ((3, 4), (5, 6)))
-- case
ALTER TABLE `employees` ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `employees` ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10
-- case
ALTER TABLE `employees` ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `employees` ADD PARTITION NO_WRITE_TO_BINLOG PARTITIONS 10
-- case
ALTER TABLE `t_n` REBUILD PARTITION ALL
-- case
ALTER TABLE `d_n`.`t_n` REBUILD PARTITION NO_WRITE_TO_BINLOG ALL
-- case
ALTER TABLE `t_n` REBUILD PARTITION NO_WRITE_TO_BINLOG `ident`
-- case
ALTER TABLE `t_n` REBUILD PARTITION NO_WRITE_TO_BINLOG `ident`,`ident`
-- case
-- error: line 1 column 39 near "" 
-- case
ALTER TABLE `t_n` REBUILD PARTITION NO_WRITE_TO_BINLOG `local`
-- case
ALTER TABLE `t_n` REBUILD PARTITION NO_WRITE_TO_BINLOG `local`,`local`
-- case
ALTER TABLE `t` DROP PARTITION `p1`
-- case
ALTER TABLE `t` DROP PARTITION `p2`
-- case
ALTER TABLE `t` DROP PARTITION IF EXISTS `p2`
-- case
ALTER TABLE `t` DROP PARTITION `p1`,`p2`
-- case
ALTER TABLE `t` DROP PARTITION IF EXISTS `p1`,`p2`
-- case
ALTER TABLE `t` CHECK PARTITION ALL
-- case
ALTER TABLE `t` CHECK PARTITION `p`
-- case
ALTER TABLE `t` CHECK PARTITION `p1`,`p2`
-- case
ALTER TABLE `employees` ADD PARTITION PARTITIONS 1
-- case
ALTER TABLE `employees` ADD PARTITION PARTITIONS 2
-- case
ALTER TABLE `clients` COALESCE PARTITION 3
-- case
ALTER TABLE `clients` COALESCE PARTITION 4
-- case
ALTER TABLE `clients` COALESCE PARTITION NO_WRITE_TO_BINLOG 4
-- case
ALTER TABLE `clients` COALESCE PARTITION NO_WRITE_TO_BINLOG 4
-- case
ALTER TABLE `t` DISABLE KEYS
-- case
ALTER TABLE `t` ENABLE KEYS
-- case
ALTER TABLE `t` MODIFY COLUMN `a` VARCHAR(255)
-- case
ALTER TABLE `t` MODIFY COLUMN IF EXISTS `a` VARCHAR(255)
-- case
ALTER TABLE `t` CHANGE COLUMN `a` `b` VARCHAR(255)
-- case
ALTER TABLE `t` CHANGE COLUMN IF EXISTS `a` `b` VARCHAR(255)
-- case
ALTER TABLE `t` CHANGE COLUMN `a` `b` VARCHAR(255) BINARY CHARACTER SET UTF8
-- case
ALTER TABLE `t` CHANGE COLUMN `a` `b` VARCHAR(255) FIRST
-- case
ALTER TABLE `db`.`t` RENAME AS `db1`.`t1`
-- case
ALTER TABLE `db`.`t` RENAME AS `db1`.`t1`
-- case
ALTER TABLE `db`.`t` RENAME AS `db1`.`t1`
-- case
ALTER TABLE `db`.`t` RENAME AS `db1`.`t1`
-- case
ALTER TABLE `t` RENAME AS `t1`
-- case
ALTER TABLE `t` RENAME AS `t1`
-- case
ALTER TABLE `t` RENAME AS `t1`
-- case
ALTER TABLE `t` RENAME AS `t1`
-- case
ALTER TABLE `t_n` ORDER BY `ident`
-- case
ALTER TABLE `t_n` ORDER BY `ident`
-- case
ALTER TABLE `t_n` ORDER BY `ident` DESC
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2` DESC
-- case
ALTER TABLE `t_n` ORDER BY `ident1` DESC,`ident2`
-- case
ALTER TABLE `t_n` ORDER BY `ident1` DESC,`ident2`
-- case
ALTER TABLE `t_n` ORDER BY `ident1` DESC,`ident2` DESC
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`,`ident3`
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`,`ident3`
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`,`ident3` DESC
-- case
ALTER TABLE `t_n` ORDER BY `ident1`,`ident2`,`ident3`
-- case
ALTER TABLE `t_n` ORDER BY `ident1` DESC,`ident2` DESC,`ident3` DESC
-- case
ALTER TABLE `t` RENAME COLUMN `a` TO `b`
-- case
-- error: line 1 column 30 near ".a TO t.b" 
-- case
-- error: line 1 column 35 near ".b" 
-- case
-- error: line 1 column 30 near ".a TO b" 
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT 1
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT 1
-- case
-- error: line 1 column 58 near "CURRENT_TIMESTAMP" 
-- case
-- error: line 1 column 44 near "NOW()" 
-- case
-- error: line 1 column 43 near "+1" 
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT (CURRENT_TIMESTAMP())
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT (NOW())
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT (1+1)
-- case
ALTER TABLE `t` ALTER COLUMN `a` SET DEFAULT 1
-- case
ALTER TABLE `t` ALTER COLUMN `a` DROP DEFAULT
-- case
ALTER TABLE `t` ALTER COLUMN `a` DROP DEFAULT
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = NONE
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = DEFAULT
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = SHARED
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = EXCLUSIVE
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = NONE
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = DEFAULT
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = SHARED
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = EXCLUSIVE
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = NONE
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = DEFAULT
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = SHARED
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, LOCK = EXCLUSIVE
-- case
ALTER TABLE `t` ADD FULLTEXT `FullText`(`name`)
-- case
ALTER TABLE `t` ADD FULLTEXT `FullText`(`name`)
-- case
ALTER TABLE `t` ADD FULLTEXT `FullText`(`name`)
-- case
ALTER TABLE `t` ADD INDEX(`a`) USING BTREE COMMENT 'a'
-- case
ALTER TABLE `t` ADD INDEX IF NOT EXISTS(`a`) USING BTREE COMMENT 'a'
-- case
ALTER TABLE `t` ADD INDEX(`a`) USING RTREE COMMENT 'a'
-- case
ALTER TABLE `t` ADD INDEX(`a`) PRE_SPLIT_REGIONS = 4
-- case
ALTER TABLE `t` ADD INDEX(`a`) PRE_SPLIT_REGIONS = 4
-- case
-- error: line 1 column 51 near "'a'" 
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`) CLUSTERED PRE_SPLIT_REGIONS = 4
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`) NONCLUSTERED PRE_SPLIT_REGIONS = 4
-- case
ALTER TABLE `t` ADD INDEX(`a`) PRE_SPLIT_REGIONS = (BETWEEN (1,_UTF8MB4'a') AND (2,_UTF8MB4'b') REGIONS 4)
-- case
ALTER TABLE `t` ADD INDEX(`a`) PRE_SPLIT_REGIONS = (BY (1,_UTF8MB4'a'),(2,_UTF8MB4'b'),(3,_UTF8MB4'c'))
-- case
ALTER TABLE `t` ADD INDEX(`a`) COMMENT 'a' PRE_SPLIT_REGIONS = (BETWEEN (1,_UTF8MB4'a') AND (2,_UTF8MB4'b') REGIONS 4)
-- case
CREATE INDEX `idx` ON `t` (`a`, `b`) PRE_SPLIT_REGIONS = 100
-- case
CREATE INDEX `idx` ON `t` (`a`, `b`) PRE_SPLIT_REGIONS = (BETWEEN (1,_UTF8MB4'a') AND (2,_UTF8MB4'b') REGIONS 4)
-- case
ALTER TABLE `t` ADD INDEX `idx`(`a`) PRE_SPLIT_REGIONS = 100, ADD INDEX `idx2`(`b`) PRE_SPLIT_REGIONS = (BY (1),(2),(3))
-- case
ALTER TABLE `t` ADD INDEX(`a`) USING HASH COMMENT 'a'
-- case
ALTER TABLE `t` ADD INDEX(`a`) USING BTREE COMMENT 'a' GLOBAL
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) GLOBAL
-- case
ALTER TABLE `t` ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD INDEX IF NOT EXISTS(`a`) USING HASH COMMENT 'a'
-- case
ALTER TABLE `t` ADD PRIMARY KEY `ident`(`a` DESC, `b`) USING RTREE
-- case
ALTER TABLE `t` ADD INDEX(`a`) USING RTREE
-- case
ALTER TABLE `t` ADD INDEX(`ident`, `ident`(123)) USING RTREE
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`) COMMENT 'a'
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) COMMENT 'a'
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) COMMENT 'a'
-- case
ALTER TABLE `t` ADD UNIQUE(`a`) COMMENT 'a'
-- case
-- error: line 1 column 26 near "(a) USING HNSW COMMENT 'a'" 
-- case
-- error: line 1 column 26 near "(VEC_COSINE_DISTANCE(a)) USING HNSW COMMENT 'a'" 
-- case
-- error: line 1 column 26 near "((VEC_COSINE_DISTANCE(a))) USING HASH COMMENT 'a'" 
-- case
-- error: line 1 column 26 near "((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'" 
-- case
-- error: line 1 column 28 near "KEY ((VEC_COSINE_DISTANCE(a, b))) USING HNSW COMMENT 'a'" 
-- case
ALTER TABLE `t` ADD VECTOR INDEX((VEC_COSINE_DISTANCE(`a`, `b`))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX((LOWER(`a`))) USING HNSW COMMENT 'a'
-- case
-- error: line 1 column 56 near ", a)) USING HNSW COMMENT 'a'" 
-- case
ALTER TABLE `t` ADD VECTOR INDEX(`a`, (VEC_COSINE_DISTANCE(`a`))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX((VEC_COSINE_DISTANCE(`a`))) USING HYPO COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX((VEC_COSINE_DISTANCE(`a`))) COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX((VEC_COSINE_DISTANCE(`a`))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX IF NOT EXISTS((VEC_COSINE_DISTANCE(`a`))) USING HNSW COMMENT 'a'
-- case
ALTER TABLE `t` ADD VECTOR INDEX IF NOT EXISTS((VEC_COSINE_DISTANCE(`a`))) ADD_COLUMNAR_REPLICA_ON_DEMAND USING HNSW COMMENT 'a'
-- case
-- error: line 1 column 28 near "(a) USING INVERTED COMMENT 'a'" 
-- case
-- error: line 1 column 28 near "((a - 1)) USING INVERTED COMMENT 'a'" 
-- case
-- error: line 1 column 28 near "(a) USING HASH COMMENT 'a'" 
-- case
-- error: line 1 column 28 near "(a, b) USING INVERTED COMMENT 'a'" 
-- case
-- error: line 1 column 30 near "KEY (a, b) USING INVERTED COMMENT 'a'" 
-- case
ALTER TABLE `t` ADD COLUMNAR INDEX(`a`, `b`) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE `t` ADD COLUMNAR INDEX(`a`) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE `t` ADD COLUMNAR INDEX(`a`) USING HYPO COMMENT 'a'
-- case
ALTER TABLE `t` ADD COLUMNAR INDEX((`a`-1)) USING HYPO COMMENT 'a'
-- case
ALTER TABLE `t` ADD COLUMNAR INDEX IF NOT EXISTS(`a`) USING INVERTED COMMENT 'a'
-- case
ALTER TABLE `t` ADD CONSTRAINT `fk_t2_id` FOREIGN KEY (`t2_id`) REFERENCES `t`(`id`)
-- case
ALTER TABLE `t` ADD CONSTRAINT `fk_t2_id` FOREIGN KEY IF NOT EXISTS (`t2_id`) REFERENCES `t`(`id`)
-- case
ALTER TABLE `t` ADD CONSTRAINT `c_1` CHECK(1+1) NOT ENFORCED, ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD CONSTRAINT `c_1` CHECK(1+1) ENFORCED, ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ADD CONSTRAINT `c_1` CHECK(1+1) ENFORCED, ADD UNIQUE(`a`)
-- case
ALTER TABLE `t` ENGINE = ''
-- case
ALTER TABLE `t` ENGINE = ''
-- case
ALTER TABLE `t` ENGINE = innodb
-- case
ALTER TABLE `t` ENGINE = innodb
-- case
ALTER TABLE `db`.`t` ENGINE = ''
-- case
ALTER TABLE `t` INSERT_METHOD = FIRST
-- case
ALTER TABLE `t` INSERT_METHOD = LAST
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT UNSIGNED, ADD COLUMN `a` SMALLINT
-- case
ALTER TABLE `t` ADD COLUMN `a` SMALLINT, ENGINE = '', DEFAULT COLLATE = UTF8_GENERAL_CI
-- case
ALTER TABLE `t` ENGINE = '', COMMENT = '', DEFAULT COLLATE = UTF8_GENERAL_CI
-- case
ALTER TABLE `t` ENGINE = '', ADD COLUMN `a` SMALLINT
-- case
ALTER TABLE `t` DEFAULT COLLATE = UTF8_GENERAL_CI, ENGINE = '', ADD COLUMN `a` SMALLINT
-- case
ALTER TABLE `t` SHARD_ROW_ID_BITS = 1
-- case
ALTER TABLE `t` AUTO_INCREMENT = 3
-- case
ALTER TABLE `t` AUTO_INCREMENT = 3
-- case
ALTER TABLE `t` FORCE AUTO_INCREMENT = 3
-- case
ALTER TABLE `t` FORCE AUTO_INCREMENT = 3
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` MEDIUMTEXT CHARACTER SET UTF8MB4 COLLATE utf8mb4_unicode_ci NOT NULL, ALGORITHM = DEFAULT
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` MEDIUMTEXT CHARACTER SET UTF8MB4 COLLATE utf8mb4_unicode_ci NOT NULL, ALGORITHM = INPLACE
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` MEDIUMTEXT CHARACTER SET UTF8MB4 COLLATE utf8mb4_unicode_ci NOT NULL, ALGORITHM = COPY
-- case
ALTER TABLE `hello-world@dev`.`User` ADD COLUMN `name` MEDIUMTEXT CHARACTER SET UTF8MB4 COLLATE utf8mb4_unicode_ci NOT NULL, ALGORITHM = INSTANT
-- case
ALTER TABLE `t` CONVERT TO CHARACTER SET UTF8
-- case
ALTER TABLE `t` CONVERT TO CHARACTER SET UTF8
-- case
ALTER TABLE `t` CONVERT TO CHARACTER SET UTF8 COLLATE UTF8_BIN
-- case
ALTER TABLE `t` CONVERT TO CHARACTER SET UTF8 COLLATE UTF8_BIN
-- case
ALTER TABLE `d_n`.`t_n` CONVERT TO CHARACTER SET DEFAULT
-- case
ALTER TABLE `d_n`.`t_n` CONVERT TO CHARACTER SET DEFAULT
-- case
ALTER TABLE `d_n`.`t_n` CONVERT TO CHARACTER SET DEFAULT
-- case
ALTER TABLE `d_n`.`t_n` CONVERT TO CHARACTER SET DEFAULT COLLATE UTF8MB4_0900_AI_CI
-- case
ALTER TABLE `t` FORCE /* AlterTableForce is not supported */ 
-- case
-- error: line 1 column 25 near ";" 
-- case
ALTER TABLE `t` DROP INDEX `a`
-- case
ALTER TABLE `t` DROP INDEX IF EXISTS `a`
-- case
ALTER TABLE `t` ALTER INDEX `a` INVISIBLE
-- case
ALTER TABLE `t` ALTER INDEX `a` VISIBLE
-- case
ALTER TABLE `t` DROP FOREIGN KEY `a`
-- case
ALTER TABLE `t` DROP COLUMN `a`
-- case
ALTER TABLE `t` DROP COLUMN IF EXISTS `a`
-- case
ALTER TABLE `testTableCompression` COMPRESSION = 'LZ4'
-- case
ALTER TABLE `t1` COMPRESSION = 'zlib'
-- case
ALTER TABLE `t1`
-- case
-- error: line 1 column 16 near "," 
-- case
ALTER TABLE `t` RENAME INDEX `a` TO `b`
-- case
ALTER TABLE `t` RENAME INDEX `a` TO `b`
-- case
ALTER TABLE `d_n`.`t_n` DROP CHECK `ident`
-- case
ALTER TABLE `t_n` LOCK = DEFAULT, DROP CHECK `ident`
-- case
ALTER TABLE `t_n` ALTER CHECK `ident` ENFORCED
-- case
ALTER TABLE `t_n` ALTER CHECK `ident` NOT ENFORCED
-- case
ALTER TABLE `t_n` DROP CHECK `ident`
-- case
ALTER TABLE `t_n` DROP CHECK `ident`
-- case
-- error: line 1 column 38 near "" 
-- case
ALTER TABLE `t_n` ALTER CHECK `ident` ENFORCED
-- case
ALTER TABLE `t_n` ALTER CHECK `ident` NOT ENFORCED
-- case
ANALYZE TABLE `t` PARTITION `a`
-- case
ANALYZE TABLE `t` PARTITION `a` WITH 4 BUCKETS
-- case
ANALYZE TABLE `t` PARTITION `a` INDEX `b`
-- case
ANALYZE TABLE `t` PARTITION `a` INDEX `b` WITH 4 BUCKETS
-- case
ALTER TABLE `t` PARTITION BY HASH (`a`) PARTITIONS 1
-- case
ALTER TABLE `t` ADD COLUMN `a` INT PARTITION BY HASH (`a`) PARTITIONS 1
-- case
ALTER TABLE `t` ADD COLUMN `a` INT PARTITION BY HASH (`a`) PARTITIONS 1 UPDATE INDEXES (`idx_a` GLOBAL)
-- case
ALTER TABLE `t` ADD COLUMN `a` INT PARTITION BY HASH (`a`) PARTITIONS 1 UPDATE INDEXES (`idx_a` GLOBAL,`idx_b` LOCAL)
-- case
-- error: line 1 column 80 near "normal)" 
-- case
-- error: line 1 column 75 near ")" 
-- case
-- error: [ddl:1492]For RANGE partitions each partition must be defined
-- case
-- error: [ddl:1492]For RANGE partitions each partition must be defined
-- case
ALTER TABLE `t` PARTITION BY RANGE (`a`) (PARTITION `x` VALUES LESS THAN (75))
-- case
-- error: line 1 column 41 near "partition by range(a) (partition x values less than (75))" 
-- case
ALTER TABLE `t` COMMENT = 'cmt' PARTITION BY HASH (`a`) PARTITIONS 1
-- case
ALTER TABLE `t` ENABLE KEYS, COMMENT = 'cmt' PARTITION BY HASH (`a`) PARTITIONS 1
-- case
-- error: line 1 column 53 near "partition by hash(a)" 
-- case
-- error: line 1 column 41 near "enable keys" 
-- case
-- error: line 1 column 35 near ", enable keys" 
-- case
ALTER TABLE `t` PARTITION BY RANGE COLUMNS (`a`) (PARTITION `x` VALUES LESS THAN (MAXVALUE))
-- case
ALTER TABLE `t` PARTITION BY LIST COLUMNS (`a`) (PARTITION `p0` VALUES IN (5, 10, 15))
-- case
ALTER TABLE `t` PARTITION BY RANGE COLUMNS (`a`,`b`,`c`) (PARTITION `p1` VALUES LESS THAN (1, 1, 1))
-- case
ALTER TABLE `t` PARTITION BY LIST COLUMNS (`a`,`b`,`c`) (PARTITION `p0` VALUES IN ((5, 10, 15)))
-- case
ALTER TABLE `t` WITH VALIDATION, ADD COLUMN `b` INT GENERATED ALWAYS AS(`a`+1) VIRTUAL
-- case
ALTER TABLE `t` WITHOUT VALIDATION, ADD COLUMN `b` INT GENERATED ALWAYS AS(`a`+1) VIRTUAL
-- case
ALTER TABLE `t` WITHOUT VALIDATION, WITH VALIDATION, ADD COLUMN `b` INT GENERATED ALWAYS AS(`a`+1) VIRTUAL
-- case
ALTER TABLE `t` WITH VALIDATION, MODIFY COLUMN `b` INT GENERATED ALWAYS AS(`a`+2) VIRTUAL
-- case
ALTER TABLE `t` WITH VALIDATION, CHANGE COLUMN `b` `c` INT GENERATED ALWAYS AS(`a`+2) VIRTUAL
-- case
ALTER TABLE `d_n`.`t_n` ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `d_n`.`t_n` ADD PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `t` WITH VALIDATION, EXCHANGE PARTITION `p` WITH TABLE `nt` WITHOUT VALIDATION
-- case
ALTER TABLE `t` EXCHANGE PARTITION `p` WITH TABLE `nt`
-- case
ALTER TABLE `t` REORGANIZE PARTITION
-- case
ALTER TABLE `t` REORGANIZE PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `t` REORGANIZE PARTITION NO_WRITE_TO_BINLOG
-- case
ALTER TABLE `members` REORGANIZE PARTITION `n0` INTO (PARTITION `s0` VALUES LESS THAN (1960), PARTITION `s1` VALUES LESS THAN (1970))
-- case
ALTER TABLE `members` REORGANIZE PARTITION NO_WRITE_TO_BINLOG `n0` INTO (PARTITION `s0` VALUES LESS THAN (1960), PARTITION `s1` VALUES LESS THAN (1970))
-- case
ALTER TABLE `members` REORGANIZE PARTITION `p1`,`p2`,`p3` INTO (PARTITION `s0` VALUES LESS THAN (1960), PARTITION `s1` VALUES LESS THAN (1970))
-- case
-- error: line 1 column 51 near "partition;" 
-- case
ALTER TABLE `t` REORGANIZE PARTITION NO_WRITE_TO_BINLOG `remove` INTO (PARTITION `p0` VALUES LESS THAN (1991))
-- case
ALTER TABLE `t` ATTRIBUTES='str'
-- case
ALTER TABLE `t` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` ATTRIBUTES=DEFAULT
-- case
ALTER TABLE `t` ATTRIBUTES=DEFAULT
-- case
ALTER TABLE `t` ATTRIBUTES=DEFAULT
-- case
-- error: line 1 column 24 near "" 
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES='str'
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES='str1,str2'
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES=DEFAULT
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES=DEFAULT
-- case
ALTER TABLE `t` PARTITION `p` ATTRIBUTES=DEFAULT
-- case
-- error: line 1 column 36 near "" 
-- case
CREATE TABLE `t1` (`attributes` INT)
-- case
CREATE INDEX `idx` ON `t` (`a`)
-- case
CREATE INDEX IF NOT EXISTS `idx` ON `t` (`a`)
-- case
CREATE UNIQUE INDEX `idx` ON `t` (`a`)
-- case
CREATE UNIQUE INDEX IF NOT EXISTS `idx` ON `t` (`a`)
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING BTREE
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING HASH
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING RTREE
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING BTREE
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING BTREE
-- case
CREATE SPATIAL INDEX `idx` ON `t` (`a`)
-- case
CREATE SPATIAL INDEX IF NOT EXISTS `idx` ON `t` (`a`)
-- case
CREATE FULLTEXT INDEX `idx` ON `t` (`a`)
-- case
CREATE FULLTEXT INDEX IF NOT EXISTS `idx` ON `t` (`a`)
-- case
CREATE FULLTEXT INDEX `idx` ON `t` (`a`) WITH PARSER `ident`
-- case
CREATE FULLTEXT INDEX `idx` ON `t` (`a`) WITH PARSER `ident` COMMENT 'string'
-- case
CREATE FULLTEXT INDEX `idx` ON `t` (`a`) WITH PARSER `ident` COMMENT 'string'
-- case
CREATE FULLTEXT INDEX `idx` ON `t` (`a`) WITH PARSER `ident` COMMENT 'string'
-- case
CREATE INDEX `idx` ON `t` (`a`) USING HASH
-- case
CREATE INDEX `idx` ON `t` (`a`) COMMENT 'foo'
-- case
CREATE INDEX `idx` ON `t` (`a`) USING HASH COMMENT 'foo'
-- case
CREATE INDEX `idx` ON `t` (`a`) LOCK = NONE
-- case
CREATE INDEX `idx` ON `t` (`a`) USING HASH COMMENT 'foo'
-- case
CREATE INDEX `idx` ON `t` (`a`) USING BTREE
-- case
CREATE INDEX `idx` ON `t` (`a`) VISIBLE
-- case
CREATE INDEX `idx` ON `t` (`a`) INVISIBLE
-- case
CREATE INDEX `idx` ON `t` (`a`) VISIBLE
-- case
CREATE INDEX `idx` ON `t` (`a`) INVISIBLE
-- case
CREATE INDEX `idx` ON `t` (`a`) USING HASH VISIBLE
-- case
CREATE INDEX `idx` ON `t` (`a`) USING HASH INVISIBLE
-- case
CREATE VECTOR INDEX `idx` ON `t` (`a`) USING HNSW
-- case
CREATE VECTOR INDEX `idx` ON `t` (`a`, `b`) USING HNSW
-- case
CREATE VECTOR INDEX `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`)))
-- case
CREATE VECTOR INDEX `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`))) USING BTREE
-- case
-- error: line 1 column 34 near "USING HNSW ((VEC_COSINE_DISTANCE(a)))" 
-- case
-- error: line 1 column 17 near "idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW" 
-- case
CREATE VECTOR INDEX `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`)), `a`) USING HNSW
-- case
CREATE VECTOR INDEX `idx` ON `t` (`a`, (VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
-- error: line 1 column 17 near "KEY idx ON t ((VEC_COSINE_DISTANCE(a))) USING HNSW" 
-- case
CREATE VECTOR INDEX `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
CREATE VECTOR INDEX IF NOT EXISTS `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
CREATE VECTOR INDEX IF NOT EXISTS `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
CREATE VECTOR INDEX `ident` ON `d_n`.`t_n` ((VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
CREATE VECTOR INDEX `idx` ON `t` ((VEC_COSINE_DISTANCE(`a`))) USING HNSW
-- case
CREATE VECTOR INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING HNSW
-- case
CREATE UNIQUE INDEX `ident` ON `d_n`.`t_n` (`ident`, `ident`) USING HNSW
-- case
CREATE INDEX `idx` ON `t` (`a`)
-- case
CREATE INDEX `idx` ON `t` (`a`)
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = INPLACE
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = INPLACE
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = COPY
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = COPY
-- case
CREATE INDEX `idx` ON `t` (`a`)
-- case
CREATE INDEX `idx` ON `t` (`a`)
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = INPLACE LOCK = EXCLUSIVE
-- case
CREATE INDEX `idx` ON `t` (`a`) ALGORITHM = INPLACE LOCK = EXCLUSIVE
-- case
-- error: [parser:1800]Unknown ALGORITHM 'ALGORITHM'
-- case
-- error: [parser:1800]Unknown ALGORITHM 'ALGORITHM'
-- case
DROP INDEX `a` ON `t`
-- case
DROP INDEX `a` ON `db`.`t`
-- case
DROP INDEX `a` ON `db`.`tb-ttb`
-- case
DROP INDEX IF EXISTS `a` ON `t`
-- case
DROP INDEX IF EXISTS `a` ON `db`.`t`
-- case
DROP INDEX IF EXISTS `a` ON `db`.`tb-ttb`
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t` ALGORITHM = INPLACE
-- case
DROP INDEX `idx` ON `t` ALGORITHM = INPLACE
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t` LOCK = SHARED
-- case
DROP INDEX `idx` ON `t` LOCK = SHARED
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t`
-- case
DROP INDEX `idx` ON `t` ALGORITHM = INPLACE LOCK = EXCLUSIVE
-- case
DROP INDEX `idx` ON `t` ALGORITHM = INPLACE LOCK = EXCLUSIVE
-- case
-- error: [parser:1800]Unknown ALGORITHM 'algorithm'
-- case
-- error: [parser:1800]Unknown ALGORITHM 'algorithm'
-- case
-- error: [parser:1801]Unknown LOCK type 'lock_type'
-- case
-- error: [parser:1801]Unknown LOCK type 'lock_type'
-- case
RENAME TABLE `t` TO `t1`
-- case
-- error: line 1 column 17 near "t1" 
-- case
RENAME TABLE `d`.`t` TO `d1`.`t1`
-- case
RENAME TABLE `t1` TO `t2`, `t3` TO `t4`
-- case
TRUNCATE TABLE `t1`
-- case
TRUNCATE TABLE `t1`
-- case
-- error: line 1 column 26 near ") " 
-- case
-- error: line 1 column 27 near ")" 
-- case
-- error: line 1 column 33 near ")" 
-- case
-- error: line 1 column 31 near ")" 
-- case
ALTER TABLE `d_n`.`t_n` SECONDARY_LOAD
-- case
ALTER TABLE `d_n`.`t_n` SECONDARY_UNLOAD
-- case
ALTER TABLE `t_n` LOCK = DEFAULT, SECONDARY_LOAD
-- case
ALTER TABLE `d_n`.`t_n` ALGORITHM = DEFAULT, SECONDARY_LOAD
-- case
ALTER TABLE `d_n`.`t_n` ALGORITHM = DEFAULT, SECONDARY_UNLOAD
-- case
CREATE TABLE `a` (`process` DOUBLE)
-- case
CREATE TABLE `t` (`a` TINYINT,`b` SMALLINT,`c` MEDIUMINT,`d` INT,`e` BIGINT)
-- case
CREATE TABLE `t` (`lv` MEDIUMTEXT NULL)
-- case
CREATE TABLE `cdp_test`.`test2-1` (`id` INT(11) DEFAULT NULL,INDEX(`id`))
-- case
CREATE TABLE `miantiao` (`扁豆焖面` INT(11))
-- case
CREATE TABLE `bar` (`m` INT) AS SELECT `n` FROM `foo`
-- case
CREATE TABLE `bar` (`m` INT) IGNORE AS SELECT `n` FROM `foo`
-- case
CREATE TABLE `bar` (`m` INT) REPLACE AS SELECT `n` FROM `foo`
-- case
-- error: [ddl:1221]Incorrect usage of ON UPDATE and generated column
-- case
-- error: [ddl:1221]Incorrect usage of AUTO_INCREMENT and generated column
-- case
-- error: [ddl:1221]Incorrect usage of DEFAULT and generated column
-- case
CREATE TABLE `t` (`a` BIGINT,`b` BIGINT GENERATED ALWAYS AS(`a`+1) VIRTUAL NOT NULL)
-- case
CREATE TABLE `t` (`a` BIGINT,`b` BIGINT GENERATED ALWAYS AS(`a`+1) VIRTUAL NOT NULL)
-- case
CREATE TABLE `t` (`a` BIGINT,`b` BIGINT GENERATED ALWAYS AS(`a`+1) VIRTUAL NOT NULL COMMENT 'ttt')
-- case
CREATE TABLE `t` (`a` INT,INDEX `idx`((CAST(`a` AS BINARY(1)))))
-- case
-- error: [ddl:1221]Incorrect usage of DEFAULT and generated column
-- case
-- error: [ddl:1221]Incorrect usage of DEFAULT and generated column
-- case
CREATE TABLE `t` (`a` INT COLUMN_FORMAT FIXED)
-- case
CREATE TABLE `t` (`a` INT COLUMN_FORMAT DEFAULT)
-- case
CREATE TABLE `t` (`a` INT COLUMN_FORMAT DYNAMIC)
-- case
ALTER TABLE `t` MODIFY COLUMN `a` BIGINT COLUMN_FORMAT DEFAULT
-- case
RECOVER TABLE BY JOB 11
-- case
-- error: line 1 column 24 near ",12,13" 
-- case
-- error: line 1 column 20 near "" 
-- case
RECOVER TABLE BY JOB 0
-- case
RECOVER TABLE `t1`
-- case
-- error: line 1 column 17 near ",t2" 
-- case
-- error: line 1 column 14 near "" 
-- case
RECOVER TABLE `t1` 100
-- case
-- error: line 1 column 20 near "abc" 
-- case
FLASHBACK TABLE `t`
-- case
FLASHBACK TABLE `t` TO `t1`
-- case
FLASHBACK TABLE `t` TO `timestamp`
-- case
FLASHBACK DATABASE `db1`
-- case
FLASHBACK DATABASE `db1`
-- case
FLASHBACK DATABASE `db1` TO `db2`
-- case
FLASHBACK DATABASE `db1` TO `db2`
-- case
FLASHBACK CLUSTER TO TIMESTAMP '2021-05-26 16:45:26'
-- case
FLASHBACK TABLE `t` TO TIMESTAMP '2021-05-26 16:45:26'
-- case
FLASHBACK TABLE `t`, `t1` TO TIMESTAMP '2021-05-26 16:45:26'
-- case
FLASHBACK DATABASE `test` TO TIMESTAMP '2021-05-26 16:45:26'
-- case
FLASHBACK DATABASE `test` TO TIMESTAMP '2021-05-26 16:45:26'
-- case
-- error: line 1 column 20 near "to timestamp TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())" 
-- case
-- error: line 1 column 20 near "to timestamp DATE_SUB(NOW(), INTERVAL 3 SECOND)" 
-- case
-- error: line 1 column 28 near "timestamp '2021-05-26 16:45:26'" 
-- case
-- error: line 1 column 31 near "timestamp '2021-05-26 16:45:26'" 
-- case
FLASHBACK CLUSTER TO TSO 445494955052105721
-- case
FLASHBACK TABLE `t` TO TSO 445494955052105722
-- case
FLASHBACK TABLE `t`, `t1` TO TSO 445494955052105723
-- case
FLASHBACK DATABASE `test` TO TSO 445494955052105724
-- case
FLASHBACK DATABASE `test` TO TSO 445494955052105725
-- case
-- error: line 1 column 22 near "tso 445494955052105726" 
-- case
-- error: line 1 column 25 near "tso 445494955052105727" 
-- case
-- error: line 1 column 30 near ""Invalid TSO value provided: 0 
-- case
-- error: line 1 column 30 near "-100" 
-- case
ALTER TABLE `t` REMOVE PARTITIONING
-- case
ALTER TABLE `db`.`ident` REMOVE PARTITIONING
-- case
ALTER TABLE `t` LOCK = DEFAULT REMOVE PARTITIONING
-- case
ALTER TABLE `t` ADD COLUMN `a` INT REMOVE PARTITIONING
-- case
ALTER TABLE `t` ADD COLUMN `a` INT, ADD INDEX(`c`) REMOVE PARTITIONING
-- case
-- error: line 1 column 38 near "remove partitioning" 
-- case
-- error: line 1 column 53 near "remove partitioning" 
-- case
-- error: line 1 column 37 near "add column a int" 
-- case
-- error: line 1 column 34 near ", add column a int" 
-- case
ALTER TABLE `t` ADD COLUMN `a` DOUBLE(4,2) UNSIGNED ZEROFILL REFERENCES `b` MATCH FULL ON UPDATE SET NULL FIRST
-- case
ALTER TABLE `d_n`.`t_n` ADD CONSTRAINT `ident` FOREIGN KEY (`ident`(1)) REFERENCES `d_n`.`t_n` MATCH FULL ON DELETE SET NULL
-- case
ALTER TABLE `t_n` ADD CONSTRAINT `ident` FOREIGN KEY (`ident`, `ident`(1)) REFERENCES `t_n` MATCH FULL ON DELETE RESTRICT ON UPDATE SET NULL
-- case
ALTER TABLE `d_n`.`t_n` ADD CONSTRAINT `ident` FOREIGN KEY (`ident`, `ident`(1)) REFERENCES `t_n` MATCH PARTIAL ON DELETE CASCADE REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` ADD CONSTRAINT FOREIGN KEY (`ident`) REFERENCES `d_n`.`t_n` MATCH SIMPLE ON DELETE CASCADE ON UPDATE CASCADE
-- case
CREATE TABLE `t` (`a` VARCHAR(1))
-- case
CREATE TABLE `t` (`a` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(1))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(1),`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` CHAR)
-- case
CREATE TABLE `t` (`a` CHAR)
-- case
CREATE TABLE `t` (`a` VARCHAR(50),`b` INT)
-- case
CREATE TABLE `t` (`a` CHAR,`b` INT)
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` CHAR)
-- case
CREATE TABLE `t` (`a` CHAR)
-- case
CREATE TABLE `t` (`a` CHAR)
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
-- error: line 1 column 35 near ");" 
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `t` (`a` VARCHAR(50))
-- case
CREATE TABLE `nchar` (`a` INT)
-- case
CREATE TABLE `nchar` (`a` INT,`b` CHAR)
-- case
CREATE TABLE `nchar` (`a` INT,`b` CHAR(50))
-- case
ALTER TABLE `t_n` STORAGE DISK, MODIFY COLUMN `ident` VARCHAR(12) COLUMN_FORMAT FIXED FIRST
-- case
CREATE TABLE `t` (`a` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE KEY NULL)
-- case
CREATE TABLE `t` (`b` INT,`a` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` INT NOT NULL AUTO_INCREMENT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` INT NOT NULL AUTO_INCREMENT UNIQUE KEY NULL)
-- case
CREATE TABLE `t` (`a` BIGINT NOT NULL AUTO_INCREMENT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` SMALLINT NOT NULL AUTO_INCREMENT UNIQUE KEY)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT,`b` MEDIUMTEXT,`c` MEDIUMTEXT,`d` MEDIUMTEXT,`e` MEDIUMTEXT,`f` MEDIUMTEXT,`g` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`a` MEDIUMBLOB)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT,`b` MEDIUMBLOB)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET UTF8)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET UTF8)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET UTF8)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET UTF8)
-- case
ALTER TABLE `d_n`.`t_n` MODIFY COLUMN `ident` MEDIUMTEXT AFTER `ident` REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` MODIFY COLUMN `ident` MEDIUMTEXT AFTER `ident` REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` MODIFY COLUMN `ident` MEDIUMTEXT AFTER `ident` REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` MODIFY COLUMN `ident` MEDIUMTEXT AFTER `ident` REMOVE PARTITIONING
-- case
ALTER TABLE `d_n`.`t_n` MODIFY COLUMN `ident` MEDIUMTEXT AFTER `ident` REMOVE PARTITIONING
-- case
ALTER TABLE `t_n` CHANGE COLUMN `ident` `ident` MEDIUMTEXT BINARY CHARACTER SET UTF8 FIRST, TABLESPACE = `ident`
-- case
ALTER TABLE `t_n` CHANGE COLUMN `ident` `ident` MEDIUMTEXT BINARY CHARACTER SET UTF8 FIRST, TABLESPACE = `ident`
-- case
-- error: line 1 column 43 near ";"The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
-- case
-- error: line 1 column 46 near ";"The value of STATS_AUTO_RECALC must be one of [0|1|DEFAULT]. 
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = 0
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = DEFAULT
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = 0
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = 1
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = DEFAULT
-- case
CREATE TABLE `t` (`a` INT) STATS_PERSISTENT = DEFAULT /* TableOptionStatsPersistent is not supported */  STATS_AUTO_RECALC = 1
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = 1 STATS_SAMPLE_PAGES = 25
-- case
ALTER TABLE `t` MODIFY COLUMN `a` BIGINT, ENGINE = InnoDB, STATS_AUTO_RECALC = 0
-- case
CREATE TABLE `stats_auto_recalc` (`a` INT)
-- case
CREATE TABLE `stats_auto_recalc` (`a` INT) STATS_AUTO_RECALC = 1
-- case
CREATE TABLE `t` (`a` INT,PRIMARY KEY `type`(`a`) USING BTREE)
-- case
-- error: line 1 column 45 near "btree (a));" 
-- case
CREATE TABLE `t` (`a` INT,PRIMARY KEY(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,PRIMARY KEY(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,PRIMARY KEY(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,UNIQUE `type`(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,UNIQUE `type`(`a`) USING BTREE)
-- case
-- error: line 1 column 46 near "btree (a));" 
-- case
CREATE TABLE `t` (`a` INT,UNIQUE(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,UNIQUE(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,UNIQUE(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,INDEX `type`(`a`) USING BTREE)
-- case
-- error: line 1 column 39 near "btree (a));" 
-- case
CREATE TABLE `t` (`a` INT,INDEX `type`(`a`) USING BTREE)
-- case
CREATE TABLE `t` (`a` INT,INDEX(`a`) USING BTREE)
-- case
ALTER TABLE `d_n`.`t_n` WITHOUT VALIDATION, ADD PARTITION (PARTITION `ident` VALUES LESS THAN (MAXVALUE) ENGINE = text_string MAX_ROWS = 12)
-- case
ALTER TABLE `d_n`.`t_n` WITH VALIDATION, ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION `ident` VALUES LESS THAN (MAXVALUE) ENGINE = text_string, PARTITION `ident` VALUES LESS THAN (MAXVALUE) (SUBPARTITION `text_string` MIN_ROWS = 11))
-- case
ALTER TABLE `d_n`.`t_n` WITHOUT VALIDATION, ADD PARTITION (PARTITION `ident` DEFAULT ENGINE = text_string MAX_ROWS = 12)
-- case
ALTER TABLE `d_n`.`t_n` WITH VALIDATION, ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION `ident` DEFAULT ENGINE = text_string MAX_ROWS = 12)
-- case
ALTER TABLE `d_n`.`t_n` ADD PARTITION (PARTITION `ident` DEFAULT, PARTITION `ptext` VALUES IN (_UTF8MB4'default'))
-- case
ALTER TABLE `d_n`.`t_n` WITH VALIDATION, ADD PARTITION NO_WRITE_TO_BINLOG (PARTITION `ident` VALUES LESS THAN (MAXVALUE) ENGINE = text_string, PARTITION `ident` DEFAULT (SUBPARTITION `text_string` MIN_ROWS = 11))
-- case
ALTER TABLE `d_n`.`t_n` ADD PARTITION (PARTITION `ident` DEFAULT)
-- case
ALTER TABLE `d_n`.`t_n` ADD PARTITION (PARTITION `ident` VALUES IN (1, DEFAULT))
-- case
ALTER TABLE `t` IMPORT TABLESPACE
-- case
ALTER TABLE `t` DISCARD TABLESPACE
-- case
ALTER TABLE `db`.`t` IMPORT TABLESPACE
-- case
ALTER TABLE `db`.`t` DISCARD TABLESPACE
-- case
ALTER TABLE `t` ADD COLUMN (CHECK(TRUE) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (CHECK(TRUE) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (CONSTRAINT `ident` CHECK(1>2) NOT ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (`b` INT, UNIQUE `c`(`b`))
-- case
ALTER TABLE `t` ADD COLUMN (CHECK(TRUE) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (CHECK(TRUE) ENFORCED, CHECK(TRUE) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, CONSTRAINT `b1` CHECK(`a1`>0) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, `a2` INT, CONSTRAINT `b1` CHECK(`a1`>0) ENFORCED, CONSTRAINT `b2` CHECK(`a2`<10) ENFORCED)
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY(`a1`))
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, PRIMARY KEY(`a1`))
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, `a2` INT, PRIMARY KEY(`a1`), UNIQUE(`a2`))
-- case
ALTER TABLE `t` ADD COLUMN (`a1` INT, `a2` INT, PRIMARY KEY(`a1`), UNIQUE `b2`(`a2`))
-- case
ALTER TABLE `ident` ADD COLUMN (CONSTRAINT `ident` FOREIGN KEY (`EXECUTE`(123)) REFERENCES `t`(`a`) MATCH SIMPLE ON DELETE CASCADE ON UPDATE SET NULL)
-- case
ALTER TABLE `t` ADD COLUMN `a` DATE CHECK(`a`>0) ENFORCED FIRST
-- case
ALTER TABLE `t` ADD COLUMN `a1` INT CONSTRAINT `ident` CHECK(`a1`>1) ENFORCED REFERENCES `b` ON DELETE CASCADE ON UPDATE CASCADE
-- case
ALTER TABLE `t` ADD COLUMN `a` DATE CHECK(`a`>0) ENFORCED FIRST
-- case
ALTER TABLE `t` ADD COLUMN `a` TINYBLOB CONSTRAINT `ident` CHECK(1>2) ENFORCED REFERENCES `b` ON DELETE CASCADE ON UPDATE CASCADE
-- case
ALTER TABLE `t` ADD COLUMN `a2` INT CONSTRAINT `ident` CHECK(`a2`>1) ENFORCED
-- case
ALTER TABLE `t` ADD COLUMN `a2` INT CONSTRAINT `ident` CHECK(`a2`>1) NOT ENFORCED
-- case
-- error: line 1 column 49 near "primary key REFERENCES b ON DELETE CASCADE ON UPDATE CASCADE;" 
-- case
-- error: line 1 column 49 near "primary key (a2))" 
-- case
-- error: line 1 column 48 near "unique key (a2))" 
-- case
ALTER TABLE `t` SET TIFLASH REPLICA 2 LOCATION LABELS 'a', 'b'
-- case
ALTER TABLE `t` SET TIFLASH REPLICA 0
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 2 LOCATION LABELS 'a', 'b'
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 0
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 1 SET TIFLASH REPLICA 2 LOCATION LABELS 'a', 'b'
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 1 SET TIFLASH REPLICA 2
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 1 LOCATION LABELS 'a', 'b' SET TIFLASH REPLICA 2
-- case
ALTER DATABASE `t` SET TIFLASH REPLICA 1 LOCATION LABELS 'a', 'b' SET TIFLASH REPLICA 2 LOCATION LABELS 'a', 'b'
-- case
CREATE TABLE IF NOT EXISTS `table_ident` (`a` YEAR(4),`b` YEAR)
-- case
CREATE TABLE IF NOT EXISTS `table_ident` (`ident1` TINYINT(1) COMMENT 'text_string' UNIQUE KEY,`ident2` YEAR(4))
-- case
CREATE TABLE `t` (`y` YEAR(4),`y1` YEAR)
-- case
CREATE TABLE `t` (`y` YEAR(4),`y1` YEAR)
-- case
INSERT INTO `t` SET `a`=DEFAULT
-- case
INSERT INTO `t` SET `a`=DEFAULT
-- case
REPLACE INTO `t` SET `a`=DEFAULT
-- case
UPDATE `t` SET `a`=DEFAULT
-- case
INSERT INTO `t` SET `a`=DEFAULT ON DUPLICATE KEY UPDATE `a`=DEFAULT
-- case
-- error: line 1 column 33 near "ascii)" 
-- case
-- error: line 1 column 35 near "charset latin1)" 
-- case
CREATE TABLE `t` (`a` LONGTEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` TINYTEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` BLOB)
-- case
CREATE TABLE `t` (`a` MEDIUMBLOB,`b` TEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` TEXT CHARACTER SET LATIN1,`b` MEDIUMTEXT CHARACTER SET LATIN1,`c` INT)
-- case
CREATE TABLE `t` (`a` INT,`b` TEXT CHARACTER SET LATIN1,`c` MEDIUMTEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET LATIN1,`b` MEDIUMTEXT CHARACTER SET LATIN1)
-- case
CREATE TABLE `t` (`a` MEDIUMTEXT CHARACTER SET UTF8MB4,`b` MEDIUMTEXT CHARACTER SET UTF8MB4,`c` MEDIUMTEXT CHARACTER SET UTF8MB4)
-- case
CREATE TABLE `t` (`a` INT STORAGE MEMORY,`b` VARCHAR(255) STORAGE MEMORY)
-- case
CREATE TABLE `t` (`a` INT STORAGE DISK,`b` VARCHAR(255) STORAGE DEFAULT)
-- case
CREATE TABLE `t` (`a` INT STORAGE DEFAULT,`b` VARCHAR(255) STORAGE DISK)
-- case
CREATE TABLE `t` (`a` DECIMAL(6,3),`b` DECIMAL PRIMARY KEY)
-- case
CREATE TABLE `t` (`a` DECIMAL,`b` DECIMAL(6))
-- case
CREATE TABLE `t` (`a` DECIMAL(65,30) UNSIGNED ZEROFILL,`b` DECIMAL,`c` DECIMAL(65) UNSIGNED ZEROFILL)
-- case
-- error: line 1 column 33 near "a)));" 
-- case
-- error: line 1 column 28 near "+1));" 
-- case
-- error: line 1 column 31 near "+1));" 
-- case
CREATE TABLE `a` (`a` INT,`b` INT,INDEX((`a`+1), (`b`+1)))
-- case
CREATE TABLE `a` (`a` INT,`b` INT,INDEX(`a`, (`b`+1)))
-- case
CREATE TABLE `a` (`a` INT,`b` INT,INDEX((`a`+1), `b`))
-- case
CREATE TABLE `a` (`a` INT,`b` INT,INDEX((`a`+1) DESC))
-- case
CREATE SEQUENCE `sequence`
-- case
CREATE SEQUENCE `seq`
-- case
CREATE SEQUENCE IF NOT EXISTS `seq`
-- case
CREATE SEQUENCE `seq`
-- case
CREATE SEQUENCE `seq`
-- case
CREATE SEQUENCE IF NOT EXISTS `seq`
-- case
CREATE SEQUENCE IF NOT EXISTS `seq`
-- case
-- error: line 1 column 43 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` INCREMENT BY 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` INCREMENT BY 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` INCREMENT BY 1
-- case
-- error: line 1 column 42 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` MINVALUE 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` MINVALUE 1
-- case
-- error: line 1 column 36 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NO MINVALUE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NO MINVALUE
-- case
-- error: line 1 column 42 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` MAXVALUE 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` MAXVALUE 1
-- case
-- error: line 1 column 36 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NO MAXVALUE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NO MAXVALUE
-- case
-- error: line 1 column 39 near "" 
-- case
-- error: line 1 column 44 near "" 
-- case
-- error: line 1 column 41 near "" 
-- case
-- error: line 1 column 44 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` START WITH 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` START WITH 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` START WITH 1
-- case
-- error: line 1 column 39 near "" 
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` CACHE 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` CACHE 1
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NOCACHE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NOCACHE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` CYCLE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NOCYCLE
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` NOCYCLE
-- case
CREATE SEQUENCE `seq` INCREMENT BY 1 START WITH 0 MINVALUE 0 MAXVALUE 1000
-- case
CREATE SEQUENCE `seq` INCREMENT BY 1 START WITH 0 MINVALUE 0 MAXVALUE 1000
-- case
CREATE SEQUENCE `seq` INCREMENT BY 10 START WITH 0 MINVALUE 0 MAXVALUE 1000
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` CACHE 1 INCREMENT BY 1 START WITH -1 MINVALUE 0 MAXVALUE 1000
-- case
CREATE SEQUENCE `sEq` START WITH 0 MINVALUE 0 MAXVALUE 1000
-- case
CREATE SEQUENCE IF NOT EXISTS `seq` INCREMENT BY 1 START WITH 0 MINVALUE -2 MAXVALUE 1000
-- case
CREATE SEQUENCE `seq` INCREMENT BY -1 START WITH -1 MINVALUE -1 MAXVALUE -1000 CACHE 10 NOCYCLE
-- case
CREATE TABLE `sequence` (`a` INT)
-- case
CREATE TABLE `t` (`sequence` INT)
-- case
-- error: line 1 column 13 near "" 
-- case
DROP SEQUENCE `seq`
-- case
DROP SEQUENCE IF EXISTS `seq`
-- case
DROP SEQUENCE `seq`
-- case
DROP SEQUENCE IF EXISTS `seq`
-- case
DROP SEQUENCE IF EXISTS `seq`, `seq2`, `seq3`
-- case
-- error: line 1 column 22 near "seq2" 
-- case
DROP SEQUENCE `seq`, `seq2`
-- case
CREATE TABLE `t` (`a` BIGINT AUTO_RANDOM(3) PRIMARY KEY,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` BIGINT AUTO_RANDOM PRIMARY KEY,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` BIGINT PRIMARY KEY AUTO_RANDOM(4),`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` BIGINT PRIMARY KEY AUTO_RANDOM(3) PRIMARY KEY UNIQUE KEY,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` BIGINT AUTO_RANDOM(5, 53) PRIMARY KEY,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` BIGINT AUTO_RANDOM(15, 32) PRIMARY KEY,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` INT) AUTO_ID_CACHE = 1
-- case
CREATE TABLE `t` (`a` INT AUTO_INCREMENT PRIMARY KEY) AUTO_ID_CACHE = 10
-- case
CREATE TABLE `t` (`a` BIGINT,`b` VARCHAR(255)) AUTO_ID_CACHE = 50
-- case
CREATE TABLE `t` (`a` BIGINT AUTO_RANDOM(3) PRIMARY KEY) AUTO_RANDOM_BASE = 10
-- case
CREATE TABLE `t` (`a` BIGINT PRIMARY KEY AUTO_RANDOM(4),`b` VARCHAR(100)) AUTO_RANDOM_BASE = 200
-- case
ALTER TABLE `t` AUTO_RANDOM_BASE = 50
-- case
ALTER TABLE `t` AUTO_INCREMENT = 30, AUTO_RANDOM_BASE = 40
-- case
ALTER TABLE `t` FORCE AUTO_RANDOM_BASE = 50
-- case
ALTER TABLE `t` AUTO_INCREMENT = 30, FORCE AUTO_RANDOM_BASE = 40
-- case
-- error: line 1 column 18 near "" 
-- case
-- error: line 1 column 26 near "comment="haha"" 
-- case
ALTER SEQUENCE `seq` START WITH 1
-- case
ALTER SEQUENCE `seq` START WITH 1 INCREMENT BY 1
-- case
ALTER SEQUENCE `seq` START WITH 1 INCREMENT BY 2 MINVALUE 0 MAXVALUE 100
-- case
ALTER SEQUENCE `seq` INCREMENT BY -1 START WITH -1 MINVALUE -1 MAXVALUE -1000 CACHE 10 NOCYCLE
-- case
ALTER SEQUENCE IF EXISTS `seq2` INCREMENT BY 2
-- case
ALTER SEQUENCE `seq` RESTART
-- case
ALTER SEQUENCE `seq` START WITH 3 RESTART WITH 5
-- case
ALTER SEQUENCE `seq` RESTART WITH 5
-- case
-- error: line 1 column 27 near "restart = 5" 
-- case
CREATE TABLE `t` (`a` INT,INDEX ``(`a`))
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255),PRIMARY KEY(`b`, `a`) CLUSTERED)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255),PRIMARY KEY(`b`, `a`) NONCLUSTERED)
-- case
CREATE TABLE `t` (`a` INT PRIMARY KEY NONCLUSTERED,`b` VARCHAR(255))
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255) PRIMARY KEY CLUSTERED)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255) DEFAULT _UTF8MB4'a' PRIMARY KEY CLUSTERED)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255) PRIMARY KEY NONCLUSTERED,PRIMARY KEY(`b`, `a`) NONCLUSTERED)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255),PRIMARY KEY(`b`, `a`) NONCLUSTERED USING RTREE)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255),PRIMARY KEY(`b`, `a`) NONCLUSTERED USING RTREE)
-- case
CREATE TABLE `t` (`a` INT,`b` VARCHAR(255),PRIMARY KEY(`b`, `a`) CLUSTERED USING RTREE)
-- case
-- error: line 1 column 47 near "clustered primary key)" 
-- case
-- error: line 1 column 72 near "clustered)" 
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`, `b`) CLUSTERED
-- case
ALTER TABLE `t` ADD PRIMARY KEY(`a`, `b`) NONCLUSTERED
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX(`b`) USING HNSW)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX(`a`, `b`) USING HNSW)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((VEC_COSINE_DISTANCE(`b`))))
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((VEC_COSINE_DISTANCE(`b`))) USING HASH)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX(`a`, (VEC_COSINE_DISTANCE(`b`))) USING HNSW)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((VEC_COSINE_DISTANCE(`b`)), `a`) USING HNSW)
-- case
-- error: line 1 column 69 near "b)) USING HNSW);" 
-- case
-- error: line 1 column 45 near "key((VEC_COSINE_DISTANCE(b))) TYPE HNSW);" 
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((`b`+1)) USING HNSW)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((VEC_COSINE_DISTANCE(`a`, `b`))) USING HNSW)
-- case
CREATE TABLE `t` (`a` INT,`b` VECTOR(3),VECTOR INDEX((VEC_COSINE_DISTANCE(`b`))) USING HNSW)
-- case
DROP PLACEMENT POLICY `x`
-- case
-- error: line 1 column 24 near ", y" 
-- case
DROP PLACEMENT POLICY IF EXISTS `x`
-- case
-- error: line 1 column 34 near ", y" 
-- case
SHOW CREATE PLACEMENT POLICY `x`
-- case
-- error: line 1 column 31 near "if exists x" 
-- case
-- error: line 1 column 31 near ", y" 
-- case
SHOW CREATE PLACEMENT POLICY `placement`
-- case
CREATE PLACEMENT POLICY `x` PRIMARY_REGION = 'us'
-- case
-- error: line 1 column 32 near "region='us, 3'" 
-- case
CREATE PLACEMENT POLICY `x` FOLLOWERS = 3
-- case
-- error: line 1 column 37 near ""FOLLOWERS must be positive 
-- case
CREATE PLACEMENT POLICY `x` VOTERS = 3
-- case
CREATE PLACEMENT POLICY `x` LEARNERS = 3
-- case
CREATE PLACEMENT POLICY `x` SCHEDULE = 'even'
-- case
CREATE PLACEMENT POLICY `x` CONSTRAINTS = 'ww'
-- case
CREATE PLACEMENT POLICY `x` LEADER_CONSTRAINTS = 'ww'
-- case
CREATE PLACEMENT POLICY `x` FOLLOWER_CONSTRAINTS = 'ww'
-- case
CREATE PLACEMENT POLICY `x` VOTER_CONSTRAINTS = 'ww'
-- case
CREATE PLACEMENT POLICY `x` LEARNER_CONSTRAINTS = 'ww'
-- case
CREATE PLACEMENT POLICY `x` PRIMARY_REGION = 'cn' REGIONS = 'us' SCHEDULE = 'even'
-- case
CREATE PLACEMENT POLICY `x` PRIMARY_REGION = 'cn' LEADER_CONSTRAINTS = 'ww' LEADER_CONSTRAINTS = 'yy'
-- case
CREATE PLACEMENT POLICY IF NOT EXISTS `x` REGIONS = 'us' FOLLOWER_CONSTRAINTS = 'yy'
-- case
CREATE OR REPLACE PLACEMENT POLICY `x` REGIONS = 'us'
-- case
-- error: line 1 column 35 near "placement policy y" 
-- case
CREATE TABLE `masking` (`id` INT)
-- case
SELECT `masking` FROM `t`
-- case
ALTER TABLE `t` ADD COLUMN `masking` INT
-- case
ALTER TABLE `t` DROP COLUMN `masking`
-- case
ALTER TABLE `t` MODIFY COLUMN `masking` INT
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS `c`
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS `c` ENABLE
-- case
CREATE MASKING POLICY IF NOT EXISTS `p` ON `t` (`c`) AS `c` DISABLE
-- case
CREATE OR REPLACE MASKING POLICY `p` ON `t` (`c`) AS `c`
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS `c`
-- case
-- error: line 1 column 61 near ""'OR REPLACE' and 'IF NOT EXISTS' are mutually exclusive 
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS CASE WHEN CURRENT_USER()=_UTF8MB4'root' THEN `c` ELSE _UTF8MB4'xxx' END ENABLE
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS `c` RESTRICT ON (INSERT_INTO_SELECT, DELETE_SELECT) ENABLE
-- case
CREATE MASKING POLICY `p` ON `t` (`c`) AS CASE WHEN CURRENT_USER() NOT IN (_UTF8MB4'root@%',_UTF8MB4'u@%') THEN _UTF8MB4'x' ELSE `c` END
-- case
ALTER TABLE `t` ADD MASKING POLICY `p` ON (`c`) AS `c`
-- case
ALTER TABLE `t` ADD MASKING POLICY `p` ON (`c`) AS `c`
-- case
ALTER TABLE `t` ADD MASKING POLICY `p` ON (`c`) AS `c` DISABLE
-- case
ALTER TABLE `t` ADD MASKING POLICY `p` ON (`c`) AS `c` RESTRICT ON (UPDATE_SELECT, CTAS) DISABLE
-- case
ALTER TABLE `t` MODIFY MASKING POLICY `p` SET EXPRESSION = CASE WHEN CURRENT_ROLE() IN (_UTF8MB4'r1',_UTF8MB4'r2') THEN `c` ELSE _UTF8MB4'x' END
-- case
-- error: line 1 column 57 near "case when current_role() in ('r1', 'r2') then c else 'x' end" 
-- case
ALTER TABLE `t` MODIFY MASKING POLICY `p` SET RESTRICT ON (INSERT_INTO_SELECT, UPDATE_SELECT, DELETE_SELECT, CTAS)
-- case
ALTER TABLE `t` MODIFY MASKING POLICY `p` SET RESTRICT ON NONE
-- case
ALTER TABLE `t` ENABLE MASKING POLICY `p`
-- case
ALTER TABLE `t` DISABLE MASKING POLICY `p`
-- case
ALTER TABLE `t` DROP MASKING POLICY `p`
-- case
SHOW MASKING POLICIES FOR `t`
-- case
SHOW MASKING POLICIES FOR `t` WHERE `column_name`=_UTF8MB4'c'
-- case
ALTER PLACEMENT POLICY `x` PRIMARY_REGION = 'us'
-- case
-- error: line 1 column 31 near "region='us, 3'" 
-- case
ALTER PLACEMENT POLICY `x` FOLLOWERS = 3
-- case
ALTER PLACEMENT POLICY `x` VOTERS = 3
-- case
ALTER PLACEMENT POLICY `x` LEARNERS = 3
-- case
ALTER PLACEMENT POLICY `x` SCHEDULE = 'even'
-- case
ALTER PLACEMENT POLICY `x` CONSTRAINTS = 'ww'
-- case
ALTER PLACEMENT POLICY `x` LEADER_CONSTRAINTS = 'ww'
-- case
ALTER PLACEMENT POLICY `x` FOLLOWER_CONSTRAINTS = 'ww'
-- case
ALTER PLACEMENT POLICY `x` VOTER_CONSTRAINTS = 'ww'
-- case
ALTER PLACEMENT POLICY `x` LEARNER_CONSTRAINTS = 'ww'
-- case
ALTER PLACEMENT POLICY `x` PRIMARY_REGION = 'cn' REGIONS = 'us' SCHEDULE = 'even'
-- case
ALTER PLACEMENT POLICY `x` PRIMARY_REGION = 'cn' LEADER_CONSTRAINTS = 'ww' LEADER_CONSTRAINTS = 'yy'
-- case
ALTER PLACEMENT POLICY IF EXISTS `x` REGIONS = 'us' FOLLOWER_CONSTRAINTS = 'yy'
-- case
-- error: line 1 column 34 near "placement policy y" 
-- case
-- error: line 1 column 27 near "cpu ='8c'" 
-- case
-- error: line 1 column 30 near "region ='us, 3'" 
-- case
-- error: line 1 column 27 near "cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'" 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 2000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 200000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED
-- case
-- error: line 1 column 42 near "'check'" 
-- case
-- error: line 1 column 33 near "followers=0" 
-- case
-- error: line 1 column 38 near "true" 
-- case
-- error: line 1 column 39 near "false" 
-- case
-- error: line 1 column 41 near "disable" 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 2000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 3000, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 4000
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = UNLIMITED, RU_PER_SEC = 4000
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 4000
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = OFF, RU_PER_SEC = 4000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 20, PRIORITY = LOW, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `default` RU_PER_SEC = 20, PRIORITY = LOW, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `default` RU_PER_SEC = UNLIMITED, PRIORITY = LOW, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = UNLIMITED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = OFF
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = UNLIMITED, RU_PER_SEC = 2000
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = OFF, RU_PER_SEC = 2000
-- case
CREATE RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 2000
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = UNLIMITED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = OFF
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, PRIORITY = LOW, BURSTABLE = UNLIMITED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, PRIORITY = LOW, BURSTABLE = OFF
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, PRIORITY = LOW, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, PRIORITY = LOW, BURSTABLE = UNLIMITED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, PRIORITY = LOW, BURSTABLE = OFF
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, PRIORITY = LOW, BURSTABLE = MODERATED
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = DRYRUN)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10m' ACTION = COOLDOWN)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (ACTION = KILL EXEC_ELAPSED = '10m')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' WATCH = SIMILAR DURATION = '10m' ACTION = COOLDOWN)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = EXACT DURATION = '10m')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '9s' ACTION = COOLDOWN WATCH = EXACT DURATION = UNLIMITED)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '8s' ACTION = COOLDOWN WATCH = EXACT DURATION = UNLIMITED)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '7s' ACTION = COOLDOWN WATCH = EXACT DURATION = UNLIMITED)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '7s' ACTION = COOLDOWN WATCH = EXACT DURATION = UNLIMITED)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' RU = 100 ACTION = DRYRUN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' PROCESSED_KEYS = 100 ACTION = DRYRUN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' WATCH = SIMILAR DURATION = '10m' ACTION = COOLDOWN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = EXACT DURATION = '10m')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BACKGROUND = (TASK_TYPES = '')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BACKGROUND = (UTILIZATION_LIMIT = 50)
-- case
-- error: line 1 column 77 near ""NAN")" 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BACKGROUND = (TASK_TYPES = 'br,lightning')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, BACKGROUND = (TASK_TYPES = 'br,lightning', UTILIZATION_LIMIT = 50)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = EXACT DURATION = '10m'), BACKGROUND = (TASK_TYPES = 'br,lightning')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = PLAN DURATION = '10m'), BACKGROUND = (TASK_TYPES = 'br,lightning')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = PLAN DURATION = '10m'), BACKGROUND = (TASK_TYPES = 'br,lightning', UTILIZATION_LIMIT = 10)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = PLAN DURATION = '10m'), BACKGROUND = (TASK_TYPES = 'br,lightning')
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s')
-- case
-- error: line 1 column 45 near "QUERY=(EXEC_ELAPSED '10s')" 
-- case
-- error: line 1 column 64 near "EXEC_ELAPSED '10s'" 
-- case
-- error: line 1 column 73 near "" 
-- case
-- error: line 1 column 45 near "LIMIT=(EXEC_ELAPSED '10s')" 
-- case
-- error: line 1 column 100 near ")"Dupliated runaway options specified 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (PROCESSED_KEYS = 100)
-- case
-- error: line 1 column 45 near "QUERY=(PROCESSED_KEYS 100)" 
-- case
-- error: line 1 column 66 near "PROCESSED_KEYS 100" 
-- case
-- error: line 1 column 73 near "" 
-- case
-- error: line 1 column 45 near "LIMIT=(PROCESSED_KEYS 100)" 
-- case
-- error: line 1 column 100 near ")"Dupliated runaway options specified 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (RU = 100)
-- case
-- error: line 1 column 45 near "QUERY=(RU 100)" 
-- case
-- error: line 1 column 54 near "RU 100" 
-- case
-- error: line 1 column 61 near "" 
-- case
-- error: line 1 column 45 near "LIMIT=(RU 100)" 
-- case
-- error: line 1 column 88 near ")"Dupliated runaway options specified 
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' PROCESSED_KEYS = 100)
-- case
CREATE RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' PROCESSED_KEYS = 100 RU = 100)
-- case
-- error: line 1 column 26 near "cpu ='8c'" 
-- case
-- error: line 1 column 29 near "region ='us, 3'" 
-- case
-- error: line 1 column 37 near "true" 
-- case
-- error: line 1 column 38 near "false" 
-- case
-- error: line 1 column 40 near "disable" 
-- case
ALTER RESOURCE GROUP `default` PRIORITY = HIGH
-- case
-- error: line 1 column 26 near "cpu='8c', io_read_bandwidth='2GB/s', io_write_bandwidth='200MB/s'" 
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 2000, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, BURSTABLE = MODERATED
-- case
-- error: line 1 column 41 near "'check', BURSTABLE" 
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 3000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 4000
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 2000, BURSTABLE = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 2000, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 2000, BURSTABLE = OFF
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, BURSTABLE = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, BURSTABLE = OFF
-- case
-- error: line 1 column 41 near "'check', BURSTABLE" 
-- case
-- error: line 1 column 53 near "yes" 
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = UNLIMITED, RU_PER_SEC = 3000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 3000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = OFF, RU_PER_SEC = 3000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = UNLIMITED, RU_PER_SEC = 4000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED, RU_PER_SEC = 4000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = OFF, RU_PER_SEC = 4000
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 200000, BURSTABLE = MODERATED
-- case
-- error: line 1 column 32 near "followers=0" 
-- case
-- error: line 1 column 49 near "MID BURSTABLE" 
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` BURSTABLE = OFF
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 200000, BURSTABLE = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 200000, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 200000, BURSTABLE = OFF
-- case
-- error: line 1 column 32 near "followers=0" 
-- case
-- error: line 1 column 49 near "MID" 
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 20, PRIORITY = HIGH, BURSTABLE = UNLIMITED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 20, PRIORITY = HIGH, BURSTABLE = MODERATED
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 20, PRIORITY = HIGH, BURSTABLE = OFF
-- case
ALTER RESOURCE GROUP `x` QUERY_LIMIT = NULL
-- case
ALTER RESOURCE GROUP `x` QUERY_LIMIT = NULL
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = DRYRUN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = NULL
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10m' ACTION = COOLDOWN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (ACTION = KILL EXEC_ELAPSED = '10m')
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' WATCH = SIMILAR DURATION = '10m' ACTION = COOLDOWN)
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = EXACT DURATION = '10m')
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = UNLIMITED, QUERY_LIMIT = (EXEC_ELAPSED = '10s' ACTION = COOLDOWN WATCH = EXACT DURATION = '10m')
-- case
ALTER RESOURCE GROUP `x` RU_PER_SEC = 1000, QUERY_LIMIT = (EXEC_ELAPSED = '10s')
-- case
-- error: line 1 column 63 near "EXEC_ELAPSED '10s'" 
-- case
-- error: line 1 column 99 near ")"Dupliated runaway options specified 
-- case
-- error: line 1 column 132 near ")"Dupliated runaway options specified 
-- case
ALTER RESOURCE GROUP `x` BACKGROUND = NULL
-- case
ALTER RESOURCE GROUP `x` BACKGROUND = NULL
-- case
ALTER RESOURCE GROUP `default` PRIORITY = LOW, BACKGROUND = (TASK_TYPES = 'ttl')
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = UNLIMITED, BACKGROUND = (TASK_TYPES = 'a,b,c')
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = MODERATED, BACKGROUND = (TASK_TYPES = 'a,b,c')
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = OFF, BACKGROUND = (TASK_TYPES = 'a,b,c')
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = UNLIMITED, BACKGROUND = (UTILIZATION_LIMIT = 20)
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = MODERATED, BACKGROUND = (UTILIZATION_LIMIT = 20)
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = OFF, BACKGROUND = (UTILIZATION_LIMIT = 20)
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = UNLIMITED, BACKGROUND = (TASK_TYPES = 'a,b,c', UTILIZATION_LIMIT = 20)
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = MODERATED, BACKGROUND = (TASK_TYPES = 'a,b,c', UTILIZATION_LIMIT = 20)
-- case
ALTER RESOURCE GROUP `default` BURSTABLE = OFF, BACKGROUND = (TASK_TYPES = 'a,b,c', UTILIZATION_LIMIT = 20)
-- case
-- error: line 1 column 87 near "'abc' )" 
-- case
-- error: line 1 column 87 near "'abc' )" 
-- case
-- error: line 1 column 81 near "'abc' )" 
-- case
DROP RESOURCE GROUP `x`
-- case
DROP RESOURCE GROUP `DEFAULT`
-- case
DROP RESOURCE GROUP IF EXISTS `x`
-- case
-- error: line 1 column 22 near ",y" 
-- case
-- error: line 1 column 32 near ",y" 
-- case
SET RESOURCE GROUP `x`
-- case
SET RESOURCE GROUP ``
-- case
SET RESOURCE GROUP `default`
-- case
SET RESOURCE GROUP `default`
-- case
-- error: line 1 column 22 near "y" 
-- case
CREATE ROLE `RESOURCE`@`%`
-- case
-- error: line 1 column 20 near "" 
-- case
CREATE TABLE `t` (`a` INT) STATS_BUCKETS = 1
-- case
-- error: line 1 column 42 near "'abc'" 
-- case
-- error: line 1 column 37 near "" 
-- case
CREATE TABLE `t` (`a` INT) STATS_TOPN = 1
-- case
-- error: line 1 column 39 near "'abc'" 
-- case
CREATE TABLE `t` (`a` INT) STATS_AUTO_RECALC = 1
-- case
-- error: line 1 column 46 near "'abc'" 
-- case
CREATE TABLE `t` (`a` INT) STATS_SAMPLE_RATE = 0.1
-- case
-- error: line 1 column 46 near "'abc'" 
-- case
CREATE TABLE `t` (`a` INT) STATS_COL_CHOICE = 'all'
-- case
CREATE TABLE `t` (`a` INT) STATS_COL_CHOICE = 'list'
-- case
-- error: line 1 column 41 near "1" 
-- case
CREATE TABLE `t` (`a` INT,`b` INT) STATS_COL_LIST = 'a,b'
-- case
-- error: line 1 column 46 near "1" 
-- case
CREATE TABLE `t` (`a` INT) STATS_BUCKETS = 1 STATS_TOPN = 1
-- case
CREATE TABLE `t` (`a` INT) STATS_BUCKETS = 1 STATS_TOPN = 1 PARTITION BY RANGE (`a`) (PARTITION `p1` VALUES LESS THAN (200))
-- case
ALTER TABLE `t` STATS_OPTIONS='str'
-- case
ALTER TABLE `t` STATS_OPTIONS='str1,str2'
-- case
ALTER TABLE `t` STATS_OPTIONS='str1,str2'
-- case
ALTER TABLE `t` STATS_OPTIONS='str1,str2'
-- case
ALTER TABLE `t` STATS_OPTIONS='str1,str2'
-- case
ALTER TABLE `t` STATS_OPTIONS=DEFAULT
-- case
ALTER TABLE `t` STATS_OPTIONS=DEFAULT
-- case
ALTER TABLE `t` STATS_OPTIONS=DEFAULT
-- case
-- error: line 1 column 27 near "" 
-- case
CREATE TABLE `t` (`a` INT) INSERT_METHOD = FIRST
