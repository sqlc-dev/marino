WITH `cte` AS (SELECT 1,2) SELECT `col1`,`col2` FROM `cte`
-- case
WITH `cte` (col1, col2) AS (SELECT 1,2 UNION ALL SELECT 3,4) SELECT col1, col2 FROM cte;
-- case
WITH `cte` AS (SELECT 1,2), cte2 as (select 3) SELECT `col1`,`col2` FROM `cte`
-- case
WITH RECURSIVE cte (n) AS (  SELECT 1  UNION ALL  SELECT n + 1 FROM cte WHERE n < 5)SELECT * FROM cte;
-- case
with cte(a) as (select 1) update t, cte set t.a=1  where t.a=cte.a;
-- case
with cte(a) as (select 1) delete t from t, cte where t.a=cte.a;
-- case
WITH cte1 AS (SELECT 1) SELECT * FROM (WITH cte2 AS (SELECT 2) SELECT * FROM cte2 JOIN cte1) AS dt;
-- case
WITH cte AS (SELECT 1) SELECT /*+ MAX_EXECUTION_TIME(1000) */ * FROM cte;
-- case
with cte as (table t) table cte;
-- case
with cte as (select 1) select 1 union with cte as (select 1) select * from cte;
-- case
with cte as (select 1) (select 1);
-- case
with cte as (select 1) (select 1 union select 1)
-- case
select * from (with cte as (select 1) select 1 union select 2) qn
-- case
select * from t where 1 > (with cte as (select 2) select * from cte)
-- case
( with cte(n) as ( select 1 )  select n+1 from cte  union select n+2 from cte) union select 1
-- case
( with cte(n) as ( select 1 )  select n+1 from cte) union select 1
-- case
( with cte(n) as ( select 1 )  (select n+1 from cte)) union select 1
