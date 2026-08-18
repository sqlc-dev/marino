CREATE TABLE `t` (`c1` INT,`c2` INT) PARTITION BY LINEAR KEY ALGORITHM = 1 (`c1`,`c2`) PARTITIONS 4
-- case
-- error: line 1 column 78 near "-1 (c1,c2) PARTITIONS 4" 
-- case
-- error: [parser:1149]You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
-- case
-- error: [parser:1149]You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use
