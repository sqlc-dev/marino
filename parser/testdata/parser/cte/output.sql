WITH `cte` AS (SELECT 1,2) SELECT `col1`,`col2` FROM `cte`
-- case
WITH `cte` (`col1`, `col2`) AS (SELECT 1,2 UNION ALL SELECT 3,4) SELECT `col1`,`col2` FROM `cte`
-- case
WITH `cte` AS (SELECT 1,2), `cte2` AS (SELECT 3) SELECT `col1`,`col2` FROM `cte`
-- case
WITH RECURSIVE `cte` (`n`) AS (SELECT 1 UNION ALL SELECT `n`+1 FROM `cte` WHERE `n`<5) SELECT * FROM `cte`
-- case
WITH `cte` (`a`) AS (SELECT 1) UPDATE (`t`) JOIN `cte` SET `t`.`a`=1 WHERE `t`.`a`=`cte`.`a`
-- case
WITH `cte` (`a`) AS (SELECT 1) DELETE `t` FROM (`t`) JOIN `cte` WHERE `t`.`a`=`cte`.`a`
-- case
WITH `cte1` AS (SELECT 1) SELECT * FROM (WITH `cte2` AS (SELECT 2) SELECT * FROM `cte2` JOIN `cte1`) AS `dt`
-- case
WITH `cte` AS (SELECT 1) SELECT /*+ MAX_EXECUTION_TIME(1000)*/ * FROM `cte`
-- case
WITH `cte` AS (TABLE `t`) TABLE `cte`
-- case
-- error: line 1 column 42 near "with cte as (select 1) select * from cte;" 
-- case
WITH `cte` AS (SELECT 1) (SELECT 1)
-- case
WITH `cte` AS (SELECT 1) (SELECT 1 UNION SELECT 1)
-- case
SELECT * FROM (WITH `cte` AS (SELECT 1) SELECT 1 UNION SELECT 2) AS `qn`
-- case
SELECT * FROM `t` WHERE 1>(WITH `cte` AS (SELECT 2) SELECT * FROM `cte`)
-- case
(WITH `cte` (`n`) AS (SELECT 1) SELECT `n`+1 FROM `cte` UNION SELECT `n`+2 FROM `cte`) UNION SELECT 1
-- case
(WITH `cte` (`n`) AS (SELECT 1) SELECT `n`+1 FROM `cte`) UNION SELECT 1
-- case
(WITH `cte` (`n`) AS (SELECT 1) (SELECT `n`+1 FROM `cte`)) UNION SELECT 1
