select * from t use index (primary)
-- case
select * from t use index (`primary`)
-- case
select * from t use index ();
-- case
select * from t use index (idx);
-- case
select * from t use index (idx1, idx2);
-- case
select * from t ignore key (idx1)
-- case
select * from t force index for join (idx1)
-- case
select * from t use index for order by (idx1)
-- case
select * from t force index for group by (idx1)
-- case
select * from t use index for group by (idx1) use index for order by (idx2), t2
