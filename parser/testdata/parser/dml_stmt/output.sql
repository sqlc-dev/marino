
-- case

-- case
INSERT INTO `foo` VALUES (1234)
-- case
INSERT INTO `foo` VALUES (1234,5678)
-- case
INSERT INTO `t1` (SELECT * FROM `t2`)
-- case
INSERT INTO `t` PARTITION(`p0`) VALUES (1234)
-- case
REPLACE INTO `t` PARTITION(`p0`) VALUES (1234)
-- case
INSERT INTO `t` PARTITION(`p0`, `p1`, `p2`) VALUES (1234)
-- case
REPLACE INTO `t` PARTITION(`p0`, `p1`, `p2`) VALUES (1234)
-- case
INSERT INTO `foo` VALUES (1 OR 2)
-- case
INSERT INTO `foo` VALUES (1|2)
-- case
INSERT INTO `foo` VALUES (FALSE OR TRUE)
-- case
INSERT INTO `foo` VALUES (BAR(5678))
-- case
INSERT INTO `foo` VALUES ()
-- case
SELECT * FROM `t`
-- case
SELECT * FROM `t` AS `u`
-- case
SELECT * FROM (`t`) JOIN `v`
-- case
SELECT * FROM (`t` AS `u`) JOIN `v`
-- case
SELECT * FROM (`t`) JOIN `v` AS `w`
-- case
SELECT * FROM (`t` AS `u`) JOIN `v` AS `w`
-- case
SELECT * FROM ((`foo`) JOIN `bar`) JOIN `foo`
-- case
-- error: line 1 column 23 near "FROM t" 
-- case
SELECT DISTINCT * FROM `t`
-- case
SELECT DISTINCT * FROM `t`
-- case
SELECT ALL * FROM `t`
-- case
-- error: [parser:1221]Incorrect usage of ALL and DISTINCT
-- case
-- error: [parser:1221]Incorrect usage of ALL and DISTINCT
-- case
INSERT INTO `foo` (`a`) VALUES (42)
-- case
-- error: line 1 column 20 near ") VALUES (42,)" 
-- case
INSERT INTO `foo` (`a`,`b`) VALUES (42,314)
-- case
-- error: line 1 column 22 near ") VALUES (42,314)" 
-- case
-- error: line 1 column 22 near ") VALUES (42,314,)" 
-- case
INSERT INTO `foo` () VALUES ()
-- case
INSERT INTO `foo` VALUES ()
-- case
INSERT INTO `tt` VALUES (1000001783)
-- case
INSERT INTO `tt` VALUES (DEFAULT)
-- case
REPLACE INTO `foo` VALUES (1 OR 2)
-- case
REPLACE INTO `foo` VALUES (1|2)
-- case
REPLACE INTO `foo` VALUES (FALSE OR TRUE)
-- case
REPLACE INTO `foo` VALUES (BAR(5678))
-- case
REPLACE INTO `foo` VALUES ()
-- case
REPLACE INTO `foo` (`a`,`b`) VALUES (42,314)
-- case
-- error: line 1 column 23 near ") VALUES (42,314)" 
-- case
-- error: line 1 column 23 near ") VALUES (42,314,)" 
-- case
REPLACE INTO `foo` () VALUES ()
-- case
REPLACE INTO `foo` VALUES ()
-- case
SELECT `stuff`.`id` FROM `stuff` WHERE `stuff`.`value`>=ALL (SELECT `stuff`.`value` FROM `stuff`)
-- case
START TRANSACTION
-- case
START TRANSACTION
-- case
COMMIT
-- case
COMMIT
-- case
COMMIT
-- case
COMMIT
-- case
COMMIT RELEASE
-- case
COMMIT AND CHAIN
-- case
-- error: line 1 column 24 near "RELEASE" 
-- case
ROLLBACK
-- case
ROLLBACK
-- case
ROLLBACK
-- case
ROLLBACK
-- case
ROLLBACK RELEASE
-- case
ROLLBACK AND CHAIN
-- case
-- error: line 1 column 26 near "RELEASE" 
-- case
START TRANSACTION; INSERT INTO `foo` VALUES (42,3.14); INSERT INTO `foo` VALUES (-1,2.78); COMMIT
-- case
START TRANSACTION; INSERT INTO `tmp` SELECT * FROM `bar`; SELECT * FROM `tmp`; ROLLBACK
-- case
SAVEPOINT x
-- case
RELEASE SAVEPOINT x
-- case
ROLLBACK TO x
-- case
ROLLBACK TO X
-- case
ROLLBACK TO x
-- case
TABLE `t`
-- case
(TABLE `t`)
-- case
-- error: line 1 column 9 near ", t2" 
-- case
TABLE `t` ORDER BY `b`
-- case
TABLE `t` LIMIT 3
-- case
TABLE `t` ORDER BY `b` LIMIT 3
-- case
TABLE `t` ORDER BY `b` LIMIT 2,3
-- case
TABLE `t` ORDER BY `b` LIMIT 2,3
-- case
INSERT INTO `ta` TABLE `tb`
-- case
INSERT INTO `t`.`a` TABLE `t`.`b`
-- case
REPLACE INTO `ta` TABLE `tb`
-- case
REPLACE INTO `t`.`a` TABLE `t`.`b`
-- case
TABLE `t1` INTO OUTFILE 'a.txt'
-- case
TABLE `t` ORDER BY `a` INTO OUTFILE '/tmp/abc'
-- case
CREATE TABLE `t`.`a` AS TABLE `t`.`b`
-- case
CREATE TABLE `ta` AS TABLE `tb`
-- case
CREATE TABLE `ta` (`x` INT) AS TABLE `tb`
-- case
CREATE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY DEFINER VIEW `v` AS TABLE `t`
-- case
CREATE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY DEFINER VIEW `v` AS (TABLE `t`)
-- case
SELECT * FROM `t1` WHERE `a` IN (TABLE `t2`)
-- case
VALUES ROW(1)
-- case
VALUES ROW()
-- case
VALUES ROW(1,DEFAULT)
-- case
VALUES ROW(1), ROW(2,3)
-- case
-- error: line 1 column 8 near "(1,2)" 
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8)
-- case
VALUES ROW(1,`s`,3.1), ROW(5,`y`,9.9)
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) LIMIT 3
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) ORDER BY `a`
-- case
VALUES ROW(1,-2,3), ROW(5,7,9), ROW(4,6,8) ORDER BY `a` LIMIT 2
-- case
VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE 'a.txt'
-- case
VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY `a` INTO OUTFILE '/tmp/abc'
-- case
CREATE TABLE `ta` AS VALUES ROW(1)
-- case
CREATE TABLE `ta` AS VALUES ROW(1)
-- case
CREATE ALGORITHM = UNDEFINED DEFINER = CURRENT_USER SQL SECURITY DEFINER VIEW `a` AS VALUES ROW(1)
-- case
SELECT `a`.`b`.`c` FROM `t`
-- case
-- error: line 1 column 13 near ".c FROM t" 
-- case
SELECT `a`.`b`.* FROM `t`
-- case
SELECT `a` FROM `t`
-- case
-- error: line 1 column 13 near ".d FROM t" 
-- case
DO 1
-- case
DO 1, SLEEP(1)
-- case
-- error: line 1 column 9 near "from t" 
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t1` FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' IGNORE 1 LINES
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t`
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` CHARACTER SET utf8
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS TERMINATED BY 'ab'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS TERMINATED BY 'ab'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` LINES STARTING BY 'ab'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` LINES STARTING BY 'ab' TERMINATED BY 'xy'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS TERMINATED BY 'ab' LINES TERMINATED BY 'xy'
-- case
-- error: line 1 column 53 near "terminated by 'xy' fields terminated by 'ab'" 
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t`
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` CHARACTER SET utf8 FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` LINES STARTING BY 'ab'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` LINES STARTING BY 'ab' TERMINATED BY 'xy'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' LINES TERMINATED BY 'xy'
-- case
-- error: line 1 column 59 near "terminated by 'xy' fields terminated by 'ab'" 
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` CHARACTER SET utf8 FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` LINES STARTING BY 'ab' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` LINES STARTING BY 'ab' TERMINATED BY 'xy' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` CHARACTER SET utf8 FIELDS TERMINATED BY 'ab' LINES TERMINATED BY 'xy' (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' LINES TERMINATED BY 'xy' (`a`,`b`)
-- case
-- error: line 1 column 61 near "fields terminated by 'ab'" 
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` IGNORE 1 LINES
-- case
-- error: line 1 column 57 near "-1 lines" 
-- case
-- error: line 1 column 103 near "ignore 1 lines" 
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` LINES STARTING BY 'ab' TERMINATED BY 'xy' IGNORE 1 LINES
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '*' IGNORE 1 LINES (`a`,`b`)
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY ''
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` FIELDS TERMINATED BY 'kk'
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` FIELDS TERMINATED BY 'kk' ENCLOSED BY ''
-- case
-- error: [parser:1083]Field separator argument is not what is expected; check the manual
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` FIELDS TERMINATED BY 'kk'
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` FIELDS TERMINATED BY 'kk' ENCLOSED BY ''
-- case
-- error: [parser:1083]Field separator argument is not what is expected; check the manual
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` LINES STARTING BY 'kk' TERMINATED BY 'kk'
-- case
LOAD DATA LOCAL INFILE '~/1.csv' IGNORE INTO TABLE `t_ascii` LINES STARTING BY 'kk' TERMINATED BY 'kk'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY ''
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t` FIELDS TERMINATED BY 'ab' ENCLOSED BY 'b' ESCAPED BY '' SET `b`=CAST(CONV(MID(@`var1`, 3, LENGTH(@`var1`)-3), 2, 10) AS UNSIGNED)
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS ENCLOSED BY ''
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS ENCLOSED BY 'a'
-- case
-- error: [parser:1083]Field separator argument is not what is expected; check the manual
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS ESCAPED BY ''
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS ESCAPED BY 'a'
-- case
-- error: [parser:1083]Field separator argument is not what is expected; check the manual
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE `t1` (`column1`,@`dummy`,`column2`,@`dummy`,`column3`)
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE `t1` (`column1`,@`var1`) SET `column2`=@`var1`/100
-- case
LOAD DATA INFILE 'file.txt' INTO TABLE `t1` (`column1`,@`var1`,@`var2`) SET `column2`=@`var1`/100, `column3`=DEFAULT, `column4`=CURRENT_TIMESTAMP(), `column5`=@`var2`+1
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t1` FIELDS TERMINATED BY ',' LINES TERMINATED BY '
'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t1` FIELDS TERMINATED BY ',' LINES TERMINATED BY '
'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' IGNORE INTO TABLE `t1` FIELDS TERMINATED BY ',' LINES TERMINATED BY '
'
-- case
LOAD DATA LOCAL INFILE '/tmp/t.csv' REPLACE INTO TABLE `t1` FIELDS TERMINATED BY ',' LINES TERMINATED BY '
'
-- case
LOAD DATA INFILE 's3://bucket-name/t.csv' INTO TABLE `t`
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS DEFINED NULL BY 'nil'
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS DEFINED NULL BY ' '
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` FIELDS DEFINED NULL BY 'NULL' OPTIONALLY ENCLOSED IGNORE 1 LINES
-- case
LOAD DATA INFILE '/tmp/t.csv' FORMAT 'delimited data' INTO TABLE `t` (`column1`,@`var1`) SET `column2`=@`var1`/100
-- case
LOAD DATA LOCAL INFILE '/tmp/t.sql' FORMAT 'sql file' REPLACE INTO TABLE `t` (`a`,`b`)
-- case
LOAD DATA INFILE '/tmp/t.parquet' FORMAT 'parquet' INTO TABLE `t` (`column1`,@`var1`) SET `column2`=@`var1`/100
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` WITH detached
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` WITH threads=10
-- case
LOAD DATA INFILE '/tmp/t.csv' INTO TABLE `t` WITH threads=10, detached
-- case
IMPORT INTO `t` FROM '/file.csv'
-- case
IMPORT INTO `t` (`a`,`b`) FROM '/file.csv'
-- case
IMPORT INTO `t` (`a`,@`1`) FROM '/file.csv'
-- case
IMPORT INTO `t` (`a`,@`1`) SET `b`=@`1`+100 FROM '/file.csv'
-- case
IMPORT INTO `t` FROM '/file.csv' FORMAT 'sql file'
-- case
IMPORT INTO `t` FROM '/file.csv' WITH detached
-- case
IMPORT INTO `t` FROM '/file.csv' WITH thread=1
-- case
IMPORT INTO `t` FROM '/file.csv' WITH detached, thread=1
-- case
SELECT * FROM `t` FOR UPDATE
-- case
SELECT * FROM `t` FOR SHARE
-- case
SELECT * FROM `t` FOR UPDATE NOWAIT
-- case
SELECT * FROM `t` FOR UPDATE WAIT 5
-- case
SELECT * FROM `t` LIMIT 1 FOR UPDATE WAIT 11
-- case
SELECT * FROM `t` FOR SHARE NOWAIT
-- case
SELECT * FROM `t` FOR UPDATE SKIP LOCKED
-- case
SELECT * FROM `t` FOR SHARE SKIP LOCKED
-- case
SELECT * FROM `t` FOR SHARE
-- case
-- error: line 1 column 41 near "nowait" 
-- case
-- error: line 1 column 39 near "skip locked" 
-- case
SELECT * FROM `t` FOR UPDATE OF `t`
-- case
SELECT * FROM `t` FOR SHARE OF `t`
-- case
SELECT * FROM `t` FOR UPDATE OF `t` NOWAIT
-- case
SELECT * FROM `t` FOR UPDATE OF `t` WAIT 5
-- case
SELECT * FROM `t` LIMIT 1 FOR UPDATE OF `t` WAIT 11
-- case
SELECT * FROM `t` FOR SHARE OF `t` NOWAIT
-- case
SELECT * FROM `t` FOR UPDATE OF `t` SKIP LOCKED
-- case
SELECT * FROM `t` FOR SHARE OF `t` SKIP LOCKED
-- case
SELECT `a`,`b` FROM `t` INTO OUTFILE '/tmp/result.txt'
-- case
SELECT `a` FROM `t` ORDER BY `a` INTO OUTFILE '/tmp/abc'
-- case
SELECT 1 INTO OUTFILE '/tmp/1.csv'
-- case
SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/1.csv'
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ','
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' ENCLOSED BY '"'
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"'
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' LINES TERMINATED BY '
'
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES TERMINATED BY ''
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY ''
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"' LINES STARTING BY 'xy' TERMINATED BY ''
-- case
SELECT `a`,`b`,`a`+`b` FROM `t` INTO OUTFILE '/tmp/result.txt' FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES STARTING BY 'xy' TERMINATED BY ''
-- case
SELECT * FROM ((`t1`) JOIN `t2`) JOIN `t3`
-- case
SELECT * FROM (`t1` JOIN `t2`) LEFT JOIN `t3` ON `t2`.`id`=`t3`.`id`
-- case
SELECT * FROM (`t1` RIGHT JOIN `t2` ON `t1`.`id`=`t2`.`id`) LEFT JOIN `t3` ON `t3`.`id`=`t2`.`id`
-- case
-- error: line 1 column 60 near "" 
-- case
SELECT * FROM (`t1` JOIN `t2`) LEFT JOIN `t3` USING (`id`)
-- case
SELECT * FROM (`t1` RIGHT JOIN `t2` USING (`id`)) LEFT JOIN `t3` USING (`id`)
-- case
-- error: line 1 column 54 near "" 
-- case
SELECT * FROM `t1` NATURAL JOIN `t2`
-- case
SELECT * FROM `t1` NATURAL RIGHT JOIN `t2`
-- case
SELECT * FROM `t1` NATURAL LEFT JOIN `t2`
-- case
-- error: line 1 column 30 near "inner join t2" 
-- case
-- error: line 1 column 30 near "cross join t2" 
-- case
SELECT * FROM `t3` JOIN (`t1` JOIN `t2` ON `t1`.`a`=`t2`.`a`) ON `t3`.`b`=`t2`.`b`
-- case
SELECT * FROM `t1` STRAIGHT_JOIN `t2` ON `t1`.`id`=`t2`.`id`
-- case
SELECT STRAIGHT_JOIN * FROM `t1` JOIN `t2` ON `t1`.`id`=`t2`.`id`
-- case
SELECT STRAIGHT_JOIN * FROM `t1` LEFT JOIN `t2` ON `t1`.`id`=`t2`.`id`
-- case
SELECT STRAIGHT_JOIN * FROM `t1` RIGHT JOIN `t2` ON `t1`.`id`=`t2`.`id`
-- case
SELECT STRAIGHT_JOIN * FROM `t1` STRAIGHT_JOIN `t2` ON `t1`.`id`=`t2`.`id`
-- case
DELETE FROM `t1`
-- case
-- error: line 1 column 16 near "" 
-- case
DELETE LOW_PRIORITY FROM `t1`
-- case
DELETE QUICK FROM `t1`
-- case
DELETE IGNORE FROM `t1`
-- case
DELETE LOW_PRIORITY QUICK IGNORE FROM `t1`
-- case
DELETE FROM `t1` WHERE `t1`.`a`>0 ORDER BY `t1`.`a`
-- case
DELETE FROM `t1` WHERE `a`=26
-- case
DELETE FROM `t1` WHERE `a`=1 LIMIT 1
-- case
DELETE FROM `t1` WHERE `t1`.`a`>0 ORDER BY `t1`.`a` LIMIT 1
-- case
DELETE FROM `x`.`y` AS `z` WHERE `z`.`a`>0
-- case
DELETE FROM `t1` AS `w` WHERE `a`>0
-- case
DELETE FROM `t1` PARTITION(`p0`, `p1`)
-- case
DELETE LOW_PRIORITY `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE QUICK `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE IGNORE `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE IGNORE `t1`,`t2` FROM (`t1` PARTITION(`p0`, `p1`)) JOIN `t2`
-- case
DELETE LOW_PRIORITY QUICK IGNORE `t1`,`t2` FROM (`t1`) JOIN `t2` WHERE `t1`.`a`>5
-- case
DELETE `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE `t1`,`t2` FROM (`t1`) JOIN `t2` WHERE `t1`.`a`=1 AND `t2`.`b`!=1
-- case
DELETE `t1` FROM (`t1`) JOIN `t2`
-- case
DELETE `t2` FROM (`t1`) JOIN `t2`
-- case
DELETE `t1` FROM `t1`
-- case
DELETE `t1`,`t2`,`t3` FROM ((`t1`) JOIN `t2`) JOIN `t3`
-- case
DELETE `t1`,`t2`,`t3` FROM ((`t1`) JOIN `t2`) JOIN `t3` WHERE `t3`.`c`<5 AND `t1`.`a`=3
-- case
DELETE `t1` FROM (`t1`) JOIN `t1` AS `t2` WHERE `t1`.`b`=`t2`.`b` AND `t1`.`a`>`t2`.`a`
-- case
DELETE `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE `t`.`t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE `t1`,`t2` FROM (`t1`) JOIN `t2`
-- case
DELETE `t11`,`t12` FROM (`t11`) JOIN `t12` WHERE `t11`.`a`=`t12`.`a` AND `t11`.`b`!=1
-- case
DELETE QUICK FROM `t1`,`t2` USING (`t1`) JOIN `t2`
-- case
DELETE LOW_PRIORITY IGNORE FROM `t1`,`t2` USING (`t1`) JOIN `t2`
-- case
DELETE LOW_PRIORITY QUICK IGNORE FROM `t1`,`t2` USING (`t1`) JOIN `t2`
-- case
DELETE FROM `t1` USING `t1` WHERE `post`=_UTF8MB4'1'
-- case
DELETE FROM `t1`,`t2` USING (`t1`) JOIN `t2`
-- case
DELETE FROM `t1`,`t2`,`t3` USING ((`t1`) JOIN `t2`) JOIN `t3` WHERE `t3`.`a`=1
-- case
DELETE FROM `t2`,`t3` USING ((`t1`) JOIN `t2`) JOIN `t3` WHERE `t1`.`a`=1
-- case
DELETE FROM `t2`,`t3` USING ((`t1`) JOIN `t2`) JOIN `t3` WHERE `t1`.`a`=1
-- case
DELETE FROM `t1`,`t2`,`t3` USING ((`t1`) JOIN `t2`) JOIN `t3` WHERE `t1`.`a`=1
-- case
DELETE `t1`,`t2` FROM (`t1` JOIN `t2`) JOIN `t3` WHERE `t1`.`id`=`t2`.`id` AND `t2`.`id`=`t3`.`id`
-- case
DELETE FROM `t1`,`t2` USING (`t1` JOIN `t2`) JOIN `t3` WHERE `t1`.`id`=`t2`.`id` AND `t2`.`id`=`t3`.`id`
-- case
DELETE /*+ TIDB_INLJ(`t1`, `t2`)*/ `t1`,`t2` FROM (`t1`) JOIN `t2` WHERE `t1`.`id`=`t2`.`id`
-- case
DELETE /*+ TIDB_HJ(`t1`, `t2`)*/ `t1`,`t2` FROM (`t1`) JOIN `t2` WHERE `t1`.`id`=`t2`.`id`
-- case
DELETE /*+ TIDB_SMJ(`t1`, `t2`)*/ `t1`,`t2` FROM (`t1`) JOIN `t2` WHERE `t1`.`id`=`t2`.`id`
-- case
DELETE FROM `t1` USE INDEX (`idx_a`) WHERE `t1`.`id`=1
-- case
DELETE `t1`,`t2` FROM `t1` USE INDEX (`idx_a`) JOIN `t2` WHERE `t1`.`id`=`t2`.`id`
-- case
DELETE `t1`,`t2` FROM `t1` USE INDEX (`idx_a`) JOIN `t2` USE INDEX (`idx_a`) WHERE `t1`.`id`=`t2`.`id`
-- case
-- error: line 1 column 89 near "limit 10;" 
-- case
-- error: line 1 column 89 near "order by t1.id;" 
-- case
INSERT INTO `t` (`a`,`b`,`c`) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE `c`=VALUES(`a`)+VALUES(`b`)
-- case
INSERT INTO `t` (`a`,`b`,`c`) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE `c`=VALUES(`a`)+VALUES(`b`)
-- case
INSERT IGNORE INTO `t` (`a`,`b`,`c`) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE `c`=VALUES(`a`)+VALUES(`b`)
-- case
INSERT IGNORE INTO `t` (`a`,`b`,`c`) VALUES (1,2,3),(4,5,6) ON DUPLICATE KEY UPDATE `c`=VALUES(`a`)+VALUES(`b`)
-- case
INSERT INTO `odk` (`id`,`v`) VALUES (1,10) AS `new` ON DUPLICATE KEY UPDATE `v`=`new`.`v`
-- case
INSERT INTO `odk` (`id`,`v`) VALUES (1,10) AS `new`(`nid`,`nv`) ON DUPLICATE KEY UPDATE `v`=`nv`
-- case
INSERT INTO `odk` VALUES (1,10) AS `new` ON DUPLICATE KEY UPDATE `v`=`new`.`v`
-- case
INSERT IGNORE INTO `odk` (`id`,`v`) VALUES (1,10) AS `new` ON DUPLICATE KEY UPDATE `v`=`new`.`v`
-- case
INSERT INTO `odk` SET `id`=1,`v`=10 AS `new` ON DUPLICATE KEY UPDATE `v`=`new`.`v`
-- case
-- error: line 1 column 57 near "AS new ON DUPLICATE KEY UPDATE v = new.v;" 
-- case
INSERT INTO `t` SET `a`=1,`b`=2
-- case
-- error: line 1 column 21 near "SET a=1" 
-- case
UPDATE LOW_PRIORITY IGNORE `t` SET `id`=`id`+1 ORDER BY `id` DESC
-- case
UPDATE `t` SET `id`=`id`+1 ORDER BY `id` DESC
-- case
UPDATE `t` SET `id`=`id`+1 ORDER BY `id` DESC LIMIT 3
-- case
UPDATE `t` SET `id`=`id`+1, `name`=_UTF8MB4'jojo'
-- case
UPDATE (`items`) JOIN `month` SET `items`.`price`=`month`.`price` WHERE `items`.`id`=`month`.`id`
-- case
UPDATE `user` AS `T0` LEFT JOIN `user_profile` AS `T1` ON `T1`.`id`=`T0`.`profile_id` SET `T0`.`profile_id`=1 WHERE `T0`.`profile_id` IN (1)
-- case
UPDATE (`t1`) JOIN `t2` SET `t1`.`profile_id`=1, `t2`.`profile_id`=1 WHERE `ta`.`a`=`t`.`ba`
-- case
UPDATE /*+ TIDB_INLJ(`t1`, `t2`)*/ (`t1`) JOIN `t2` SET `t1`.`profile_id`=1, `t2`.`profile_id`=1 WHERE `ta`.`a`=`t`.`ba`
-- case
UPDATE /*+ TIDB_SMJ(`t1`, `t2`)*/ (`t1`) JOIN `t2` SET `t1`.`profile_id`=1, `t2`.`profile_id`=1 WHERE `ta`.`a`=`t`.`ba`
-- case
UPDATE /*+ TIDB_HJ(`t1`, `t2`)*/ (`t1`) JOIN `t2` SET `t1`.`profile_id`=1, `t2`.`profile_id`=1 WHERE `ta`.`a`=`t`.`ba`
-- case
-- error: line 1 column 76 near "LIMIT 10;" 
-- case
-- error: line 1 column 76 near "order by month.id;" 
-- case
UPDATE `t1` USE INDEX (`idx_a`) SET `t1`.`price`=3.25 WHERE `t1`.`id`=1
-- case
UPDATE `t1` USE INDEX (`idx_a`) JOIN `t2` SET `t1`.`price`=`t2`.`price` WHERE `t1`.`id`=`t2`.`id`
-- case
UPDATE `t1` USE INDEX (`idx_a`) JOIN `t2` USE INDEX (`idx_a`) SET `t1`.`price`=`t2`.`price` WHERE `t1`.`id`=`t2`.`id`
-- case
SELECT * FROM `t` WHERE 1=1
-- case
SELECT * FROM `t` LIMIT 5
-- case
SELECT * FROM `t` LIMIT 5
-- case
SELECT * FROM `t` LIMIT 5
-- case
SELECT * FROM `t` LIMIT 5
-- case
SELECT * FROM `t` LIMIT 1
-- case
SELECT * FROM `t` LIMIT 1
-- case
SELECT 1
-- case
SELECT 1 LIMIT 1
-- case
SELECT 1 FROM DUAL WHERE EXISTS (SELECT 2)
-- case
SELECT 1 FROM DUAL WHERE NOT EXISTS (SELECT 2)
-- case
SELECT 1 AS `a` ORDER BY `a`
-- case
SELECT 1 AS `a` FROM DUAL WHERE 1<ANY (SELECT 2) ORDER BY `a`
-- case
SELECT 1 ORDER BY 1
-- case
(SELECT 1)
-- case
SELECT 1 FROM DUAL WHERE 1=1
-- case
SELECT 1 GROUP BY 1
-- case
SELECT 1 GROUP BY 1
-- case
SELECT MIN(`b`) AS `b` FROM (SELECT MIN(`t`.`b`) AS `b` FROM `t` WHERE `t`.`a`=_UTF8MB4'')
-- case
SELECT MIN(`b`) AS `b` FROM (SELECT MIN(`t`.`b`) AS `b` FROM `t` WHERE `t`.`a`=_UTF8MB4'') AS `t1`
-- case
SELECT SQL_NO_CACHE * FROM `test` WHERE 1 LIMIT 0,2000
-- case
ANALYZE TABLE `t`
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SELECT 1
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
SHOW DATABASES
-- case
BINLOG '
BxSFVw8JAAAA8QAAAPUAAAAAAAQANS41LjQ0LU1hcmlhREItbG9nAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAEzgNAAgAEgAEBAQEEgAA2QAEGggAAAAICAgCAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAA5gm5Mg==
'
-- case
SELECT * FROM `t1` PARTITION(`p1`)
-- case
SELECT * FROM `t1` PARTITION(`p1`, `p2`)
-- case
SELECT * FROM `t1` PARTITION(`p1`, `p2`, `p3`)
-- case
-- error: line 1 column 29 near ")" 
-- case
SPLIT TABLE `t1` INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT TABLE `t1` INDEX `idx1` BY (1)
-- case
SPLIT TABLE `t1` INDEX `idx1` BY (_UTF8MB4'abc',123),(_UTF8MB4'xyz'),(_UTF8MB4'yz',1000)
-- case
-- error: line 1 column 29 near "" 
-- case
SPLIT TABLE `t1` INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT TABLE `t1` INDEX `idx1` BETWEEN (_UTF8MB4'a',1) AND (_UTF8MB4'z',2) REGIONS 10
-- case
SPLIT TABLE `t1` INDEX `idx1` BETWEEN () AND () REGIONS 10
-- case
-- error: line 1 column 23 near "by (1)" 
-- case
SPLIT REGION FOR TABLE `t1` INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT PARTITION TABLE `t1` INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR PARTITION TABLE `t1` INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR TABLE `t1` INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT PARTITION TABLE `t1` INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT REGION FOR PARTITION TABLE `t1` INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT REGION FOR TABLE `t1` PARTITION(`p0`, `p1`) INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT PARTITION TABLE `t1` PARTITION(`p0`) INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR PARTITION TABLE `t1` PARTITION(`p0`) INDEX `idx1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR TABLE `t1` PARTITION(`p0`) INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT PARTITION TABLE `t1` PARTITION(`p0`) INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT REGION FOR PARTITION TABLE `t1` PARTITION(`p0`) INDEX `idx1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT TABLE `t1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT TABLE `t1` BY (1)
-- case
SPLIT TABLE `t1` BY (_UTF8MB4'abc',123),(_UTF8MB4'xyz'),(_UTF8MB4'yz',1000)
-- case
-- error: line 1 column 18 near "" 
-- case
SPLIT TABLE `t1` BETWEEN (_UTF8MB4'a') AND (_UTF8MB4'z') REGIONS 10
-- case
SPLIT TABLE `t1` BETWEEN (_UTF8MB4'a',1) AND (_UTF8MB4'z',2) REGIONS 10
-- case
SPLIT TABLE `t1` BETWEEN () AND () REGIONS 10
-- case
SPLIT REGION FOR TABLE `t1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT PARTITION TABLE `t1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR PARTITION TABLE `t1` BY (_UTF8MB4'a'),(_UTF8MB4'b'),(_UTF8MB4'c')
-- case
SPLIT REGION FOR TABLE `t1` BETWEEN (1) AND (1000) REGIONS 10
-- case
SPLIT PARTITION TABLE `t1` BETWEEN (1) AND (1000) REGIONS 10
-- case
SPLIT REGION FOR PARTITION TABLE `t1` BETWEEN (1) AND (1000) REGIONS 10
-- case
SHOW TABLE `t1` REGIONS
-- case
SHOW TABLE `t1` REGIONS WHERE `a`=1
-- case
-- error: line 1 column 13 near "" 
-- case
SHOW TABLE `t1` INDEX `idx1` REGIONS
-- case
SHOW TABLE `t1` INDEX `idx1` REGIONS WHERE `a`=2
-- case
-- error: line 1 column 24 near "" 
-- case
SHOW TABLE `t1` PARTITION(`p0`, `p1`) REGIONS
-- case
SHOW TABLE `t1` PARTITION(`p0`) REGIONS WHERE `a`=1
-- case
-- error: line 1 column 23 near "" 
-- case
SHOW TABLE `t1` PARTITION(`p0`) INDEX `idx1` REGIONS
-- case
SHOW TABLE `t1` PARTITION(`p0`, `p1`) INDEX `idx1` REGIONS WHERE `a`=2
-- case
-- error: line 1 column 29 near "index idx1" 
-- case
SHOW TABLE `t1` DISTRIBUTIONS
-- case
SHOW TABLE `t1` DISTRIBUTIONS WHERE `a`=1
-- case
SHOW TABLE `t1` PARTITION(`p0`, `p1`) DISTRIBUTIONS
-- case
SHOW TABLE `t1` PARTITION(`p0`, `p1`) DISTRIBUTIONS WHERE `a`=1
-- case
-- error: line 1 column 19 near "" 
-- case
-- error: line 1 column 33 near "" 
-- case
-- error: line 1 column 36 near "" 
-- case
-- error: line 1 column 43 near "engine = tikv" 
-- case
DISTRIBUTE TABLE `t1` RULE = 'leader-scatter' ENGINE = 'tikv'
-- case
DISTRIBUTE TABLE `t1` RULE = 'leader-scatter' ENGINE = 'tikv'
-- case
DISTRIBUTE TABLE `t1` PARTITION(`p0`, `p1`) RULE = 'learner-scatter' ENGINE = 'tikv'
-- case
DISTRIBUTE TABLE `t1` PARTITION(`p0`) RULE = 'peer-scatter' ENGINE = 'tiflash'
-- case
DISTRIBUTE TABLE `t1` PARTITION(`p0`) RULE = 'peer-scatter' ENGINE = 'tiflash' TIMEOUT = '30m'
-- case
-- error: line 1 column 24 near "1" 
-- case
SHOW DISTRIBUTION JOBS
-- case
SHOW DISTRIBUTION JOBS WHERE `id`>0
-- case
-- error: line 1 column 29 near "where id > 0" 
-- case
SHOW DISTRIBUTION JOB 1
-- case
-- error: line 1 column 23 near "" 
-- case
CANCEL DISTRIBUTION JOB 1
-- case
SHOW TABLE `t1`.`t1` NEXT_ROW_ID
-- case
SHOW TABLE `t1` NEXT_ROW_ID
-- case
-- error: line 1 column 22 near "" 
-- case
BEGIN PESSIMISTIC
-- case
BEGIN OPTIMISTIC
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`a` INT)
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`a` CHAR(1))
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`a` CHAR(1),`b` INT)
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`c1` TIME(2),`c2` DATETIME(2),`c3` TIMESTAMP(2))
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`a` TINYINT UNSIGNED)
-- case
ADMIN REPAIR TABLE `t` CREATE TABLE `t` (`name` CHAR(50) CHARACTER SET UTF8)
-- case
ALTER INSTANCE RELOAD TLS
-- case
ALTER INSTANCE RELOAD TLS NO ROLLBACK ON ERROR
-- case
ALTER RANGE `global` PLACEMENT POLICY = `mypolicy`
-- case
ALTER RANGE `global` PLACEMENT POLICY = `default`
-- case
ALTER RANGE `meta` PLACEMENT POLICY = `mypolicy`
-- case
CREATE SEQUENCE `seq` INCREMENT BY -9223372036854775807
-- case
CREATE SEQUENCE `seq` INCREMENT BY -9223372036854775808
-- case
CREATE SEQUENCE `seq` INCREMENT BY -9223372036854775808
-- case
-- error: line 1 column 50 near ""the Signed Value should be at the range of [-9223372036854775808, 9223372036854775807]. 
-- case
SELECT `t`.`1a`.`1` FROM `t`
-- case
SELECT * FROM `1db`.`1table`
-- case
SELECT * FROM `t` WHERE `t`.`status`=1
-- case
SHOW PLACEMENT
-- case
SHOW PLACEMENT LIKE _UTF8MB4'POLICY foo%'
-- case
SHOW PLACEMENT WHERE `Target`=_UTF8MB4'TABLE test.t1'
-- case
SHOW PLACEMENT FOR DATABASE `db1`
-- case
SHOW PLACEMENT FOR DATABASE `db1`
-- case
SHOW PLACEMENT FOR TABLE `tb1`
-- case
SHOW PLACEMENT FOR TABLE `db1`.`tb1`
-- case
SHOW PLACEMENT FOR TABLE `tb1` PARTITION `p1`
-- case
SHOW PLACEMENT FOR TABLE `db1`.`tb1` PARTITION `p1`
-- case
-- error: line 1 column 18 near "" 
-- case
-- error: line 1 column 23 near "DATABASE db1" 
-- case
-- error: line 1 column 21 near "DB db1" 
-- case
-- error: line 1 column 37 near "TABLE tb1" 
-- case
-- error: line 1 column 28 near "PARTITION p1" 
-- case
-- error: line 1 column 21 near "DB LIKE '%'" 
-- case
-- error: line 1 column 21 near "DB db1 LIKE '%'" 
-- case
SHOW PLACEMENT LABELS
-- case
SHOW PLACEMENT LABELS LIKE _UTF8MB4'%zone%'
-- case
SHOW PLACEMENT LABELS WHERE `label`=_UTF8MB4'l123'
-- case
SHOW SESSION_STATES
-- case
SET SESSION_STATES 'x'
-- case
-- error: line 1 column 18 near "" 
-- case
-- error: line 1 column 20 near "1" 
-- case
-- error: line 1 column 22 near "now()" 
-- case
CALIBRATE RESOURCE
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2021-04-15 00:00:00'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' END_TIME _UTF8MB4'2023-04-01 16:00:00'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' DURATION '20m'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' END_TIME _UTF8MB4'2023-04-01 16:00:00' DURATION '20m'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' END_TIME _UTF8MB4'2023-04-01 16:00:00'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' DURATION '20m'
-- case
CALIBRATE RESOURCE DURATION '20m' START_TIME _UTF8MB4'2023-04-01 13:00:00'
-- case
CALIBRATE RESOURCE START_TIME _UTF8MB4'2023-04-01 13:00:00' END_TIME _UTF8MB4'2023-04-01 16:00:00' DURATION '20m'
-- case
CALIBRATE RESOURCE START_TIME CURRENT_TIMESTAMP() END_TIME CURRENT_TIMESTAMP()
-- case
CALIBRATE RESOURCE END_TIME NOW()
-- case
CALIBRATE RESOURCE START_TIME NOW()
-- case
CALIBRATE RESOURCE START_TIME NOW() END_TIME NOW()
-- case
CALIBRATE RESOURCE START_TIME DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 10 MINUTE) END_TIME NOW()
-- case
CALIBRATE RESOURCE START_TIME NOW()-1000 END_TIME CURRENT_TIMESTAMP()
-- case
CALIBRATE RESOURCE START_TIME DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 20 MINUTE) DURATION INTERVAL 15 MINUTE
-- case
CALIBRATE RESOURCE START_TIME DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 20 MINUTE) DURATION '15m'
-- case
CALIBRATE RESOURCE END_TIME NOW() START_TIME DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL 20 MINUTE)
-- case
-- error: line 1 column 27 near "" 
-- case
CALIBRATE RESOURCE WORKLOAD TPCC
-- case
CALIBRATE RESOURCE WORKLOAD OLTP_READ_WRITE
-- case
CALIBRATE RESOURCE WORKLOAD OLTP_READ_ONLY
-- case
CALIBRATE RESOURCE WORKLOAD OLTP_WRITE_ONLY
-- case
-- error: line 1 column 29 near "= oltp_read_write START_TIME '2023-04-01 13:00:00'" 
-- case
QUERY WATCH ADD SQL DIGEST `b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377`
-- case
QUERY WATCH ADD SQL DIGEST `b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377`
-- case
QUERY WATCH ADD SQL DIGEST _UTF8MB4'b13858789fce00208f9a262c99621b7045f4869807cd4e6568008ae7ca19a377'
-- case
QUERY WATCH ADD PLAN DIGEST `5e3ddd388f6012e328233dbcdda5d48f404e0536c6c54d9618233210f3d5762a`
-- case
QUERY WATCH ADD PLAN DIGEST @`digest1`
-- case
QUERY WATCH ADD SQL TEXT SIMILAR TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD SQL TEXT EXACT TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD SQL TEXT PLAN TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD RESOURCE GROUP `default` SQL TEXT SIMILAR TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD RESOURCE GROUP @`rg` SQL TEXT SIMILAR TO @`sql1`
-- case
QUERY WATCH ADD RESOURCE GROUP `rg1` SQL TEXT SIMILAR TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD SQL TEXT SIMILAR TO _UTF8MB4'select 1' RESOURCE GROUP `rg1`
-- case
QUERY WATCH ADD ACTION = KILL SQL TEXT SIMILAR TO _UTF8MB4'select 1'
-- case
QUERY WATCH ADD ACTION = COOLDOWN RESOURCE GROUP `rg1` SQL TEXT SIMILAR TO _UTF8MB4'select 1'
-- case
-- error: line 1 column 65 near "SQL TEXT SIMILAR to 'select 1'"Dupliated options specified 
-- case
-- error: line 1 column 27 near "SIMILAR to 'select 1'" 
-- case
-- error: line 1 column 43 near "'select 1'" 
-- case
QUERY WATCH REMOVE 1
-- case
QUERY WATCH REMOVE RESOURCE GROUP `rg1`
-- case
QUERY WATCH REMOVE RESOURCE GROUP @`rg`
-- case
-- error: line 1 column 18 near "" 
-- case
REPLACE /*+ SET_VAR(sql_mode = 'ALLOW_INVALID_DATES')*/ INTO `t` VALUES (_UTF8MB4'2004-04-31')
