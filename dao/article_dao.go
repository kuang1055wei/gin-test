package dao

import (
	"gin-test/model"
	"gorm.io/gorm"
)

var ArticleDao = neArticleDao()

func neArticleDao() *articleDao {
	return &articleDao{}
}

type articleDao struct {
}

//查询兑换码
func (c *articleDao) Get(db *gorm.DB, id int) *model.Article {
	code := &model.Article{}
	res := db.First(code, id)
	if res.Error != nil {
		return nil
	}
	return code
}

func (c *articleDao) Take(db *gorm.DB, where ...interface{}) *model.Article {
	ret := &model.Article{}
	if err := db.Take(ret, where...).Error; err != nil {
		return nil
	}
	return ret
}

func (c *articleDao) Create(db *gorm.DB, t *model.Article) (err error) {
	err = db.Create(t).Error
	return
}

func (c *articleDao) Update(db *gorm.DB, t *model.Article) (err error) {
	err = db.Save(t).Error
	return
}

func (c *articleDao) Updates(db *gorm.DB, id int64, columns map[string]interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).Updates(columns).Error
	return
}

func (c *articleDao) UpdateColumn(db *gorm.DB, id int64, name string, value interface{}) (err error) {
	err = db.Model(&model.Article{}).Where("id = ?", id).UpdateColumn(name, value).Error
	return
}

func (c *articleDao) Delete(db *gorm.DB, id int64) {
	db.Delete(&model.Article{}, "id = ?", id)
}

// BatchSave 批量插入数据
//func (c *articleDao) BatchSave(db *gorm.DB, courseCode []*model.Article) error {
//	cc := model.Article{}
//	var buffer bytes.Buffer
//	sql := "insert into `%s` (`course_id`,`goods_id`,`ship_code`,`exchange_code`) VALUES "
//	sql = fmt.Sprintf(sql, cc.TableName())
//
//	if _, err := buffer.WriteString(sql); err != nil {
//		return err
//	}
//	for key, val := range courseCode {
//		if val == nil {
//			continue
//		}
//		if len(courseCode)-1 == key {
//			buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s');", val.CourseID, val.GoodsID, val.ShipCode, val.ExchangeCode))
//		} else {
//			buffer.WriteString(fmt.Sprintf("('%d','%d','%s','%s'),", val.CourseID, val.GoodsID, val.ShipCode, val.ExchangeCode))
//		}
//	}
//	return db.Exec(buffer.String()).Error
//}

//获取商品的解锁码数量
func (c *articleDao) GetGoodsCodeNum(db *gorm.DB, id int) int64 {
	var total int64
	if err := db.Model(&model.Article{}).Where("goods_id = ?", id).Count(&total); err != nil {
		return 0
	}
	return total
}

func (c *articleDao) FindPageByParams(db *gorm.DB, params *model.QueryParams) (list []model.Article, paging *model.Paging) {
	return c.FindPageByCnd(db, &params.SqlCnd)
}

func (c *articleDao) FindPageByCnd(db *gorm.DB, cnd *model.SqlCnd) (list []model.Article, paging *model.Paging) {
	cnd.Find(db, &list)
	count := cnd.Count(db, &model.Article{})

	paging = &model.Paging{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return
}
