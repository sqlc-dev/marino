CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 3.1415 YEAR
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL _UTF8MB4'1 1:1:1' DAY_SECOND
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 1 YEAR
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'OFF'
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'OFF'
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'OFF' TTL_JOB_INTERVAL = '8h'
-- case
CREATE TABLE `t` (`created_at` DATETIME) TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'ON'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 MONTH
-- case
ALTER TABLE `t` TTL_ENABLE = 'ON'
-- case
ALTER TABLE `t` TTL_ENABLE = 'OFF'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 MONTH TTL_ENABLE = 'OFF'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 MONTH TTL_ENABLE = 'OFF' TTL_JOB_INTERVAL = '1h'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'ON'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'ON' TTL_JOB_INTERVAL = '8h'
-- case
ALTER TABLE `t` TTL = `created_at` + INTERVAL 1 YEAR TTL_ENABLE = 'ON' TTL_JOB_INTERVAL = '8.645124531235h'
-- case
ALTER TABLE `t` REMOVE TTL
-- case
-- error: line 1 column 61 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
-- case
-- error: line 1 column 74 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
-- case
-- error: line 1 column 51 near ""The TTL_ENABLE option has to be set 'ON' or 'OFF' 
-- case
-- error: line 1 column 66 near ""The TTL_JOB_INTERVAL option is not a valid duration: fail to read an integer 
-- case
-- error: line 1 column 66 near ""The TTL_JOB_INTERVAL option is not a valid duration: fail to read an integer 
-- case
-- error: line 1 column 68 near ""The TTL_JOB_INTERVAL option is not a valid duration: strconv.ParseFloat: parsing "10.10.255": invalid syntax 
