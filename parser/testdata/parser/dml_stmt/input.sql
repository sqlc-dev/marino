
-- case
;
-- case
INSERT INTO foo VALUES (1234)
-- case
INSERT INTO foo VALUES (1234, 5678)
-- case
INSERT INTO t1 (SELECT * FROM t2)
-- case
INSERT INTO t partition (p0) values(1234)
-- case
REPLACE INTO t partition (p0) values(1234)
-- case
INSERT INTO t partition (p0, p1, p2) values(1234)
-- case
REPLACE INTO t partition (p0, p1, p2) values(1234)
-- case
INSERT INTO foo VALUES (1 || 2)
-- case
INSERT INTO foo VALUES (1 | 2)
-- case
INSERT INTO foo VALUES (false || true)
-- case
INSERT INTO foo VALUES (bar(5678))
-- case
INSERT INTO foo VALUES ()
-- case
SELECT * FROM t
-- case
SELECT * FROM t AS u
-- case
SELECT * FROM t, v
-- case
SELECT * FROM t AS u, v
-- case
SELECT * FROM t, v AS w
-- case
SELECT * FROM t AS u, v AS w
-- case
SELECT * FROM foo, bar, foo
-- case
SELECT DISTINCTS * FROM t
-- case
SELECT DISTINCT * FROM t
-- case
SELECT DISTINCTROW * FROM t
-- case
SELECT ALL * FROM t
-- case
SELECT DISTINCT ALL * FROM t
-- case
SELECT DISTINCTROW ALL * FROM t
-- case
INSERT INTO foo (a) VALUES (42)
-- case
INSERT INTO foo (a,) VALUES (42,)
-- case
INSERT INTO foo (a,b) VALUES (42,314)
-- case
INSERT INTO foo (a,b,) VALUES (42,314)
-- case
INSERT INTO foo (a,b,) VALUES (42,314,)
-- case
INSERT INTO foo () VALUES ()
-- case
INSERT INTO foo VALUE ()
-- case
INSERT INTO tt VALUES (01000001783);
-- case
INSERT INTO tt VALUES (default);
-- case
REPLACE INTO foo VALUES (1 || 2)
-- case
REPLACE INTO foo VALUES (1 | 2)
-- case
REPLACE INTO foo VALUES (false || true)
-- case
REPLACE INTO foo VALUES (bar(5678))
-- case
REPLACE INTO foo VALUES ()
-- case
REPLACE INTO foo (a,b) VALUES (42,314)
-- case
REPLACE INTO foo (a,b,) VALUES (42,314)
-- case
REPLACE INTO foo (a,b,) VALUES (42,314,)
-- case
REPLACE INTO foo () VALUES ()
-- case
REPLACE INTO foo VALUE ()
-- case
SELECT stuff.id
			FROM stuff
			WHERE stuff.value >= ALL (SELECT stuff.value
			FROM stuff)
-- case
BEGIN
-- case
START TRANSACTION
-- case
COMMIT
-- case
COMMIT AND NO CHAIN
-- case
COMMIT NO RELEASE
-- case
COMMIT AND NO CHAIN NO RELEASE
-- case
COMMIT AND NO CHAIN RELEASE
-- case
COMMIT AND CHAIN NO RELEASE
-- case
COMMIT AND CHAIN RELEASE
-- case
ROLLBACK
-- case
ROLLBACK AND NO CHAIN
-- case
ROLLBACK NO RELEASE
-- case
ROLLBACK AND NO CHAIN NO RELEASE
-- case
ROLLBACK AND NO CHAIN RELEASE
-- case
ROLLBACK AND CHAIN NO RELEASE
-- case
ROLLBACK AND CHAIN RELEASE
-- case
BEGIN;
			INSERT INTO foo VALUES (42, 3.14);
			INSERT INTO foo VALUES (-1, 2.78);
		COMMIT;
-- case
BEGIN;
			INSERT INTO tmp SELECT * from bar;
			SELECT * from tmp;
		ROLLBACK;
