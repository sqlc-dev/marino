create global binding for select * from t using select * from t use index(a)
-- case
create session binding for select * from t using select * from t use index(a)
-- case
drop global binding for select * from t
-- case
drop session binding for select * from t
-- case
drop global binding for select * from t using select * from t use index(a)
-- case
drop session binding for select * from t using select * from t use index(a)
-- case
show global bindings
-- case
show session bindings
-- case
set binding enabled for select * from t
-- case
set binding enabled for select * from t using select * from t use index(a)
-- case
set binding disabled for select * from t
-- case
set binding disabled for select * from t using select * from t use index(a)
-- case
create global binding for select * from t union all select * from t using select * from t use index(a) union all select * from t use index(a)
-- case
create session binding for select * from t union all select * from t using select * from t use index(a) union all select * from t use index(a)
-- case
drop global binding for select * from t union all select * from t using select * from t use index(a) union all select * from t use index(a)
-- case
drop session binding for select * from t union all select * from t using select * from t use index(a) union all select * from t use index(a)
-- case
drop global binding for select * from t union all select * from t
-- case
create session binding for select 1 union select 2 intersect select 3 using select 1 union select 2 intersect select 3
-- case
drop session binding for select 1 union select 2 intersect select 3 using select 1 union select 2 intersect select 3
-- case
drop session binding for select 1 union select 2 intersect select 3
-- case
create global binding using select * from *.t1
-- case
create global binding using select * from *.t1 where t1.a > (select max(a) from t2)
-- case
create session binding using select * from *.t1
-- case
create binding using select * from *.t1
-- case
CREATE GLOBAL BINDING FOR UPDATE `t` SET `a`=1 WHERE `b`=1 USING UPDATE /*+ USE_INDEX(`t` `b`)*/ `t` SET `a`=1 WHERE `b`=1
-- case
CREATE SESSION BINDING FOR UPDATE `t` SET `a`=1 WHERE `b`=1 USING UPDATE /*+ USE_INDEX(`t` `b`)*/ `t` SET `a`=1 WHERE `b`=1
-- case
drop global binding for update t set a = 1 where b = 1
-- case
drop session binding for update t set a = 1 where b = 1
-- case
DROP GLOBAL BINDING FOR UPDATE `t` SET `a`=1 WHERE `b`=1 USING UPDATE /*+ USE_INDEX(`t` `b`)*/ `t` SET `a`=1 WHERE `b`=1
-- case
DROP SESSION BINDING FOR UPDATE `t` SET `a`=1 WHERE `b`=1 USING UPDATE /*+ USE_INDEX(`t` `b`)*/ `t` SET `a`=1 WHERE `b`=1
-- case
CREATE GLOBAL BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b` USING UPDATE /*+ INL_JOIN(`t1`)*/ `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
CREATE SESSION BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b` USING UPDATE /*+ INL_JOIN(`t1`)*/ `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
DROP GLOBAL BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
DROP SESSION BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
DROP GLOBAL BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b` USING UPDATE /*+ INL_JOIN(`t1`)*/ `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
DROP SESSION BINDING FOR UPDATE `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b` USING UPDATE /*+ INL_JOIN(`t1`)*/ `t1` JOIN `t2` SET `t1`.`a`=1 WHERE `t1`.`b`=`t2`.`b`
-- case
CREATE GLOBAL BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHERE `a`=1
-- case
CREATE SESSION BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHERE `a`=1
-- case
drop global binding for delete from t where a = 1
-- case
drop session binding for delete from t where a = 1
-- case
DROP GLOBAL BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHERE `a`=1
-- case
DROP SESSION BINDING FOR DELETE FROM `t` WHERE `a`=1 USING DELETE /*+ USE_INDEX(`t` `a`)*/ FROM `t` WHERE `a`=1
-- case
CREATE GLOBAL BINDING FOR DELETE `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1 USING DELETE /*+ HASH_JOIN(`t1`, `t2`)*/ `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1
-- case
CREATE SESSION BINDING FOR DELETE `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1 USING DELETE /*+ HASH_JOIN(`t1`, `t2`)*/ `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1
-- case
drop global binding for delete t1, t2 from t1 inner join t2 on t1.b = t2.b where t1.a = 1
-- case
drop session binding for delete t1, t2 from t1 inner join t2 on t1.b = t2.b where t1.a = 1
-- case
DROP GLOBAL BINDING FOR DELETE `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1 USING DELETE /*+ HASH_JOIN(`t1`, `t2`)*/ `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1
-- case
DROP SESSION BINDING FOR DELETE `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1 USING DELETE /*+ HASH_JOIN(`t1`, `t2`)*/ `t1`,`t2` FROM `t1` JOIN `t2` ON `t1`.`b`=`t2`.`b` WHERE `t1`.`a`=1
-- case
CREATE GLOBAL BINDING FOR INSERT INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING INSERT INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
CREATE SESSION BINDING FOR INSERT INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING INSERT INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
drop global binding for insert into t1 select * from t2 where t1.a=1
-- case
drop session binding for insert into t1 select * from t2 where t1.a=1
-- case
DROP GLOBAL BINDING FOR INSERT INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING INSERT INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
DROP SESSION BINDING FOR INSERT INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING INSERT INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
CREATE GLOBAL BINDING FOR REPLACE INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING REPLACE INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
CREATE SESSION BINDING FOR REPLACE INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING REPLACE INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
drop global binding for replace into t1 select * from t2 where t1.a=1
-- case
drop session binding for replace into t1 select * from t2 where t1.a=1
-- case
DROP GLOBAL BINDING FOR REPLACE INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING REPLACE INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
DROP SESSION BINDING FOR REPLACE INTO `t1` SELECT * FROM `t2` WHERE `t2`.`a`=1 USING REPLACE INTO `t1` SELECT /*+ USE_INDEX(`t2` `a`)*/ * FROM `t2` WHERE `t2`.`a`=1
-- case
DROP SESSION BINDING FOR SQL DIGEST 'a'
-- case
drop global binding for sql digest 's'
-- case
drop global binding for sql digest @a, @b, 'test1,test2', @c, 'test333'
-- case
create session binding from history using plan digest 'sss'
-- case
create session binding from history using plan digest @a, @b, 'test1,test2', @c, 'test333'
-- case
CREATE GLOBAL BINDING FROM HISTORY USING PLAN DIGEST 'sss'
-- case
set binding enabled for sql digest '1'
-- case
set binding disabled for sql digest '1'
-- case
explain explore 'select a from t'
-- case
explain explore '23adc8e6f62'
