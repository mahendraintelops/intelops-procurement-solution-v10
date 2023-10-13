package models

type Payment struct {
	Id int64 `json:"id,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	InvoiceNumber string `json:"invoiceNumber,omitempty"`
}
