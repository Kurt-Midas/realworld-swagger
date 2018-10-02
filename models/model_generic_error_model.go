package models

type GenericErrorModel struct {
	Errors *GenericErrorModelErrors `json:"errors"`
}
