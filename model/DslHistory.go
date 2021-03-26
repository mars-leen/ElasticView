package model

import "ElasticView/engine/db"

//http://sql2struct.atotoa.com/

type GmDslHistoryModel struct {
	ID         int      `gorm:"column:id" json:"id" form:"id" db:"id"`
	Uid        int      `gorm:"column:uid" json:"uid" form:"uid" db:"uid"`
	Method     string   `gorm:"column:method" json:"method" form:"method" db:"method"`
	Path       string   `gorm:"column:path" json:"path" form:"path" db:"path"`
	Body       string   `gorm:"column:body" json:"body" form:"body" db:"body"`
	Created    string   `gorm:"column:created" json:"created" form:"created" db:"created"`
	FilterDate []string `json:"date"`
	IndexName  string   `json:"indexName"`
}

func (this *GmDslHistoryModel) TableName() string {
	return "gm_dsl_history"
}

func (this *GmDslHistoryModel) Insert() (err error) {
	_, err = db.SqlBuilder.
		Insert(this.TableName()).
		SetMap(map[string]interface{}{
			"uid":    this.Uid,
			"method": this.Method,
			"path":   this.Path,
			"body":   this.Body,
		}).RunWith(db.Sqlx).Exec()
	return
}

func (this *GmDslHistoryModel) Clean() (err error) {
	_, err = db.SqlBuilder.
		Delete(this.TableName()).
		Where(db.Eq{"uid": this.Uid}).
		RunWith(db.Sqlx).
		Exec()
	return
}

func (this *GmDslHistoryModel) List() (gmDslHistoryModelList []GmDslHistoryModel, err error) {
	builder := db.SqlBuilder.
		Select("*").
		From(this.TableName()).
		Where(db.Eq{"uid": this.Uid}).
		OrderBy("id desc").
		Limit(20)

	if this.IndexName != "" {
		builder = builder.Where(db.Like{"path": db.CreateLike(this.IndexName)})
	}

	if len(this.FilterDate) > 0 {
		builder = builder.Where(db.Gte{"created": this.FilterDate[0]}).Where(db.Lte{"created": this.FilterDate[1]})
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return
	}
	err = db.Sqlx.Select(&gmDslHistoryModelList, sql, args...)
	return
}
