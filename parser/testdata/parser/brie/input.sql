BACKUP DATABASE a TO 'local:///tmp/archive01/'
-- case
BACKUP SCHEMA a TO 'local:///tmp/archive01/'
-- case
BACKUP DATABASE a,b,c TO 'noop://'
-- case
BACKUP DATABASE a.b TO 'noop://'
-- case
BACKUP DATABASE * TO 'noop://'
-- case
BACKUP DATABASE *, a TO 'noop://'
-- case
BACKUP DATABASE a, * TO 'noop://'
-- case
BACKUP DATABASE TO 'noop://'
-- case
BACKUP TABLE a TO 'noop://' checksum_concurrency 4 compression_level 4 ignore_stats 1 compression_type 'lz4'
-- case
RESTORE TABLE a FROM 'noop://' checksum_concurrency 4 wait_tiflash_ready 1 with_sys_table 1
-- case
BACKUP TABLE a.b TO 'noop://'
-- case
BACKUP TABLE a.b,c.d,e TO 'noop://'
-- case
BACKUP TABLE a.* TO 'noop://'
-- case
BACKUP TABLE * TO 'noop://'
-- case
BACKUP TABLE TO 'noop://'
-- case
RESTORE DATABASE * FROM 's3://bucket/path/'
-- case
BACKUP DATABASE * TO 'noop://' LAST_BACKUP = '2020-02-02 14:14:14'
-- case
BACKUP DATABASE * TO 'noop://' LAST_BACKUP = 1234567890
-- case
backup database * to 'noop://' rate_limit 500 MB/second snapshot 5 minute ago
-- case
backup database * to 'noop://' snapshot = '2020-03-18 18:13:54'
-- case
backup database * to 'noop://' snapshot = 1234567890
-- case
restore table g from 'noop://' concurrency 40 checksum 0 online 1
-- case
backup table x to 's3://bucket/path/?endpoint=https://test-cluster-s3.local&access-key=aaaaaaaaa&secret-access-key=bbbbbbbb&force-path-style=1'
-- case
backup database * to 's3://bucket/path/?provider=alibaba&region=us-west-9&storage-class=glacier&sse=AES256&acl=authenticated-read&use-accelerate-endpoint=1' send_credentials_to_tikv = 1
-- case
restore database * from 'gcs://bucket/path/?endpoint=https://test-cluster.gcs.local&storage-class=coldline&predefined-acl=OWNER&credentials-file=/data/private/creds.json'
-- case
restore table g from 'noop://' checksum off
-- case
restore table g from 'noop://' checksum optional
-- case
backup logs to 'noop://'
-- case
backup logs to 'noop://' start_ts='20220304'
-- case
pause backup logs
-- case
pause backup logs gc_ttl='20220304'
-- case
resume backup logs
-- case
show backup logs status
-- case
show backup logs metadata from 'noop://'
-- case
show br job 1234
-- case
show br job query 1234
-- case
cancel br job 1234
-- case
purge backup logs from 'noop://'
-- case
purge backup logs from 'noop://' until_ts='2012122304'
-- case
restore point from 'noop://log_backup'
-- case
restore point from 'noop://log_backup' full_backup_storage='noop://full_log'
-- case
restore point from 'noop://log_backup' full_backup_storage='noop://full_log' restored_ts='20230123'
-- case
restore point from 'noop://log_backup' full_backup_storage='noop://full_log' start_ts='20230101' restored_ts='20230123'
