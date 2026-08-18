create statistics stats1 (cardinality) on t(a,b,c)
-- case
create statistics stats2 (dependency) on t(a,b)
-- case
create statistics stats3 (correlation) on t(a,b)
-- case
create statistics stats3 on t(a,b)
-- case
create statistics if not exists stats1 (cardinality) on t(a,b,c)
-- case
create statistics if not exists stats2 (dependency) on t(a,b)
-- case
create statistics if not exists stats3 (correlation) on t(a,b)
-- case
create statistics if not exists stats3 on t(a,b)
-- case
create statistics stats1(cardinality) on t(a,b,c)
-- case
drop statistics stats1
