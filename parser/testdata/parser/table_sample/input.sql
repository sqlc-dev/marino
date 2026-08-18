select * from tbl tablesample system (50);
-- case
select * from tbl tablesample system (50 percent);
-- case
select * from tbl tablesample system (49.9 percent);
-- case
select * from tbl tablesample system (120 rows);
-- case
select * from tbl tablesample bernoulli (50);
-- case
select * from tbl tablesample (50);
-- case
select * from tbl tablesample (50) repeatable (123456789);
-- case
select * from tbl as a tablesample (50);
-- case
select * from tbl `tablesample` tablesample (50);
-- case
select * from tbl tablesample (50) where id > 20;
-- case
select * from tbl partition (p0) tablesample (50);
-- case
select * from tbl tablesample (0 percent);
-- case
select * from tbl tablesample (100 percent);
-- case
select * from tbl tablesample (0 rows);
-- case
select * from tbl tablesample ('34');
-- case
select * from tbl1 tablesample (10), tbl2 tablesample (20);
-- case
select * from tbl1 a tablesample (10) join tbl2 b tablesample (20) on a.id <> b.id;
-- case
select * from demo tablesample bernoulli(50) limit 1 into outfile '/tmp/sample.csv';
-- case
select * from demo tablesample bernoulli(50) order by a, b into outfile '/tmp/sample.csv';
-- case
select * from tbl tablesample system(50) a;
-- case
select * from tbl tablesample (50) partition (p0);
-- case
select * from tbl where id > 20 tablesample system(50);
-- case
select * from (select * from tbl) a tablesample system(50);
-- case
select * from tbl tablesample system(50) tablesample system(50);
-- case
select * from tbl tablesample system(50, 50);
-- case
select * from tbl tablesample dhfksdlfljcoew(50);
-- case
select * from tbl tablesample system;
-- case
select * from tbl tablesample system (33) repeatable;
-- case
select 1 from dual tablesample system (50);
