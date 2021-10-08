package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"sync"
	"time"
)

var db *gorm.DB

const (
	BookTableName   = "fiction_books"
	writerTableName = "fiction_book_writers"
	tagTableName    = "fiction_book_tags"
)

type BannerBook struct {
	Name         string `json:"name" gorm:"column:book_name"`
	ActionType   int    `json:"actionType gorm:"column:action_type"`
	BannerUrl    string `json:"bannerUrl" gorm:"column:url"`
	Introduction string `json:"introduction" gorm:"column:introduction"`
	BookId       int    `json:"bookId" gorm:"column:book_id"`
}

type BookItem struct {
	BookId                uint      `json:"bookId"`
	BookName              string    `json:"bookName"`
	Author                string    `json:"author"`
	Introduction          string    `json:"introduction"`
	Ratings               float32   `json:"ratings"`
	Cover                 string    `json:"cover"`
	Tags                  []string  `json:"tags"`
	ViewCount             int       `json:"viewCount"`
	WriteStatus           string    `json:"writeStatus"`
	LastUpdateTime        time.Time `json:"lastUpdateTime"`
	TypeOneNames          []string  `json:"typeOneNames"`
	TypeTwoNames          []string  `json:"typeTwoNames"`
	ViewCountDisplay      string    `json:"viewCountDisplay"`
	LastUpdateTimeDisplay string    `json:"lastUpdateTimeDisplay"`
}
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FictionBook struct {
	Model
	BookId            uint    `gorm:"column:book_id"`
	BookName          string  `gorm:"column:book_name"`
	ChapterCount      int     `gorm:"column:chapter_count"`
	Introduction      string  `gorm:"column:introduction"`
	CoverUrl          string  `gorm:"column:cover_url"`
	Labels            string  `gorm:"column:labels"`
	Language          string  `gorm:"column:language"`
	PseudonymWriterId int     `gorm:"column:pseudonym_writer_id"`
	ToTalWords        int     `gorm:"column:total_words"`
	Ratings           float32 `gorm:"column:ratings"`
	ViewCount         int     `gorm:"column:view_count"`
	CommentCount      int     `gorm:"column:comment_count"`
	CategoryOne       string  `gorm:"column:category_one"`
	CategoryTwo       string  `gorm:"column:category_two"`
	ContractStatus    string  `gorm:"column:contract_status"`
	ContractType      string  `gorm:"column:contract_type"`
	WritingStatus     string  `gorm:"column:writing_status"`
	SourceType        string  `gorm:"column:source_type"`
	SourceUrl         string  `gorm:"column:source_url"`
}

type BookTagQuery struct {
	FictionBook
	TagNames string
	Author   string
}

func GetTagNames(labelIds string) (string, error) {
	var names []string
	err := db.Table(tagTableName).
		Select("group_concat(fiction_book_tags.label_name separator ',') as tag_names").
		Where("find_in_set(fiction_book_tags.id, (?))", labelIds).
		Pluck("tag_names", &names).Error
	if err != nil {
		return "", err
	}
	if len(names) <= 0 {
		//sugar.Warnf("GetTagNames: get tags from label string result is empty")
		return "", gorm.ErrRecordNotFound
	}
	return names[0], nil
}

func (l *BookTagQueryList) GetTagsFromLabelStr() (err error) {
	var wg sync.WaitGroup

	for _, query := range *l {
		wg.Add(1)
		go func(query *BookTagQuery) {
			tagNames, innerErr := GetTagNames(query.Labels)
			if innerErr != nil {
				//sugar.Errorf("GetTagsFromLabelStr: get tag from book label string failed, innerErr: %v", innerErr)
				//_ = multierr.Append(err, innerErr)
			} else {
				query.TagNames = tagNames
			}
			wg.Done()
		}(query)
	}

	wg.Wait()

	return
}
func (l *BookTagQueryList) GetPopularBooks(orderBy string, limit int, sources []string) error {
	writerNameSubQuery := fmt.Sprintf(
		`(select w.writer_name from %s w 
					where w.id = b.pseudonym_writer_id) as author`,
		writerTableName,
	)
	if err := db.Table(BookTableName+" b").
		Where("b.deleted_at is null AND b.source_type in (?)", sources).
		Select([]string{"b.*", writerNameSubQuery}).
		Order("b." + orderBy + " desc").
		Limit(limit).
		Find(l).
		Error; err != nil {
		return err
	}

	err := l.GetTagsFromLabelStr()
	if err != nil {
		fmt.Printf("GetPopularBooks: can not get tag names from book label field, err: %v", err)
		//sugar.Errorf("GetPopularBooks: can not get tag names from book label field, err: %v", err)
	}

	return nil
}

type BookTagQueryList []*BookTagQuery

func main() {
	fmt.Print("数据库开始连接\n")
	db, err := gorm.Open(
		"mysql",
		"root:peng1365451772@/novel?charset=utf8&parseTime=True&loc=Local",
	)
	if err != nil {
		fmt.Println("出现异常异常原因", err)
	}
	defer db.Close()
	var banner []BannerBook
	db.Table("read_banners").Select("book_name,action_type,url,introduction,book_id").Where("app_id = 'com.fantasy.best.novel'").Order("sort").Scan(&banner)
	writerNameSubQuery := fmt.Sprintf(
		`(select w.writer_name from %s w 
					where w.id = b.pseudonym_writer_id) as author`,
		writerTableName,
	)
	book := &BookTagQueryList{}
	err = db.Table("fiction_books b").
		Select([]string{"b.*", writerNameSubQuery}).
		Joins("right join read_banners r on b.book_id = r.book_id").
		Where("r.app_id = ?", "com.fantasy.best.novel").
		Order("r.sort").
		Find(book).Error
	if err != nil {
		fmt.Println("出新异常，异常原因", err)
	}
	fmt.Printf("打印从数据库取出的值:\n")
	for _, v := range *book {
		fmt.Println(v)
	}

}
