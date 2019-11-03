package repository

import (
  "../model"
  "../util"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "gopkg.in/gorp.v1"
  "log"
  "os"
)

func setupDatabase() *gorp.DbMap {

  log.Print("user", os.Getenv("MYSQL_USER"))
  log.Print("mysql database", os.Getenv("MYSQL_DATABASE"))
  connection_string := "event" + ":" + "event_password" + "@/" + "events"
  db, err := sql.Open("mysql", connection_string)
  util.CheckErr(err, "sql.Open failed")

  dbmap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

  return dbmap
}

func createTable() *gorp.DbMap {

  database := setupDatabase()
  database.AddTableWithName(model.Event{}, "event").SetKeys(true, "Id")

  err := database.CreateTablesIfNotExists()
  util.CheckErr(err, "Create table failed")
  return database
}

func LoadInitialData() *gorp.DbMap {

  var db = createTable()
  err := db.TruncateTables()

  if err == nil {
    log.Print("Table truncated")
  }

  var events []model.Event
  messageReceived := model.Event{1, "enabled", "MESSAGE_RECEIVED"}
  messageLost := model.Event{2, "enabled", "MESSAGE_LOST"}
  events = append(events, messageReceived)
  events = append(events, messageLost)

  for _, dataInstruction := range events {
    result, _ := db.Exec(`INSERT INTO event (event_status, event_name) VALUES (?, ?)`, dataInstruction.EventStatus, dataInstruction.EventName)
    if result != nil {
      preInsertedId, err := result.LastInsertId()
      if err == nil {
        log.Print("Succesfully pre inserted Data ", preInsertedId)
      }
    }
  }

  return db;

}
