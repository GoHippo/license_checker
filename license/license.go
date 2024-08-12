package license

import (
	"github.com/GoHippo/license_checker/license/points"
)

type LicenseService struct {
	DataFromServer string
	*points.Points
}

type coreRunner interface {
	Run()
}

type Runner interface {
	Run()
}

func (ls LicenseService) CheckLicense() Runner {
	return ls.CheckPoint(ls.DataFromServer)
}
