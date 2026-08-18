create view v as select * from t
-- case
create or replace view v as select * from t
-- case
create or replace algorithm = undefined view v as select * from t
-- case
create or replace algorithm = merge view v as select * from t
-- case
create or replace algorithm = temptable view v as select * from t
-- case
create or replace algorithm = merge definer = 'root' view v as select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security definer view v as select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v as select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t with local check option
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t with cascaded check option
-- case
create or replace algorithm = merge definer = current_user view v as select * from t
-- case
create view v as (select * from t)
-- case
create or replace view v as (select * from t)
-- case
create or replace algorithm = undefined view v as (select * from t)
-- case
create or replace algorithm = merge view v as (select * from t)
-- case
create or replace algorithm = temptable view v as (select * from t)
-- case
create or replace algorithm = merge definer = 'root' view v as (select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security definer view v as (select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v as (select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t) with local check option
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t) with cascaded check option
-- case
create or replace algorithm = merge definer = current_user view v as (select * from t)
-- case
create view v as select * from t union select * from t
-- case
create or replace view v as select * from t union select * from t
-- case
create or replace algorithm = undefined view v as select * from t union select * from t
-- case
create or replace algorithm = merge view v as select * from t union select * from t
-- case
create or replace algorithm = temptable view v as select * from t union select * from t
-- case
create or replace algorithm = merge definer = 'root' view v as select * from t union select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security definer view v as select * from t union select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v as select * from t union select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union select * from t with local check option
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union select * from t with cascaded check option
-- case
create or replace algorithm = merge definer = current_user view v as select * from t union select * from t
-- case
create view v as select * from t union all select * from t
-- case
create or replace view v as select * from t union all select * from t
-- case
create or replace algorithm = undefined view v as select * from t union all select * from t
-- case
create or replace algorithm = merge view v as select * from t union all select * from t
-- case
create or replace algorithm = temptable view v as select * from t union all select * from t
-- case
create or replace algorithm = merge definer = 'root' view v as select * from t union all select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security definer view v as select * from t union all select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v as select * from t union all select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union all select * from t
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union all select * from t with local check option
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as select * from t union all select * from t with cascaded check option
-- case
create or replace algorithm = merge definer = current_user view v as select * from t union all select * from t
-- case
create view v as (select * from t union all select * from t)
-- case
create or replace view v as (select * from t union all select * from t)
-- case
create or replace algorithm = undefined view v as (select * from t union all select * from t)
-- case
create or replace algorithm = merge view v as (select * from t union all select * from t)
-- case
create or replace algorithm = temptable view v as (select * from t union all select * from t)
-- case
create or replace algorithm = merge definer = 'root' view v as (select * from t union all select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security definer view v as (select * from t union all select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v as (select * from t union all select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t union all select * from t)
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t union all select * from t) with local check option
-- case
create or replace algorithm = merge definer = 'root' sql security invoker view v(a,b) as (select * from t union all select * from t) with cascaded check option
-- case
create or replace algorithm = merge definer = current_user view v as select * from t union all select * from t
