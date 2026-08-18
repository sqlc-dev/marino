LOCK INSTANCE FOR BACKUP
-- case
UNLOCK INSTANCE
-- case
XA START 'xid1'
-- case
XA END 'xid1'
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
