package gorm

import (
	"reflect"
	"sort"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/yushaona/gjson"
)

func replace(s string, start, length int, new string) string {
	t := make([]byte, len(s)+length)
	w := 0
	w += copy(t[0:], s[0:start])
	w += copy(t[w:], new)
	start += length
	w += copy(t[w:], s[start:])
	return string(t[0:w])
}

type QueryDB struct {
	o orm.Ormer
}

func NewQuery() *QueryDB {
	t := new(QueryDB)
	t.o = orm.NewOrm()
	return t
}

//RawQuery Query Data
func (t *QueryDB) RawQuery(query string, params *gjson.GJSON) (result *gjson.GJSON, num int64, err error) {
	//分配一个新数组

	result = gjson.NewGJSON(gjson.TypeArray)
	o := t.o

	if params == nil {
		params = gjson.NewGJSON(gjson.TypeObject)
	}
	var args []interface{}
	posKey := make(map[int]string)
	var sortPos []int
	// convert to Upper word
	sql := strings.ToUpper(query)
	// params.MustMap
	ps := params.Interface().(map[string]interface{})
	if err != nil {
		return
	}
	for name, _ := range ps { //random sort
		nameUpper := strings.ToUpper(name)
		var last int = 0
		for {
			if last >= len(sql) {
				break
			}
			pos := strings.Index(sql[last:], ":"+nameUpper)
			if pos != -1 {
				pos += last
				last = pos + len(":"+nameUpper)
				if last < len(sql) {
					a := sql[last]
					if (a >= '0' && a <= '9') || (a >= 'A' && a <= 'z') {
						continue
					}
				}
				posKey[pos] = name
				sortPos = append(sortPos, pos)
			} else { //说明提供的字段在sql中没有，就不需要管了
				break
			}
		}
	}
	sort.Ints(sortPos)
	//从后往前替换字符串
	for i := len(sortPos) - 1; i >= 0; i-- {
		pos := sortPos[i]
		name := posKey[pos]
		nameUpper := strings.ToUpper(name)
		sql = replace(sql, pos, len(":"+nameUpper), "?")
		query = replace(query, pos, len(":"+nameUpper), "?")
	}
	//加入参数
	for _, pos := range sortPos {
		name := posKey[pos]
		val := ps[name]
		args = append(args, val)
	}
	var temp []orm.Params
	num, err = o.Raw(query, args).Values(&temp)
	if err != nil {
		return
	}
	for _, v := range temp {
		item := result.AddItem()
		for mk, mv := range v {
			switch reflect.ValueOf(mv).Kind() {
			case reflect.Invalid:
				item.SetString(mk, "")
			default:
				item.SetString(mk, mv.(string))
			}
		}
	}
	return
}

func (t *QueryDB) RawExec(query string, params *gjson.GJSON) (lastInsertId int64, err error) {

	o := t.o

	if params == nil {
		params = gjson.NewGJSON(gjson.TypeObject)
	}
	var args []interface{}
	posKey := make(map[int]string)
	var sortPos []int
	// convert to Upper word
	sql := strings.ToUpper(query)
	// params.MustMap
	ps := params.Interface().(map[string]interface{})
	if err != nil {
		return
	}
	for name, _ := range ps { //random sort
		nameUpper := strings.ToUpper(name)
		var last int = 0
		for {
			if last >= len(sql) {
				break
			}
			pos := strings.Index(sql[last:], ":"+nameUpper)
			if pos != -1 {
				pos += last
				last = pos + len(":"+nameUpper)
				if last < len(sql) {
					a := sql[last]
					if (a >= '0' && a <= '9') || (a >= 'A' && a <= 'z') {
						continue
					}
				}
				posKey[pos] = name
				sortPos = append(sortPos, pos)
			} else { //说明提供的字段在sql中没有，就不需要管了
				break
			}
		}
	}
	sort.Ints(sortPos)
	//从后往前替换字符串
	for i := len(sortPos) - 1; i >= 0; i-- {
		pos := sortPos[i]
		name := posKey[pos]
		nameUpper := strings.ToUpper(name)
		sql = replace(sql, pos, len(":"+nameUpper), "?")
		query = replace(query, pos, len(":"+nameUpper), "?")
	}
	//加入参数
	for _, pos := range sortPos {
		name := posKey[pos]
		val := ps[name]
		args = append(args, val)
	}
	r, err := o.Raw(query, args).Exec()
	if err != nil {
		return
	}
	lastInsertId, err = r.LastInsertId()
	if err != nil {
		return
	}
	return
}
