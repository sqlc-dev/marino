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
