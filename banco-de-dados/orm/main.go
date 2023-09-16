package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	ID    int `gorm:"primaryKey"`
	Name  string
	Price float64
	gorm.Model
}

func main() {
	dsn := "root:root@tcp(localhost:3306)/nfsociety?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Product{})

	//create
	db.Create(&Product{
		Name:  "Teste",
		Price: 1000,
	})

	//create batch
	db.Create(&[]Product{
		Product{
			Name:  "Iphone",
			Price: 1000,
		},
		Product{
			Name:  "Samsung",
			Price: 2000,
		},
	})

	//select one
	var product Product
	db.First(&product, "name = ?", "samsung")
	db.First(&product, 2)
	fmt.Println(product)

	//read
	var products []Product
	db.Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}

	//limite offset
	db.Limit(2).Offset(1).Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}

	//where
	db.Where("name LIKE ?", "%iphone").Find(&products)
	for _, product := range products {
		fmt.Println(product)
	}

	//update
	var productToUpdate Product
	if err := db.First(&productToUpdate, 2).Error; err != nil {
		fmt.Printf("Erro ao buscar o produto: %v\n", err)
		return
	}

	productToUpdate.Price = 2999
	productToUpdate.Name = "Iphone"
	if err := db.Save(&productToUpdate).Error; err != nil {
		fmt.Printf("Erro ao salvar o produto: %v\n", err)
		return
	}

	fmt.Printf("Produto atualizado com sucesso: %v\n", productToUpdate)

	//delete
	var productToDelete Product
	if err := db.First(&productToDelete, 7).Error; err != nil {
		fmt.Printf("Erro ao buscar o produto: %v\n", err)
		return
	} else {
		if err := db.Delete(&productToDelete).Error; err != nil {
			fmt.Printf("Erro ao deletar o produto: %v\n", err)
			return
		}
		fmt.Printf("Produto deletado com sucesso: %v\n", productToDelete)
	}

}
