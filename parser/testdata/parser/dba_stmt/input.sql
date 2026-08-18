SHOW VARIABLES LIKE 'character_set_results'
-- case
SHOW GLOBAL VARIABLES LIKE 'character_set_results'
-- case
SHOW SESSION VARIABLES LIKE 'character_set_results'
-- case
SHOW VARIABLES
-- case
SHOW GLOBAL VARIABLES
-- case
SHOW GLOBAL VARIABLES WHERE Variable_name = 'autocommit'
-- case
SHOW STATUS
-- case
SHOW GLOBAL STATUS
-- case
SHOW SESSION STATUS
-- case
SHOW STATUS LIKE 'Up%'
-- case
SHOW STATUS WHERE Variable_name
-- case
SHOW STATUS WHERE Variable_name LIKE 'Up%'
-- case
SHOW FULL TABLES FROM icar_qa LIKE play_evolutions
-- case
SHOW FULL TABLES WHERE Table_Type != 'VIEW'
-- case
SHOW GRANTS
-- case
SHOW GRANTS FOR 'test'@'localhost'
-- case
SHOW GRANTS FOR 'test'@'LOCALHOST'
-- case
SHOW GRANTS FOR current_user()
-- case
SHOW GRANTS FOR current_user
-- case
SHOW GRANTS FOR 'u1'@'localhost' USING 'r1'
-- case
SHOW GRANTS FOR 'u1'@'localhost' USING 'r1', 'r2'
-- case
SHOW COLUMNS FROM City;
-- case
SHOW COLUMNS FROM tv189.1_t_1_x;
-- case
SHOW FIELDS FROM City;
-- case
SHOW TRIGGERS LIKE 't'
-- case
SHOW DATABASES LIKE 'test2'
-- case
SHOW PROCEDURE STATUS WHERE Db='test'
-- case
SHOW FUNCTION STATUS WHERE Db='test'
-- case
SHOW INDEX FROM t;
-- case
SHOW KEYS FROM t;
-- case
SHOW INDEX IN t;
-- case
SHOW KEYS IN t;
-- case
SHOW INDEXES IN t where true;
-- case
SHOW KEYS FROM t FROM test where true;
-- case
SHOW EVENTS FROM test_db WHERE definer = 'current_user'
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
SHOW PROFILE CPU, MEMORY, BLOCK IO, CONTEXT SWITCHES, PAGE FAULTS, IPC, SWAPS, SOURCE FOR QUERY 1 limit 100
-- case
SHOW MASTER STATUS
-- case
SHOW BINARY LOG STATUS
-- case
SHOW PRIVILEGES
-- case
show character set;
-- case
show charset
-- case
show collation
-- case
show collation like 'utf8%'
-- case
show collation where Charset = 'utf8' and Collation = 'utf8_bin'
-- case
show columns in t;
-- case
show full columns in t;
-- case
SHOW COLUMNS FROM City;
-- case
SHOW EXTENDED COLUMNS FROM City;
-- case
SHOW EXTENDED FIELDS FROM City;
-- case
SHOW EXTENDED FULL COLUMNS FROM City;
-- case
SHOW EXTENDED FULL FIELDS FROM City;
-- case
show create table test.t
-- case
show create table t
-- case
show create view test.t
-- case
show create view t
-- case
show create database d1
-- case
show create database if not exists d1
-- case
show create sequence seq
-- case
show create sequence test.seq
-- case
show stats_extended
-- case
show stats_extended where table_name = 't'
-- case
show stats_meta
-- case
show stats_meta where table_name = 't'
-- case
show stats_locked
-- case
show stats_locked where table_name = 't'
-- case
show stats_histograms
-- case
show stats_histograms where col_name = 'a'
-- case
show stats_buckets
-- case
show stats_buckets where col_name = 'a'
-- case
show stats_healthy
-- case
show stats_healthy where table_name = 't'
-- case
show stats_topn
-- case
show stats_topn where table_name = 't'
-- case
show histograms_in_flight
-- case
show column_stats_usage
-- case
show column_stats_usage where table_name = 't'
-- case
show binding_cache status
-- case
show analyze status
-- case
show analyze status where table_name = 't'
-- case
show analyze status where table_name like '%'
-- case
show builtins
-- case
show backups
-- case
show restores like 'r0001'
-- case
show backups where start_time > now() - interval 10 hour
-- case
show backup
-- case
show restore
-- case
show replica status
-- case
show slave status
-- case
load stats '/tmp/stats.json'
-- case
lock stats test.t
-- case
lock stats t, t2
-- case
lock stats t partition (p0, p1)
-- case
unlock stats test.t
-- case
unlock stats t, t2
-- case
unlock stats t partition (p0, p1)
-- case
SET @ = 1
-- case
SET @' ' = 1
-- case
SET @! = 1
-- case
SET @1 = 1
-- case
SET @a = 1
-- case
SET @b := 1
-- case
SET @.c = 1
-- case
SET @_d = 1
-- case
SET @_e._$. = 1
-- case
SET @~f = 1
-- case
SET @`g,` = 1
-- case
SET
-- case
SET @a = 1, @b := 2
-- case
SET SESSION autocommit = 1
-- case
SET @@session.autocommit = 1
-- case
SET @@SESSION.autocommit = 1
-- case
SET @@GLOBAL.GTID_PURGED = '123'
-- case
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN
-- case
SET LOCAL autocommit = 1
-- case
SET @@local.autocommit = 1
-- case
SET @@autocommit = 1
-- case
SET autocommit = 1
-- case
SET GLOBAL autocommit = 1
-- case
SET @@global.autocommit = 1
-- case
SET autocommit := 1
-- case
SET @@session.autocommit := 1
-- case
SET @MYSQLDUMP_TEMP_LOG_BIN := @@SESSION.SQL_LOG_BIN
-- case
SET LOCAL autocommit := 1
-- case
SET @@global.autocommit := default
-- case
SET @@global.autocommit = default
-- case
SET @@session.autocommit = default
-- case
SET @@character_set_results = binary
-- case
SET CHARACTER SET utf8mb4;
-- case
SET CHARACTER SET 'utf8mb4';
-- case
SET PASSWORD = 'password';
-- case
SET PASSWORD FOR 'root'@'localhost' = 'password';
-- case
SET SESSION TRANSACTION ISOLATION LEVEL REPEATABLE READ
-- case
SET GLOBAL TRANSACTION ISOLATION LEVEL REPEATABLE READ
-- case
SET SESSION TRANSACTION READ WRITE
-- case
SET SESSION TRANSACTION READ ONLY
-- case
SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED
-- case
SET SESSION TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
-- case
SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE
-- case
SET TRANSACTION ISOLATION LEVEL REPEATABLE READ
-- case
SET TRANSACTION READ WRITE
-- case
SET TRANSACTION READ ONLY
-- case
SET TRANSACTION ISOLATION LEVEL READ COMMITTED
-- case
SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED
-- case
SET TRANSACTION ISOLATION LEVEL SERIALIZABLE
-- case
set names utf8
-- case
set names utf8 collate utf8_unicode_ci
-- case
set names binary
-- case
set names default
-- case
set character set default
-- case
set charset default
-- case
set char set default
-- case
set role `role1`
-- case
SET ROLE DEFAULT
-- case
SET ROLE ALL
-- case
SET ROLE ALL EXCEPT `role1`, `role2`
-- case
SET DEFAULT ROLE administrator, developer TO `joe`@`10.0.0.1`
-- case
set names utf8, @@session.sql_mode=1;
-- case
set @@session.sql_mode=1, names utf8, charset utf8;
-- case
set config TIKV LOG.LEVEL='info'
-- case
set config PD LOG.LEVEL='info'
-- case
set config TIDB LOG.LEVEL='info'
-- case
set config '127.0.0.1:3306' LOG.LEVEL='info'
-- case
set config '127.0.0.1:3306' AUTO-COMPACTION-MODE=TRUE
-- case
set config '127.0.0.1:3306' LABEL-PROPERTY.REJECT-LEADER.KEY='zone'
-- case
show config
-- case
show config where type='tidb'
-- case
show config where instance='127.0.0.1:3306'
-- case
create table CONFIG (a int)
-- case
flush no_write_to_binlog tables tbl1 with read lock
-- case
flush table
-- case
flush tables
-- case
flush tables tbl1
-- case
flush no_write_to_binlog tables tbl1
-- case
flush local tables tbl1
-- case
flush table with read lock
-- case
flush tables tbl1, tbl2, tbl3
-- case
flush tables tbl1, tbl2, tbl3 with read lock
-- case
flush privileges
-- case
flush status
-- case
flush tidb plugins plugin1
-- case
flush tidb plugins plugin1, plugin2
-- case
flush hosts
-- case
flush logs
-- case
flush binary logs
-- case
flush engine logs
-- case
flush error logs
-- case
flush general logs
-- case
flush slow logs
-- case
flush client_errors_summary
-- case
flush stats_delta
-- case
flush stats_delta cluster
-- case
flush stats_delta cluster cluster
-- case
flush stats_delta *.*
-- case
flush stats_delta *.* cluster
-- case
flush stats_delta db1.*
-- case
flush stats_delta db1.* cluster
-- case
flush stats_delta t1
-- case
flush stats_delta db1.t1
-- case
flush stats_delta db1.t1 cluster
-- case
flush stats_delta db1.t1, db2.*
-- case
flush stats_delta db1.t1, db2.* cluster
-- case
call 
-- case
call test
-- case
call test()
-- case
call test(1, 'test', true)
-- case
call x.y;
-- case
call x.y();
-- case
call x.y('p', 'q', 'r');
-- case
call `x`.`y`;
-- case
call `x`.`y`();
-- case
call `x`.`y`('p', 'q', 'r');
