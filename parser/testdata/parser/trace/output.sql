TRACE START TRANSACTION
-- case
TRACE COMMIT
-- case
TRACE ROLLBACK
-- case
TRACE SET @@SESSION.`a`=1
-- case
TRACE SELECT `c1` FROM `t1`
-- case
TRACE DELETE `t1`,`t2` FROM (`t1` JOIN `t2`) JOIN `t3` WHERE `t1`.`id`=`t2`.`id` AND `t2`.`id`=`t3`.`id`
-- case
TRACE INSERT INTO `t` VALUES (1),(2),(3)
-- case
TRACE REPLACE INTO `foo` VALUES (1 OR 2)
-- case
TRACE UPDATE `t` SET `id`=`id`+1 ORDER BY `id` DESC
-- case
TRACE SELECT `c1` FROM `t1` UNION (SELECT `c2` FROM `t2`) LIMIT 1,1
-- case
TRACE SELECT `c1` FROM `t1` UNION (SELECT `c2` FROM `t2`) LIMIT 1,1
-- case
TRACE FORMAT = 'json' UPDATE `t` SET `id`=`id`+1 ORDER BY `id` DESC
-- case
TRACE PLAN SELECT `c1` FROM `t1`
-- case
TRACE PLAN TARGET = 'estimation' SELECT `c1` FROM `t1`
-- case
TRACE PLAN TARGET = 'arandomstring' SELECT `c1` FROM `t1`
