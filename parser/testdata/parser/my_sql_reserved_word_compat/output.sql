-- error: line 1 column 17 near "cube (a int)" 
-- case
-- error: line 1 column 21 near "external (a int)" 
-- case
-- error: line 1 column 20 near "qualify (a int)" 
-- case
-- error: line 1 column 11 near "cube from t" 
-- case
-- error: line 1 column 15 near "external from t" 
-- case
-- error: line 1 column 14 near "qualify from t" 
-- case
CREATE TABLE `cube` (`a` INT)
-- case
CREATE TABLE `external` (`a` INT)
-- case
CREATE TABLE `qualify` (`a` INT)
-- case
SELECT `t`.`cube` FROM `t`
-- case
SELECT `t`.`external` FROM `t`
-- case
SELECT `t`.`qualify` FROM `t`
-- case
CREATE TABLE `manual` (`parallel` INT)
-- case
SELECT `manual`,`parallel` FROM `t`
