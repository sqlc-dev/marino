SELECT ++1
-- case
-- error: line 1 column 9 near "*1" 
-- case
SELECT -+1
-- case
SELECT -1
-- case
SELECT --1
-- case
SELECT _UTF8MB4'''a''',_UTF8MB4'"a"'
-- case
-- error: line 1 column 12 near "''" 
-- case
-- error: line 1 column 12 near """" 
-- case
SELECT _UTF8MB4'''a'''
-- case
SELECT _UTF8MB4'''a'''
-- case
SELECT _UTF8MB4'"a"'
-- case
SELECT _UTF8MB4'"a"'
-- case
SELECT _UTF8'string'
-- case
SELECT _BINARY'string'
-- case
SELECT _UTF8'string'
-- case
SELECT _UTF8'string'
-- case
SELECT _UTF8 x'd0b1'
-- case
SELECT _UTF8 x'd0b1'
-- case
SELECT _UTF8 b'1101000010110001'
-- case
SELECT _UTF8 b'1101000010110001'
-- case
SELECT 1<=>0,1<=>NULL,1=NULL
-- case
SELECT DATE '1989-09-10'
-- case
-- error: line 1 column 20 near "19890910" 
-- case
SELECT TIME '00:00:00.111'
-- case
-- error: line 1 column 20 near "19890910" 
-- case
SELECT TIMESTAMP '1989-09-10 11:11:11'
-- case
-- error: line 1 column 25 near "19890910" 
-- case
SELECT TIMESTAMP '1989-09-10 11:11:11'
-- case
SELECT DATE '1989-09-10'
-- case
SELECT TIME '00:00:00.111'
-- case
SELECT * FROM `t` WHERE `a`>TIMESTAMP '1989-09-10 11:11:11'
-- case
SELECT * FROM `t` WHERE `a`>TIMESTAMP '1989-09-10 11:11:11'
-- case
SELECT _UTF8MB4'1989-09-10 11:11:11'
-- case
SELECT 123
-- case
SELECT 1 XOR 1
-- case
SELECT * FROM `t` WHERE `a`>_UTF8MB4'1989-09-10 11:11:11'
-- case
-- error: line 1 column 8 near ".t.a from t" 
