select * from t group by a, b rollup
-- case
select * from t group by a, b with rollup
-- case
select * from t group by (a, b) with rollup
-- case
select * from t group by (a+b) with rollup
