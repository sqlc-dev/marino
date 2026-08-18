SELECT CUME_DIST() OVER `w` FROM `t`
-- case
SELECT DENSE_RANK() OVER (`w`) FROM `t`
-- case
SELECT FIRST_VALUE(`val`) OVER `w` FROM `t`
-- case
SELECT FIRST_VALUE(`val`) OVER `w` FROM `t`
-- case
SELECT FIRST_VALUE(`val`) IGNORE NULLS OVER `w` FROM `t`
-- case
SELECT LAG(`val`) OVER (`w`) FROM `t`
-- case
SELECT LAG(`val`, 1) OVER (`w`) FROM `t`
-- case
SELECT LAG(`val`, 1, `def`) OVER (`w`) FROM `t`
-- case
SELECT LAST_VALUE(`val`) OVER (`w`) FROM `t`
-- case
SELECT LEAD(`val`) OVER `w` FROM `t`
-- case
SELECT LEAD(`val`, 1) OVER `w` FROM `t`
-- case
SELECT LEAD(`val`, 1, `def`) OVER `w` FROM `t`
-- case
SELECT NTH_VALUE(`val`, 233) OVER `w` FROM `t`
-- case
SELECT NTH_VALUE(`val`, 233) OVER `w` FROM `t`
-- case
SELECT NTH_VALUE(`val`, 233) FROM LAST OVER `w` FROM `t`
-- case
SELECT NTH_VALUE(`val`, 233) FROM LAST IGNORE NULLS OVER `w` FROM `t`
-- case
-- error: line 1 column 21 near ") OVER w FROM t;" 
-- case
SELECT NTILE(233) OVER (`w`) FROM `t`
-- case
SELECT PERCENT_RANK() OVER (`w`) FROM `t`
-- case
SELECT RANK() OVER (`w`) FROM `t`
-- case
SELECT ROW_NUMBER() OVER (`w`) FROM `t`
-- case
SELECT `n`,LAG(`n`, 1, 0) OVER (`w`),LEAD(`n`, 1, 0) OVER `w`,`n`+LAG(`n`, 1, 0) OVER (`w`) FROM `fib`
-- case
SELECT SUM(`profit`) OVER (PARTITION BY `country`) AS `country_profit` FROM `sales`
-- case
SELECT SUM(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT AVG(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT BIT_XOR(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT COUNT(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT COUNT(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT COUNT(1) OVER () AS `country_profit` FROM `sales`
-- case
SELECT MAX(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT MIN(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT SUM(`profit`) OVER () AS `country_profit` FROM `sales`
-- case
SELECT ROW_NUMBER() OVER (PARTITION BY `country`) AS `row_num1` FROM `sales`
-- case
SELECT ROW_NUMBER() OVER (PARTITION BY `country`, `d` ORDER BY `year`,`product`) AS `row_num2` FROM `sales`
-- case
SELECT SUM(`val`) OVER (PARTITION BY `subject` ORDER BY `time` ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM `t`
-- case
SELECT AVG(`val`) OVER (PARTITION BY `subject` ORDER BY `time` ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM `t`
-- case
SELECT AVG(`val`) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM `t`
-- case
SELECT AVG(`val`) OVER (ROWS BETWEEN 1 PRECEDING AND UNBOUNDED FOLLOWING) FROM `t`
-- case
SELECT AVG(`val`) OVER (RANGE BETWEEN INTERVAL 5 DAY PRECEDING AND INTERVAL _UTF8MB4'2:30' MINUTE_SECOND FOLLOWING) FROM `t`
-- case
SELECT AVG(`val`) OVER (RANGE BETWEEN CURRENT ROW AND CURRENT ROW) FROM `t`
-- case
SELECT AVG(`val`) OVER (RANGE BETWEEN CURRENT ROW AND CURRENT ROW) FROM `t`
-- case
SELECT RANK() OVER (`w`) FROM `t` WINDOW `w` AS (ORDER BY `val`)
-- case
SELECT RANK() OVER `w` FROM `t` WINDOW `w` AS ()
-- case
SELECT FIRST_VALUE(`year`) OVER (`w` ORDER BY `year`) AS `first` FROM `sales` WINDOW `w` AS (PARTITION BY `country`)
-- case
SELECT RANK() OVER (`w1`) FROM `t` WINDOW `w1` AS (`w2`),`w2` AS (),`w3` AS (`w1`)
-- case
SELECT RANK() OVER `w1` FROM `t` WINDOW `w1` AS (`w2`),`w2` AS (`w3`),`w3` AS (`w1`)
-- case
SELECT TIDB_PARSE_TSO(1)
-- case
SELECT TIDB_PARSE_TSO_LOGICAL(1)
-- case
SELECT TIDB_BOUNDED_STALENESS(_UTF8MB4'2015-09-21 00:07:01', NOW())
-- case
SELECT TIDB_BOUNDED_STALENESS(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())
-- case
SELECT TIDB_BOUNDED_STALENESS(_UTF8MB4'2015-09-21 00:07:01', _UTF8MB4'2021-04-27 11:26:13')
-- case
SELECT FROM_UNIXTIME(404411537129996288)
-- case
SELECT FROM_UNIXTIME(404411537129996288.22)
