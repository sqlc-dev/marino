admin show ddl;
-- case
admin show ddl jobs;
-- case
admin show ddl jobs where id > 0;
-- case
admin show ddl jobs 20 where id=0;
-- case
admin show ddl jobs -1;
-- case
admin show ddl job queries 1
-- case
admin show ddl job queries 1, 2, 3, 4
-- case
admin show ddl job queries limit 5
-- case
admin show ddl job queries limit 5, 10
-- case
admin show ddl job queries limit 3 offset 2
-- case
admin show ddl job queries limit 22 offset 0
-- case
admin show t1 next_row_id
-- case
admin create workload snapshot;
-- case
admin check table t1, t2;
-- case
admin check index tableName idxName;
-- case
admin check index tableName idxName (1, 2), (4, 5);
-- case
admin checksum table t1, t2;
-- case
admin cancel ddl jobs 1
-- case
admin cancel ddl jobs 1, 2
-- case
admin pause ddl jobs 1, 3
-- case
admin pause ddl jobs 5
-- case
admin pause ddl jobs
-- case
admin pause ddl jobs str_not_num
-- case
admin resume ddl jobs 1, 2
-- case
admin resume ddl jobs 3
-- case
admin resume ddl jobs
-- case
admin resume ddl jobs str_not_num
-- case
admin recover index t1 idx_a
-- case
admin cleanup index t1 idx_a
-- case
admin show slow top 3
-- case
admin show slow top internal 7
-- case
admin show slow top all 9
-- case
admin show slow recent 11
-- case
admin reload expr_pushdown_blacklist
-- case
admin plugins disable audit, whitelist
-- case
admin plugins enable audit, whitelist
-- case
admin flush bindings
-- case
admin capture bindings
-- case
admin evolve bindings
-- case
admin reload bindings
-- case
admin reload cluster bindings
-- case
admin reload statistics
-- case
admin reload stats_extended
-- case
admin flush instance plan_cache
-- case
admin flush session plan_cache
-- case
admin flush global plan_cache
-- case
admin set bdr role primary
-- case
admin set bdr role secondary
-- case
admin unset bdr role
-- case
admin show bdr role
-- case
admin alter ddl jobs 1 thread = 2
-- case
admin alter ddl jobs 1 thread = 
-- case
admin alter ddl jobs 1 thread
-- case
admin alter ddl jobs 1 batch_size = 3
-- case
admin alter ddl jobs 1 batch_size = 
-- case
admin alter ddl jobs 1 batch_size
-- case
admin alter ddl jobs 1 max_write_speed = 4
-- case
admin alter ddl jobs 1 max_write_speed = _UTF8MB4'4MiB'
-- case
admin alter ddl jobs 1 max_write_speed = 
-- case
admin alter ddl jobs 1 max_write_speed
