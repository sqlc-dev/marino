CREATE TABLE `t` (`c1` TIME(2),`c2` DATETIME(2),`c3` TIMESTAMP(2))
-- case
SELECT x'0a',x'11',x'11'
-- case
SELECT x'13181c76734725455a'
-- case
-- error: line 1 column 10 near "x'0xaa'" 
-- case
-- error: line 1 column 11 near "0X11"hex literal: invalid hexadecimal format 0X11 
-- case
SELECT x'4920616d2061206c6f6e672068657820737472696e67'
-- case
SELECT b'1',b'0',b'11',b'11'
-- case
CREATE TABLE `t` (`c1` ENUM('a','b'),`c2` SET('a','b'))
-- case
CREATE TABLE `t` (`c1` ENUM('a','b	'),`c2` SET('a','b	'))
-- case
CREATE TABLE `t` (`c1` ENUM('a','b') BINARY,`c2` SET('a','b') BINARY)
-- case
CREATE TABLE `t` (`c1` ENUM('a','b'),`c2` SET('a','b'))
-- case
CREATE TABLE `t` (`c1` ENUM('a','b'),`c2` SET('a','b'))
-- case
-- error: line 1 column 24 near ")" 
-- case
-- error: line 1 column 23 near ")" 
-- case
CREATE TABLE `t` (`c1` BLOB(1024),`c2` TEXT(1024))
-- case
CREATE TABLE `t` (`y` YEAR(4),`y1` YEAR)
-- case
CREATE TABLE `t` (`y` YEAR(4),`y1` YEAR)
-- case
CREATE TABLE `t` (`c1` CHAR(2),`c2` VARCHAR(2))
-- case
CREATE TABLE `t` (`a` JSON)
