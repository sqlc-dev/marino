EXPLAIN FORMAT = 'TREE' SELECT 1
-- case
EXPLAIN ANALYZE FORMAT = 'TREE' SELECT 1
-- case
EXPLAIN FORMAT = 'JSON' INTO @`plan` SELECT 1
-- case
EXPLAIN FORMAT = 'JSON' INTO @`plan` SELECT `t`.`id`,COUNT(1) FROM `information_schema`.`tables` AS `t` GROUP BY `t`.`id`
-- case
DESC `t` 'a%'
