CREATE TABLE `t` (`g` GEOMETRY)
-- case
CREATE TABLE `t` (`pt` POINT)
-- case
CREATE TABLE `t` (`ls` LINESTRING)
-- case
CREATE TABLE `t` (`pg` POLYGON)
-- case
CREATE TABLE `t` (`mpt` MULTIPOINT)
-- case
CREATE TABLE `t` (`mls` MULTILINESTRING)
-- case
CREATE TABLE `t` (`mpg` MULTIPOLYGON)
-- case
CREATE TABLE `t` (`gc` GEOMETRYCOLLECTION)
-- case
CREATE TABLE `t` (`gc` GEOMETRYCOLLECTION)
-- case
CREATE TABLE `t` (`g` GEOMETRY NOT NULL)
-- case
CREATE TABLE `t` (`g` GEOMETRY NOT NULL SRID 4326)
-- case
CREATE TABLE `t` (`pt` POINT SRID 0)
-- case
CREATE TABLE `t` (`g` GEOMETRY NOT NULL SRID 4326,SPATIAL `idx_g`(`g`))
-- case
CREATE TABLE `t` (`g` GEOMETRY NOT NULL,SPATIAL(`g`))
-- case
CREATE TABLE `t` (`g` GEOMETRY NOT NULL,SPATIAL(`g`))
-- case
ALTER TABLE `t` ADD SPATIAL `idx_g`(`g`)
-- case
CREATE SPATIAL INDEX `idx_g` ON `t` (`g`)
-- case
CREATE TABLE `geom` (`g` GEOMETRY,`p` POINT NOT NULL SRID 4326,`ll` LINESTRING COMMENT 'path')
-- case
ALTER TABLE `t` ADD COLUMN `g` GEOMETRY NOT NULL SRID 4326
-- case
ALTER TABLE `t` MODIFY COLUMN `pt` POINT SRID 4326
-- case
CREATE TABLE `t` (`geometry` INT,`point` INT,`linestring` INT,`polygon` INT,`multipoint` INT,`multilinestring` INT,`multipolygon` INT,`geometrycollection` INT,`srid` INT)
-- case
SELECT `geometry`,`point`,`srid` FROM `polygon`
-- case
CREATE TABLE `t` (`g` GEOMETRY SRID 4326 COMMENT 'geo column' NOT NULL)
-- case
SELECT CAST(`g` AS POINT) FROM `t`
-- case
SELECT CAST(`x` AS GEOMETRY),CAST(`x` AS GEOMETRYCOLLECTION),CAST(`x` AS GEOMETRYCOLLECTION) FROM `t`
-- case
SELECT CONVERT(`x`, MULTIPOLYGON),CONVERT(`x`, LINESTRING) FROM `t`
