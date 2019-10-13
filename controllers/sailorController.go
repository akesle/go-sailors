package controllers

import (
  "context"
  "database/sql"
  "github.com/akesle/sailors/models"
  "github.com/gin-gonic/gin"
  "github.com/golang-sql/sqlexp"
  "log"
  "net/http"
  "strings"
)

type SailorController struct {
  DBSrc func() sqlexp.Querier
}

func (sc *SailorController) AddSailor(c *gin.Context) {

  var
  (
    sailor models.Sailor
    err    error
    db     sqlexp.Querier
    result sql.Result
  )

  if err = c.ShouldBind(&sailor); err == nil {
    log.Printf("%v", sailor)

    db = sc.DBSrc()

    result, err = db.ExecContext(context.Background(), "INSERT INTO Sailor (FirstName, LastName, Age) VALUES(?,?,?)",
      sailor.FirstName, sailor.LastName, sailor.Age)

    if err != nil {
      log.Printf("Failed to create sailor: %v", err)
      c.String(http.StatusInternalServerError, "Issue creating sailor; user not at fault")
      return
    }

    if count, rowsError := result.RowsAffected(); count != 1 || rowsError != nil {
      log.Printf("Failed to create sailor: %v", rowsError)
      c.String(http.StatusInternalServerError, "Issue creating sailor; user not at fault")
      return
    }

    c.String(http.StatusCreated, "Success")
    return
  }
  log.Printf("Error parsing input: %v", err)
  c.String(http.StatusBadRequest, "Failed to parse input")
}

func (sc *SailorController) FindSailor(c *gin.Context) {

  var
  (
    queriedSailor models.Sailor
    foundSailor   models.Sailor
    err           error
    db            sqlexp.Querier
    rows          *sql.Rows
  )

  if err = c.ShouldBind(&queriedSailor); err == nil {
    log.Printf("%v", queriedSailor)

    db = sc.DBSrc()

    handleQueryErr := func(err error) {
      log.Printf("Failed to query sailor: %v", err)
      c.String(http.StatusInternalServerError, "Issue querying sailor; user not at fault")
    }

    baseQuery := "SELECT FirstName, LastName, Age FROM Sailor WHERE "
    var criteria []string
    var args []interface{}
    if queriedSailor.LastName != "" {
      criteria = append(criteria, "LastName = ?")
      args = append(args, queriedSailor.LastName)
    }
    if queriedSailor.FirstName != "" {
      criteria = append(criteria, "FirstName = ?")
      args = append(args, queriedSailor.FirstName)
    }
    if queriedSailor.Age != 0 {
      criteria = append(criteria, "Age = ?")
      args = append(args, queriedSailor.Age)
    }
    query := baseQuery + strings.Join(criteria, " AND ")

    rows, err = db.QueryContext(context.Background(), query, args...)
    if err != nil {
      handleQueryErr(err)
      return
    }

    if err = rows.Err(); err != nil {
      handleQueryErr(err)
      return
    }
    defer func() {
      err = rows.Close()
      if err != nil {
        log.Printf("Failed closing rows for queried %v: %v", queriedSailor, err)
      }
    }()
    sailors := make([]models.Sailor, 0)
    for rows.Next() {
      err = rows.Scan(&foundSailor.FirstName, &foundSailor.LastName, &foundSailor.Age)
      if err != nil {
        log.Printf("Failed to scan sailor database fields: %v", err)
        c.String(http.StatusInternalServerError, "Issue querying sailor; user not at fault")
        return
      }
      sailors = append(sailors, foundSailor)
    }

    c.JSON(http.StatusOK, sailors)

    return
  }
  log.Printf("Error parsing input: %v", err)
  c.String(http.StatusBadRequest, "Failed to parse input")
}

func (sc *SailorController) RemoveSailor(c *gin.Context) {

  var
  (
    sailor models.AffectedSailor
    err    error
    db     sqlexp.Querier
    result sql.Result
  )

  if err = c.ShouldBind(&sailor); err == nil {
    log.Printf("%v", sailor)

    db = sc.DBSrc()

    result, err = db.ExecContext(context.Background(), "DELETE FROM Sailor WHERE LastName = ? AND FirstName = ? AND Age = ?",
      sailor.LastName, sailor.FirstName, sailor.Age)

    if err != nil {
      log.Printf("Failed to remove sailor: %v", err)
      c.String(http.StatusInternalServerError, "Issue removing sailor; user not at fault")
      return
    }

    if count, rowsError := result.RowsAffected(); count < 1 || rowsError != nil {
      log.Printf("Failed to remove any matching sailor: %v", rowsError)
    }

    c.String(http.StatusAccepted, "Success")
    return
  }
  log.Printf("Error parsing input: %v", err)
  c.String(http.StatusBadRequest, "Failed to parse input")
}

func (sc *SailorController) ModifySailor(c *gin.Context) {

  var
  (
    requested models.SailorUpdates
    err       error
    db        sqlexp.Querier
    result    sql.Result
  )

  if err = c.ShouldBind(&requested); err == nil {
    log.Printf("%v", requested)

    db = sc.DBSrc()

    baseQuery := "UPDATE Sailor SET "
    querySuffix := " WHERE LastName = ? AND FirstName = ? AND Age = ?"
    var updates []string
    var args []interface{}
    if requested.UpdatedFirstName != "" {
      updates = append(updates, "FirstName = ?")
      args = append(args, requested.UpdatedFirstName)
    }
    if requested.UpdatedLastName != "" {
      updates = append(updates, "LastName = ?")
      args = append(args, requested.UpdatedLastName)
    }
    if requested.UpdatedAge != 0 {
      updates = append(updates, "Age = ?")
      args = append(args, requested.UpdatedAge)
    }
    args = append(args,
      requested.AffectedSailor.LastName, requested.AffectedSailor.FirstName, requested.AffectedSailor.Age)
    query := baseQuery + strings.Join(updates, ", ") + querySuffix

    result, err = db.ExecContext(context.Background(), query, args...)

    if err != nil {
      log.Printf("Failed to update sailor: %v", err)
      c.String(http.StatusInternalServerError, "Issue updating sailor; user not at fault")
      return
    }

    if count, rowsError := result.RowsAffected(); count < 1 || rowsError != nil {
      log.Printf("Failed to update any matching sailor: %v", rowsError)
      c.String(http.StatusBadRequest, "Failed to find any matching sailors")
    }

    c.String(http.StatusOK, "Success")
    return
  }
  log.Printf("Error parsing input: %v", err)
  c.String(http.StatusBadRequest, "Failed to parse input")
}
