package service

import (
  "../model"
  "../repository"
  "../util"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
  "strconv"
)

var dbmap = repository.LoadInitialData()

func GetEvents(c *gin.Context) {
  var instructions []model.Event
  _, err := dbmap.Select(&instructions, "SELECT * FROM event")
  if err == nil {
    c.JSON(200, instructions)
  } else {
    c.JSON(404, gin.H{"error": "no event(s) into the table"})
  }
}

func GetEvent(c *gin.Context) {
  id := c.Params.ByName("id")
  var instruction model.Event
  
  err := dbmap.SelectOne(&instruction, "SELECT * FROM event WHERE id=?", id)
  if err == nil {
    instruction_id, _ := strconv.ParseInt(id, 0, 64)
  
    content := &model.Event{
      Id: instruction_id,
      EventStatus: instruction.EventStatus,
      EventName: instruction.EventName,
    }
 
    c.JSON(200, content)
  } else {
    c.JSON(404, gin.H{"error": "instruction not found"})
  }
}

func PostEvent(c *gin.Context) {
  var instruction model.Event
  _ = c.Bind(&instruction)

  if instruction.EventStatus != "" && instruction.EventName != "" {
    if insert, _ := dbmap.Exec(`INSERT INTO event (event_status, event_name) VALUES (?, ?)`, instruction.EventStatus, instruction.EventName); insert != nil {
      instruction_id, err := insert.LastInsertId()
      if err == nil {
        content := &model.Event{
          Id: instruction_id,
          EventStatus: instruction.EventStatus,
          EventName: instruction.EventName,
        }
        c.JSON(201, content)
      } else {
        util.CheckErr(err, "Insert failed")
      }
    }
  } else {
    c.JSON(422, gin.H{"error": "fields are empty"})
  }
}

func UpdateEvent(c *gin.Context) {
  id := c.Params.ByName("id")
  var instruction model.Event
  err := dbmap.SelectOne(&instruction, "SELECT * FROM event WHERE id=?", id)
  
  if err == nil {
    var json model.Event
    _ = c.Bind(&json)
    instruction_id, _ := strconv.ParseInt(id, 0, 64)
    instruction := model.Event{
      Id: instruction_id,
      EventStatus: json.EventStatus,
      EventName: json.EventName,
    }

    if instruction.EventStatus != "" && instruction.EventName != ""{
    _, err = dbmap.Update(&instruction)

      if err == nil {
        c.JSON(200, instruction)
      } else {
        util.CheckErr(err, "Updated failed")
      }
    } else {
      c.JSON(422, gin.H{"error": "fields are empty"})
    }
  } else {
    c.JSON(404, gin.H{"error": "instruction not found"})
  }
}

func DeleteEvent(c *gin.Context) {
  id := c.Params.ByName("id")
  var instruction model.Event
  err := dbmap.SelectOne(&instruction, "SELECT id FROM event WHERE id=?", id)
  
  if err == nil {
    _, err = dbmap.Delete(&instruction)
    
    if err == nil {
      c.JSON(200, gin.H{"id #" + id: " deleted"})
    } else {
      util.CheckErr(err, "Delete failed")
    }
  } else {
    c.JSON(404, gin.H{"error": "event not found"})
  }
}
