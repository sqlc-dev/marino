create table t (a int) AFFINITY = 'table'
-- case
create table t (a int) affinity 'TABLE'
-- case
create table t (a int) affinity 'partition'
-- case
create table t (a int) AFFINITY = ''
-- case
create table t (a int) AFFINITY 'none'
-- case
create table t (a int) AFFINITY 'PARTITION' partition by hash ( a ) PARTITIONS 1
-- case
create table t (a int) /*T![affinity] AFFINITY = 'table'*/
-- case
create table t (a int) AFFINITY 'abcd'
-- case
alter table t AFFINITY = 'table'
-- case
alter table t affinity 'TABLE'
-- case
alter table t /*T![affinity] affinity 'table'*/
-- case
create table t (a int) AFFINITY 1
-- case
create table t (a int) AFFINITY = 1
-- case
create table t (a int) AFFINITY
