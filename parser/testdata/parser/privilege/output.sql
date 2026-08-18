CREATE USER `ttt`@`%` REQUIRE X509
-- case
CREATE USER `ttt`@`%` REQUIRE SSL
-- case
CREATE USER `ttt`@`%` REQUIRE NONE
-- case
CREATE USER `ttt`@`%` REQUIRE ISSUER '/C=SE/ST=Stockholm/L=Stockholm/O=MySQL/CN=CA/emailAddress=ca@example.com' AND CIPHER 'EDH-RSA-DES-CBC3-SHA'
-- case
CREATE USER `ttt`@`%` REQUIRE ISSUER '/C=SE/ST=Stockholm/L=Stockholm/O=MySQL/CN=CA/emailAddress=ca@example.com' AND CIPHER 'EDH-RSA-DES-CBC3-SHA' AND SUBJECT '/C=SE/ST=Stockholm/L=Stockholm/O=MySQL/CN=CA/emailAddress=ca@example.com'
-- case
CREATE USER `ttt`@`%` REQUIRE SAN 'DNS:mysql-user, URI:spiffe://example.org/myservice'
-- case
CREATE USER `ttt`@`%` WITH MAX_QUERIES_PER_HOUR 2
-- case
CREATE USER `ttt`@`localhost` REQUIRE NONE WITH MAX_QUERIES_PER_HOUR 1 MAX_UPDATES_PER_HOUR 10 PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
-- case
CREATE USER `u1`@`%` IDENTIFIED WITH 'mysql_native_password' AS '' REQUIRE NONE PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
-- case
CREATE USER `test`@`%`
-- case
CREATE USER `test`@`%`
-- case
CREATE USER `test`@`%`
-- case
-- error: line 1 column 17 near "-user" 
-- case
-- error: line 1 column 17 near ".user" 
-- case
CREATE USER `test-user`@`%`
-- case
CREATE USER `test-user`@`%`
-- case
-- error: line 1 column 17 near ".user" 
-- case
CREATE USER `test.user`@`%`
-- case
CREATE USER `test.user`@`%`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `uesr1`@`localhost`
-- case
CREATE USER `test@localhost`@`%` PASSWORD EXPIRE
-- case
CREATE USER `test@localhost`@`%` PASSWORD EXPIRE NEVER
-- case
CREATE USER `test@localhost`@`%` PASSWORD EXPIRE DEFAULT
-- case
CREATE USER `test@localhost`@`%` PASSWORD EXPIRE INTERVAL 3 DAY
-- case
CREATE USER `test@localhost`@`%` IDENTIFIED BY 'password' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME 3
-- case
CREATE USER `test@localhost`@`%` IDENTIFIED BY 'password' FAILED_LOGIN_ATTEMPTS 3 PASSWORD_LOCK_TIME UNBOUNDED
-- case
CREATE USER `test@localhost`@`%` IDENTIFIED BY 'password' FAILED_LOGIN_ATTEMPTS 3
-- case
CREATE USER `test@localhost`@`%` IDENTIFIED BY 'password' PASSWORD_LOCK_TIME 3
-- case
CREATE USER `test@localhost`@`%` IDENTIFIED BY 'password' PASSWORD_LOCK_TIME UNBOUNDED
-- case
CREATE USER `sha_test`@`localhost` IDENTIFIED WITH 'caching_sha2_password' BY 'sha_test'
-- case
CREATE USER `sha_test3`@`localhost` IDENTIFIED WITH 'caching_sha2_password' AS '$A$005$%[Ilf,U{;eJ/ }HRvp6DpKv02GGEdAjQfb4lXhdm2vLkQOCXZG7yVYGQO4'
-- case
CREATE USER `sha_test4`@`localhost` IDENTIFIED WITH 'caching_sha2_password' AS '$A$005$%[Ilf,U{;eJ/ }HRvp6DpKv02GGEdAjQfb4lXhdm2vLkQOCXZG7yVYGQO4'
-- case
CREATE USER `user@pingcap.com`@`localhost` IDENTIFIED WITH 'tidb_auth_token' REQUIRE TOKEN_ISSUER 'issuer-abc' ATTRIBUTE '{"email": "user@pingcap.com"}'
-- case
CREATE USER `nopwd_native`@`localhost` IDENTIFIED WITH 'mysql_native_password'
-- case
CREATE USER `nopwd_sha`@`localhost` IDENTIFIED WITH 'caching_sha2_password'
-- case
CREATE ROLE `test-role`@`%`, `role1`@`localhost`
-- case
CREATE ROLE `test-role`@`%`
-- case
CREATE ROLE `role1`@`%`
-- case
CREATE ROLE `role1`@`localhost`
-- case
CREATE USER `bug19354014user`@`%` IDENTIFIED WITH 'mysql_native_password'
-- case
CREATE USER `bug19354014user`@`%` IDENTIFIED WITH 'mysql_native_password' BY 'new-password'
-- case
CREATE USER `bug19354014user`@`%` IDENTIFIED WITH 'mysql_native_password' AS 'hashstring'
-- case
CREATE USER IF NOT EXISTS `root`@`localhost` IDENTIFIED BY 'new-password'
-- case
CREATE USER `root`@`localhost` IDENTIFIED BY 'new-password'
-- case
CREATE USER `root`@`localhost` IDENTIFIED WITH 'mysql_native_password' AS 'hashstring'
-- case
CREATE USER `root`@`localhost` IDENTIFIED BY 'new-password', `root`@`127.0.0.1` IDENTIFIED WITH 'mysql_native_password' AS 'hashstring'
-- case
CREATE USER `root`@`127.0.0.1` IDENTIFIED BY 'hashstring' RESOURCE GROUP `rg1`
-- case
ALTER USER IF EXISTS `root`@`localhost` IDENTIFIED BY 'new-password'
-- case
ALTER USER `root`@`localhost` IDENTIFIED BY 'new-password'
-- case
ALTER USER `root`@`localhost` RESOURCE GROUP `rg2`
-- case
ALTER USER `root`@`localhost` IDENTIFIED WITH 'mysql_native_password' AS 'hashstring'
-- case
ALTER USER `root`@`localhost` IDENTIFIED BY 'new-password', `root`@`127.0.0.1` IDENTIFIED WITH 'mysql_native_password' AS 'hashstring'
-- case
ALTER USER USER() IDENTIFIED BY 'new-password'
-- case
ALTER USER IF EXISTS USER() IDENTIFIED BY 'new-password'
-- case
ALTER USER `test@localhost`@`%` PASSWORD EXPIRE
-- case
ALTER USER `test@localhost`@`%` PASSWORD EXPIRE NEVER
-- case
ALTER USER `test@localhost`@`%` PASSWORD EXPIRE DEFAULT
-- case
ALTER USER `test@localhost`@`%` PASSWORD EXPIRE INTERVAL 3 DAY
-- case
ALTER USER `ttt`@`%` REQUIRE X509
-- case
ALTER USER `ttt`@`%` REQUIRE SSL
-- case
ALTER USER `ttt`@`%` REQUIRE NONE
-- case
ALTER USER `ttt`@`%` REQUIRE ISSUER '/C=SE/ST=Stockholm/L=Stockholm/O=MySQL/CN=CA/emailAddress=ca@example.com' AND CIPHER 'EDH-RSA-DES-CBC3-SHA'
-- case
ALTER USER `ttt`@`%` WITH MAX_QUERIES_PER_HOUR 2
-- case
ALTER USER `ttt`@`%` WITH MAX_UPDATES_PER_HOUR 2
-- case
ALTER USER `ttt`@`%` WITH MAX_CONNECTIONS_PER_HOUR 2
-- case
ALTER USER `ttt`@`%` WITH MAX_USER_CONNECTIONS 2
-- case
ALTER USER `ttt`@`localhost` REQUIRE NONE WITH MAX_QUERIES_PER_HOUR 1 MAX_UPDATES_PER_HOUR 10 PASSWORD EXPIRE DEFAULT ACCOUNT UNLOCK
-- case
DROP USER `root`@`localhost`, `root1`@`localhost`
-- case
DROP USER IF EXISTS `root`@`localhost`
-- case
RENAME USER `root`@`localhost` TO `root`@`%`
-- case
RENAME USER `fred`@`%` TO `barry`@`%`
-- case
RENAME USER `u1`@`%` TO `u2`@`%`, `u3`@`%` TO `u4`@`%`
-- case
DROP ROLE `role`@`localhost`, `role1`@`localhost`
-- case
DROP ROLE `administrator`@`%`, `developer`@`%`
-- case
DROP ROLE IF EXISTS `role`@`localhost`
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost` REQUIRE X509
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost` REQUIRE SSL
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost` REQUIRE NONE
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost` REQUIRE ISSUER '/C=SE/ST=Stockholm/L=Stockholm/O=MySQL/CN=CA/emailAddress=ca@example.com' AND CIPHER 'EDH-RSA-DES-CBC3-SHA'
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost`
-- case
GRANT ALL ON TABLE `db1`.* TO `jeffrey`@`localhost`
-- case
GRANT ALL ON `db1`.* TO `jeffrey`@`localhost` WITH GRANT OPTION
-- case
GRANT SELECT ON `db2`.`invoice` TO `jeffrey`@`localhost`
-- case
GRANT ALL ON *.* TO `someuser`@`somehost`
-- case
GRANT ALL ON *.* TO `SOMEuser`@`somehost`
-- case
GRANT SELECT, INSERT ON *.* TO `someuser`@`somehost`
-- case
GRANT ALL ON `mydb`.* TO `someuser`@`somehost`
-- case
GRANT SELECT, INSERT ON `mydb`.* TO `someuser`@`somehost`
-- case
GRANT ALL ON `mydb`.`mytbl` TO `someuser`@`somehost`
-- case
GRANT SELECT, INSERT ON `mydb`.`mytbl` TO `someuser`@`somehost`
-- case
GRANT SELECT (`col1`), INSERT (`col1`,`col2`) ON `mydb`.`mytbl` TO `someuser`@`somehost`
-- case
GRANT ALL ON `zabbix`.* TO `zabbix`@`localhost` IDENTIFIED BY 'password'
-- case
GRANT SELECT ON `test`.* TO `test`@`%`
-- case
GRANT PROCESS, USAGE, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO `xxxxxxxxxx`@`%` IDENTIFIED WITH 'mysql_native_password' AS 'xxxxxxxxxxxxxxxxxxxxxxxxxxxx'
-- case
GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, REFERENCES, RELOAD, PROCESS, INDEX, ALTER, CREATE TEMPORARY TABLES, LOCK TABLES, EXECUTE, REPLICATION SLAVE, REPLICATION CLIENT, CREATE VIEW, SHOW VIEW, CREATE ROUTINE, ALTER ROUTINE, CREATE USER, EVENT, TRIGGER ON *.* TO `root2`@`%` IDENTIFIED WITH 'mysql_native_password' AS '*sdsadsdsadssadsadsadsadsada' WITH GRANT OPTION
-- case
GRANT `role1`@`%`, `role2`@`%` TO `user1`@`localhost`, `user2`@`localhost`
-- case
GRANT `u1`@`%` TO `u1`@`%`
-- case
GRANT `app_read`@`%`, `app_write`@`%` TO `rw_user1`@`localhost`
-- case
GRANT `app_developer`@`%` TO `dev1`@`localhost`
-- case
GRANT SHUTDOWN ON *.* TO `dev1`@`localhost`
-- case
GRANT CONFIG ON *.* TO `dev1`@`localhost`
-- case
GRANT CREATE ON *.* TO `dev1`@`localhost`
-- case
GRANT CREATE TABLESPACE ON *.* TO `dev1`@`localhost`
-- case
GRANT EXECUTE ON FUNCTION `db1`.`anomaly_score` TO `user1`@`domain-or-ip-address1`
-- case
GRANT EXECUTE ON PROCEDURE `mydb`.`myproc` TO `someuser`@`somehost`
-- case
GRANT APPLICATION_PASSWORD_ADMIN, AUDIT_ADMIN ON *.* TO `root`@`localhost`
-- case
GRANT LOAD FROM S3, SELECT INTO S3, INVOKE LAMBDA, INVOKE SAGEMAKER, INVOKE COMPREHEND ON *.* TO `root`@`localhost`
-- case
GRANT PROXY ON `localuser`@`localhost` TO `externaluser`@`somehost`
-- case
GRANT PROXY ON ``@`` TO `root`@`localhost` WITH GRANT OPTION
-- case
GRANT PROXY ON `proxied_user`@`%` TO `proxy_user1`@`%`, `proxy_user2`@`%`
-- case
GRANT GRANT OPTION ON *.* TO `u1`@`%`
-- case
REVOKE ALL ON `db1`.* FROM `jeffrey`@`localhost`
-- case
REVOKE SELECT ON `db2`.`invoice` FROM `jeffrey`@`localhost`
-- case
REVOKE ALL ON *.* FROM `someuser`@`somehost`
-- case
REVOKE SELECT, INSERT ON *.* FROM `someuser`@`somehost`
-- case
REVOKE ALL ON `mydb`.* FROM `someuser`@`somehost`
-- case
REVOKE SELECT, INSERT ON `mydb`.* FROM `someuser`@`somehost`
-- case
REVOKE ALL ON `mydb`.`mytbl` FROM `someuser`@`somehost`
-- case
REVOKE SELECT, INSERT ON `mydb`.`mytbl` FROM `someuser`@`somehost`
-- case
REVOKE SELECT (`col1`), INSERT (`col1`,`col2`) ON `mydb`.`mytbl` FROM `someuser`@`somehost`
-- case
REVOKE ALL ON `zabbix`.* FROM `zabbix`@`localhost` IDENTIFIED BY 'password'
-- case
REVOKE `role1`@`%`, `role2`@`%` FROM `user1`@`localhost`, `user2`@`localhost`
-- case
REVOKE SHUTDOWN ON *.* FROM `dev1`@`localhost`
-- case
REVOKE CONFIG ON *.* FROM `dev1`@`localhost`
-- case
REVOKE EXECUTE ON FUNCTION `db`.`func` FROM `user`@`localhost`
-- case
REVOKE EXECUTE ON PROCEDURE `db`.`func` FROM `user`@`localhost`
-- case
REVOKE APPLICATION_PASSWORD_ADMIN, AUDIT_ADMIN ON *.* FROM `root`@`localhost`
-- case
REVOKE ALL, GRANT OPTION ON *.* FROM `u1`@`%`
-- case
REVOKE ALL, GRANT OPTION ON *.* FROM `u1`@`%`, `u2`@`%`, `u3`@`%`
