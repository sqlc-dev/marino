SELECT 1>(SELECT 1)
-- case
SELECT 1>ANY (SELECT 1)
-- case
SELECT 1>ALL (SELECT 1)
-- case
SELECT 1>ANY (SELECT 1)
-- case
-- error: line 1 column 20 near "select 1" 
-- case
SELECT EXISTS (SELECT 1)
-- case
SELECT +EXISTS (SELECT 1)
-- case
SELECT -EXISTS (SELECT 1)
-- case
SELECT NOT EXISTS (SELECT 1)
-- case
-- error: line 1 column 12 near "NOT EXISTS (select 1)" 
-- case
-- error: line 1 column 12 near "NOT EXISTS (select 1)" 
-- case
SELECT * FROM `t` WHERE `t`.`a` IN (SELECT `a` FROM `t` LIMIT 1,10)
-- case
SELECT * FROM `t` WHERE `t`.`a` IN ((SELECT `a` FROM `t` LIMIT 1,10))
-- case
SELECT * FROM `t` WHERE `t`.`a` IN ((SELECT `a` FROM `t` LIMIT 1,10),1)
-- case
SELECT * FROM ((SELECT `a` FROM `t`) AS `t1` JOIN `t` AS `t2`) JOIN `t3`
-- case
SELECT `t1`.`a` AS `a` FROM (SELECT `a` FROM `t`) AS `t1`
-- case
SELECT COUNT(1) FROM (SELECT `a`,`b` FROM `x1` UNION ALL SELECT `a`,`b` FROM `x3` UNION ALL (SELECT `x1`.`a`,`x3`.`b` FROM (SELECT * FROM `x3` UNION ALL SELECT * FROM `x2`) AS `x3` LEFT JOIN `x1` ON `x3`.`a`=`x1`.`b`))
-- case
(SELECT 1 AS `a`,3 AS `b`) UNION (SELECT 2,1) ORDER BY (SELECT 2)
-- case
(SELECT * FROM `t1`) UNION (SELECT * FROM `t1`)
-- case
(SELECT * FROM `t1`) UNION (SELECT * FROM `t1`)
-- case
SELECT * FROM ((SELECT * FROM `t1`) UNION (SELECT * FROM `t1`) UNION (SELECT * FROM `t1`)) AS `a`
-- case
SELECT COUNT(1) FROM `plan_executions` WHERE (EXISTS (SELECT * FROM `triggers` WHERE `plan_executions`.`trigger_id`=`triggers`.`id` AND `triggers`.`type`=_UTF8MB4'CRON'))
-- case
SELECT EXISTS (SELECT 1)
-- case
SELECT * FROM ((SELECT 1 AS `a`,3 AS `b`) UNION (SELECT 2,1) ORDER BY (SELECT 2)) AS `t` ORDER BY `a`,`b`
-- case
SELECT (SELECT * FROM `t1` WHERE `a`!=`t`.`a` UNION ALL (SELECT * FROM `t2` WHERE `a`!=`t`.`a`) ORDER BY `a` LIMIT 1) FROM `t1` AS `t`
-- case
WITH `v0` AS (SELECT TRUE) (SELECT _UTF8MB4'abc' EXCEPT (SELECT TRUE))
