package main

import "testing"

func Test_Execute(t *testing.T) {
	Execute()
	t.Log("one test passed.")
}

func TestGetTableNames(t *testing.T) {
	tableNames := GetTableNames("id_value")
	t.Log(tableNames)
}

func TestStart(t *testing.T) {
	sql := `
update id_value(id_code, current_value) values
('go_test1', 110),
('go_test2', 120),
('go_test3', 130),
('go_test4', 140);


replace into id_value(id_code, current_value) values
('go_test1', 110),
('go_test2', 120),
('go_test3', 130),
('go_test4', 140);

replace into id_value(id_code, current_value) values
('go_test1', 110),
('go_test5', 150),
('go_test3', 130),
('go_test4', 140);

select *  from id_value where id_code like 'go\_test%'
delete from id_value where id_code like 'go\_test%'

desc ent.id_value
`
	t.Log(sql)
}
