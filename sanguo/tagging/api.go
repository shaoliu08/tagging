package tagging

import (
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	success bool
	content string
}

type Api struct {
	dao *DAO
}

func (api *Api) Init() {
	api.dao = &DAO{}
	api.dao.InitDb("tagging")
}

// resume
func (api *Api) ListResumes(w rest.ResponseWriter, r *rest.Request) {
	start := 0
	size := 10
	// the parameter's format will be start_sizeperpage
	if query := r.PathParam("start"); query != "" {
		ary := strings.Split(query, "_")
		if len(ary) == 2 {
			if v, err := strconv.Atoi(ary[1]); err == nil {
				size = v
			}
		}
		if v, err := strconv.Atoi(ary[0]); err == nil {
			start = v * size
		}
	}

	var resumes []Resume
	log.Printf("select * from resume limit %d, %d\n", start, size)
	_, err := api.dao.dbmap.Select(&resumes, "select * from resume limit ?,?", start, size)
	if err != nil {
		api.dao.CheckErr(err, "list resumes failed")
		w.WriteJson(&Message{false, err.Error()})
		return
	}

	for i, resume := range resumes {
		var tags []Tag
		_, err := api.dao.dbmap.Select(&tags, "select * from tag where rid=?", resume.Id)
		if err != nil {
			api.dao.CheckErr(err, "list resumes failed")
			w.WriteJson(&Message{false, err.Error()})
			return
		}
		resumes[i].Tags = tags
	}
	w.WriteJson(&resumes)
}

func (api *Api) GetResumeById(w rest.ResponseWriter, r *rest.Request) {
	var resume Resume

	if id, err := strconv.Atoi(r.PathParam("id")); err == nil {
		err = api.dao.dbmap.SelectOne(&resume, "select * from resume where id=?", id)
		if err != nil {
			api.dao.CheckErr(err, "GetResumeById failed")
			w.WriteJson(&Message{false, err.Error()})
			return
		}

		var tags []Tag
		_, err = api.dao.dbmap.Select(&tags, "select * from tag where rid=?", resume.Id)
		if err != nil {
			api.dao.CheckErr(err, "GetResumeById failed")
			w.WriteJson(&Message{false, err.Error()})
			return
		}
		resume.Tags = tags

	}
	w.WriteJson(&resume)

}

// Tag
func (api *Api) SaveTag(w rest.ResponseWriter, r *rest.Request) {
	tag := Tag{}
	err := r.DecodeJsonPayload(&tag)
	if err != nil {
		api.dao.CheckErr(err, "SaveTag failed")
		w.WriteJson(&Message{false, err.Error()})
		return
	}
	if tag.Tag == "" {
		api.dao.CheckErr(err, "SaveTag failed")
		w.WriteJson(&Message{false, "tag required\n" + err.Error()})
		return
	}
	tag.TagDate = time.Now().UnixNano()
	api.dao.dbmap.Insert(&tag)
	w.WriteJson(&tag)
}

func (api *Api) DeleteTagById(w rest.ResponseWriter, r *rest.Request) {
	msg := Message{}

	msg.success = true
	if id, err := strconv.Atoi(r.PathParam("id")); err == nil {
		_, err = api.dao.dbmap.Exec("delete from tag where id=?", id)
		if err != nil {
			api.dao.CheckErr(err, "delete tag failed")
			msg.success = false
		}
	} else {
		msg.success = false
	}

	w.WriteJson(&msg)
}

// tagformat
func (api *Api) ListTagFormats(w rest.ResponseWriter, r *rest.Request) {
	var tagformats []TagFormat

	_, err := api.dao.dbmap.Select(&tagformats, "select * from tagformat")
	if err != nil {
		api.dao.CheckErr(err, "ListTagFormats failed")
		w.WriteJson(&Message{false, err.Error()})
		return

	}
	w.WriteJson(&tagformats)
}

func (api *Api) SaveTagFormat(w rest.ResponseWriter, r *rest.Request) {
	tagformat := TagFormat{}
	err := r.DecodeJsonPayload(&tagformat)
	if err != nil {
		api.dao.CheckErr(err, "SaveTagFormat failed")
		w.WriteJson(&Message{false, err.Error()})
		return
	}
	if tagformat.Tag == "" {
		api.dao.CheckErr(err, "SaveTagFormat failed")
		w.WriteJson(&Message{false, "tag required\n" + err.Error()})
		return
	}
	api.dao.dbmap.Insert(&tagformat)
	w.WriteJson(&tagformat)
}

func (api *Api) DeleteTagFormatByTag(w rest.ResponseWriter, r *rest.Request) {
	msg := Message{}

	if tag := strings.Replace(r.PathParam("tag"), "%20", " ", -1); tag != "" {
		log.Println("DeleteTagFormatByTag:\t" + tag)
		_, err := api.dao.dbmap.Exec("delete from tagformat where tag=?", tag)
		if err != nil {
			msg.success = false
		} else {
			msg.success = true
		}
	}

	w.WriteJson(&msg)
}
