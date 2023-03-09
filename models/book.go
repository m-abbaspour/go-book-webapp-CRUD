package models

import ( 
	"fmt" // 1

	"github.com/jinzhu/gorm" // 3
	_ "github.com/jinzhu/gorm/dialects/mysql" // 5
)

type Book struct { // 6
	gorm.Model
	Title string
	Author string
	Publication string
}

var db *gorm.DB // 2

func init() { // 4
	d, err := gorm.Open("mysql", "abbas:abbas@/book?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db = d
	fmt.Println("Connected to DB")
	db.AutoMigrate(&Book{}) // Any change in struct will Migrate to the DB
}

func GetAllBooks() []Book { // 7 // Got to webapp.go and write interactHandler and rerun and check DB for table
	var books []Book
	db.Find(&books)

	return books
}

func GetBook(id int64) (Book,int) {
	var book Book
	resault:=db.Where("ID=?", id).Find(&book)
	return book,int(resault.RowsAffected)
}

func DeleteBook(id int64) Book {
	var book Book
	db.Where("ID=?", id).Delete(book)

	return book
}

func (b *Book) CreateBook () (*Book){
	db.NewRecord(b)
	db.Create(&b)
	return b
}

func (b *Book) UpdateBook() (*Book){
	db.Save(&b)
	return b
}