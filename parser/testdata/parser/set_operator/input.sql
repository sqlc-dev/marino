select c1 from t1 union select c2 from t2
-- case
select c1 from t1 union (select c2 from t2)
-- case
select c1 from t1 union (select c2 from t2) order by c1
-- case
select c1 from t1 union select c2 from t2 order by c2
-- case
select c1 from t1 union (select c2 from t2) limit 1
-- case
select c1 from t1 union (select c2 from t2) limit 1, 1
-- case
select c1 from t1 union (select c2 from t2) order by c1 limit 1
-- case
(select c1 from t1) union distinct select c2 from t2
-- case
(select c1 from t1) union distinctrow select c2 from t2
-- case
(select c1 from t1) union all select c2 from t2
-- case
(select c1 from t1) union distinct all select c2 from t2
-- case
(select c1 from t1) union distinctrow all select c2 from t2
-- case
(select c1 from t1) union (select c2 from t2) order by c1 union select c3 from t3
-- case
(select c1 from t1) union (select c2 from t2) limit 1 union select c3 from t3
-- case
(select c1 from t1) union select c2 from t2 union (select c3 from t3) order by c1 limit 1
-- case
select (select 1 union select 1) as a
-- case
select * from (select 1 union select 2) as a
-- case
insert into t select c1 from t1 union select c2 from t2
-- case
insert into t (c) select c1 from t1 union select c2 from t2
-- case
select 2 as a from dual union select 1 as b from dual order by a
-- case
table t1 union table t2
-- case
table t1 union (table t2)
-- case
table t1 union select * from t2
-- case
select * from t1 union table t2
-- case
table t1 union (select c2 from t2) order by c1 limit 1
-- case
select c1 from t1 union (table t2) order by c1 limit 1
-- case
(select c1 from t1) union table t2 union (select c3 from t3) order by c1 limit 1
-- case
(table t1) union select c2 from t2 union (table t3) order by c1 limit 1
-- case
values row(1,-2,3), row(5,7,9) union values row(1,-2,3), row(5,7,9)
-- case
values row(1,-2,3), row(5,7,9) union (values row(1,-2,3), row(5,7,9))
-- case
values row(1,-2,3), row(5,7,9) union select * from t
-- case
values row(1,-2,3), row(5,7,9) union table t
-- case
select * from t union values row(1,-2,3), row(5,7,9)
-- case
table t union values row(1,-2,3), row(5,7,9)
-- case
select c1 from t1 except select c2 from t2
-- case
select c1 from t1 except (select c2 from t2)
-- case
select c1 from t1 except (select c2 from t2) order by c1
-- case
select c1 from t1 except select c2 from t2 order by c2
-- case
select c1 from t1 except (select c2 from t2) limit 1
-- case
select c1 from t1 except (select c2 from t2) limit 1, 1
-- case
select c1 from t1 except (select c2 from t2) order by c1 limit 1
-- case
(select c1 from t1) except (select c2 from t2) order by c1 except select c3 from t3
-- case
(select c1 from t1) except (select c2 from t2) limit 1 except select c3 from t3
-- case
(select c1 from t1) except select c2 from t2 except (select c3 from t3) order by c1 limit 1
-- case
select (select 1 except select 1) as a
-- case
select * from (select 1 except select 2) as a
-- case
insert into t select c1 from t1 except select c2 from t2
-- case
insert into t (c) select c1 from t1 except select c2 from t2
-- case
select 2 as a from dual except select 1 as b from dual order by a
-- case
table t1 except table t2
-- case
table t1 except (table t2)
-- case
table t1 except select * from t2
-- case
select * from t1 except table t2
-- case
table t1 except (select c2 from t2) order by c1 limit 1
-- case
select c1 from t1 except (table t2) order by c1 limit 1
-- case
(select c1 from t1) except table t2 except (select c3 from t3) order by c1 limit 1
-- case
(table t1) except select c2 from t2 except (table t3) order by c1 limit 1
-- case
values row(1,-2,3), row(5,7,9) except values row(1,-2,3), row(5,7,9)
-- case
values row(1,-2,3), row(5,7,9) except (values row(1,-2,3), row(5,7,9))
-- case
values row(1,-2,3), row(5,7,9) except select * from t
-- case
values row(1,-2,3), row(5,7,9) except table t
-- case
select * from t except values row(1,-2,3), row(5,7,9)
-- case
table t except values row(1,-2,3), row(5,7,9)
-- case
select c1 from t1 intersect select c2 from t2
-- case
select c1 from t1 intersect (select c2 from t2)
-- case
select c1 from t1 intersect (select c2 from t2) order by c1
-- case
select c1 from t1 intersect select c2 from t2 order by c2
-- case
select c1 from t1 intersect (select c2 from t2) limit 1
-- case
select c1 from t1 intersect (select c2 from t2) limit 1, 1
-- case
select c1 from t1 intersect (select c2 from t2) order by c1 limit 1
-- case
(select c1 from t1) intersect (select c2 from t2) order by c1 intersect select c3 from t3
-- case
(select c1 from t1) intersect (select c2 from t2) limit 1 intersect select c3 from t3
-- case
(select c1 from t1) intersect select c2 from t2 intersect (select c3 from t3) order by c1 limit 1
-- case
select (select 1 intersect select 1) as a
-- case
select * from (select 1 intersect select 2) as a
-- case
insert into t select c1 from t1 intersect select c2 from t2
-- case
insert into t (c) select c1 from t1 intersect select c2 from t2
-- case
select 2 as a from dual intersect select 1 as b from dual order by a
-- case
table t1 intersect table t2
-- case
table t1 intersect (table t2)
-- case
table t1 intersect select * from t2
-- case
select * from t1 intersect table t2
-- case
table t1 intersect (select c2 from t2) order by c1 limit 1
-- case
select c1 from t1 intersect (table t2) order by c1 limit 1
-- case
(select c1 from t1) intersect table t2 intersect (select c3 from t3) order by c1 limit 1
-- case
(table t1) intersect select c2 from t2 intersect (table t3) order by c1 limit 1
-- case
values row(1,-2,3), row(5,7,9) intersect values row(1,-2,3), row(5,7,9)
-- case
values row(1,-2,3), row(5,7,9) intersect (values row(1,-2,3), row(5,7,9))
-- case
values row(1,-2,3), row(5,7,9) intersect select * from t
-- case
values row(1,-2,3), row(5,7,9) intersect table t
-- case
select * from t intersect values row(1,-2,3), row(5,7,9)
-- case
table t intersect values row(1,-2,3), row(5,7,9)
-- case
(select c1 from t1) intersect select c2 from t2 union (select c3 from t3) order by c1 limit 1
-- case
(select c1 from t1) union all select c2 from t2 except (select c3 from t3) order by c1 limit 1
-- case
(select c1 from t1) except select c2 from t2 intersect (select c3 from t3) order by c1 limit 1
-- case
select 1 union distinct select 1 except select 1 intersect select 1
-- case
(select c1 from t1) intersect all (select c2 from t2 union (select c3 from t3)) order by c1 limit 1
-- case
(select c1 from t1) union all (select c2 from t2 except select c3 from t3) order by c1 limit 1
-- case
((select c1 from t1) except select c2 from t2) intersect all (select c3 from t3) order by c1 limit 1
-- case
select 1 union distinct (select 1 except all select 1 intersect select 1)
-- case
select * from a where PK = 0 union all (select * from b where PK = 0 union all (select * from b where PK != 0) order by pk limit 1)
-- case
select * from a where PK = 0 union all (select * from b where PK = 0 union all (select * from b where PK != 0) order by pk limit 1) order by pk limit 2
-- case
(select * from b where pk= 0 union all (select * from b where pk !=0) order by pk limit 1) order by pk limit 2
-- case
(select * from b where pk= 0 union all (select * from b where pk !=0) order by pk limit 1) order by pk
