package check

import (
	"fmt"
	"github.com/GoHippo/license_checker/license"
	"github.com/GoHippo/license_checker/license/points"
	"github.com/GoHippo/license_checker/license_server_request"
	"log/slog"
)

type LicenseService interface {
	CheckLicense() license.Runner
}

type Runner interface {
	Run()
}

type CheckLicenseOptions struct {
	BoxRun             Runner // интерфейс с запуском
	PublicLicenseKeyLk string
	DataServerOptions  license_server_request.DataServerOptions
	Log                *slog.Logger
}

func CheckLicense(opt CheckLicenseOptions) (license.Runner, error) {
	data, err := license_server_request.GetDataFromServer(opt.DataServerOptions)
	if err != nil {
		if err.Error() == license_server_request.INVALID_AUTHORIZATION {
			return nil, fmt.Errorf("Your license key is invalid - %v or %v ", "The key has expired", "The UUID of the linked machine does not match")
		} else {
			return nil, fmt.Errorf("Error checking license key - %w", err)
		}
	}

	lc := LicenseService(&license.LicenseService{data, points.NewPointsLicense(opt.BoxRun, opt.PublicLicenseKeyLk, opt.Log)})
	r := lc.CheckLicense()

	if r.Run == nil {
		return nil, fmt.Errorf("license failed!")
	}

	return r, nil
}
