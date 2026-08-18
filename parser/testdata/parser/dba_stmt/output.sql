SHOW SESSION VARIABLES LIKE _UTF8MB4'character_set_results'
-- case
SHOW GLOBAL VARIABLES LIKE _UTF8MB4'character_set_results'
-- case
SHOW SESSION VARIABLES LIKE _UTF8MB4'character_set_results'
-- case
SHOW SESSION VARIABLES
-- case
SHOW GLOBAL VARIABLES
-- case
SHOW GLOBAL VARIABLES WHERE `Variable_name`=_UTF8MB4'autocommit'
-- case
SHOW SESSION STATUS
-- case
SHOW GLOBAL STATUS
-- case
SHOW SESSION STATUS
-- case
SHOW SESSION STATUS LIKE _UTF8MB4'Up%'
-- case
SHOW SESSION STATUS WHERE `Variable_name`
-- case
SHOW SESSION STATUS WHERE `Variable_name` LIKE _UTF8MB4'Up%'
-- case
SHOW FULL TABLES IN `icar_qa` LIKE `play_evolutions`
-- case
SHOW FULL TABLES WHERE `Table_Type`!=_UTF8MB4'VIEW'
-- case
SHOW GRANTS
-- case
SHOW GRANTS FOR `test`@`localhost`
-- case
SHOW GRANTS FOR `test`@`localhost`
-- case
SHOW GRANTS FOR CURRENT_USER
-- case
SHOW GRANTS FOR CURRENT_USER
-- case
SHOW GRANTS FOR `u1`@`localhost` USING `r1`@`%`
-- case
SHOW GRANTS FOR `u1`@`localhost` USING `r1`@`%`, `r2`@`%`
-- case
SHOW COLUMNS IN `City`
-- case
SHOW COLUMNS IN `tv189`.`1_t_1_x`
-- case
SHOW COLUMNS IN `City`
-- case
SHOW TRIGGERS LIKE _UTF8MB4't'
-- case
SHOW DATABASES LIKE _UTF8MB4'test2'
-- case
SHOW PROCEDURE STATUS WHERE `Db`=_UTF8MB4'test'
-- case
SHOW FUNCTION STATUS WHERE `Db`=_UTF8MB4'test'
-- case
SHOW INDEX IN `t`
-- case
SHOW INDEX IN `t`
-- case
SHOW INDEX IN `t`
-- case
SHOW INDEX IN `t`
-- case
SHOW INDEX IN `t` WHERE TRUE
-- case
SHOW INDEX IN `test`.`t` WHERE TRUE
-- case
SHOW EVENTS IN `test_db` WHERE `definer`=_UTF8MB4'current_user'
-- case
SHOW PLUGINS
-- case
SHOW PROFILES
-- case
SHOW PROFILE
-- case
SHOW PROFILE FOR QUERY 1
-- case
SHOW PROFILE CPU FOR QUERY 2
-- case
SHOW PROFILE CPU FOR QUERY 2 LIMIT 1,1
-- case
SHOW PROFILE CPU, MEMORY, BLOCK IO, CONTEXT SWITCHES, PAGE FAULTS, IPC, SWAPS, SOURCE FOR QUERY 1 LIMIT 100
-- case
SHOW MASTER STATUS
-- case
SHOW BINARY LOG STATUS
-- case
SHOW PRIVILEGES
-- case
SHOW CHARSET
-- case
SHOW CHARSET
-- case
SHOW COLLATION
-- case
SHOW COLLATION LIKE _UTF8MB4'utf8%'
-- case
SHOW COLLATION WHERE `Charset`=_UTF8MB4'utf8' AND `Collation`=_UTF8MB4'utf8_bin'
-- case
SHOW COLUMNS IN `t`
-- case
SHOW FULL COLUMNS IN `t`
-- case
SHOW COLUMNS IN `City`
-- case
SHOW EXTENDED COLUMNS IN `City`
-- case
SHOW EXTENDED COLUMNS IN `City`
-- case
SHOW EXTENDED FULL COLUMNS IN `City`
-- case
SHOW EXTENDED FULL COLUMNS IN `City`
-- case
SHOW CREATE TABLE `test`.`t`
-- case
SHOW CREATE TABLE `t`
-- case
SHOW CREATE VIEW `test`.`t`
-- case
SHOW CREATE VIEW `t`
-- case
SHOW CREATE DATABASE `d1`
-- case
SHOW CREATE DATABASE IF NOT EXISTS `d1`
-- case
SHOW CREATE SEQUENCE `seq`
-- case
SHOW CREATE SEQUENCE `test`.`seq`
-- case
SHOW STATS_EXTENDED
-- case
SHOW STATS_EXTENDED WHERE `table_name`=_UTF8MB4't'
-- case
SHOW STATS_META
-- case
SHOW STATS_META WHERE `table_name`=_UTF8MB4't'
-- case
SHOW STATS_LOCKED
-- case
SHOW STATS_LOCKED WHERE `table_name`=_UTF8MB4't'
-- case
SHOW STATS_HISTOGRAMS
-- case
SHOW STATS_HISTOGRAMS WHERE `col_name`=_UTF8MB4'a'
-- case
SHOW STATS_BUCKETS
-- case
SHOW STATS_BUCKETS WHERE `col_name`=_UTF8MB4'a'
-- case
SHOW STATS_HEALTHY
-- case
SHOW STATS_HEALTHY WHERE `table_name`=_UTF8MB4't'
-- case
SHOW STATS_TOPN
-- case
SHOW STATS_TOPN WHERE `table_name`=_UTF8MB4't'
-- case
SHOW HISTOGRAMS_IN_FLIGHT
-- case
SHOW COLUMN_STATS_USAGE
-- case
SHOW COLUMN_STATS_USAGE WHERE `table_name`=_UTF8MB4't'
-- case
SHOW BINDING_CACHE STATUS
-- case
SHOW ANALYZE STATUS
-- case
SHOW ANALYZE STATUS WHERE `table_name`=_UTF8MB4't'
-- case
SHOW ANALYZE STATUS WHERE `table_name` LIKE _UTF8MB4'%'
-- case
SHOW BUILTINS
-- case
SHOW BACKUPS
-- case
SHOW RESTORES LIKE _UTF8MB4'r0001'
-- case
SHOW BACKUPS WHERE `start_time`>DATE_SUB(NOW(), INTERVAL 10 HOUR)
-- case
-- error: line 1 column 11 near "" 
-- case
-- error: line 1 column 12 near "restore" 
-- case
SHOW REPLICA STATUS
-- case
SHOW REPLICA STATUS
-- case
LOAD STATS '/tmp/stats.json'
-- case
LOCK STATS `test`.`t`
-- case
LOCK STATS `t`, `t2`
-- case
LOCK STATS `t` PARTITION(`p0`, `p1`)
-- case
UNLOCK STATS `test`.`t`
-- case
UNLOCK STATS `t`, `t2`
-- case
UNLOCK STATS `t` PARTITION(`p0`, `p1`)
-- case
SET @``=1
-- case
SET @` `=1
-- case
-- error: line 1 column 6 near "! = 1" 
-- case
SET @`1`=1
-- case
SET @`a`=1
-- case
SET @`b`=1
-- case
SET @`.c`=1
-- case
SET @`_d`=1
-- case
SET @`_e._$.`=1
-- case
-- error: line 1 column 6 near "~f = 1" 
-- case
SET @`g,`=1
-- case
-- error: line 1 column 3 near "" 
-- case
SET @`a`=1, @`b`=2
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@GLOBAL.`gtid_purged`=_UTF8MB4'123'
-- case
SET @`MYSQLDUMP_TEMP_LOG_BIN`=@@SESSION.`sql_log_bin`
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@GLOBAL.`autocommit`=1
-- case
SET @@GLOBAL.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @`MYSQLDUMP_TEMP_LOG_BIN`=@@SESSION.`sql_log_bin`
-- case
SET @@SESSION.`autocommit`=1
-- case
SET @@GLOBAL.`autocommit`=DEFAULT
-- case
SET @@GLOBAL.`autocommit`=DEFAULT
-- case
SET @@SESSION.`autocommit`=DEFAULT
-- case
SET @@SESSION.`character_set_results`=_UTF8MB4'BINARY'
-- case
SET CHARSET 'utf8mb4'
-- case
SET CHARSET 'utf8mb4'
-- case
SET PASSWORD='password'
-- case
SET PASSWORD FOR `root`@`localhost`='password'
-- case
SET @@SESSION.`tx_isolation`=_UTF8MB4'REPEATABLE-READ'
-- case
SET @@GLOBAL.`tx_isolation`=_UTF8MB4'REPEATABLE-READ'
-- case
SET @@SESSION.`tx_read_only`=_UTF8MB4'0'
-- case
SET @@SESSION.`tx_read_only`=_UTF8MB4'1'
-- case
SET @@SESSION.`tx_isolation`=_UTF8MB4'READ-COMMITTED'
-- case
SET @@SESSION.`tx_isolation`=_UTF8MB4'READ-UNCOMMITTED'
-- case
SET @@SESSION.`tx_isolation`=_UTF8MB4'SERIALIZABLE'
-- case
SET @@SESSION.`tx_isolation_one_shot`=_UTF8MB4'REPEATABLE-READ'
-- case
SET @@SESSION.`tx_read_only`=_UTF8MB4'0'
-- case
SET @@SESSION.`tx_read_only`=_UTF8MB4'1'
-- case
SET @@SESSION.`tx_isolation_one_shot`=_UTF8MB4'READ-COMMITTED'
-- case
SET @@SESSION.`tx_isolation_one_shot`=_UTF8MB4'READ-UNCOMMITTED'
-- case
SET @@SESSION.`tx_isolation_one_shot`=_UTF8MB4'SERIALIZABLE'
-- case
SET NAMES 'utf8'
-- case
SET NAMES 'utf8' COLLATE 'utf8_unicode_ci'
-- case
SET NAMES 'binary'
-- case
SET NAMES DEFAULT
-- case
SET CHARSET DEFAULT
-- case
SET CHARSET DEFAULT
-- case
SET CHARSET DEFAULT
-- case
SET ROLE `role1`@`%`
-- case
SET ROLE DEFAULT
-- case
SET ROLE ALL
-- case
SET ROLE ALL EXCEPT `role1`@`%`, `role2`@`%`
-- case
SET DEFAULT ROLE `administrator`@`%`, `developer`@`%` TO `joe`@`10.0.0.1`
-- case
SET NAMES 'utf8', @@SESSION.`sql_mode`=1
-- case
SET @@SESSION.`sql_mode`=1, NAMES 'utf8', CHARSET 'utf8'
-- case
SET CONFIG TIKV LOG.LEVEL = _UTF8MB4'info'
-- case
SET CONFIG PD LOG.LEVEL = _UTF8MB4'info'
-- case
SET CONFIG TIDB LOG.LEVEL = _UTF8MB4'info'
-- case
SET CONFIG '127.0.0.1:3306' LOG.LEVEL = _UTF8MB4'info'
-- case
SET CONFIG '127.0.0.1:3306' AUTO-COMPACTION-MODE = TRUE
-- case
SET CONFIG '127.0.0.1:3306' LABEL-PROPERTY.REJECT-LEADER.KEY = _UTF8MB4'zone'
-- case
SHOW CONFIG
-- case
SHOW CONFIG WHERE `type`=_UTF8MB4'tidb'
-- case
SHOW CONFIG WHERE `instance`=_UTF8MB4'127.0.0.1:3306'
-- case
CREATE TABLE `CONFIG` (`a` INT)
-- case
FLUSH NO_WRITE_TO_BINLOG TABLES `tbl1` WITH READ LOCK
-- case
FLUSH TABLES
-- case
FLUSH TABLES
-- case
FLUSH TABLES `tbl1`
-- case
FLUSH NO_WRITE_TO_BINLOG TABLES `tbl1`
-- case
FLUSH NO_WRITE_TO_BINLOG TABLES `tbl1`
-- case
FLUSH TABLES WITH READ LOCK
-- case
FLUSH TABLES `tbl1`, `tbl2`, `tbl3`
-- case
FLUSH TABLES `tbl1`, `tbl2`, `tbl3` WITH READ LOCK
-- case
FLUSH PRIVILEGES
-- case
FLUSH STATUS
-- case
FLUSH TIDB PLUGINS plugin1
-- case
FLUSH TIDB PLUGINS plugin1, plugin2
-- case
FLUSH HOSTS
-- case
FLUSH LOGS
-- case
FLUSH BINARY LOGS
-- case
FLUSH ENGINE LOGS
-- case
FLUSH ERROR LOGS
-- case
FLUSH GENERAL LOGS
-- case
FLUSH SLOW LOGS
-- case
FLUSH CLIENT_ERRORS_SUMMARY
-- case
-- error: line 1 column 17 near "" 
-- case
FLUSH STATS_DELTA `cluster`
-- case
FLUSH STATS_DELTA `cluster` CLUSTER
-- case
FLUSH STATS_DELTA *.*
-- case
FLUSH STATS_DELTA *.* CLUSTER
-- case
FLUSH STATS_DELTA `db1`.*
-- case
FLUSH STATS_DELTA `db1`.* CLUSTER
-- case
FLUSH STATS_DELTA `t1`
-- case
FLUSH STATS_DELTA `db1`.`t1`
-- case
FLUSH STATS_DELTA `db1`.`t1` CLUSTER
-- case
FLUSH STATS_DELTA `db1`.`t1`, `db2`.*
-- case
FLUSH STATS_DELTA `db1`.`t1`, `db2`.* CLUSTER
-- case
-- error: line 1 column 5 near "" 
-- case
CALL `test`()
-- case
CALL `test`()
-- case
CALL `test`(1, _UTF8MB4'test', TRUE)
-- case
CALL `x`.`y`()
-- case
CALL `x`.`y`()
-- case
CALL `x`.`y`(_UTF8MB4'p', _UTF8MB4'q', _UTF8MB4'r')
-- case
CALL `x`.`y`()
-- case
CALL `x`.`y`()
-- case
CALL `x`.`y`(_UTF8MB4'p', _UTF8MB4'q', _UTF8MB4'r')
