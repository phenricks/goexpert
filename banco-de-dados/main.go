package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

// dev em go sempre tentam utilizar ao maximo o que ja vem com a linguagem. Maioria dos devs não utilizam ORM

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/nfsociety")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//create product
	product := NewProduct("SmartWatch", 299.99)
	err = createProduct(db, product)
	if err != nil {
		panic(err)
	}

	//update product
	toUpdate := &Product{
		ID:    "6a209aed-a7d5-4500-8dae-2da777a0975a",
		Price: 2999.99,
	}

	updated, err := updateProduct(db, toUpdate)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)

	//getByID
	found, err := findProductById(db, "6a209aed-a7d5-4500-8dae-2da777a0975a")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Product: %v, possui o preço de %.2f\n", found.Name, found.Price)

	//listAll
	products, err := listAll(db)
	if err != nil {
		panic(err)
	}

	for _, product := range products {
		fmt.Printf("Product: %v - Price: %.2f\n", product.Name, product.Price)
	}

	//delete
	err = deleteProduct(db, "16f02318-6bbd-43ed-84ad-f48f91781357")
	if err != nil {
		panic(err)
	}
}

func createProduct(db *sql.DB, product *Product) error {

	stmt, err := db.Prepare("INSERT INTO products (id, name, price) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(product.ID, product.Name, product.Price)
	if err != nil {
		return err
	}

	return nil
}

func updateProduct(db *sql.DB, product *Product) (*Product, error) {

	query := "UPDATE products SET "
	params := []interface{}{}

	if product.Name != "" {
		query += "name = ?, "
		params = append(params, product.Name)
	}
	if product.Price != 0 {
		query += "price = ?, "
		params = append(params, product.Price)
	}

	// Remova a vírgula extra no final da consulta.
	query = query[:len(query)-2]
	query += " WHERE id = ?"
	params = append(params, product.ID)

	// Prepare a consulta SQL.
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	// Execute a consulta SQL preparada.
	_, err = stmt.Exec(params...)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func findProductById(db *sql.DB, id string) (*Product, error) {

	query := "SELECT id, name, price FROM products WHERE id =?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var product Product
	// err = stmt.QueryRowContext(ctx, id).Scan(&product.ID, &product.Name, &product.Price) // <- caso eu precise de contexto
	err = stmt.QueryRow(id).Scan(&product.ID, &product.Name, &product.Price) // scan olha os valores de cada coluna e atribui em cada atributo do product
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func listAll(db *sql.DB) ([]Product, error) {

	query := "SELECT id, name, price FROM products"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err = rows.Scan(&product.ID, &product.Name, &product.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func deleteProduct(db *sql.DB, id string) error {
	query := "DELETE FROM products WHERE id =?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}
