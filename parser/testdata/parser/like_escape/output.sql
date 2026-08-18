SELECT _UTF8MB4'abc_' LIKE _UTF8MB4'abc\_' ESCAPE ''
-- case
SELECT _UTF8MB4'abc_' LIKE _UTF8MB4'abc\_'
-- case
-- error: [parser:1210]Incorrect arguments to ESCAPE
-- case
SELECT _UTF8MB4'abc' LIKE _UTF8MB4'escape' ESCAPE '+'
-- case
SELECT _UTF8MB4'''_' LIKE _UTF8MB4'''_' ESCAPE ''''
