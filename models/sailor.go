package models

type Sailor struct {
  FirstName string `uri:"first_name" form:"first_name"`
  LastName  string `uri:"last_name" form:"last_name"`
  Age       int    `uri:"age" form:"age"`
}

type AffectedSailor struct {
  FirstName string `uri:"first_name" form:"first_name" binding:"required"`
  LastName  string `uri:"last_name" form:"last_name" binding:"required"`
  Age       int    `uri:"age" form:"age" binding:"required"`
}

type SailorUpdates struct {
  AffectedSailor
  UpdatedFirstName string `uri:"updated_first_name" form:"updated_first_name"`
  UpdatedLastName  string `uri:"updated_last_name" form:"updated_last_name"`
  UpdatedAge       int    `uri:"updated_age" form:"updated_age"`
}
