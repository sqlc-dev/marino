CHANGE REPLICATION SOURCE TO APPLIER_VERSION = 2
-- case
CHANGE REPLICATION SOURCE TO APPLIER_VERSION = 1, APPLIER_WORKER_COUNT = 64, APPLIER_EVENT_MEMORY_LIMIT = 1073741824 FOR CHANNEL 'channel_1'
-- case
CHANGE REPLICATION SOURCE TO SOURCE_HOST = 'replica.example.com', SOURCE_PORT = 3306
-- case
CHANGE REPLICATION SOURCE TO SOURCE_AUTO_POSITION = 1 FOR CHANNEL 'group_replication_recovery'
-- case
CHANGE REPLICATION SOURCE TO SOURCE_HEARTBEAT_PERIOD = 60.5
-- case
-- error: line 1 column 28 near "" 
-- case
-- error: line 1 column 41 near "applier_version = 2" 
-- case
-- error: line 1 column 13 near "master to master_host = 'h'" 
-- case
-- error: line 1 column 44 near "" 
-- case
-- error: line 1 column 49 near "" 
-- case
-- error: line 1 column 60 near "" 
-- case
-- error: line 1 column 32 near "for channel 'ch'" 
