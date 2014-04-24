package tagging

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"log"
)

type Resume struct {
	Id          int    `db:"id"`
	Content     string `db:"content"`
	CreatedDate int64  `db:"createdDate"`
	Tags        []Tag  `db:"-"`
}

type Tag struct {
	Id       int    `db:"id"`
	Rid      int    `db:"rid"`
	Selstart int    `db:"selstart"`
	Selend   int    `db:"selend"`
	Tag      string `db:"tag"`
	Tagger   string `db:"tagger"`
	TagDate  int64  `db:"tagDate"`
}

type TagFormat struct {
	Tag       string `db:"tag"`
	Formatstr string `db:"formatstr"`
}

type DAO struct {
	dbmap *gorp.DbMap
}

func (dao *DAO) InitDb(dbname string) {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/"+dbname)
	dao.CheckErr(err, "sql.Open failed")

	// construct a gorp DbMap
	dao.dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	dao.dbmap.AddTableWithName(Resume{}, "resume").SetKeys(true, "Id")
	dao.dbmap.AddTableWithName(Tag{}, "tag").SetKeys(true, "Id")
	dao.dbmap.AddTableWithName(TagFormat{}, "tagformat").SetKeys(false, "Tag")
}

func (dao *DAO) CloseDb() {
	dao.dbmap.Db.Close()
}

func (dao *DAO) CheckErr(err error, msg string) {
	if err != nil {
		log.Println(msg, err)
	}
}
