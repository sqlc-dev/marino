create table cube (a int)
-- case
create table external (a int)
-- case
create table qualify (a int)
-- case
select cube from t
-- case
select external from t
-- case
select qualify from t
-- case
create table `cube` (a int)
-- case
create table `external` (a int)
-- case
create table `qualify` (a int)
-- case
select t.cube from t
-- case
select t.external from t
-- case
select t.qualify from t
-- case
create table manual (parallel int)
-- case
select manual, parallel from t
