CREATE STATISTICS `stats1` (CARDINALITY) ON `t`(`a`, `b`, `c`)
-- case
CREATE STATISTICS `stats2` (DEPENDENCY) ON `t`(`a`, `b`)
-- case
CREATE STATISTICS `stats3` (CORRELATION) ON `t`(`a`, `b`)
-- case
-- error: line 1 column 27 near "on t(a,b)" 
-- case
CREATE STATISTICS IF NOT EXISTS `stats1` (CARDINALITY) ON `t`(`a`, `b`, `c`)
-- case
CREATE STATISTICS IF NOT EXISTS `stats2` (DEPENDENCY) ON `t`(`a`, `b`)
-- case
CREATE STATISTICS IF NOT EXISTS `stats3` (CORRELATION) ON `t`(`a`, `b`)
-- case
-- error: line 1 column 41 near "on t(a,b)" 
-- case
CREATE STATISTICS `stats1` (CARDINALITY) ON `t`(`a`, `b`, `c`)
-- case
DROP STATISTICS `stats1`