-- case
SAVEPOINT x
-- case
RELEASE SAVEPOINT x
-- case
ROLLBACK TO x
-- case
ROLLBACK TO X
-- case
ROLLBACK TO SAVEPOINT x
-- case
TABLE t
-- case
(TABLE t)
-- case
TABLE t1, t2
-- case
TABLE t ORDER BY b
-- case
TABLE t LIMIT 3
-- case
TABLE t ORDER BY b LIMIT 3
-- case
TABLE t ORDER BY b LIMIT 3 OFFSET 2
-- case
TABLE t ORDER BY b LIMIT 2,3
-- case
INSERT INTO ta TABLE tb
-- case
INSERT INTO t.a TABLE t.b
-- case
REPLACE INTO ta TABLE tb
-- case
REPLACE INTO t.a TABLE t.b
-- case
TABLE t1 INTO OUTFILE 'a.txt'
-- case
TABLE t ORDER BY a INTO OUTFILE '/tmp/abc'
-- case
CREATE TABLE t.a TABLE t.b
-- case
CREATE TABLE ta TABLE tb
-- case
CREATE TABLE ta (x INT) TABLE tb
-- case
CREATE VIEW v AS TABLE t
-- case
CREATE VIEW v AS (TABLE t)
-- case
SELECT * FROM t1 WHERE a IN (TABLE t2)
-- case
VALUES ROW(1)
-- case
VALUES ROW()
-- case
VALUES ROW(1, default)
-- case
VALUES ROW(1), ROW(2,3)
-- case
VALUES (1,2)
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8)
-- case
VALUES ROW(1,s,3.1), ROW(5,y,9.9)
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) LIMIT 3
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) ORDER BY a
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) ORDER BY a LIMIT 2
-- case
VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE 'a.txt'
-- case
VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/abc'
-- case
CREATE TABLE ta VALUES ROW(1)
-- case
CREATE TABLE ta AS VALUES ROW(1)
-- case
CREATE VIEW a AS VALUES ROW(1)
-- case
SELECT a.b.c FROM t
-- case
SELECT a.b.*.c FROM t
-- case
SELECT a.b.* FROM t
-- case
SELECT a FROM t
-- case
SELECT a.b.c.d FROM t
-- case
DO 1
-- case
DO 1, sleep(1)
-- case
DO 1 from t
-- case
load data local infile '/tmp/t.csv' into table t1 fields terminated by ',' optionally enclosed by '"' ignore 1 lines
-- case
load data infile '/tmp/t.csv' into table t
-- case
load data infile '/tmp/t.csv' into table t character set utf8
-- case
load data infile '/tmp/t.csv' into table t fields terminated by 'ab'
-- case
load data infile '/tmp/t.csv' into table t columns terminated by 'ab'
-- case
load data infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b'
-- case
load data infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' escaped by '*'
-- case
load data infile '/tmp/t.csv' into table t lines starting by 'ab'
-- case
load data infile '/tmp/t.csv' into table t lines starting by 'ab' terminated by 'xy'
-- case
load data infile '/tmp/t.csv' into table t fields terminated by 'ab' lines terminated by 'xy'
-- case
load data infile '/tmp/t.csv' into table t terminated by 'xy' fields terminated by 'ab'
-- case
load data local infile '/tmp/t.csv' into table t
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab'
-- case
load data local infile '/tmp/t.csv' into table t columns terminated by 'ab'
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b'
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' escaped by '*'
-- case
load data local infile '/tmp/t.csv' into table t character set utf8 fields terminated by 'ab' enclosed by 'b' escaped by '*'
-- case
load data local infile '/tmp/t.csv' into table t lines starting by 'ab'
-- case
load data local infile '/tmp/t.csv' into table t lines starting by 'ab' terminated by 'xy'
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' lines terminated by 'xy'
-- case
load data local infile '/tmp/t.csv' into table t terminated by 'xy' fields terminated by 'ab'
-- case
load data infile '/tmp/t.csv' into table t (a,b)
-- case
load data local infile '/tmp/t.csv' into table t (a,b)
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t columns terminated by 'ab' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' escaped by '*' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t character set utf8 fields terminated by 'ab' enclosed by 'b' escaped by '*' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t lines starting by 'ab' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t lines starting by 'ab' terminated by 'xy' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t character set utf8 fields terminated by 'ab' lines terminated by 'xy' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' lines terminated by 'xy' (a,b)
-- case
load data local infile '/tmp/t.csv' into table t (a,b) fields terminated by 'ab'
-- case
load data local infile '/tmp/t.csv' into table t ignore 1 lines
-- case
load data local infile '/tmp/t.csv' into table t ignore -1 lines
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' (a,b) ignore 1 lines
-- case
load data local infile '/tmp/t.csv' into table t lines starting by 'ab' terminated by 'xy' ignore 1 lines
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' escaped by '*' ignore 1 lines (a,b)
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' escaped by ''
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by X'6B6B';
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by X'6B6B' enclosed by X'0D';
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by X'6B6B' enclosed by X'0D0D';
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by B'110101101101011';
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by B'110101101101011' enclosed by B'1101';
-- case
load data local infile '~/1.csv' into table `t_ascii` fields terminated by B'110101101101011' enclosed by B'110100001101';
-- case
load data local infile '~/1.csv' into table `t_ascii` lines starting by B'110101101101011' terminated by B'110101101101011';
-- case
load data local infile '~/1.csv' into table `t_ascii` lines starting by X'6B6B' terminated by X'6B6B';
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' enclosed by 'b' enclosed by 'b'
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' escaped by '' enclosed by 'b'
-- case
load data local infile '/tmp/t.csv' into table t fields terminated by 'ab' escaped by '' enclosed by 'b' SET b = CAST(CONV(MID(@var1, 3, LENGTH(@var1)-3), 2, 10) AS UNSIGNED)
-- case
load data infile '/tmp/t.csv' into table t fields enclosed by ''
-- case
load data infile '/tmp/t.csv' into table t fields enclosed by 'a'
-- case
load data infile '/tmp/t.csv' into table t fields enclosed by 'aa'
-- case
load data infile '/tmp/t.csv' into table t fields escaped by ''
-- case
load data infile '/tmp/t.csv' into table t fields escaped by 'a'
-- case
load data infile '/tmp/t.csv' into table t fields escaped by 'aa'
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE t1 (column1, @dummy, column2, @dummy, column3)
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE t1 (column1, @var1) SET column2 = @var1/100
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE t1 (column1, @var1, @var2) SET column2 = @var1/100, column3 = DEFAULT, column4=CURRENT_TIMESTAMP, column5=@var2+1
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE t1 FIELDS TERMINATED BY ',' LINES TERMINATED BY '
';
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' INTO TABLE t1 FIELDS TERMINATED BY ',' LINES TERMINATED BY '
';
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE t1 FIELDS TERMINATED BY ',' LINES TERMINATED BY '
';
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' REPLACE INTO TABLE t1 FIELDS TERMINATED BY ',' LINES TERMINATED BY '
';
-- case
load data infile 's3://bucket-name/t.csv' into table t
-- case
load data infile '/tmp/t.csv' into table t fields defined null by 'nil'
-- case
load data infile '/tmp/t.csv' into table t fields defined null by X'00'
-- case
load data infile '/tmp/t.csv' into table t fields defined null by 'NULL' optionally enclosed ignore 1 lines
-- case
load data infile '/tmp/t.csv' format 'delimited data' into table t (column1, @var1) SET column2 = @var1/100
-- case
load data local infile '/tmp/t.sql' format 'sql file' replace into table t (a,b)
-- case
load data infile '/tmp/t.parquet' format 'parquet' into table t (column1, @var1) SET column2 = @var1/100
-- case
load data infile '/tmp/t.csv' into table t with detached
-- case
load data infile '/tmp/t.csv' into table `t` with threads=10
-- case
load data infile '/tmp/t.csv' into table `t` with threads=10, detached
-- case
import into t from '/file.csv'
-- case
import into t (a,b) from '/file.csv'
-- case
import into t (a,@1) from '/file.csv'
-- case
import into t (a,@1) set b=@1+100 from '/file.csv'
-- case
import into t from '/file.csv' format 'sql file'
-- case
import into t from '/file.csv' with detached
-- case
import into `t` from '/file.csv' with thread=1
-- case
import into `t` from '/file.csv' with detached, thread=1
-- case
select * from t for update
-- case
select * from t for share
-- case
select * from t for update nowait
-- case
select * from t for update wait 5
-- case
select * from t limit 1 for update wait 11
-- case
select * from t for share nowait
-- case
select * from t for update skip locked
-- case
select * from t for share skip locked
-- case
select * from t lock in share mode
-- case
select * from t lock in share mode nowait
-- case
select * from t lock in share mode skip locked
-- case
select * from t for update of t
-- case
select * from t for share of t
-- case
select * from t for update of t nowait
-- case
select * from t for update of t wait 5
-- case
select * from t limit 1 for update of t wait 11
-- case
select * from t for share of t nowait
-- case
select * from t for update of t skip locked
-- case
select * from t for share of t skip locked
-- case
select a, b from t into outfile '/tmp/result.txt'
-- case
select a from t order by a into outfile '/tmp/abc'
-- case
select 1 into outfile '/tmp/1.csv'
-- case
select 1 for update into outfile '/tmp/1.csv'
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ','
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' enclosed BY '"'
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' optionally enclosed BY '"'
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' lines terminated BY '
'
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' optionally enclosed BY '"' lines terminated BY ''
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' enclosed BY '"' lines terminated BY ''
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' optionally enclosed BY '"' lines starting by 'xy' terminated BY ''
-- case
select a,b,a+b from t into outfile '/tmp/result.txt' fields terminated BY ',' enclosed BY '"' lines starting by 'xy' terminated BY ''
-- case
SELECT * from t1, t2, t3
-- case
select * from t1 join t2 left join t3 on t2.id = t3.id
-- case
select * from t1 right join t2 on t1.id = t2.id left join t3 on t3.id = t2.id
-- case
select * from t1 right join t2 on t1.id = t2.id left join t3
-- case
select * from t1 join t2 left join t3 using (id)
-- case
select * from t1 right join t2 using (id) left join t3 using (id)
-- case
select * from t1 right join t2 using (id) left join t3
-- case
select * from t1 natural join t2
-- case
select * from t1 natural right join t2
-- case
select * from t1 natural left outer join t2
-- case
select * from t1 natural inner join t2
-- case
select * from t1 natural cross join t2
-- case
select * from t3 join t1 join t2 on t1.a=t2.a on t3.b=t2.b
-- case
select * from t1 straight_join t2 on t1.id = t2.id
-- case
select straight_join * from t1 join t2 on t1.id = t2.id
-- case
select straight_join * from t1 left join t2 on t1.id = t2.id
-- case
select straight_join * from t1 right join t2 on t1.id = t2.id
-- case
select straight_join * from t1 straight_join t2 on t1.id = t2.id
-- case
DELETE from t1
-- case
DELETE from t1.*
-- case
DELETE LOW_priORITY from t1
-- case
DELETE quick from t1
-- case
DELETE ignore from t1
-- case
DELETE low_priority quick ignore from t1
-- case
DELETE FROM t1 WHERE t1.a > 0 ORDER BY t1.a
-- case
delete from t1 where a=26
-- case
DELETE from t1 where a=1 limit 1
-- case
DELETE FROM t1 WHERE t1.a > 0 ORDER BY t1.a LIMIT 1
-- case
DELETE FROM x.y z WHERE z.a > 0
-- case
DELETE FROM t1 AS w WHERE a > 0
-- case
DELETE from t1 partition (p0,p1)
-- case
delete low_priority t1, t2 from t1, t2
-- case
delete quick t1, t2 from t1, t2
-- case
delete ignore t1, t2 from t1, t2
-- case
delete ignore t1, t2 from t1 partition (p0,p1), t2
-- case
delete low_priority quick ignore t1, t2 from t1, t2 where t1.a > 5
-- case
delete t1, t2 from t1, t2
-- case
delete t1, t2 from t1, t2 where t1.a = 1 and t2.b <> 1
-- case
delete t1 from t1, t2
-- case
delete t2 from t1, t2
-- case
delete t1 from t1
-- case
delete t1,t2,t3 from t1, t2, t3
-- case
delete t1,t2,t3 from t1, t2, t3 where t3.c < 5 and t1.a = 3
-- case
delete t1 from t1, t1 as t2 where t1.b = t2.b and t1.a > t2.a
-- case
delete t1.*,t2 from t1, t2
-- case
delete t.t1.*,t2 from t1, t2
-- case
delete t1.*, t2.* from t1, t2
-- case
delete t11.*, t12.* from t11, t12 where t11.a = t12.a and t11.b <> 1
-- case
DELETE quick FROM t1,t2 USING t1,t2
-- case
DELETE low_priority ignore FROM t1,t2 USING t1,t2
-- case
DELETE low_priority quick ignore FROM t1,t2 USING t1,t2
-- case
DELETE FROM t1 USING t1 WHERE post='1'
-- case
DELETE FROM t1,t2 USING t1,t2
-- case
DELETE FROM t1,t2,t3 USING t1,t2,t3 where t3.a = 1
-- case
DELETE FROM t2,t3 USING t1,t2,t3 where t1.a = 1
-- case
DELETE FROM t2.*,t3.* USING t1,t2,t3 where t1.a = 1
-- case
DELETE FROM t1,t2.*,t3.* USING t1,t2,t3 where t1.a = 1
-- case
DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;
-- case
DELETE FROM t1, t2 USING t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id;
-- case
DELETE /*+ TiDB_INLJ(t1, t2) */ t1, t2 from t1, t2 where t1.id=t2.id;
-- case
DELETE /*+ TiDB_HJ(t1, t2) */ t1, t2 from t1, t2 where t1.id=t2.id
-- case
DELETE /*+ TiDB_SMJ(t1, t2) */ t1, t2 from t1, t2 where t1.id=t2.id
-- case
DELETE FROM t1 USE INDEX(idx_a) WHERE t1.id=1;
-- case
DELETE t1, t2 FROM t1 USE INDEX(idx_a) JOIN t2 WHERE t1.id=t2.id;
-- case
DELETE t1, t2 FROM t1 USE INDEX(idx_a) JOIN t2 USE INDEX(idx_a) WHERE t1.id=t2.id;
-- case
DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id limit 10;
-- case
DELETE t1, t2 FROM t1 INNER JOIN t2 INNER JOIN t3 WHERE t1.id=t2.id AND t2.id=t3.id order by t1.id;
-- case
INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c=VALUES(a)+VALUES(b);
-- case
INSERT INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b);
-- case
INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c=VALUES(a)+VALUES(b);
-- case
INSERT IGNORE INTO t (a,b,c) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE c:=VALUES(a)+VALUES(b);
-- case
INSERT INTO odk (id, v) VALUES (1, 10) AS new ON DUPLICATE KEY UPDATE v = new.v;
-- case
INSERT INTO odk (id, v) VALUES (1, 10) AS new(nid, nv) ON DUPLICATE KEY UPDATE v = nv;
-- case
INSERT INTO odk VALUES (1, 10) AS new ON DUPLICATE KEY UPDATE v = new.v;
-- case
INSERT IGNORE INTO odk (id, v) VALUES (1, 10) AS new ON DUPLICATE KEY UPDATE v = new.v;
-- case
INSERT INTO odk SET id = 1, v = 10 AS new ON DUPLICATE KEY UPDATE v = new.v;
-- case
INSERT INTO odk (id, v) SELECT id, v FROM staging AS s AS new ON DUPLICATE KEY UPDATE v = new.v;
-- case
INSERT INTO t SET a=1,b=2
-- case
INSERT INTO t (a) SET a=1
-- case
UPDATE LOW_PRIORITY IGNORE t SET id = id + 1 ORDER BY id DESC;
-- case
UPDATE t SET id = id + 1 ORDER BY id DESC;
-- case
UPDATE t SET id = id + 1 ORDER BY id DESC limit 3 ;
-- case
UPDATE t SET id = id + 1, name = 'jojo';
-- case
UPDATE items,month SET items.price=month.price WHERE items.id=month.id;
-- case
UPDATE user T0 LEFT OUTER JOIN user_profile T1 ON T1.id = T0.profile_id SET T0.profile_id = 1 WHERE T0.profile_id IN (1);
-- case
UPDATE t1, t2 set t1.profile_id = 1, t2.profile_id = 1 where ta.a=t.ba
-- case
UPDATE /*+ TiDB_INLJ(t1, t2) */ t1, t2 set t1.profile_id = 1, t2.profile_id = 1 where ta.a=t.ba
-- case
UPDATE /*+ TiDB_SMJ(t1, t2) */ t1, t2 set t1.profile_id = 1, t2.profile_id = 1 where ta.a=t.ba
-- case
UPDATE /*+ TiDB_HJ(t1, t2) */ t1, t2 set t1.profile_id = 1, t2.profile_id = 1 where ta.a=t.ba
-- case
UPDATE items,month SET items.price=month.price WHERE items.id=month.id LIMIT 10;
-- case
UPDATE items,month SET items.price=month.price WHERE items.id=month.id order by month.id;
-- case
UPDATE t1 USE INDEX(idx_a) SET t1.price=3.25 WHERE t1.id=1;
-- case
UPDATE t1 USE INDEX(idx_a) JOIN t2 SET t1.price=t2.price WHERE t1.id=t2.id;
-- case
UPDATE t1 USE INDEX(idx_a) JOIN t2 USE INDEX(idx_a) SET t1.price=t2.price WHERE t1.id=t2.id;
-- case
SELECT * FROM t WHERE 1 = 1
-- case
SELECT * FROM t FETCH FIRST 5 ROW ONLY
-- case
SELECT * FROM t FETCH NEXT 5 ROW ONLY
-- case
SELECT * FROM t FETCH FIRST 5 ROWS ONLY
-- case
SELECT * FROM t FETCH NEXT 5 ROWS ONLY
-- case
SELECT * FROM t FETCH FIRST ROW ONLY
-- case
SELECT * FROM t FETCH NEXT ROW ONLY
-- case
select 1 from dual
-- case
select 1 from dual limit 1
-- case
select 1 where exists (select 2)
-- case
select 1 from dual where not exists (select 2)
-- case
select 1 as a from dual order by a
-- case
select 1 as a from dual where 1 < any (select 2) order by a
-- case
select 1 order by 1
-- case
(select 1);
-- case
select 1 where 1=1
-- case
select 1 group by 1
-- case
select 1 from dual group by 1
-- case
select min(b) b from (select min(t.b) b from t where t.a = '');
-- case
select min(b) b from (select min(t.b) b from t where t.a = '') as t1;
-- case
SELECT /*!40001 SQL_NO_CACHE */ * FROM test WHERE 1 limit 0, 2000;
-- case
ANALYZE TABLE t
-- case
/** 20180417 **/ show databases;
-- case
/* 20180417 **/ show databases;
-- case
/** 20180417 */ show databases;
-- case
/** 20180417 ******/ show databases;
-- case
/**/show databases;
-- case
/*+*/show databases;
-- case
select/*+*/1;
-- case
/*T*/show databases;
-- case
/*M*/show databases;
-- case
/*!*/show databases;
-- case
/*T!*/show databases;
-- case
/*M!*/show databases;
-- case
BINLOG '
BxSFVw8JAAAA8QAAAPUAAAAAAAQANS41LjQ0LU1hcmlhREItbG9nAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAEzgNAAgAEgAEBAQEEgAA2QAEGggAAAAICAgCAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAA5gm5Mg==
'/*!*/;
-- case
select * from t1 partition (p1)
-- case
select * from t1 partition (p1,p2)
-- case
select * from t1 partition (`p1`, p2, p3)
-- case
select * from t1 partition ()
-- case
split table t1 index idx1 by ('a'),('b'),('c')
-- case
split table t1 index idx1 by (1)
-- case
split table t1 index idx1 by ('abc',123), ('xyz'), ('yz', 1000)
-- case
split table t1 index idx1 by 
-- case
split table t1 index idx1 between ('a') and ('z') regions 10
-- case
split table t1 index idx1 between ('a',1) and ('z',2) regions 10
-- case
split table t1 index idx1 between () and () regions 10
-- case
split table t1 index by (1)
-- case
split region for table t1 index idx1 by ('a'),('b'),('c')
-- case
split partition table t1 index idx1 by ('a'),('b'),('c')
-- case
split region for partition table t1 index idx1 by ('a'),('b'),('c')
-- case
split region for table t1 index idx1 between ('a') and ('z') regions 10
-- case
split partition table t1 index idx1 between ('a') and ('z') regions 10
-- case
split region for partition table t1 index idx1 between ('a') and ('z') regions 10
-- case
split region for table t1 partition (p0,p1) index idx1 by ('a'),('b'),('c')
-- case
split partition table t1 partition (p0) index idx1 by ('a'),('b'),('c')
-- case
split region for partition table t1 partition (p0) index idx1 by ('a'),('b'),('c')
-- case
split region for table t1 partition (p0) index idx1 between ('a') and ('z') regions 10
-- case
split partition table t1 partition (p0) index idx1 between ('a') and ('z') regions 10
-- case
split region for partition table t1 partition (p0) index idx1 between ('a') and ('z') regions 10
-- case
split table t1 by ('a'),('b'),('c')
-- case
split table t1 by (1)
-- case
split table t1 by ('abc',123), ('xyz'), ('yz', 1000)
-- case
split table t1 by 
-- case
split table t1 between ('a') and ('z') regions 10
-- case
split table t1 between ('a',1) and ('z',2) regions 10
-- case
split table t1 between () and () regions 10
-- case
split region for table t1 by ('a'),('b'),('c')
-- case
split partition table t1 by ('a'),('b'),('c')
-- case
split region for partition table t1 by ('a'),('b'),('c')
-- case
split region for table t1 between (1) and (1000) regions 10
-- case
split partition table t1 between (1) and (1000) regions 10
-- case
split region for partition table t1 between (1) and (1000) regions 10
-- case
show table t1 regions
-- case
show table t1 regions where a=1
-- case
show table t1
-- case
show table t1 index idx1 regions
-- case
show table t1 index idx1 regions where a=2
-- case
show table t1 index idx1
-- case
show table t1 partition (p0,p1) regions
-- case
show table t1 partition (p0) regions where a=1
-- case
show table t1 partition
-- case
show table t1 partition (p0) index idx1 regions
-- case
show table t1 partition (p0,p1) index idx1 regions where a=2
-- case
show table t1 partition index idx1
-- case
show table t1 distributions
-- case
show table t1 distributions where a=1
-- case
show table t1 partition (p0,p1) distributions
-- case
show table t1 partition (p0,p1) distributions where a=1
-- case
distribute table t1
-- case
distribute table t1 partition(p0)
-- case
distribute table t1 partition(p0,p1)
-- case
distribute table t1 partition(p0,p1) engine = tikv
-- case
distribute table t1 rule = 'leader-scatter' engine = 'tikv'
-- case
distribute table t1 rule = "leader-scatter" engine = "tikv"
-- case
distribute table t1 partition(p0,p1) rule = 'learner-scatter' engine = 'tikv'
-- case
distribute table t1 partition(p0) rule = 'peer-scatter' engine = 'tiflash'
-- case
distribute table t1 partition(p0) rule = 'peer-scatter' engine = 'tiflash' timeout = '30m'
-- case
show distribution jobs 1
-- case
show distribution jobs
-- case
show distribution jobs where id > 0
-- case
show distribution job 1 where id > 0
-- case
show distribution job 1
-- case
cancel distribution job
-- case
cancel distribution job 1
-- case
show table t1.t1 next_row_id
-- case
show table t1 next_row_id
-- case
show table next_row_id
-- case
begin pessimistic
-- case
begin optimistic
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (a int)
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (a char(1))
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (a char(1), b int)
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (c1 TIME(2), c2 DATETIME(2), c3 TIMESTAMP(2));
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (a TINYINT UNSIGNED);
-- case
ADMIN REPAIR TABLE t CREATE TABLE t (name CHAR(50) CHARACTER SET UTF8)
-- case
ALTER INSTANCE RELOAD TLS
-- case
ALTER INSTANCE RELOAD TLS NO ROLLBACK ON ERROR
-- case
ALTER RANGE global PLACEMENT POLICY mypolicy
-- case
ALTER RANGE global PLACEMENT POLICY default
-- case
ALTER RANGE meta PLACEMENT POLICY mypolicy
-- case
CREATE SEQUENCE seq INCREMENT - 9223372036854775807
-- case
CREATE SEQUENCE seq INCREMENT - 9223372036854775808
-- case
CREATE SEQUENCE seq INCREMENT -9223372036854775808
-- case
CREATE SEQUENCE seq INCREMENT -9223372036854775809
-- case
select `t`.`1a`.1 from t;
-- case
select * from 1db.1table;
-- case
select * from t where t. status = 1;
-- case
SHOW PLACEMENT
-- case
SHOW PLACEMENT LIKE 'POLICY foo%'
-- case
SHOW PLACEMENT WHERE Target='TABLE test.t1'
-- case
SHOW PLACEMENT FOR DATABASE db1
-- case
SHOW PLACEMENT FOR SCHEMA db1
-- case
SHOW PLACEMENT FOR TABLE tb1
-- case
SHOW PLACEMENT FOR TABLE db1.tb1
-- case
SHOW PLACEMENT FOR TABLE tb1 PARTITION p1
-- case
SHOW PLACEMENT FOR TABLE db1.tb1 PARTITION p1
-- case
SHOW PLACEMENT FOR
-- case
SHOW PLACEMENT DATABASE db1
-- case
SHOW PLACEMENT FOR DB db1
-- case
SHOW PLACEMENT FOR DATABASE db1 TABLE tb1
-- case
SHOW PLACEMENT FOR PARTITION p1
-- case
SHOW PLACEMENT FOR DB LIKE '%'
-- case
SHOW PLACEMENT FOR DB db1 LIKE '%'
-- case
SHOW PLACEMENT LABELS
-- case
SHOW PLACEMENT LABELS LIKE '%zone%'
-- case
SHOW PLACEMENT LABELS WHERE label='l123'
-- case
SHOW SESSION_STATES
-- case
SET SESSION_STATES 'x'
-- case
SET SESSION_STATES
-- case
SET SESSION_STATES 1
-- case
SET SESSION_STATES now()
-- case
calibrate resource
-- case
calibrate resource START_TIME '2021-04-15 00:00:00'
-- case
calibrate resource START_TIME '2023-04-01 13:00:00' END_TIME '2023-04-01 16:00:00'
-- case
calibrate resource START_TIME '2023-04-01 13:00:00' DURATION '20m'
-- case
calibrate resource START_TIME '2023-04-01 13:00:00' END_TIME '2023-04-01 16:00:00' DURATION '20m'
-- case
calibrate resource START_TIME '2023-04-01 13:00:00',END_TIME='2023-04-01 16:00:00'
-- case
calibrate resource START_TIME '2023-04-01 13:00:00',DURATION='20m'
-- case
calibrate resource DURATION='20m' START_TIME '2023-04-01 13:00:00'
-- case
calibrate resource   START_TIME '2023-04-01 13:00:00' END_TIME='2023-04-01 16:00:00',DURATION '20m'
-- case
calibrate resource START_TIME CURRENT_TIMESTAMP() END_TIME current_timestamp()
-- case
calibrate resource END_TIME now()
-- case
calibrate resource START_TIME now()
-- case
calibrate resource START_TIME NOW() END_TIME now()
-- case
calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 10 minute END_TIME now()
-- case
calibrate resource START_TIME now() - 1000 END_TIME current_timestamp()
-- case
calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute DURATION interval 15 minute
-- case
calibrate resource START_TIME CURRENT_TIMESTAMP() - interval 20 minute DURATION '15m'
-- case
calibrate resource END_TIME now() START_TIME CURRENT_TIMESTAMP() - interval 20 minute
-- case
calibrate resource workload
-- case
calibrate resource workload tpcc
-- case
calibrate resource workload oltp_read_write
-- case
calibrate resource workload oltp_read_only
-- case
calibrate resource workload oltp_write_only
-- case
calibrate resource workload = oltp_read_write START_TIME '2023-04-01 13:00:00'
-- case
query watch add SQL DIGEST b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377 
-- case
query watch add SQL DIGEST b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377 
-- case
query watch add SQL DIGEST 'b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377' 
-- case
query watch add PLAN DIGEST `5e3ddd388f6012e328233dbcdda5d48f404e0536c6c54d9618233210f3d5762a` 
-- case
query watch add PLAN DIGEST @digest1 
-- case
query watch add SQL TEXT SIMILAR to 'select 1'
-- case
query watch add SQL TEXT EXACT to 'select 1'
-- case
query watch add SQL TEXT PLAN to 'select 1'
-- case
query watch add resource group `default` SQL TEXT SIMILAR to 'select 1'
-- case
query watch add resource group @rg SQL TEXT SIMILAR to @sql1
-- case
query watch add resource group rg1 SQL TEXT SIMILAR to 'select 1'
-- case
query watch add SQL TEXT SIMILAR to 'select 1' resource group rg1
-- case
query watch add ACTION = KILL SQL TEXT SIMILAR to 'select 1'
-- case
query watch add ACTION COOLDOWN resource group rg1 SQL TEXT SIMILAR to 'select 1'
-- case
query watch add resource group `default` resource group `rg1` SQL TEXT SIMILAR to 'select 1'
-- case
query watch add SQL SIMILAR to 'select 1'
-- case
query watch add SQL TEXT SIMILAR 'select 1'
-- case
query watch remove 1
-- case
query watch remove resource group rg1
-- case
query watch remove resource group @rg
-- case
query watch remove
-- case
replace /*+ SET_VAR(sql_mode='ALLOW_INVALID_DATES') */ into t values ('2004-04-31');
