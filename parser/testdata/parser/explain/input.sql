explain select c1 from t1
-- case
explain delete t1, t2 from t1 inner join t2 inner join t3 where t1.id=t2.id and t2.id=t3.id;
-- case
explain insert into t values (1), (2), (3)
-- case
explain replace into foo values (1 || 2)
-- case
explain update t set id = id + 1 order by id desc;
-- case
explain select c1 from t1 union (select c2 from t2) limit 1, 1
-- case
explain format = "row" select c1 from t1 union (select c2 from t2) limit 1, 1
-- case
explain format = 'brief' select * from t
-- case
DESC SCHE.TABL
-- case
DESC SCHE.TABL COLUM
-- case
DESCRIBE SCHE.TABL COLUM
-- case
EXPLAIN ANALYZE SELECT 1
-- case
EXPLAIN ANALYZE format=VERBOSE SELECT 1
-- case
EXPLAIN ANALYZE format=TRUE_CARD_COST SELECT 1
-- case
EXPLAIN ANALYZE format='VERBOSE' SELECT 1
-- case
EXPLAIN ANALYZE format='TRUE_CARD_COST' SELECT 1
-- case
EXPLAIN FORMAT = 'dot' SELECT 1
-- case
EXPLAIN FORMAT = DOT SELECT 1
-- case
EXPLAIN FORMAT = 'row' SELECT 1
-- case
EXPLAIN FORMAT = 'ROW' SELECT 1
-- case
EXPLAIN FORMAT = 'BRIEF' SELECT 1
-- case
EXPLAIN FORMAT = BRIEF SELECT 1
-- case
EXPLAIN FORMAT = 'verbose' SELECT 1
-- case
EXPLAIN FORMAT = 'VERBOSE' SELECT 1
-- case
EXPLAIN FORMAT = VERBOSE SELECT 1
-- case
EXPLAIN SELECT 1
-- case
EXPLAIN FOR CONNECTION 1
-- case
EXPLAIN FOR connection 42
-- case
EXPLAIN FORMAT = 'dot' FOR CONNECTION 1
-- case
EXPLAIN FORMAT = DOT FOR CONNECTION 1
-- case
EXPLAIN FORMAT = 'row' FOR connection 1
-- case
EXPLAIN FORMAT = ROW FOR connection 1
-- case
EXPLAIN FORMAT = TRADITIONAL FOR CONNECTION 1
-- case
EXPLAIN FORMAT = TRADITIONAL SELECT 1
-- case
EXPLAIN FORMAT = BRIEF SELECT 1
-- case
EXPLAIN FORMAT = 'brief' SELECT 1
-- case
EXPLAIN FORMAT = DOT SELECT 1
-- case
EXPLAIN FORMAT = 'dot' SELECT 1
-- case
EXPLAIN FORMAT = VERBOSE SELECT 1
-- case
EXPLAIN FORMAT = 'verbose' SELECT 1
-- case
EXPLAIN FORMAT = JSON FOR CONNECTION 1
-- case
EXPLAIN FORMAT = JSON SELECT 1
-- case
EXPLAIN FORMAT = 'hint' SELECT 1
-- case
EXPLAIN ANALYZE FORMAT = 'verbose' SELECT 1
-- case
EXPLAIN ANALYZE FORMAT = 'binary' SELECT 1
-- case
EXPLAIN ALTER TABLE t1 ADD INDEX (a)
-- case
EXPLAIN ALTER TABLE t1 ADD a varchar(255)
-- case
EXPLAIN FORMAT = TIDB_JSON FOR CONNECTION 1
-- case
EXPLAIN FORMAT = tidb_json SELECT 1
-- case
EXPLAIN ANALYZE FORMAT = tidb_json SELECT 1
-- case
EXPLAIN 'sqldigest'
-- case
EXPLAIN ANALYZE 'sqldigest'
-- case
EXPLAIN format='json' 'sqldigest'
-- case
EXPLAIN ANALYZE format='json' 'sqldigest'
