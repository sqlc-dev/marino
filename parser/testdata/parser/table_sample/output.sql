SELECT * FROM `tbl` TABLESAMPLE SYSTEM (50)
-- case
SELECT * FROM `tbl` TABLESAMPLE SYSTEM (50 PERCENT)
-- case
SELECT * FROM `tbl` TABLESAMPLE SYSTEM (49.9 PERCENT)
-- case
SELECT * FROM `tbl` TABLESAMPLE SYSTEM (120 ROWS)
-- case
SELECT * FROM `tbl` TABLESAMPLE BERNOULLI (50)
-- case
SELECT * FROM `tbl` TABLESAMPLE (50)
-- case
SELECT * FROM `tbl` TABLESAMPLE (50) REPEATABLE(123456789)
-- case
SELECT * FROM `tbl` AS `a` TABLESAMPLE (50)
-- case
SELECT * FROM `tbl` AS `tablesample` TABLESAMPLE (50)
-- case
SELECT * FROM `tbl` TABLESAMPLE (50) WHERE `id`>20
-- case
SELECT * FROM `tbl` PARTITION(`p0`) TABLESAMPLE (50)
-- case
SELECT * FROM `tbl` TABLESAMPLE (0 PERCENT)
-- case
SELECT * FROM `tbl` TABLESAMPLE (100 PERCENT)
-- case
SELECT * FROM `tbl` TABLESAMPLE (0 ROWS)
-- case
SELECT * FROM `tbl` TABLESAMPLE (_UTF8MB4'34')
-- case
SELECT * FROM (`tbl1` TABLESAMPLE (10)) JOIN `tbl2` TABLESAMPLE (20)
-- case
SELECT * FROM `tbl1` AS `a` TABLESAMPLE (10) JOIN `tbl2` AS `b` TABLESAMPLE (20) ON `a`.`id`!=`b`.`id`
-- case
SELECT * FROM `demo` TABLESAMPLE BERNOULLI (50) LIMIT 1 INTO OUTFILE '/tmp/sample.csv'
-- case
SELECT * FROM `demo` TABLESAMPLE BERNOULLI (50) ORDER BY `a`,`b` INTO OUTFILE '/tmp/sample.csv'
-- case
-- error: line 1 column 42 near "a;" 
-- case
-- error: line 1 column 44 near "partition (p0);" 
-- case
-- error: line 1 column 43 near "tablesample system(50);" 
-- case
-- error: line 1 column 47 near "tablesample system(50);" 
-- case
-- error: line 1 column 52 near "tablesample system(50);" 
-- case
-- error: line 1 column 40 near ", 50);" 
-- case
-- error: line 1 column 44 near "dhfksdlfljcoew(50);" 
-- case
-- error: line 1 column 37 near ";" 
-- case
-- error: line 1 column 53 near ";" 
-- case
-- error: line 1 column 30 near "tablesample system (50);" 
