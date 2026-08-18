BACKUP DATABASE `a` TO 'local:///tmp/archive01/'
-- case
BACKUP DATABASE `a` TO 'local:///tmp/archive01/'
-- case
BACKUP DATABASE `a`, `b`, `c` TO 'noop://'
-- case
-- error: line 1 column 18 near ".b TO 'noop://'" 
-- case
BACKUP DATABASE * TO 'noop://'
-- case
-- error: line 1 column 18 near ", a TO 'noop://'" 
-- case
-- error: line 1 column 20 near "* TO 'noop://'" 
-- case
-- error: line 1 column 18 near "TO 'noop://'" 
-- case
BACKUP TABLE `a` TO 'noop://' CHECKSUM_CONCURRENCY = 4 COMPRESSION_LEVEL = 4 IGNORE_STATS = 1 COMPRESSION_TYPE = 'lz4'
-- case
RESTORE TABLE `a` FROM 'noop://' CHECKSUM_CONCURRENCY = 4 WAIT_TIFLASH_READY = 1 WITH_SYS_TABLE = 1
-- case
BACKUP TABLE `a`.`b` TO 'noop://'
-- case
BACKUP TABLE `a`.`b`, `c`.`d`, `e` TO 'noop://'
-- case
-- error: line 1 column 16 near "* TO 'noop://'" 
-- case
-- error: line 1 column 17 near "TO 'noop://'" 
-- case
-- error: line 1 column 15 near "TO 'noop://'" 
-- case
RESTORE DATABASE * FROM 's3://bucket/path/'
-- case
BACKUP DATABASE * TO 'noop://' LAST_BACKUP = '2020-02-02 14:14:14'
-- case
BACKUP DATABASE * TO 'noop://' LAST_BACKUP = 1234567890
-- case
BACKUP DATABASE * TO 'noop://' RATE_LIMIT = 500 MB/SECOND SNAPSHOT = 300000000 MICROSECOND AGO
-- case
BACKUP DATABASE * TO 'noop://' SNAPSHOT = '2020-03-18 18:13:54'
-- case
BACKUP DATABASE * TO 'noop://' SNAPSHOT = 1234567890
-- case
RESTORE TABLE `g` FROM 'noop://' CONCURRENCY = 40 CHECKSUM = OFF ONLINE = 1
-- case
BACKUP TABLE `x` TO 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret-access-key=bbbbbbbb&force-path-style=1'
-- case
BACKUP DATABASE * TO 's3://bucket/path/?provider=alibaba&region=us-west-9&storage-class=glacier&sse=AES256&acl=authenticated-read&use-accelerate-endpoint=1' SEND_CREDENTIALS_TO_TIKV = 1
-- case
RESTORE DATABASE * FROM 'gcs://bucket/path/?endpoint=https://test-cluster.gcs.local&storage-class=coldline&predefined-acl=OWNER&credentials-file=/data/private/creds.json'
-- case
RESTORE TABLE `g` FROM 'noop://' CHECKSUM = OFF
-- case
RESTORE TABLE `g` FROM 'noop://' CHECKSUM = OPTIONAL
-- case
BACKUP LOGS TO 'noop://'
-- case
BACKUP LOGS TO 'noop://' START_TS = '20220304'
-- case
PAUSE BACKUP LOGS
-- case
PAUSE BACKUP LOGS GC_TTL = '20220304'
-- case
RESUME BACKUP LOGS
-- case
SHOW BACKUP LOGS STATUS
-- case
SHOW BACKUP LOGS METADATA FROM 'noop://'
-- case
SHOW BR JOB 1234
-- case
SHOW BR JOB QUERY 1234
-- case
CANCEL BR JOB 1234
-- case
PURGE BACKUP LOGS FROM 'noop://'
-- case
PURGE BACKUP LOGS FROM 'noop://' UNTIL_TS = '2012122304'
-- case
RESTORE POINT FROM 'noop://log_backup'
-- case
RESTORE POINT FROM 'noop://log_backup' FULL_BACKUP_STORAGE = 'noop://full_log'
-- case
RESTORE POINT FROM 'noop://log_backup' FULL_BACKUP_STORAGE = 'noop://full_log' RESTORED_TS = '20230123'
-- case
RESTORE POINT FROM 'noop://log_backup' FULL_BACKUP_STORAGE = 'noop://full_log' START_TS = '20230101' RESTORED_TS = '20230123'
