package main

import (
	sb "github.com/samtech09/gosql"
)

func main() {
	fw := sb.NewFileWriter(5)
    
    os.Setenv("SQL_PARAM_FORMAT", ParamPostgreSQL)

	stmt := sb.InsertBuilder().Table("dbusers").
		Columns("id", "name").Returning("id").
		Build(true)
	fw.Queue(stmt, "user", "Create", "Creates new user.")

	stmt = sb.UpdateBuilder().Table("dbusers").
		Columns("name").
		Where(sb.C().EQ("id", "?")).
		Build(true)
	fw.Queue(stmt, "user", "Update", "Update existing user by id.")

	stmt = sb.DeleteBuilder().Table("dbusers").
		Where(sb.C().EQ("id", "?")).
		Build(true)
	fw.Queue(stmt, "user", "Delete", "Delete a user by id.")

	stmt = sb.SelectBuilder().Select("id", "name").
		From("dbusers", "").
		RowCount().
		Build(true)
	fw.Queue(stmt, "user", "ListAll", "List all users.")

	fw.Write("../", "sqlbuilder", "sqls", sb.WriteJSONandJSONLoaderGoCode)
}
