HANDLER t OPEN
-- case
HANDLER t OPEN AS h
-- case
HANDLER db1.t OPEN h
-- case
HANDLER t READ FIRST
-- case
HANDLER t READ NEXT
-- case
HANDLER t READ FIRST WHERE c > 1 LIMIT 10
-- case
HANDLER t READ idx FIRST
-- case
HANDLER t READ idx NEXT
-- case
HANDLER t READ idx PREV
-- case
HANDLER t READ idx LAST
-- case
HANDLER t READ idx = (1, 'a')
-- case
HANDLER t READ idx <= (1)
-- case
HANDLER t READ idx >= (1) WHERE c < 100
-- case
HANDLER t READ idx < (1)
-- case
HANDLER t READ idx > (1) LIMIT 2, 5
-- case
HANDLER t READ PRIMARY = (1)
-- case
HANDLER t CLOSE
-- case
IMPORT TABLE FROM 't.sdi'
-- case
IMPORT TABLE FROM 't.sdi', 'u.sdi'
-- case
LOAD XML INFILE 'f.xml' INTO TABLE t
-- case
LOAD XML LOW_PRIORITY LOCAL INFILE 'f.xml' REPLACE INTO TABLE db1.t
-- case
LOAD XML CONCURRENT INFILE 'f.xml' IGNORE INTO TABLE t CHARACTER SET utf8mb4 ROWS IDENTIFIED BY '<row>' IGNORE 2 LINES (a, @b) SET c = 1
-- case
SELECT 1 INTO @a
-- case
SELECT a, b INTO @x, @y FROM t LIMIT 1
-- case
SELECT 1 INTO DUMPFILE '/tmp/out'
-- case
SELECT prev, xml, dumpfile FROM concurrent
