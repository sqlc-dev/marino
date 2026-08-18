change replication source to applier_version = 2
-- case
CHANGE REPLICATION SOURCE TO APPLIER_VERSION = 1, APPLIER_WORKER_COUNT = 64, APPLIER_EVENT_MEMORY_LIMIT = 1073741824 FOR CHANNEL 'channel_1'
-- case
change replication source to source_host = 'replica.example.com', source_port = 3306
-- case
change replication source to source_auto_position = 1 for channel 'group_replication_recovery'
-- case
change replication source to source_heartbeat_period = 60.5
-- case
change replication source to
-- case
change replication source applier_version = 2
-- case
change master to master_host = 'h'
-- case
change replication source to applier_version
-- case
change replication source to applier_version = 2,
-- case
change replication source to applier_version = 2 for channel
-- case
change replication source to for channel 'ch'
