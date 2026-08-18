select "abc_" like "abc\\_" escape ''
-- case
select "abc_" like "abc\\_" escape '\\'
-- case
select "abc_" like "abc\\_" escape '||'
-- case
select "abc" like "escape" escape '+'
-- case
select '''_' like '''_' escape ''''
