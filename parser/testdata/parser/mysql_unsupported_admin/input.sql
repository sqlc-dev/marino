CREATE RESOURCE GROUP rg1 TYPE = USER VCPU = 0-3
-- case
ALTER RESOURCE GROUP rg1 VCPU = 0-3
-- case
CHECK TABLE t
-- case
CHECKSUM TABLE t
-- case
REPAIR TABLE t
-- case
CREATE FUNCTION metaphon RETURNS STRING SONAME 'udf.so'
-- case
INSTALL COMPONENT 'file://component_validate_password'
-- case
INSTALL PLUGIN myplugin SONAME 'plugin.so'
-- case
UNINSTALL COMPONENT 'file://component_validate_password'
-- case
UNINSTALL PLUGIN myplugin
-- case
CLONE LOCAL DATA DIRECTORY = '/tmp/clone'
-- case
CLONE INSTANCE FROM 'user'@'host':3306 IDENTIFIED BY 'password'
-- case
CACHE INDEX t IN hot_cache
-- case
LOAD INDEX INTO CACHE t
-- case
RESET PERSIST
