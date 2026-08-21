CREATE USER `u`@`%` IDENTIFIED BY RANDOM PASSWORD
-- case
CREATE USER `u`@`%` IDENTIFIED WITH 'caching_sha2_password' BY RANDOM PASSWORD
-- case
ALTER USER `u`@`%` IDENTIFIED BY RANDOM PASSWORD
-- case
ALTER USER `u`@`%` IDENTIFIED BY RANDOM PASSWORD RETAIN CURRENT PASSWORD
-- case
CREATE USER `u`@`%` IDENTIFIED BY 'p1' AND IDENTIFIED WITH 'authentication_ldap_simple'
-- case
CREATE USER `u`@`%` IDENTIFIED BY 'p1' AND IDENTIFIED WITH 'authentication_fido' AND IDENTIFIED WITH 'authentication_ldap_simple' BY 'p3'
-- case
ALTER USER `u`@`%` ADD 2 FACTOR IDENTIFIED WITH 'authentication_ldap_simple'
-- case
ALTER USER `u`@`%` MODIFY 2 FACTOR IDENTIFIED WITH 'authentication_ldap_simple' BY 'p'
-- case
ALTER USER `u`@`%` DROP 2 FACTOR
-- case
ALTER USER `u`@`%` 2 FACTOR INITIATE REGISTRATION
-- case
ALTER USER `u`@`%` 2 FACTOR FINISH REGISTRATION SET CHALLENGE_RESPONSE AS 'blob'
-- case
ALTER USER `u`@`%` 2 FACTOR UNREGISTER
-- case
CREATE USER `u`@`%` DEFAULT ROLE `r1`@`%`, `r2`@`%`
-- case
CREATE USER `u`@`%` PASSWORD REQUIRE CURRENT
-- case
CREATE USER `u`@`%` PASSWORD REQUIRE CURRENT OPTIONAL
-- case
ALTER USER `u`@`%` IDENTIFIED BY 'p' REPLACE 'old'
-- case
ALTER USER `u`@`%` IDENTIFIED BY 'p' RETAIN CURRENT PASSWORD
-- case
ALTER USER `u`@`%` IDENTIFIED BY 'p' REPLACE 'old' RETAIN CURRENT PASSWORD
-- case
ALTER USER `u`@`%` DISCARD OLD PASSWORD
-- case
ALTER USER `u`@`%` DEFAULT ROLE ALL
-- case
ALTER USER `u`@`%` DEFAULT ROLE NONE
-- case
ALTER USER `u`@`%` DEFAULT ROLE `r1`@`%`, `r2`@`%`
-- case
GRANT `r1`@`%` TO `u`@`%` WITH ADMIN OPTION
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost`
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost` WITH ROLE DEFAULT
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost` WITH ROLE NONE
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost` WITH ROLE ALL
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost` WITH ROLE ALL EXCEPT `r1`@`%`, `r2`@`%`
-- case
GRANT SELECT ON *.* TO `u`@`%` AS `root`@`localhost` WITH ROLE `r1`@`%`, `r2`@`%`
-- case
REVOKE IF EXISTS SELECT ON `db`.`t` FROM `u`@`%`
-- case
REVOKE SELECT ON `db`.`t` FROM `u`@`%` IGNORE UNKNOWN USER
-- case
REVOKE IF EXISTS SELECT ON `db`.`t` FROM `u`@`%` IGNORE UNKNOWN USER
-- case
REVOKE PROXY ON `pu`@`%` FROM `u`@`%`
-- case
SET PASSWORD='p' REPLACE 'old'
-- case
SET PASSWORD FOR `u`@`%`='p' RETAIN CURRENT PASSWORD
-- case
SET PASSWORD TO RANDOM
-- case
SET PASSWORD FOR `u`@`%` TO RANDOM
-- case
SET PASSWORD FOR `u`@`%` TO RANDOM REPLACE 'old'
