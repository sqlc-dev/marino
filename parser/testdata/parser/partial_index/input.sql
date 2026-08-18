create table `t` (`id` int primary key,`col` int,index(`col`) where `col`>100)
-- case
create index `idx` on `t` (`col`) where `col`>100
-- case
alter table `t` add index `idx`(`col`) where `col`>100
