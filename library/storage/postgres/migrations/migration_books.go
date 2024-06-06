package migrations

import (
	. "practice/migrator"
)

func addBooksMigrations(mg *Migrator) {
	books := Table{
		Name: "books",
		Columns: []*Column{
			{Name: "id", Type: DB_BigInt, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "name", Type: DB_NVarchar, Length: 100, Nullable: false},
			{Name: "created_at", Type: DB_DateTime, Nullable: false},
			{Name: "updated_at", Type: DB_DateTime, Nullable: true},
			{Name: "created_by", Type: DB_NVarchar, Nullable: false},
			{Name: "updated_by", Type: DB_NVarchar, Nullable: true},
		},
	}

	mg.AddMigration("create books table", NewAddTableMigration(books))
}
