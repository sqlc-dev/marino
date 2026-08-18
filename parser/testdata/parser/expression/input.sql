SELECT ++1
-- case
SELECT -*1
-- case
SELECT -+1
-- case
SELECT -1
-- case
SELECT --1
-- case
select '''a''', """a"""
-- case
select ''a''
-- case
select ""a""
-- case
select '''a''';
-- case
select '\'a\'';
-- case
select "\"a\"";
-- case
select """a""";
-- case
select _utf8"string";
-- case
select _binary"string";
-- case
select N'string'
-- case
select n'string'
-- case
select _utf8 0xD0B1;
-- case
select _utf8 X'D0B1';
-- case
select _utf8 0b1101000010110001;
-- case
select _utf8 B'1101000010110001';
-- case
select 1 <=> 0, 1 <=> null, 1 = null
-- case
select date'1989-09-10'
-- case
select date 19890910
-- case
select time '00:00:00.111'
-- case
select time 19890910
-- case
select timestamp '1989-09-10 11:11:11'
-- case
select timestamp 19890910
-- case
select {ts '1989-09-10 11:11:11'}
-- case
select {d '1989-09-10'}
-- case
select {t '00:00:00.111'}
-- case
select * from t where a > {ts '1989-09-10 11:11:11'}
-- case
select * from t where a > {ts {abc '1989-09-10 11:11:11'}}
-- case
select {ts123 '1989-09-10 11:11:11'}
-- case
select {ts123 123}
-- case
select {ts123 1 xor 1}
-- case
select * from t where a > {ts123 '1989-09-10 11:11:11'}
-- case
select .t.a from t
