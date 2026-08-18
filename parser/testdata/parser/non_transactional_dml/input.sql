batch on c limit 10 delete from t where c = 10
-- case
batch on c limit 10 dry run delete from t where c = 10
-- case
batch on c limit 10 dry run query delete from t where c = 10
-- case
batch limit 10 delete from t where c = 10
-- case
batch limit 10 dry run delete from t where c = 10
-- case
batch limit 10 dry run query delete from t where c = 10
-- case
batch on c limit 10 update t set c = 10
-- case
batch on c limit 10 dry run update t set c = 10
-- case
batch on c limit 10 dry run query update t set c = 10
-- case
batch limit 10 update t set c = 10
-- case
batch limit 10 dry run update t set c = 10
-- case
batch limit 10 dry run query update t set c = 10
-- case
batch on c limit 10 insert into t1 select * from t2 where c = 10
-- case
batch on c limit 10 dry run insert into t1 select * from t2 where c = 10
-- case
batch on c limit 10 dry run query insert into t1 select * from t2 where c = 10
-- case
batch limit 10 insert into t1 select * from t2 where c = 10
-- case
batch limit 10 dry run insert into t1 select * from t2 where c = 10
-- case
batch limit 10 dry run query insert into t1 select * from t2 where c = 10
-- case
batch on c limit 10 insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
-- case
batch on c limit 10 dry run insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
-- case
batch on c limit 10 dry run query insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
-- case
batch limit 10 insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
-- case
batch limit 10 dry run insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
-- case
batch limit 10 dry run query insert into t1 select * from t2 where c = 10 on duplicate key update t1.val = t2.val
