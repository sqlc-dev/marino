CREATE TABLE `t` (`a` INT) AFFINITY = 'table'
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'TABLE'
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'partition'
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = ''
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'none'
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'PARTITION' PARTITION BY HASH (`a`) PARTITIONS 1
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'table'
-- case
CREATE TABLE `t` (`a` INT) AFFINITY = 'abcd'
-- case
ALTER TABLE `t` AFFINITY = 'table'
-- case
ALTER TABLE `t` AFFINITY = 'TABLE'
-- case
ALTER TABLE `t` AFFINITY = 'table'
-- case
-- error: line 1 column 33 near "1" 
-- case
-- error: line 1 column 35 near "1" 
-- case
-- error: line 1 column 31 near "" 
