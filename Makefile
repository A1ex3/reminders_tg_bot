PATH-TO-FILE-SQL = schema.sql
PATH-TO-FILE-DATABASE = remindersTgBot.db

makedb:
	cat $(PATH-TO-FILE-SQL) | sqlite3 $(PATH-TO-FILE-DATABASE)

build:
	go mod tidy
	go build .