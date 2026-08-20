-- error: line 1 column 29 near "COMMENT 'c' SELECT 1" 
-- case
-- error: line 1 column 30 near "LANGUAGE SQL SELECT 1" 
-- case
-- error: line 1 column 25 near "NOT DETERMINISTIC SELECT 1" 
-- case
-- error: line 1 column 30 near "CONTAINS SQL SELECT 1" 
-- case
-- error: line 1 column 24 near "NO SQL SELECT 1" 
-- case
-- error: line 1 column 27 near "READS SQL DATA SELECT 1" 
-- case
-- error: line 1 column 30 near "MODIFIES SQL DATA SELECT 1" 
-- case
-- error: line 1 column 25 near "SQL SECURITY DEFINER SELECT 1" 
-- case
-- error: line 1 column 25 near "SQL SECURITY INVOKER SELECT 1" 
-- case
-- error: line 1 column 86 near "COMMENT 'c' LANGUAGE SQL NOT DETERMINISTIC CONTAINS SQL SQL SECURITY DEFINER SELECT 1" 
-- case
-- error: line 1 column 92 near "AS $$
  return a + b;
$$" 
-- case
-- error: line 1 column 53 near "LANGUAGE JAVASCRIPT AS $$
  console.log(msg);
$$" 
-- case
-- error: line 1 column 100 near "AS $body$
  return "got: " + s;
$body$" 
-- case
-- error: line 1 column 85 near "AS 'return a'" 
-- case
-- error: line 1 column 49 near "$$
  export function inc(n) { return n + 1; }
$$" 
-- case
-- error: line 1 column 68 near "$$
  export const PI = 3.14159;
$$" 
-- case
-- error: line 1 column 88 near "USING (lib_demo) AS $$
  return lib_demo.inc(a);
$$" 
-- case
-- error: line 1 column 34 near "LANGUAGE JAVASCRIPT USING (test.lib_demo AS d, other_lib) AS $$
  d.inc(1);
$$" 
