ADMIN SHOW DDL
-- case
ADMIN SHOW DDL JOBS
-- case
ADMIN SHOW DDL JOBS WHERE `id`>0
-- case
ADMIN SHOW DDL JOBS 20 WHERE `id`=0
-- case
-- error: line 1 column 21 near "-1;" 
-- case
ADMIN SHOW DDL JOB QUERIES 1
-- case
ADMIN SHOW DDL JOB QUERIES 1, 2, 3, 4
-- case
ADMIN SHOW DDL JOB QUERIES LIMIT 0, 5
-- case
ADMIN SHOW DDL JOB QUERIES LIMIT 5, 10
-- case
ADMIN SHOW DDL JOB QUERIES LIMIT 2, 3
-- case
ADMIN SHOW DDL JOB QUERIES LIMIT 0, 22
-- case
ADMIN SHOW `t1` NEXT_ROW_ID
-- case
ADMIN CREATE WORKLOAD SNAPSHOT
-- case
ADMIN CHECK TABLE `t1`, `t2`
-- case
ADMIN CHECK INDEX `tableName` idxName
-- case
ADMIN CHECK INDEX `tableName` idxName (1,2), (4,5)
-- case
ADMIN CHECKSUM TABLE `t1`, `t2`
-- case
ADMIN CANCEL DDL JOBS 1
-- case
ADMIN CANCEL DDL JOBS 1, 2
-- case
ADMIN PAUSE DDL JOBS 1, 3
-- case
ADMIN PAUSE DDL JOBS 5
-- case
-- error: line 1 column 20 near "" 
-- case
-- error: line 1 column 32 near "str_not_num" 
-- case
ADMIN RESUME DDL JOBS 1, 2
-- case
ADMIN RESUME DDL JOBS 3
-- case
-- error: line 1 column 21 near "" 
-- case
-- error: line 1 column 33 near "str_not_num" 
-- case
ADMIN RECOVER INDEX `t1` idx_a
-- case
ADMIN CLEANUP INDEX `t1` idx_a
-- case
ADMIN SHOW SLOW TOP 3
-- case
ADMIN SHOW SLOW TOP INTERNAL 7
-- case
ADMIN SHOW SLOW TOP ALL 9
-- case
ADMIN SHOW SLOW RECENT 11
-- case
ADMIN RELOAD EXPR_PUSHDOWN_BLACKLIST
-- case
ADMIN PLUGINS DISABLE audit, whitelist
-- case
ADMIN PLUGINS ENABLE audit, whitelist
-- case
ADMIN FLUSH BINDINGS
-- case
ADMIN CAPTURE BINDINGS
-- case
ADMIN EVOLVE BINDINGS
-- case
ADMIN RELOAD BINDINGS
-- case
ADMIN RELOAD CLUSTER BINDINGS
-- case
ADMIN RELOAD STATS_EXTENDED
-- case
ADMIN RELOAD STATS_EXTENDED
-- case
ADMIN FLUSH INSTANCE PLAN_CACHE
-- case
ADMIN FLUSH SESSION PLAN_CACHE
-- case
ADMIN FLUSH GLOBAL PLAN_CACHE
-- case
ADMIN SET BDR ROLE PRIMARY
-- case
ADMIN SET BDR ROLE SECONDARY
-- case
ADMIN UNSET BDR ROLE
-- case
ADMIN SHOW BDR ROLE
-- case
ADMIN ALTER DDL JOBS 1 thread = 2
-- case
-- error: line 1 column 32 near "" 
-- case
-- error: line 1 column 29 near "" 
-- case
ADMIN ALTER DDL JOBS 1 batch_size = 3
-- case
-- error: line 1 column 36 near "" 
-- case
-- error: line 1 column 33 near "" 
-- case
ADMIN ALTER DDL JOBS 1 max_write_speed = 4
-- case
ADMIN ALTER DDL JOBS 1 max_write_speed = _UTF8MB4'4MiB'
-- case
-- error: line 1 column 41 near "" 
-- case
-- error: line 1 column 38 near "" 
