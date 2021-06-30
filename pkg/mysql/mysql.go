package mysql

import (
	"sync"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/memcache-go"
	"github.com/jmoiron/sqlx"
)

var (
	mySQLClient     *sqlx.DB
	mySQLClientLock sync.Mutex

	memcacheClient = memcache.NewClient("localhost:11211")
)

func getClient() (*sqlx.DB, error) {

	mySQLClientLock.Lock()
	defer mySQLClientLock.Unlock()

	c, err := sqlx.Connect("mysql", config.MySQLUser+":"+config.MySQLPass+"@"+config.MySQLIP+"/"+config.MySQLDB+"?parseTime=True")
	if err != nil {
		return nil, err
	}

	mySQLClient = c

	return mySQLClient, nil
}
