CREATE TABLE `t` (`id` INT PRIMARY KEY,`col` INT,INDEX(`col`) WHERE `col`>100)
-- case
CREATE INDEX `idx` ON `t` (`col`) WHERE `col`>100
-- case
ALTER TABLE `t` ADD INDEX `idx`(`col`) WHERE `col`>100
