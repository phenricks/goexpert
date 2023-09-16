package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Category struct {
	ID       int `gorm:"primaryKey"`
	Name     string
	Products []Product `gorm:"many2Many:products_categories;"`
}

type Product struct {
	ID           int `gorm:"primaryKey"`
	Name         string
	Price        float64
	Categories   []Category `gorm:"many2Many:products_categories;"`
	SerialNumber SerialNumber
	gorm.Model
}

type SerialNumber struct {
	ID        int `gorm:"primaryKey"`
	Number    string
	ProductID int
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/nfsociety?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{}, &Category{}, &SerialNumber{})

	//create category
	category := Category{Name: "Cozinha"}
	db.Create(&category)

	category2 := Category{Name: "Eletronicos"}
	db.Create(&category2)

	//create product
	db.Create(&Product{
		Name:       "Panela",
		Price:      99.90,
		Categories: []Category{category, category2},
	})

	db.Create(&SerialNumber{
		Number:    "123456789",
		ProductID: 1,
	})

	// belongs to
	var products []Product
	db.Preload("Category").Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}

	//has one
	var product Product
	db.Preload("Category").Preload("SerialNumber").First(&product)
	fmt.Println(product)

	//has many
	var categories []Category
	err = db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		fmt.Println(category.Name, ":")
		for _, product := range category.Products {
			fmt.Println("-", product.Name, " - Serial Number:", product.SerialNumber.Number)
		}
	}

	//many to many
	err = db.Model(&Category{}).Preload("Products.SerialNumber").Find(&categories).Error
	if err != nil {
		panic(err)
	}
	for _, category := range categories {
		fmt.Printf("%s:\n", category.Name)
		for _, product := range category.Products {
			if product.SerialNumber.Number != "" {
				fmt.Printf(" - Serial Number: %s", product.SerialNumber.Number)
			}
			fmt.Printf("- %s\n", product.Name)
		}
	}

	//lock pessimista ele locka a tabela, a linha do bd que você esta trabalhando. Exemplo abaixo:
	tx := db.Begin()
	var categoryLock Category
	err = tx.Debug().Clauses(clause.Locking{Strength: "UPDATE"}).First(&categoryLock, 1).Error
	if err != nil {
		panic(err)
	}
	category.Name = "Eletronicos"
	tx.Debug().Save(categoryLock)
	tx.Commit()

	//lock otimista versiona quando algo faz alteração no sistema

}
