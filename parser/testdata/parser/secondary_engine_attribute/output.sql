CREATE TABLE `t` (`id` INT) PARTITION BY RANGE (`id`) (PARTITION `p0` VALUES LESS THAN (10) SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value"}',PARTITION `p1` VALUES LESS THAN (20) SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value2"}')
-- case
CREATE TABLE `t` (`id` INT) SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value"}'
-- case
CREATE TABLE `t` (`id` INT) SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value"}' PARTITION BY RANGE (`id`) (PARTITION `p0` VALUES LESS THAN (10) SECONDARY_ENGINE_ATTRIBUTE = '{"key":"partition_value"}')
-- case
CREATE TABLE `t` (`id` INT SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value"}')
-- case
CREATE TABLE `t` (`id` INT) TABLESPACE = `ts1` SECONDARY_ENGINE_ATTRIBUTE = '{"key":"value"}'
-- case
CREATE TABLE `t` (`id` INT,INDEX `idx`(`id`) INVISIBLE SECONDARY_ENGINE_ATTRIBUTE = '{"key1":"value1"}')
-- case
-- error: line 1 column 112 near ")" 
-- case
-- error: line 1 column 51 near "" 
-- case
-- error: line 1 column 51 near ")" 
-- case
-- error: line 1 column 66 near "" 
-- case
-- error: line 1 column 67 near ")" 
-- case
-- error: line 1 column 111 near ")" 
-- case
-- error: line 1 column 50 near "" 
-- case
-- error: line 1 column 50 near ")" 
-- case
-- error: line 1 column 65 near "" 
-- case
-- error: line 1 column 66 near ")" 
-- case
CREATE INDEX `i` ON `t` (`a`) SECONDARY_ENGINE_ATTRIBUTE = '{}'
-- case
CREATE INDEX `i` ON `t` (`a`) SECONDARY_ENGINE_ATTRIBUTE = '{}'
-- case
-- error: line 1 column 50 near "" 
