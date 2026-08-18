SELECT * FROM `t` USE INDEX (`primary`)
-- case
SELECT * FROM `t` USE INDEX (`primary`)
-- case
SELECT * FROM `t` USE INDEX ()
-- case
SELECT * FROM `t` USE INDEX (`idx`)
-- case
SELECT * FROM `t` USE INDEX (`idx1`, `idx2`)
-- case
SELECT * FROM `t` IGNORE INDEX (`idx1`)
-- case
SELECT * FROM `t` FORCE INDEX FOR JOIN (`idx1`)
-- case
SELECT * FROM `t` USE INDEX FOR ORDER BY (`idx1`)
-- case
SELECT * FROM `t` FORCE INDEX FOR GROUP BY (`idx1`)
-- case
SELECT * FROM (`t` USE INDEX FOR GROUP BY (`idx1`) USE INDEX FOR ORDER BY (`idx2`)) JOIN `t2`
