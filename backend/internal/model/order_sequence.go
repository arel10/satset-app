package model

// OrderSequence tracks the daily order number counter.
// Uses SELECT FOR UPDATE for concurrent safety.
type OrderSequence struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	Date         string `gorm:"uniqueIndex;size:10;not null" json:"date"` // YYYY-MM-DD
	LastSequence int    `gorm:"not null;default:0" json:"last_sequence"`
}

// TableName specifies the database table name for OrderSequence.
func (OrderSequence) TableName() string {
	return "order_sequences"
}
