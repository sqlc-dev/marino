trace begin
-- case
trace commit
-- case
trace rollback
-- case
trace set a = 1
-- case
trace select c1 from t1
-- case
trace delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;
-- case
trace insert into t values (1), (2), (3)
-- case
trace replace into foo values (1 || 2)
-- case
trace update t set id = id + 1 order by id desc;
-- case
trace select c1 from t1 union (select c2 from t2) limit 1, 1
-- case
trace format = 'row' select c1 from t1 union (select c2 from t2) limit 1, 1
-- case
trace format = 'json' update t set id = id + 1 order by id desc;
-- case
trace plan select c1 from t1
-- case
trace plan target = 'estimation' select c1 from t1
-- case
trace plan target = 'arandomstring' select c1 from t1
