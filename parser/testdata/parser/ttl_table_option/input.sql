create table t (created_at datetime) TTL = created_at + INTERVAL 3.1415 YEAR
-- case
create table t (created_at datetime) TTL = created_at + INTERVAL '1 1:1:1' DAY_SECOND
-- case
create table t (created_at datetime) TTL = created_at + INTERVAL 1 YEAR
-- case
create table t (created_at datetime) TTL = created_at + INTERVAL 1 YEAR TTL_ENABLE = 'OFF'
-- case
create table t (created_at datetime) TTL created_at + INTERVAL 1 YEAR TTL_ENABLE 'OFF'
-- case
create table t (created_at datetime) TTL created_at + INTERVAL 1 YEAR TTL_ENABLE 'OFF' TTL_JOB_INTERVAL='8h'
-- case
create table t (created_at datetime) /*T![ttl] ttl=created_at + INTERVAL 1 YEAR ttl_enable='ON'*/
-- case
alter table t TTL = created_at + INTERVAL 1 MONTH
-- case
alter table t TTL_ENABLE = 'ON'
-- case
alter table t TTL_ENABLE = 'OFF'
-- case
alter table t TTL = created_at + INTERVAL 1 MONTH TTL_ENABLE 'OFF'
-- case
alter table t TTL = created_at + INTERVAL 1 MONTH TTL_ENABLE 'OFF' TTL_JOB_INTERVAL '1h'
-- case
alter table t /*T![ttl] ttl=created_at + INTERVAL 1 YEAR ttl_enable='ON'*/
-- case
alter table t /*T![ttl] ttl=created_at + INTERVAL 1 YEAR ttl_enable='ON' TTL_JOB_INTERVAL='8h'*/
-- case
alter table t /*T![ttl] ttl=created_at + INTERVAL 1 YEAR ttl_enable='ON' TTL_JOB_INTERVAL='8.645124531235h'*/
-- case
alter table t remove ttl
-- case
create table t (created_at datetime) TTL_ENABLE = 'test_case'
-- case
create table t (created_at datetime) /*T![ttl] TTL_ENABLE = 'test_case' */
-- case
alter table t /*T![ttl] TTL_ENABLE = 'test_case' */
-- case
create table t (created_at datetime) TTL_JOB_INTERVAL = '@monthly'
-- case
create table t (created_at datetime) TTL_JOB_INTERVAL = '10hourxx'
-- case
create table t (created_at datetime) TTL_JOB_INTERVAL = '10.10.255h'
