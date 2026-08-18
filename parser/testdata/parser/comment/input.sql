create table t (c int comment 'comment')
-- case
create table t (c int) comment = 'comment'
-- case
create table t (c int) comment 'comment'
-- case
create table t (c int) comment comment
-- case
create table t (comment text)
-- case
START TRANSACTION /*!40108 WITH CONSISTENT SNAPSHOT */
-- case
/*comment*/ /*comment*/ select c /* this is a comment */ from t;
-- case
delete from t where a = 7 or 1=1/*' and b = 'p'
-- case
create table t (ssl int)
-- case
create table t (require int)
-- case
create table t (account int)
-- case
create table t (expire int)
-- case
create table t (cipher int)
-- case
create table t (issuer int)
-- case
create table t (never int)
-- case
create table t (subject int)
-- case
create table t (x509 int)
-- case
create user commentUser COMMENT '123456' '{"name": "Tom", "age", 19}
-- case
alter user commentUser COMMENT '123456' '{"name": "Tom", "age", 19}
-- case
create user commentUser COMMENT '123456'
-- case
alter user commentUser COMMENT '123456'
-- case
create user commentUser ATTRIBUTE '{"name": "Tom", "age", 19}'
-- case
alter user commentUser ATTRIBUTE '{"name": "Tom", "age", 19}'
