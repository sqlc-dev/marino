-- flags: real_as_float
select cast(1 as float);
-- case
select cast(1 as float(0));
-- case
select cast(1 as float(24));
-- case
select cast(1 as float(25));
-- case
select cast(1 as float(53));
-- case
select cast(1 as float(54));
-- case
select cast(1 as real);
