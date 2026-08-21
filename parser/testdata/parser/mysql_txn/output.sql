XA START 'xid1'
-- case
XA START 'xid1'
-- case
XA START 'gtrid', 'bqual'
-- case
XA START 'gtrid', 'bqual', 7
-- case
XA START 'xid1' JOIN
-- case
XA START 'xid1' RESUME
-- case
XA END 'xid1'
-- case
XA END 'gtrid', 'bqual', 7
-- case
XA END 'xid1' SUSPEND
-- case
XA END 'xid1' SUSPEND FOR MIGRATE
-- case
XA PREPARE 'xid1'
-- case
XA COMMIT 'xid1'
-- case
XA COMMIT 'xid1' ONE PHASE
-- case
XA ROLLBACK 'xid1'
-- case
XA RECOVER
-- case
XA RECOVER CONVERT XID
-- case
LOCK INSTANCE FOR BACKUP
-- case
UNLOCK INSTANCE
-- case
SELECT `xa`,`xid`,`one`,`phase` FROM `suspend` WHERE `migrate`=1
-- case
START TRANSACTION READ ONLY
-- case
START TRANSACTION
-- case
COMMIT
-- case
ROLLBACK RELEASE
-- case
ROLLBACK TO s1
-- case
LOCK TABLES `t1` AS `a1` WRITE
-- case
LOCK TABLES `t1` AS `a1` READ
-- case
LOCK TABLES `t2` LOW_PRIORITY WRITE
