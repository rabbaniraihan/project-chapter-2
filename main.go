package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Book struct {
	Id     int    `gorm:"primaryKey" json:"id"`
	Title  string `gorm:"not null;unique;type:varchar(255)" json:"name_book"`
	Author string `gorm:"not null;type:varchar(255)" json:"author"`
	Desc   string `gorm:"not null;type:varchar(255)" json:"desc"`
}

var db *gorm.DB

func init() {
	var err error

	db, err = gorm.Open(postgres.Open("host=localhost port=5432 user=postgres password=rabbani11 dbname=golang-db sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(Book{})
}

func main() {
	g := gin.Default()

	g.GET("/book", getAllBook)
	g.POST("/book", addBook)
	g.DELETE("/book/:id", deleteBook)
	g.GET("/book/:id", getBookById)
	g.PUT("/book/:id", updateBook)

	g.Run(":8080")
}

func getAllBook(ctx *gin.Context) {
	var books []Book
	tx := db.Find(&books)
	if tx.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func addBook(ctx *gin.Context) {
	var newBook Book

	err := ctx.ShouldBindJSON(&newBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx := db.Create(&newBook)
	if tx.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	row := tx.Row()
	err = row.Scan(&newBook.Id, &newBook.Title, &newBook.Author, &newBook.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, newBook)
}

func deleteBook(ctx *gin.Context) {
	//Ambil id dari param
	stringId := ctx.Param("id")

	//Convert string -> int
	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var deletedBook Book
	deletedBook.Id = id

	tx := db.Clauses(clause.Returning{}).Delete(&deletedBook)
	if tx.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, "message : Book deleted sucesfully")
}

func getBookById(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var getBook Book
	getBook.Id = id

	tx := db.Find(&getBook)
	if tx.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	row := tx.Row()
	err = row.Scan(&getBook.Id, &getBook.Title, &getBook.Author, &getBook.Desc)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, getBook)
}

func updateBook(ctx *gin.Context) {
	stringId := ctx.Param("id")

	id, err := strconv.Atoi(stringId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var updatedBook Book

	err = ctx.ShouldBindJSON(&updatedBook)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	tx := db.Clauses(clause.Returning{
		Columns: []clause.Column{
			{
				Name: "id",
			},
		}}).Where("id = ?", id).Updates(&updatedBook)
	if tx.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	row := db.Row()
	err = row.Scan(&updatedBook.Id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, updatedBook)
}
