CREATE TABLE t( c1 TIME(2), c2 DATETIME(2), c3 TIMESTAMP(2) );
-- case
select x'0a', X'11', 0x11
-- case
select x'13181C76734725455A'
-- case
select x'0xaa'
-- case
select 0X11
-- case
select 0x4920616D2061206C6F6E672068657820737472696E67
-- case
select 0b01, 0b0, b'11', B'11'
-- case
create table t (c1 enum('a', 'b'), c2 set('a', 'b'))
-- case
create table t (c1 enum('a  ', 'b	'), c2 set('a  ', 'b	'))
-- case
create table t (c1 enum('a', 'b') binary, c2 set('a', 'b') binary)
-- case
create table t (c1 enum(0x61, 'b'), c2 set(0x61, 'b'))
-- case
create table t (c1 enum(0b01100001, 'b'), c2 set(0b01100001, 'b'))
-- case
create table t (c1 enum)
-- case
create table t (c1 set)
-- case
create table t (c1 blob(1024), c2 text(1024))
-- case
create table t (y year(4), y1 year)
-- case
create table t (y year(4) unsigned zerofill zerofill, y1 year signed unsigned zerofill)
-- case
create table t (c1 national char(2), c2 national varchar(2))
-- case
create table t (a JSON);
