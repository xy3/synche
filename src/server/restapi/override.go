package restapi

import c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"

func OverrideFlags() {
	port = c.Config.Server.Port
}