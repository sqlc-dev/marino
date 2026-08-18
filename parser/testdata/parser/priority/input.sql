select high_priority * from t
-- case
select low_priority * from t
-- case
select delayed * from t
-- case
insert high_priority into t values (1)
-- case
insert LOW_PRIORITY into t values (1)
-- case
insert delayed into t values (1)
-- case
update low_priority t set a = 2
-- case
update high_priority t set a = 2
-- case
update delayed t set a = 2
-- case
delete low_priority from t where a = 2
-- case
delete high_priority from t where a = 2
-- case
delete delayed from t where a = 2
-- case
replace high_priority into t values (1)
-- case
replace LOW_PRIORITY into t values (1)
-- case
replace delayed into t values (1)
