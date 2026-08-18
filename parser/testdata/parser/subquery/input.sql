SELECT 1 > (select 1)
-- case
SELECT 1 > ANY (select 1)
-- case
SELECT 1 > ALL (select 1)
-- case
SELECT 1 > SOME (select 1)
-- case
SELECT EXISTS select 1
-- case
SELECT EXISTS (select 1)
-- case
SELECT + EXISTS (select 1)
-- case
SELECT - EXISTS (select 1)
-- case
SELECT NOT EXISTS (select 1)
-- case
SELECT + NOT EXISTS (select 1)
-- case
SELECT - NOT EXISTS (select 1)
-- case
SELECT * FROM t where t.a in (select a from t limit 1, 10)
-- case
SELECT * FROM t where t.a in ((select a from t limit 1, 10))
-- case
SELECT * FROM t where t.a in ((select a from t limit 1, 10), 1)
-- case
select * from ((select a from t) t1 join t t2) join t3
-- case
SELECT t1.a AS a FROM ((SELECT a FROM t) AS t1)
-- case
select count(*) from (select a, b from x1 union all select a, b from x3 union all (select x1.a, x3.b from (select * from x3 union all select * from x2) x3 left join x1 on x3.a = x1.b))
-- case
(SELECT 1 a,3 b) UNION (SELECT 2,1) ORDER BY (SELECT 2)
-- case
((select * from t1)) union (select * from t1)
-- case
(((select * from t1))) union (select * from t1)
-- case
select * from (((select * from t1)) union (select * from t1) union (select * from t1)) a
-- case
SELECT COUNT(*) FROM plan_executions WHERE (EXISTS((SELECT * FROM triggers WHERE plan_executions.trigger_id=triggers.id AND triggers.type='CRON')))
-- case
select exists((select 1));
-- case
select * from ((SELECT 1 a,3 b) UNION (SELECT 2,1) ORDER BY (SELECT 2)) t order by a,b
-- case
select (select * from t1 where a != t.a union all (select * from t2 where a != t.a) order by a limit 1) from t1 t
-- case
(WITH v0 AS (SELECT TRUE) (SELECT 'abc' EXCEPT (SELECT TRUE)))
