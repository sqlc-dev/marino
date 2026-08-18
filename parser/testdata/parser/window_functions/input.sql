-- flags: window_func
SELECT CUME_DIST() OVER w FROM t;
-- case
SELECT DENSE_RANK() OVER (w) FROM t;
-- case
SELECT FIRST_VALUE(val) OVER w FROM t;
-- case
SELECT FIRST_VALUE(val) RESPECT NULLS OVER w FROM t;
-- case
SELECT FIRST_VALUE(val) IGNORE NULLS OVER w FROM t;
-- case
SELECT LAG(val) OVER (w) FROM t;
-- case
SELECT LAG(val, 1) OVER (w) FROM t;
-- case
SELECT LAG(val, 1, def) OVER (w) FROM t;
-- case
SELECT LAST_VALUE(val) OVER (w) FROM t;
-- case
SELECT LEAD(val) OVER w FROM t;
-- case
SELECT LEAD(val, 1) OVER w FROM t;
-- case
SELECT LEAD(val, 1, def) OVER w FROM t;
-- case
SELECT NTH_VALUE(val, 233) OVER w FROM t;
-- case
SELECT NTH_VALUE(val, 233) FROM FIRST OVER w FROM t;
-- case
SELECT NTH_VALUE(val, 233) FROM LAST OVER w FROM t;
-- case
SELECT NTH_VALUE(val, 233) FROM LAST IGNORE NULLS OVER w FROM t;
-- case
SELECT NTH_VALUE(val) OVER w FROM t;
-- case
SELECT NTILE(233) OVER (w) FROM t;
-- case
SELECT PERCENT_RANK() OVER (w) FROM t;
-- case
SELECT RANK() OVER (w) FROM t;
-- case
SELECT ROW_NUMBER() OVER (w) FROM t;
-- case
SELECT n, LAG(n, 1, 0) OVER (w), LEAD(n, 1, 0) OVER w, n + LAG(n, 1, 0) OVER (w) FROM fib;
-- case
SELECT SUM(profit) OVER(PARTITION BY country) AS country_profit FROM sales;
-- case
SELECT SUM(profit) OVER() AS country_profit FROM sales;
-- case
SELECT AVG(profit) OVER() AS country_profit FROM sales;
-- case
SELECT BIT_XOR(profit) OVER() AS country_profit FROM sales;
-- case
SELECT COUNT(profit) OVER() AS country_profit FROM sales;
-- case
SELECT COUNT(ALL profit) OVER() AS country_profit FROM sales;
-- case
SELECT COUNT(*) OVER() AS country_profit FROM sales;
-- case
SELECT MAX(profit) OVER() AS country_profit FROM sales;
-- case
SELECT MIN(profit) OVER() AS country_profit FROM sales;
-- case
SELECT SUM(profit) OVER() AS country_profit FROM sales;
-- case
SELECT ROW_NUMBER() OVER(PARTITION BY country) AS row_num1 FROM sales;
-- case
SELECT ROW_NUMBER() OVER(PARTITION BY country, d ORDER BY year, product) AS row_num2 FROM sales;
-- case
SELECT SUM(val) OVER (PARTITION BY subject ORDER BY time ROWS UNBOUNDED PRECEDING) FROM t;
-- case
SELECT AVG(val) OVER (PARTITION BY subject ORDER BY time ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t;
-- case
SELECT AVG(val) OVER (ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t;
-- case
SELECT AVG(val) OVER (ROWS BETWEEN 1 PRECEDING AND UNBOUNDED FOLLOWING) FROM t;
-- case
SELECT AVG(val) OVER (RANGE BETWEEN INTERVAL 5 DAY PRECEDING AND INTERVAL '2:30' MINUTE_SECOND FOLLOWING) FROM t;
-- case
SELECT AVG(val) OVER (RANGE BETWEEN CURRENT ROW AND CURRENT ROW) FROM t;
-- case
SELECT AVG(val) OVER (RANGE CURRENT ROW) FROM t;
-- case
SELECT RANK() OVER (w) FROM t WINDOW w AS (ORDER BY val);
-- case
SELECT RANK() OVER w FROM t WINDOW w AS ();
-- case
SELECT FIRST_VALUE(year) OVER (w ORDER BY year) AS first FROM sales WINDOW w AS (PARTITION BY country);
-- case
SELECT RANK() OVER (w1) FROM t WINDOW w1 AS (w2), w2 AS (), w3 AS (w1);
-- case
SELECT RANK() OVER w1 FROM t WINDOW w1 AS (w2), w2 AS (w3), w3 AS (w1);
-- case
select tidb_parse_tso(1)
-- case
select tidb_parse_tso_logical(1)
-- case
select tidb_bounded_staleness('2015-09-21 00:07:01', NOW())
-- case
select tidb_bounded_staleness(DATE_SUB(NOW(), INTERVAL 3 SECOND), NOW())
-- case
select tidb_bounded_staleness('2015-09-21 00:07:01', '2021-04-27 11:26:13')
-- case
select from_unixtime(404411537129996288)
-- case
select from_unixtime(404411537129996288.22)
