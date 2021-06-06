package sqlite

import (
	"os"
	"testing"
)

func TestSqlite(t *testing.T) {
    var my SqliteClient
	my.Open("./test.db")

	cloums := []string{"ID INT PRIMARY KEY NOT NULL","NAME TEXT NOT NULL"}
	my.CreateTable("test",cloums)

	values := []string{"111","bbb"}
	my.Insert("test", values)

	aa := my.Select("test", "*", "")
	t.Log(aa)
	
	my.Delete("test", "ID=111")

	my.DeleteTable("test")

    os.Remove("./test.db")
}