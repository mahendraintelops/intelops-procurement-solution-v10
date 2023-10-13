package models

type Invoice struct {
	Id int64 `json:"id,omitempty"`

	Amount float32 `json:"amount,omitempty"`

	InvoiceDate string `json:"invoiceDate,omitempty"`

	Items string `json:"items,omitempty"`

	PaymentTerms string `json:"paymentTerms,omitempty"`
}
