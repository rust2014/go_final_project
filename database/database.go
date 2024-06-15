package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3" // импорт драйвера SQLite
)

func ConnectDB() *sql.DB {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err) // лоогируем ошибку
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile) // проверка существования файла

	var install bool
	if err != nil {
		install = true // если файл базы данных не существует, установить флаг установки
	} else if os.IsNotExist(err) {
		log.Fatal(err)
	}
	db, err := sql.Open("sqlite3", dbFile) // подключение к БД
	if err != nil {
		log.Fatal(err)
	}
	if install {
		createTableAndIndex(db) // создание таблицы и индекса если надо
	}
	return db
}

func createTableAndIndex(db *sql.DB) { // функция для создания таблицы и индекса
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date CHAR(8) NOT NULL DEFAULT '',
		title TEXT NOT NULL DEFAULT '',
		comment TEXT NOT NULL DEFAULT '',
		repeat TEXT NOT NULL DEFAULT ''
	);
	CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`
	_, err := db.Exec(sqlStmt) // выполнение запроса sql
	if err != nil {
		return
	}
}
