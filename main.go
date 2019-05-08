package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robertkrimen/otto"
	"net/http"
)

const (
	DB_Driver = "b2b_5:123456@tcp(10.3.5.30:3306)/b2b_third?charset=utf8"
)

var scriptLoader *ScriptLoader

func main() {
	scriptLoader = New(".")
	scriptLoader.Init()
	//
	//router := gin.Default()
	//router.GET("/ok", ok)
	//router.GET("/domain", domain)
	//router.GET("/db", openDB)
	//router.GET("/js", js)
	//router.GET("/reload", reload)
	//
	//_ = router.Run(":3000")
	//fmt.Println("start finished")
	Start()
}

func reload(c *gin.Context) {
	value, _ := scriptLoader.Load("/conf.js")
	fv, _ := value.Call(otto.Value{})
	vv, _ := fv.Export()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vv,
	})
}

func js(c *gin.Context) {
	vm := otto.New()

	b, err := ReadAll("1.js")
	checkErr(err)

	js := string(b)
	r, re := vm.Run(js)
	checkErr(re)
	fmt.Println(r.ToInteger())

	value, error := vm.Get("abc")
	checkErr(error)

	v, _ := value.ToInteger()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    v,
	})
	return
}

func ok(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func openDB(c *gin.Context) {
	db, err := sql.Open("mysql", DB_Driver)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		checkErr(err)
		return
	}

	stmt, err := db.Prepare(`SELECT
activity_id, 
activity_name, 
activity_content, 
create_by,
create_at
FROM tb_activity where site_id =?`)
	checkErr(err)
	rows, err := stmt.Query("1111111112")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		checkErr(err)
		return
	}

	var l []map[string]string
	for rows.Next() {
		var activity_id string
		var activity_name string
		var activity_content string
		var create_by string
		var create_at string
		vv := map[string]string{}
		err := rows.Scan(&vv)
		checkErr(err)

		l = append(l, map[string]string{
			"activity_id":      activity_id,
			"activity_name":    activity_name,
			"activity_content": activity_content,
			"create_by":        create_by,
			"create_at":        create_at,
		})
	}
	defer db.Close()
	defer stmt.Close()
	defer rows.Close()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "获取成功",
		"data":    l,
	})
}

func domain(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "domain",
	})
}
