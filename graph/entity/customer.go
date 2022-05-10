package entity

import (
	"database/sql"
	uuid "github.com/gofrs/uuid"
)

type CustomerContact struct {
	ID                      *uuid.UUID
	ContactName             string
	Designation             string
	ContactNumber           string
	ContactImage            *string
	CustomerID              uuid.UUID
	EmailID                 *string
	AdditionalContactNumber string
}

type CustomerContactDelete struct {
	ID uuid.UUID
}

type CustomerContactData struct {
	ID                      sql.NullString `json:"id"`
	ContactName             sql.NullString `json:"contactName"`
	Designation             sql.NullString `json:"designation"`
	ContactNumber           sql.NullString `json:"contactNumber"`
	CustomerID              sql.NullString `json:"customerId"`
	ContactImage            sql.NullString `json:"ContactImage"`
	CustomerName            sql.NullString `json:"customerName"`
	HasConsent              sql.NullBool   `json:"hasConsent"`
	EmailID                 sql.NullString `json:"emailId"`
	AdditionalContactNumber sql.NullString `json:"AdditionalContactNumber"`
}
