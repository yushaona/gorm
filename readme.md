


## Create Test Table
```
CREATE TABLE `userinfo` (
  `uid` int(10) NOT NULL AUTO_INCREMENT,
  `username` varchar(64) CHARACTER SET utf8 COLLATE utf8_unicode_ci DEFAULT NULL,
  `department` varchar(64) DEFAULT NULL,
  `created` datetime DEFAULT NULL,
  `num` decimal(10,2) DEFAULT NULL,
  `time` date DEFAULT NULL,
  `timestamp` timestamp NULL DEFAULT NULL,
  `count` varchar(45) DEFAULT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;
```

## Test code

```
package main

import (
    "github.com/yushaona/gorm"
)
func main() {

	q := gorm.NewQuery()
	var a gjson.GJSON
	a.SetString("id", "0")
	r, n, e := q.RawQuery(" select * from userinfo where uid> :id order by uid desc ", &a)
	fmt.Println(r.ToString())
	fmt.Println(n)
	if e != nil {
		fmt.Println(e.Error())
	}

	var b gjson.GJSON
	b.SetString("username", "von")
	b.SetString("department", "pm")
	b.SetString("created", time.Now().Format("2006-01-02 15:04:05"))
	b.SetString("num", "15.22")
	b.SetString("time", time.Now().Format("2006-01-02 15:04:05"))
	b.SetString("timestamp", time.Now().Format("2006-01-02 15:04:05"))
	b.SetString("count", "10")
	nn, err := q.RawExec(" insert into userinfo (username,department,created,num,time,timestamp) values (:username,:department,:created,:num,:time,:timestamp) ", &b)
	if err != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println(nn)
	}
}

```