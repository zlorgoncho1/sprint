// finance.go
package models

import (
	"errors"
	"fmt"
	"log"
)

// AddIncome adds an income record to the database
func AddIncome(amount float64, category string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if category == "" {
		return errors.New("category cannot be empty")
	}

	stmt, err := DB.Prepare("INSERT INTO income(category, amount) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category, amount)
	if err != nil {
		return err
	}

	fmt.Println("Income added successfully")
	return nil
}

// AddExpense adds an expense record to the database
func AddExpense(amount float64, category string) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}
	if category == "" {
		return errors.New("category cannot be empty")
	}

	stmt, err := DB.Prepare("INSERT INTO expense(category, amount) VALUES(?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(category, amount)
	if err != nil {
		return err
	}

	fmt.Println("Expense added successfully")
	return nil
}

// GetSummary retrieves a summary of income and expenses by category
func GetSummary() map[string]interface{} {
	summary := make(map[string]interface{})
	incomes := make(map[string]float64)
	expenses := make(map[string]float64)

	// Retrieve income data
	rows, err := DB.Query("SELECT category, SUM(amount) as total_amount FROM income GROUP BY category")
	if err != nil {
		log.Println("Error fetching income data: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var totalAmount float64
		err := rows.Scan(&category, &totalAmount)
		if err != nil {
			log.Println("Error scanning income row: ", err)
			continue
		}
		incomes[category] = totalAmount
	}

	// Retrieve expense data
	rows, err = DB.Query("SELECT category, SUM(amount) as total_amount FROM expense GROUP BY category")
	if err != nil {
		log.Println("Error fetching expense data: ", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var totalAmount float64
		err := rows.Scan(&category, &totalAmount)
		if err != nil {
			log.Println("Error scanning expense row: ", err)
			continue
		}
		expenses[category] = totalAmount
	}

	summary["incomes"] = incomes
	summary["expenses"] = expenses

	return summary
}
