package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type BannerBook struct {
	Name         string `json:"name" gorm:"column:book_name"`
	ActionType   int    `json:"actionType gorm:"column:action_type"`
	BannerUrl    string `json:"bannerUrl" gorm:"column:url"`
	Introduction string `json:"introduction" gorm:"column:introduction"`
	BookId       int    `json:"bookId" gorm:"column:book_id"`
}

func main() {
	fmt.Print("数据库开始连接\n")
	db, err := gorm.Open(
		"mysql",
		"root:peng1365451772@/novel?charset=utf8&parseTime=True&loc=Local",
	)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	var banner []BannerBook
	db.Table("read_ws_banner").Select("book_name,action_type,url,introduction,book_id").Where("app_id = 'com.fantasy.best.novel'").Order("sort").Scan(&banner)
	fmt.Println(db)
	fmt.Println(banner)
	var names []string
	db.Table("read_ws_banner").Pluck("book_name", &names)

	fmt.Println(names)

}
