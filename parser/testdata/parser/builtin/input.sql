SELECT POW(1, 2)
-- case
SELECT POW(1, 2, 1)
-- case
SELECT POW(1, 0.5)
-- case
SELECT POW(1, -1)
-- case
SELECT POW(-1, 1)
-- case
SELECT RAND();
-- case
SELECT RAND(1);
-- case
SELECT MOD(10, 2);
-- case
SELECT ROUND(-1.23);
-- case
SELECT ROUND(1.23, 1);
-- case
SELECT ROUND(1.23, 1, 1);
-- case
SELECT CEIL(-1.23);
-- case
SELECT CEILING(1.23);
-- case
SELECT FLOOR(-1.23);
-- case
SELECT LN(1);
-- case
SELECT LN(1, 2);
-- case
SELECT LOG(-2);
-- case
SELECT LOG(2, 65536);
-- case
SELECT LOG(2, 65536, 1);
-- case
SELECT LOG2(2);
-- case
SELECT LOG2(2, 2);
-- case
SELECT LOG10(10);
-- case
SELECT LOG10(10, 1);
-- case
SELECT ABS(10, 1);
-- case
SELECT ABS(10);
-- case
SELECT ABS();
-- case
SELECT CONV(10+'10'+'10'+X'0a',10,10);
-- case
SELECT CONV();
-- case
SELECT CRC32('MySQL');
-- case
SELECT CRC32();
-- case
SELECT SIGN();
-- case
SELECT SIGN(0);
-- case
SELECT SQRT(0);
-- case
SELECT SQRT();
-- case
SELECT ACOS();
-- case
SELECT ACOS(1);
-- case
SELECT ACOS(1, 2);
-- case
SELECT ASIN();
-- case
SELECT ASIN(1);
-- case
SELECT ASIN(1, 2);
-- case
SELECT ATAN(0), ATAN(1), ATAN(1, 2);
-- case
SELECT ATAN2(), ATAN2(1,2);
-- case
SELECT COS(0);
-- case
SELECT COS(1);
-- case
SELECT COS(1, 2);
-- case
SELECT COT();
-- case
SELECT COT(1);
-- case
SELECT COT(1, 2);
-- case
SELECT DEGREES();
-- case
SELECT DEGREES(0);
-- case
SELECT EXP();
-- case
SELECT EXP(1);
-- case
SELECT PI();
-- case
SELECT PI(1);
-- case
SELECT RADIANS();
-- case
SELECT RADIANS(1);
-- case
SELECT SIN();
-- case
SELECT SIN(1);
-- case
SELECT TAN(1);
-- case
SELECT TAN();
-- case
SELECT TRUNCATE(1.223,1);
-- case
SELECT TRUNCATE();
-- case
SELECT SUBSTR('Quadratically',5);
-- case
SELECT SUBSTR('Quadratically',5, 3);
-- case
SELECT SUBSTR('Quadratically' FROM 5);
-- case
SELECT SUBSTR('Quadratically' FROM 5 FOR 3);
-- case
SELECT SUBSTRING('Quadratically',5);
-- case
SELECT SUBSTRING('Quadratically',5, 3);
-- case
SELECT SUBSTRING('Quadratically' FROM 5);
-- case
SELECT SUBSTRING('Quadratically' FROM 5 FOR 3);
-- case
SELECT CONVERT('111', SIGNED);
-- case
SELECT LEAST(), LEAST(1, 2, 3);
-- case
SELECT INTERVAL(1, 0, 1, 2)
-- case
SELECT (INTERVAL(1, 0, 1, 2)+5)*7+INTERVAL(1, 0, 1, 2)/2
-- case
SELECT INTERVAL(0, (1*5)/2)+INTERVAL(5, 4, 3)
-- case
SELECT DATE_ADD('2008-01-02', INTERVAL INTERVAL(1, 0, 1) DAY);
-- case
SELECT DATABASE();
-- case
SELECT SCHEMA();
-- case
SELECT USER();
-- case
SELECT USER(1);
-- case
SELECT CURRENT_USER();
-- case
SELECT CURRENT_ROLE();
-- case
SELECT CURRENT_USER;
-- case
SELECT CONNECTION_ID();
-- case
SELECT VERSION();
-- case
SELECT CURRENT_RESOURCE_GROUP();
-- case
SELECT BENCHMARK(1000000, AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3')));
-- case
SELECT BENCHMARK(AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3')));
-- case
SELECT CHARSET('abc');
-- case
SELECT COERCIBILITY('abc');
-- case
SELECT COERCIBILITY('abc', 'a');
-- case
SELECT COLLATION('abc');
-- case
SELECT ROW_COUNT();
-- case
SELECT SESSION_USER();
-- case
SELECT SYSTEM_USER();
-- case
SELECT FORMAT_BYTES(512);
-- case
SELECT FORMAT_NANO_TIME(3501);
-- case
SELECT SUBSTRING_INDEX('www.mysql.com', '.', 2);
-- case
SELECT SUBSTRING_INDEX('www.mysql.com', '.', -2);
-- case
SELECT ASCII(), ASCII(""), ASCII("A"), ASCII(1);
-- case
SELECT LOWER("A"), UPPER("a")
-- case
SELECT LCASE("A"), UCASE("a")
-- case
SELECT REPLACE('www.mysql.com', 'w', 'Ww')
-- case
SELECT LOCATE('bar', 'foobarbar');
-- case
SELECT LOCATE('bar', 'foobarbar', 5);
-- case
SELECT tidb_version();
-- case
SELECT tidb_is_ddl_owner();
-- case
SELECT tidb_decode_plan();
-- case
SELECT tidb_decode_key('abc');
-- case
SELECT tidb_decode_base64_key('abc');
-- case
SELECT tidb_decode_sql_digests('[]');
-- case
CREATE TABLE t( c1 TIME(2), c2 DATETIME(2), c3 TIMESTAMP(2) );
-- case
select row(1)
-- case
select row(1, 1,)
-- case
select (1, 1,)
-- case
select row(1, 1) > row(1, 1), row(1, 1, 1) > row(1, 1, 1)
-- case
Select (1, 1) > (1, 1)
-- case
create table t (`row` int)
-- case
create table t (row int)
-- case
SELECT *, CAST(data AS CHAR CHARACTER SET utf8) FROM t;
-- case
SELECT CAST(data AS CHARACTER);
-- case
SELECT CAST(data AS CHARACTER(10) CHARACTER SET utf8);
-- case
SELECT CAST(data AS BINARY)
-- case
SELECT *, CAST(data AS JSON) FROM t;
-- case
SELECT *, JSON_SUM_CRC32(data AS UNSIGNED ARRAY) FROM t;
-- case
SELECT *, JSON_SUM_CRC32(data AS DOUBLE ARRAY) FROM t;
-- case
SELECT *, JSON_SUM_CRC32(data AS DOUBLE) FROM t;
-- case
SELECT *, JSON_SUM_CRC32(data) FROM t;
-- case
select cast(1 as signed int);
-- case
select cast(1 as double);
-- case
select cast(1 as float);
-- case
select cast(1 as float(0));
-- case
select cast(1 as float(24));
-- case
select cast(1 as float(25));
-- case
select cast(1 as float(53));
-- case
select cast(1 as float(54));
-- case
select cast(1 as real);
-- case
select cast('2000' as year);
-- case
select cast(time '2000' as year);
-- case
select cast(b as signed array);
-- case
select cast(b as char(10) array);
-- case
SELECT last_insert_id();
-- case
SELECT last_insert_id(1);
-- case
SELECT binary 'a';
-- case
SELECT BIT_COUNT(1);
-- case
select current_timestamp
-- case
select current_timestamp()
-- case
select current_timestamp(6)
-- case
select current_timestamp(null)
-- case
select current_timestamp(-1)
-- case
select current_timestamp(1.0)
-- case
select current_timestamp('2')
-- case
select now()
-- case
select now(6)
-- case
select sysdate(), sysdate(6)
-- case
SELECT time('01:02:03');
-- case
SELECT time('01:02:03.1')
-- case
SELECT time('20.1')
-- case
SELECT TIMEDIFF('2000:01:01 00:00:00', '2000:01:01 00:00:00.000001');
-- case
SELECT TIMESTAMPDIFF(MONTH,'2003-02-01','2003-05-01');
-- case
SELECT TIMESTAMPDIFF(YEAR,'2002-05-01','2001-01-01');
-- case
SELECT TIMESTAMPDIFF(MINUTE,'2003-02-01','2003-05-01 12:05:55');
-- case
select current_time
-- case
select current_time()
-- case
select current_time(6)
-- case
select current_time(-1)
-- case
select current_time(1.0)
-- case
select current_time('1')
-- case
select current_time(null)
-- case
select curtime()
-- case
select curtime(6)
-- case
select curtime(-1)
-- case
select curtime(1.0)
-- case
select curtime('1')
-- case
select curtime(null)
-- case
select utc_timestamp
-- case
select utc_timestamp()
-- case
select utc_timestamp(6)
-- case
select utc_timestamp(-1)
-- case
select utc_timestamp(1.0)
-- case
select utc_timestamp('1')
-- case
select utc_timestamp(null)
-- case
select utc_time
-- case
select utc_time()
-- case
select utc_time(6)
-- case
select utc_time(-1)
-- case
select utc_time(1.0)
-- case
select utc_time('1')
-- case
select utc_time(null)
-- case
SELECT MICROSECOND('2009-12-31 23:59:59.000010');
-- case
SELECT SECOND('10:05:03');
-- case
SELECT MINUTE('2008-02-03 10:05:03');
-- case
SELECT HOUR(), HOUR('10:05:03');
-- case
SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE()
-- case
SELECT CURRENT_DATE, CURRENT_DATE(), CURDATE(1)
-- case
SELECT DATEDIFF('2003-12-31', '2003-12-30');
-- case
SELECT DATE('2003-12-31 01:02:03');
-- case
SELECT DATE();
-- case
SELECT DATE('2003-12-31 01:02:03', '');
-- case
SELECT DATE_FORMAT('2003-12-31 01:02:03', '%W %M %Y');
-- case
SELECT DAY('2007-02-03');
-- case
SELECT DAYOFMONTH('2007-02-03');
-- case
SELECT DAYOFWEEK('2007-02-03');
-- case
SELECT DAYOFYEAR('2007-02-03');
-- case
SELECT DAYNAME('2007-02-03');
-- case
SELECT FROM_DAYS(1423);
-- case
SELECT WEEKDAY('2007-02-03');
-- case
SELECT UTC_DATE, UTC_DATE();
-- case
SELECT UTC_DATE(), UTC_DATE()+0
-- case
SELECT WEEK();
-- case
SELECT WEEK('2007-02-03');
-- case
SELECT WEEK('2007-02-03', 0);
-- case
SELECT WEEKOFYEAR('2007-02-03');
-- case
SELECT MONTH('2007-02-03');
-- case
SELECT MONTHNAME('2007-02-03');
-- case
SELECT YEAR('2007-02-03');
-- case
SELECT YEARWEEK('2007-02-03');
-- case
SELECT YEARWEEK('2007-02-03', 0);
-- case
SELECT ADDTIME('01:00:00.999999', '02:00:00.999998');
-- case
SELECT ADDTIME('02:00:00.999998');
-- case
SELECT ADDTIME();
-- case
SELECT SUBTIME('01:00:00.999999', '02:00:00.999998');
-- case
SELECT CONVERT_TZ();
-- case
SELECT CONVERT_TZ('2004-01-01 12:00:00','+00:00','+10:00');
-- case
SELECT CONVERT_TZ('2004-01-01 12:00:00','+00:00','+10:00', '+10:00');
-- case
SELECT GET_FORMAT(DATE, 'USA');
-- case
SELECT GET_FORMAT(DATETIME, 'USA');
-- case
SELECT GET_FORMAT(TIME, 'USA');
-- case
SELECT GET_FORMAT(TIMESTAMP, 'USA');
-- case
SELECT LOCALTIME(), LOCALTIME(1)
-- case
SELECT LOCALTIMESTAMP(), LOCALTIMESTAMP(2)
-- case
SELECT MAKEDATE(2011,31);
-- case
SELECT MAKETIME(12,15,30);
-- case
SELECT MAKEDATE();
-- case
SELECT MAKETIME();
-- case
SELECT PERIOD_ADD(200801,2)
-- case
SELECT PERIOD_DIFF(200802,200703)
-- case
SELECT QUARTER('2008-04-01');
-- case
SELECT SEC_TO_TIME(2378)
-- case
SELECT TIME_FORMAT('100:00:00', '%H %k %h %I %l')
-- case
SELECT TIME_TO_SEC('22:23:00')
-- case
SELECT TIMESTAMPADD(WEEK,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_SECOND,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_MINUTE,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_HOUR,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_DAY,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_WEEK,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_MONTH,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_QUARTER,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_YEAR,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(SQL_TSI_MICROSECOND,1,'2003-01-02');
-- case
SELECT TIMESTAMPADD(MICROSECOND,1,'2003-01-02');
-- case
SELECT TO_DAYS('2007-10-07')
-- case
SELECT TO_SECONDS('2009-11-29')
-- case
SELECT LAST_DAY('2003-02-05');
-- case
SELECT UTC_TIME(), UTC_TIME(1)
-- case
select extract(microsecond from "2011-11-11 10:10:10.123456")
-- case
select extract(second from "2011-11-11 10:10:10.123456")
-- case
select extract(minute from "2011-11-11 10:10:10.123456")
-- case
select extract(hour from "2011-11-11 10:10:10.123456")
-- case
select extract(day from "2011-11-11 10:10:10.123456")
-- case
select extract(week from "2011-11-11 10:10:10.123456")
-- case
select extract(month from "2011-11-11 10:10:10.123456")
-- case
select extract(quarter from "2011-11-11 10:10:10.123456")
-- case
select extract(year from "2011-11-11 10:10:10.123456")
-- case
select extract(second_microsecond from "2011-11-11 10:10:10.123456")
-- case
select extract(minute_microsecond from "2011-11-11 10:10:10.123456")
-- case
select extract(minute_second from "2011-11-11 10:10:10.123456")
-- case
select extract(hour_microsecond from "2011-11-11 10:10:10.123456")
-- case
select extract(hour_second from "2011-11-11 10:10:10.123456")
-- case
select extract(hour_minute from "2011-11-11 10:10:10.123456")
-- case
select extract(day_microsecond from "2011-11-11 10:10:10.123456")
-- case
select extract(day_second from "2011-11-11 10:10:10.123456")
-- case
select extract(day_minute from "2011-11-11 10:10:10.123456")
-- case
select extract(day_hour from "2011-11-11 10:10:10.123456")
-- case
select extract(year_month from "2011-11-11 10:10:10.123456")
-- case
select from_unixtime(1447430881)
-- case
select from_unixtime(1447430881.123456)
-- case
select from_unixtime(1447430881.1234567)
-- case
select from_unixtime(1447430881.9999999)
-- case
select from_unixtime(1447430881, "%Y %D %M %h:%i:%s %x")
-- case
select from_unixtime(1447430881.123456, "%Y %D %M %h:%i:%s %x")
-- case
select from_unixtime(1447430881.1234567, "%Y %D %M %h:%i:%s %x")
-- case
SELECT CAST('test collated returns' AS CHAR CHARACTER SET utf8) COLLATE utf8_bin;
-- case
SELECT TRIM('  bar   ');
-- case
SELECT TRIM(LEADING 'x' FROM 'xxxbarxxx');
-- case
SELECT TRIM(BOTH 'x' FROM 'xxxbarxxx');
-- case
SELECT TRIM(TRAILING 'xyz' FROM 'barxxyz');
-- case
SELECT LTRIM(' foo ');
-- case
SELECT RTRIM(' bar ');
-- case
SELECT RPAD('hi', 6, 'c');
-- case
SELECT BIT_LENGTH('hi');
-- case
SELECT CHAR(65);
-- case
SELECT CHAR_LENGTH('abc');
-- case
SELECT CHARACTER_LENGTH('abc');
-- case
SELECT FIELD('ej', 'Hej', 'ej', 'Heja', 'hej', 'foo');
-- case
SELECT FIND_IN_SET('foo', 'foo,bar')
-- case
SELECT FIND_IN_SET('foo')
-- case
SELECT MAKE_SET(1,'a'), MAKE_SET(1,'a','b','c')
-- case
SELECT MID('Sakila', -5, 3)
-- case
SELECT OCT(12)
-- case
SELECT OCTET_LENGTH('text')
-- case
SELECT ORD('2')
-- case
SELECT POSITION('bar' IN 'foobarbar')
-- case
SELECT QUOTE('Don\'t!')
-- case
SELECT BIN(12)
-- case
SELECT ELT(1, 'ej', 'Heja', 'hej', 'foo')
-- case
SELECT EXPORT_SET(5,'Y','N'), EXPORT_SET(5,'Y','N',','), EXPORT_SET(5,'Y','N',',',4)
-- case
SELECT FORMAT(), FORMAT(12332.2,2,'de_DE'), FORMAT(12332.123456, 4)
-- case
SELECT FROM_BASE64('abc')
-- case
SELECT TO_BASE64('abc')
-- case
SELECT INSERT(), INSERT('Quadratic', 3, 4, 'What'), INSTR('foobarbar', 'bar')
-- case
SELECT LOAD_FILE('/tmp/picture')
-- case
SELECT LPAD('hi',4,'??')
-- case
SELECT LEFT("foobar", 3)
-- case
SELECT RIGHT("foobar", 3)
-- case
SELECT REPEAT("a", 10);
-- case
SELECT SLEEP(10);
-- case
SELECT ANY_VALUE(@arg);
-- case
SELECT INET_ATON('10.0.5.9');
-- case
SELECT INET_NTOA(167773449);
-- case
SELECT INET6_ATON('fdfe::5a55:caff:fefa:9089');
-- case
SELECT INET6_NTOA(INET_NTOA(167773449));
-- case
SELECT IS_FREE_LOCK(@str);
-- case
SELECT IS_IPV4('10.0.5.9');
-- case
SELECT IS_IPV4_COMPAT(INET6_ATON('::10.0.5.9'));
-- case
SELECT IS_IPV4_MAPPED(INET6_ATON('::10.0.5.9'));
-- case
SELECT IS_IPV6('10.0.5.9');
-- case
SELECT IS_USED_LOCK(@str);
-- case
SELECT NAME_CONST('myname', 14);
-- case
SELECT RELEASE_ALL_LOCKS();
-- case
SELECT UUID();
-- case
SELECT UUID_SHORT()
-- case
SELECT UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db')
-- case
SELECT UUID_TO_BIN('6ccd780c-baba-1026-9564-5b8c656024db', 1)
-- case
SELECT BIN_TO_UUID(0x6ccd780cbaba102695645b8c656024db)
-- case
SELECT BIN_TO_UUID(0x6ccd780cbaba102695645b8c656024db, 1)
-- case
SELECT SLEEP();
-- case
SELECT ANY_VALUE();
-- case
SELECT INET_ATON();
-- case
SELECT INET_NTOA();
-- case
SELECT INET6_ATON();
-- case
SELECT INET6_NTOA(INET_NTOA());
-- case
SELECT IS_FREE_LOCK();
-- case
SELECT IS_IPV4();
-- case
SELECT IS_IPV4_COMPAT(INET6_ATON());
-- case
SELECT IS_IPV4_MAPPED(INET6_ATON());
-- case
SELECT IS_IPV6()
-- case
SELECT IS_USED_LOCK();
-- case
SELECT NAME_CONST();
-- case
SELECT RELEASE_ALL_LOCKS(1);
-- case
SELECT UUID(1);
-- case
SELECT UUID_SHORT(1)
-- case
select "2011-11-11 10:10:10.123456" + interval 10 second
-- case
select "2011-11-11 10:10:10.123456" - interval 10 second
-- case
select  interval 10 second + "2011-11-11 10:10:10.123456"
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 second)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 minute)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 hour)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 day)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 week)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 month)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 quarter)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 year)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10:10" minute_second)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "10:10" hour_minute)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10.10 hour_minute)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "11 10" day_hour)
-- case
select date_add("2011-11-11 10:10:10.123456", interval "11-11" year_month)
-- case
select date_add("2011-11-11 10:10:10.123456", 10)
-- case
select date_add("2011-11-11 10:10:10.123456", 0.10)
-- case
select date_add("2011-11-11 10:10:10.123456", "11,11")
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_microsecond)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_second)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_minute)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_hour)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 10 sql_tsi_day)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_week)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_month)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_quarter)
-- case
select date_add("2011-11-11 10:10:10.123456", interval 1 sql_tsi_year)
-- case
select strcmp('abc', 'def')
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10 microsecond)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10 second)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10 minute)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10 hour)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10 day)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 1 week)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 1 month)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 1 quarter)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 1 year)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10:10" minute_second)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute)
-- case
select adddate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "11 10" day_hour)
-- case
select adddate("2011-11-11 10:10:10.123456", interval "11-11" year_month)
-- case
select adddate("2011-11-11 10:10:10.123456", 10)
-- case
select adddate("2011-11-11 10:10:10.123456", 0.10)
-- case
select adddate("2011-11-11 10:10:10.123456", "11,11")
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10 microsecond)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10 second)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10 minute)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10 hour)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10 day)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 1 week)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 1 month)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 1 quarter)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 1 year)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10:10" minute_second)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "10:10" hour_minute)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval 10.10 hour_minute)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "11 10" day_hour)
-- case
select date_sub("2011-11-11 10:10:10.123456", interval "11-11" year_month)
-- case
select date_sub("2011-11-11 10:10:10.123456", 10)
-- case
select date_sub("2011-11-11 10:10:10.123456", 0.10)
-- case
select date_sub("2011-11-11 10:10:10.123456", "11,11")
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10 microsecond)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10 second)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10 minute)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10 hour)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10 day)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 1 week)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 1 month)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 1 quarter)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 1 year)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10.10" second_microsecond)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10:10.10" minute_microsecond)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10:10" minute_second)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10:10:10.10" hour_microsecond)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10:10:10" hour_second)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "10:10" hour_minute)
-- case
select subdate("2011-11-11 10:10:10.123456", interval 10.10 hour_minute)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10.10" day_microsecond)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10:10" day_second)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "11 10:10" day_minute)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "11 10" day_hour)
-- case
select subdate("2011-11-11 10:10:10.123456", interval "11-11" year_month)
-- case
select subdate("2011-11-11 10:10:10.123456", 10)
-- case
select subdate("2011-11-11 10:10:10.123456", 0.10)
-- case
select subdate("2011-11-11 10:10:10.123456", "11,11")
-- case
select unix_timestamp()
-- case
select unix_timestamp('2015-11-13 10:20:19.012')
-- case
SELECT GET_LOCK('lock1',10);
-- case
SELECT RELEASE_LOCK('lock1');
-- case
select avg(), avg(c1,c2) from t;
-- case
select avg(distinct c1) from t;
-- case
select avg(distinctrow c1) from t;
-- case
select avg(distinct all c1) from t;
-- case
select avg(distinctrow all c1) from t;
-- case
select avg(c2) from t;
-- case
select bit_and(c1) from t;
-- case
select bit_and(all c1) from t;
-- case
select bit_and(distinct c1) from t;
-- case
select bit_and(distinctrow c1) from t;
-- case
select bit_and(distinctrow all c1) from t;
-- case
select bit_and(distinct all c1) from t;
-- case
select bit_and(), bit_and(distinct c1) from t;
-- case
select bit_and(), bit_and(distinctrow c1) from t;
-- case
select bit_and(), bit_and(all c1) from t;
-- case
select bit_or(c1) from t;
-- case
select bit_or(all c1) from t;
-- case
select bit_or(distinct c1) from t;
-- case
select bit_or(distinctrow c1) from t;
-- case
select bit_or(distinctrow all c1) from t;
-- case
select bit_or(distinct all c1) from t;
-- case
select bit_or(), bit_or(distinct c1) from t;
-- case
select bit_or(), bit_or(distinctrow c1) from t;
-- case
select bit_or(), bit_or(all c1) from t;
-- case
select bit_xor(c1) from t;
-- case
select bit_xor(all c1) from t;
-- case
select bit_xor(distinct c1) from t;
-- case
select bit_xor(distinctrow c1) from t;
-- case
select bit_xor(distinctrow all c1) from t;
-- case
select bit_xor(), bit_xor(distinct c1) from t;
-- case
select bit_xor(), bit_xor(distinctrow c1) from t;
-- case
select bit_xor(), bit_xor(all c1) from t;
-- case
select max(c1,c2) from t;
-- case
select max(distinct c1) from t;
-- case
select max(distinctrow c1) from t;
-- case
select max(distinct all c1) from t;
-- case
select max(distinctrow all c1) from t;
-- case
select max(c2) from t;
-- case
select min(c1,c2) from t;
-- case
select min(distinct c1) from t;
-- case
select min(distinctrow c1) from t;
-- case
select min(distinct all c1) from t;
-- case
select min(distinctrow all c1) from t;
-- case
select min(c2) from t;
-- case
select sum(c1,c2) from t;
-- case
select sum(distinct c1) from t;
-- case
select sum(distinctrow c1) from t;
-- case
select sum(distinct all c1) from t;
-- case
select sum(distinctrow all c1) from t;
-- case
select sum(c2) from t;
-- case
select count(c1) from t;
-- case
select count(distinct *) from t;
-- case
select count(distinctrow *) from t;
-- case
select count(*) from t;
-- case
select count(distinct c1, c2) from t;
-- case
select count(distinctrow c1, c2) from t;
-- case
select count(c1, c2) from t;
-- case
select count(all c1) from t;
-- case
select count(distinct all c1) from t;
-- case
select count(distinctrow all c1) from t;
-- case
select approx_count_distinct(c1) from t;
-- case
select approx_count_distinct(c1, c2) from t;
-- case
select approx_count_distinct(c1, 123) from t;
-- case
select approx_percentile(c1) from t;
-- case
select approx_percentile(c1, c2) from t;
-- case
select approx_percentile(c1, 123) from t;
-- case
select group_concat(c2,c1) from t group by c1;
-- case
select group_concat(c2,c1 SEPARATOR ';') from t group by c1;
-- case
select group_concat(distinct c2,c1) from t group by c1;
-- case
select group_concat(distinctrow c2,c1) from t group by c1;
-- case
SELECT student_name, GROUP_CONCAT(DISTINCT test_score ORDER BY test_score DESC SEPARATOR ' ') FROM student GROUP BY student_name;
-- case
select std(c1), std(all c1), std(distinct c1) from t
-- case
select std(c1, c2) from t
-- case
select stddev(c1), stddev(all c1), stddev(distinct c1) from t
-- case
select stddev(c1, c2) from t
-- case
select stddev_pop(c1), stddev_pop(all c1), stddev_pop(distinct c1) from t
-- case
select stddev_pop(c1, c2) from t
-- case
select stddev_samp(c1), stddev_samp(all c1), stddev_samp(distinct c1) from t
-- case
select stddev_samp(c1, c2) from t
-- case
select variance(c1), variance(all c1), variance(distinct c1) from t
-- case
select variance(c1, c2) from t
-- case
select var_pop(c1), var_pop(all c1), var_pop(distinct c1) from t
-- case
select var_pop(c1, c2) from t
-- case
select var_samp(c1), var_samp(all c1), var_samp(distinct c1) from t
-- case
select var_samp(c1, c2) from t
-- case
select json_arrayagg(c2) from t group by c1
-- case
select json_arrayagg(c1, c2) from t group by c1
-- case
select json_arrayagg(distinct c2) from t group by c1
-- case
select json_arrayagg(all c2) from t group by c1
-- case
select json_objectagg(c1, c2) from t group by c1
-- case
select json_objectagg(c1, c2, c3) from t group by c1
-- case
select json_objectagg(distinct c1, c2) from t group by c1
-- case
select json_objectagg(c1, distinct c2) from t group by c1
-- case
select json_objectagg(distinct c1, distinct c2) from t group by c1
-- case
select json_objectagg(all c1, c2) from t group by c1
-- case
select json_objectagg(c1, all c2) from t group by c1
-- case
select json_objectagg(all c1, all c2) from t group by c1
-- case
select AES_ENCRYPT('text',UNHEX('F3229A0B371ED2D9441B830D21A390C3'))
-- case
select AES_DECRYPT(@crypt_str,@key_str)
-- case
select AES_DECRYPT(@crypt_str,@key_str,@init_vector);
-- case
SELECT COMPRESS('');
-- case
SELECT DECODE(@crypt_str, @pass_str);
-- case
SELECT DES_DECRYPT(@crypt_str), DES_DECRYPT(@crypt_str, @key_str);
-- case
SELECT DES_ENCRYPT(@str), DES_ENCRYPT(@key_num);
-- case
SELECT ENCODE('cleartext', CONCAT('my_random_salt','my_secret_password'));
-- case
SELECT ENCRYPT('hello'), ENCRYPT('hello', @salt);
-- case
SELECT MD5('testing');
-- case
SELECT OLD_PASSWORD(@str);
-- case
SELECT PASSWORD(@str);
-- case
SELECT RANDOM_BYTES(@len);
-- case
SELECT SHA1('abc');
-- case
SELECT SHA('abc');
-- case
SELECT SHA2('abc', 224);
-- case
SELECT SM3('abc');
-- case
SELECT UNCOMPRESS('any string');
-- case
SELECT UNCOMPRESSED_LENGTH(@compressed_string);
-- case
SELECT VALIDATE_PASSWORD_STRENGTH(@str);
-- case
SELECT JSON_EXTRACT();
-- case
SELECT JSON_UNQUOTE();
-- case
SELECT JSON_TYPE('[123]');
-- case
SELECT JSON_TYPE();
-- case
SELECT a->'$.a' FROM t
-- case
SELECT a->>'$.a' FROM t
-- case
SELECT '{}'->'$.a' FROM t
-- case
SELECT '{}'->>'$.a' FROM t
-- case
SELECT a->3 FROM t
-- case
SELECT a->>3 FROM t
-- case
SELECT 1 member of (a)
-- case
SELECT 1 member of a
-- case
SELECT 1 member a
-- case
SELECT 1 not member of a
-- case
SELECT 1 member of (1+1)
-- case
SELECT concat('a') member of (cast(1 as char(1)))
-- case
SELECT `uuid`()
-- case
select nextval(seq)
-- case
select lastval(seq)
-- case
select setval(seq, 100)
-- case
select next value for seq
-- case
select next value for sequence
-- case
select NeXt vAluE for seQuEncE2
-- case
select regexp_like('aBc', 'abc', 'im');
-- case
select regexp_substr('aBc', 'abc', 1, 1, 'im');
-- case
select regexp_instr('aBc', 'abc', 1, 1, 0, 'im');
-- case
select regexp_replace('aBc', 'abc', 'def', 1, 1, 'i');
-- case
select 'aBc' ilike 'abc';
