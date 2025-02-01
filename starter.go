package mysqlStarter

import (
	"github.com/dennesshen/photon-core-starter/core"
	"photon-mysql-starter/mysql"
)

func init() {
	core.RegisterAddModule(mysql.Start)
}
