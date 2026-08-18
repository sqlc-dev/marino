-- error: line 1 column 36 near "rollup" 
-- case
SELECT * FROM `t` GROUP BY `a`,`b` WITH ROLLUP
-- case
SELECT * FROM `t` GROUP BY ROW(`a`,`b`) WITH ROLLUP
-- case
SELECT * FROM `t` GROUP BY (`a`+`b`) WITH ROLLUP
