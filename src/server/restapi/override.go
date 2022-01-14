package restapi

import (
	c "github.com/xy3/synche/src/server"
)

func OverrideFlags() {
	port = c.Config.Server.Port
}
