select SQL_BIG_RESULT c1 from t group by c1
-- case
select SQL_SMALL_RESULT c1 from t group by c1
-- case
select SQL_BUFFER_RESULT * from t
-- case
select sql_small_result sql_big_result sql_buffer_result 1
-- case
select STRAIGHT_JOIN SQL_SMALL_RESULT * from t
-- case
select SQL_CALC_FOUND_ROWS DISTINCT * from t
