package main

import (
	"gorm.io/gorm"
)

type CustomerAccount struct {
	ID   uint
	Name string
}

// Vulnerable: gorm Delete without primary key or Where clause can wipe entire table
func deleteCustomerUnscoped(db *gorm.DB, account CustomerAccount) error {
	// If account.ID is 0/empty, gorm Delete(&account) generates "DELETE FROM customer_accounts" wiping all rows!
	return db.Delete(&account).Error
}
