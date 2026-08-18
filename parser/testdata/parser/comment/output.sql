CREATE TABLE `t` (`c` INT COMMENT 'comment')
-- case
CREATE TABLE `t` (`c` INT) COMMENT = 'comment'
-- case
CREATE TABLE `t` (`c` INT) COMMENT = 'comment'
-- case
-- error: line 1 column 38 near "comment" 
-- case
CREATE TABLE `t` (`comment` TEXT)
-- case
START TRANSACTION
-- case
SELECT `c` FROM `t`
-- case
-- error: near '/*' and b = 'p'' at line 1
-- case
-- error: line 1 column 19 near "ssl int)" 
-- case
-- error: line 1 column 23 near "require int)" 
-- case
CREATE TABLE `t` (`account` INT)
-- case
CREATE TABLE `t` (`expire` INT)
-- case
CREATE TABLE `t` (`cipher` INT)
-- case
CREATE TABLE `t` (`issuer` INT)
-- case
CREATE TABLE `t` (`never` INT)
-- case
CREATE TABLE `t` (`subject` INT)
-- case
CREATE TABLE `t` (`x509` INT)
-- case
-- error: line 1 column 68 near "'{"name": "Tom", "age", 19}" 
-- case
-- error: line 1 column 67 near "'{"name": "Tom", "age", 19}" 
-- case
CREATE USER `commentUser`@`%` COMMENT '123456'
-- case
ALTER USER `commentUser`@`%` COMMENT '123456'
-- case
CREATE USER `commentUser`@`%` ATTRIBUTE '{"name": "Tom", "age", 19}'
-- case
ALTER USER `commentUser`@`%` ATTRIBUTE '{"name": "Tom", "age", 19}'
