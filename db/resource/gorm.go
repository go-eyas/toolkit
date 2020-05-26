package resource

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

type Resource struct {
	db        *gorm.DB
	tableName string
	model     *gorm.DB
	scope     *gorm.Scope
	pk        string
	sample    interface{}
}

func NewGormResource(db *gorm.DB, v interface{}) *Resource {
	scope := db.NewScope(v)

	return &Resource{
		sample:    v,
		db:        db,
		tableName: scope.TableName(),
		model:     db.Table(scope.TableName()),
		scope:     scope,
		pk:        scope.PrimaryKey(),
	}
}

func (r *Resource) Row(pk interface{}) *gorm.DB {
	return r.model.Where(r.pk+" = ?", pk)
}

func (r *Resource) Create(v interface{}) error {
	return r.model.Create(v).Error
}
func (r *Resource) Update(pk interface{}, v interface{}) error {
	return r.Row(pk).Updates(v).Error
}
func (r *Resource) Query(v interface{}, slice interface{}) (int64, error) {
	return 0, nil
}
func (r *Resource) Delete(pk interface{}) error {
	return r.model.Delete(r.sample, r.pk+" = ?", pk).Error
}
func (r *Resource) CreateHandle(res http.ResponseWriter, req http.Request) {

}
func (r *Resource) UpdateHandle(res http.ResponseWriter, req http.Request) {}
func (r *Resource) QueryHandle(res http.ResponseWriter, req http.Request)  {}
func (r *Resource) DeleteHandle(res http.ResponseWriter, req http.Request) {}
