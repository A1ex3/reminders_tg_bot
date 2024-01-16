package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	config "reminders_tg_got/config"
)

type Repository struct {
	SqLite *sql.DB
	Config *config.Config
}

func (repo *Repository) Open() {
	db, err := sql.Open("sqlite3", repo.Config.PathToDataBase)
	if err != nil {
		log.Println(err)
		return
	}
	repo.SqLite = db
}

func (repo *Repository) Close(){
	repo.SqLite.Close()
}