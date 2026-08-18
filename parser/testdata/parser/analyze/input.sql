analyze table t1
-- case
analyze table t1.*
-- case
analyze table t,t1
-- case
analyze table t1 index
-- case
analyze table t1 index a
-- case
analyze table t1 index a,b
-- case
analyze table t with 4 buckets
-- case
analyze table t with 4 topn
-- case
analyze table t with 4 cmsketch width
-- case
analyze table t with 4 cmsketch depth
-- case
analyze table t with 4 samples
-- case
analyze table t with 4 buckets, 4 topn, 4 cmsketch width, 4 cmsketch depth, 4 samples
-- case
analyze table t index a with 4 buckets
-- case
analyze table t partition a
-- case
analyze table t partition a with 4 buckets
-- case
analyze table t partition a index b
-- case
analyze table t partition a index b with 4 buckets
-- case
analyze incremental table t index
-- case
analyze incremental table t index idx
-- case
analyze table t update histogram on b with 1024 buckets
-- case
analyze table t drop histogram on b
-- case
analyze table t update histogram on c1, c2;
-- case
analyze table t drop histogram on c1, c2;
-- case
analyze table t update histogram on t.c1, t.c2
-- case
analyze table t drop histogram on t.c1, t.c2
-- case
analyze table t1,t2 all columns
-- case
analyze table t partition a all columns
-- case
analyze table t1,t2 all columns with 4 topn
-- case
analyze table t partition a all columns with 1024 buckets
-- case
analyze table t1,t2 predicate columns
-- case
analyze table t partition a predicate columns
-- case
analyze table t1,t2 predicate columns with 4 topn
-- case
analyze table t partition a predicate columns with 1024 buckets
-- case
analyze table t columns c1,c2
-- case
analyze table t partition a columns c1,c2
-- case
analyze table t columns t.c1,t.c2
-- case
analyze table t partition a columns t.c1,t.c2
-- case
analyze table t columns c1,c2 with 4 topn
-- case
analyze table t partition a columns c1,c2 with 1024 buckets
-- case
analyze table t index a columns c
-- case
analyze table t index a all columns
-- case
analyze table t index a predicate columns
-- case
analyze table t with 10 samplerate
-- case
analyze table t with 0.1 samplerate
-- case
analyze table t with 0.05 ndvrate
-- case
analyze table t with 0.05 ndvrate 0.00001 samplerate
-- case
analyze no_write_to_binlog table t1
-- case
analyze local table t,t1
