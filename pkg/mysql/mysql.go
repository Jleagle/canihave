package mysql

import (
	"sync"

	"github.com/Jleagle/canihave/pkg/config"
	"github.com/Jleagle/canihave/pkg/logger"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	client     *sqlx.DB
	clientLock sync.Mutex
)

func getClient() {

	var err error
	client, err = sqlx.Connect("mysql", "Uid="+config.MySQLUser+" Pwd="+config.MySQLPass+" Server="+config.MySQLIP+" Database="+config.MySQLDB)
	if err != nil {
		logger.Logger.Error("Connecting to MySQL", zap.Error(err))
	}
}
