explain explore 'digestxxx'
-- case
explain explore select 1 from t
-- case
explain explore select 1 from t1, t2
-- case
explain explore select 1 from t where t1.a > (select max(a) from t2)
