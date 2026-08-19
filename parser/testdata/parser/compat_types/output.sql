CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` TINYINT(1))
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` TINYINT(1))
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` VARCHAR(0))
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` DECIMAL)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` FLOAT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` DOUBLE)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` TINYINT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` SMALLINT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` MEDIUMINT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` INT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` BIGINT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` MEDIUMBLOB)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` MEDIUMTEXT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` MEDIUMINT)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` DECIMAL)
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` BINARY(5))
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` CHAR(5) CHARACTER SET LATIN1)
-- case
-- error: [parser:1115]Unknown character set: 'ucs2'
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` VARBINARY(10))
-- case
CREATE TABLE `t` (`id` INT PRIMARY KEY,`c1` VARCHAR(10) CHARACTER SET LATIN1)
-- case
-- error: [parser:1115]Unknown character set: 'ucs2'
