package mgo

import (
	"github.com/hero1s/golib/log"
	"gopkg.in/mgo.v2"
	"time"
)

var DefaultMgoSession *mgo.Session

func InitMgo(addr []string, database, username, password string, timeout time.Duration, poolLimit int) error {
	dialInfo := &mgo.DialInfo{
		Addrs:     addr,
		Timeout:   timeout,
		Database:  database,
		Username:  username,
		Password:  password,
		PoolLimit: poolLimit,
	}
	var err error
	DefaultMgoSession, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Errorf("InitMgo Err:%v", err)
	}
	return err
}

func getSession() *mgo.Session {
	return DefaultMgoSession.Copy()
}

func connect(db, collection string) (*mgo.Session, *mgo.Collection) {
	ms := getSession()
	c := ms.DB(db).C(collection)
	ms.SetMode(mgo.Monotonic, true)
	return ms, c
}

func getDb(db string) (*mgo.Session, *mgo.Database) {
	ms := getSession()
	return ms, ms.DB(db)
}

/*
	设置索引
*/
func SetKeyIndex(db, collection string, keys []string) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	err := c.EnsureIndexKey(keys...)
	return err
}

func GetKeyIndexs(db, collection string) ([]mgo.Index, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Indexes()
}

/*
	每次操作之后都要主动关闭 Session defer Session.Close()
	db:操作的数据库
	collection:操作的文档(表)
	doc:要插入的数据
*/
func Insert(db, collection string, doc interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Insert(doc)
}

/*
data := &Data{
    Id:      bson.NewObjectId().Hex(),
    Title:   "标题",
    Des:     "博客描述信息",
    Content: "博客的内容信息",
    Img:     "https://upload-images.jianshu.io/upload_images/8679037-67456031925afca6.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/700",
    Date:    time.Now(),
}

err := db.Insert("Test", "TestModel", data)
*/

/*
	db:操作的数据库
	collection:操作的文档(表)
	query:查询条件
	selector:需要过滤的数据(projection)
	result:查询到的结果
*/
func FindOne(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).One(result)
}

func FindAll(db, collection string, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).All(result)
}

func FindDistinct(db, collection, key string, query, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Distinct(key, result)
}

/*
    查询title="标题",并且返回结果中去除`_id`字段
	var result Data
	err = db.FindOne(database, collection, bson.M{"title": "标题"}, bson.M{"_id":0}, &result)
*/

/*
db:操作的数据库
collection:操作的文档(表)
selector:更新条件
update:更新的操作
*/
func Update(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Update(selector, update)
}

// 更新，如果不存在就插入一个新的数据 `upsert:true`
func Upsert(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.Upsert(selector, update)
	return err
}

// `multi:true`
func UpdateAll(db, collection string, selector, update interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.UpdateAll(selector, update)
	return err
}

/*
	test
	err = db.Update(database, collection, bson.M{"_id": "5b3c30639d5e3e24b8786540"}, bson.M{"$set": bson.M{"title": "更新标题"}})
*/

func Remove(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Remove(selector)
}

func RemoveAll(db, collection string, selector interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	_, err := c.RemoveAll(selector)
	return err
}

/*
	err = db.Remove(database,collection,bson.M{"_id":"5b3c30639d5e3e24b8786540"})
*/
func FindPage(db, collection string, page, limit int, query, selector, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Select(selector).Skip(page * limit).Limit(limit).All(result)
}

func FindIter(db, collection string, query interface{}) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()

	return c.Find(query).Iter()
}

func IsEmpty(db, collection string) bool {
	ms, c := connect(db, collection)
	defer ms.Close()
	count, err := c.Count()
	if err != nil {
		log.Error(err)
	}
	return count == 0
}

func Count(db, collection string, query interface{}) (int, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	return c.Find(query).Count()
}

//insert one or multi documents
func BulkInsert(db, collection string, docs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Insert(docs...)
	return bulk.Run()
}

func BulkRemove(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()

	bulk := c.Bulk()
	bulk.Remove(selector...)
	return bulk.Run()
}

func BulkRemoveAll(db, collection string, selector ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.RemoveAll(selector...)
	return bulk.Run()
}

func BulkUpdate(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Update(pairs...)
	return bulk.Run()
}

func BulkUpdateAll(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.UpdateAll(pairs...)
	return bulk.Run()
}

func BulkUpsert(db, collection string, pairs ...interface{}) (*mgo.BulkResult, error) {
	ms, c := connect(db, collection)
	defer ms.Close()
	bulk := c.Bulk()
	bulk.Upsert(pairs...)
	return bulk.Run()
}

func PipeAll(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.All(result)
}

func PipeOne(db, collection string, pipeline, result interface{}, allowDiskUse bool) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}
	return pipe.One(result)
}

func PipeIter(db, collection string, pipeline interface{}, allowDiskUse bool) *mgo.Iter {
	ms, c := connect(db, collection)
	defer ms.Close()
	var pipe *mgo.Pipe
	if allowDiskUse {
		pipe = c.Pipe(pipeline).AllowDiskUse()
	} else {
		pipe = c.Pipe(pipeline)
	}

	return pipe.Iter()

}

func Explain(db, collection string, pipeline, result interface{}) error {
	ms, c := connect(db, collection)
	defer ms.Close()
	pipe := c.Pipe(pipeline)
	return pipe.Explain(result)
}

func GridFSCreate(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Create(name)
}

func GridFSFindOne(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).One(result)
}

func GridFSFindAll(db, prefix string, query, result interface{}) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Find(query).All(result)
}

func GridFSOpen(db, prefix, name string) (*mgo.GridFile, error) {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Open(name)
}

func GridFSRemove(db, prefix, name string) error {
	ms, d := getDb(db)
	defer ms.Close()
	gridFs := d.GridFS(prefix)
	return gridFs.Remove(name)
}
