package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type TableData struct {
	Name        string
	RecordCount int64
	Size        int64
}

func main() {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.POST("/query", func(c *gin.Context) {
		address := c.PostForm("address")
		port := c.PostForm("port")
		username := c.PostForm("username")
		password := c.PostForm("password")

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, address, port))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		rows, err := db.Query("SELECT table_name, table_rows, data_length FROM information_schema.tables WHERE table_schema = 'mdm'")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var tableData []TableData

		for rows.Next() {
			var tableName string
			var recordCount int64
			var size int64

			err = rows.Scan(&tableName, &recordCount, &size)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			tableData = append(tableData, TableData{Name: tableName, RecordCount: recordCount, Size: size})
		}

		c.HTML(http.StatusOK, "result.html", gin.H{"tableData": tableData})
	})

	router.Run(":8080")
}
