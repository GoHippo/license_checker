package points

import (
	"encoding/json"
	"fmt"
	"github.com/GoHippo/license_checker/pkg/sign"
	"github.com/GoHippo/license_checker/pkg/uuid"
	"github.com/GoHippo/slogpretty/sl"
	"github.com/hyperboloide/lk"
	"log/slog"
	"time"
)

type coreRunner interface {
	Run()
}

type Run interface {
	Run()
}

type Points struct {
	log                *slog.Logger
	PublicLicenseKeyLk string
	License            chan string
	LicenseData        chan LicenseData
	Ok                 chan Run
	runCore            coreRunner
}

type LicenseData struct {
	PublicKey string `json:"PublicKey"`
	Sign      string `json:"Sign"`
}

func NewPointsLicense(runCore coreRunner, publicKeyLK string, log *slog.Logger) *Points {
	//log.Debug("Start NewPointsLicense")

	p := &Points{
		log:                log,
		PublicLicenseKeyLk: publicKeyLK,
		License:            make(chan string, 0),
		LicenseData:        make(chan LicenseData, 0),
		Ok:                 make(chan Run, 0),
		runCore:            runCore,
	}

	p.goLicense()
	p.goKeyPC()
	return p
}

func (p Points) CheckPoint(lic string) Run {
	p.License <- lic
	return <-p.Ok
}

func (p Points) goLicense() {
	//p.log.Debug("Start goLicense")
	go func() {
		for {
			// continue вместо os.exit()
			select {
			case data := <-p.License:

				pub, err := lk.PublicKeyFromHexString(p.PublicLicenseKeyLk)
				if err != nil {
					p.errprint("PublicKeyFromHexString", err)
					continue
				}

				lic, err := lk.LicenseFromHexString(data)
				if err != nil {
					p.errprint("LicenseFromHexString", err)
					continue
				}

				ok, err := lic.Verify(pub)
				if err != nil {
					p.errprint("Verify lc", err)
					continue
				}

				if len(lic.Data) == 0 {
					p.errprint("Verify lc", fmt.Errorf("license data is empty"))
					continue
				}

				if ok {
					//p.log.Debug("Completed license check")

					var licenseData = LicenseData{}
					//p.log.Debug("Data from license:", slog.String("data", string(data)))
					//p.log.Debug("Data from license:", slog.String("lic.data", string(lic.Data)))

					if err := json.Unmarshal(lic.Data, &licenseData); err != nil {
						p.errprint("Unmarshal lc", err)
						continue
					}

					p.LicenseData <- licenseData

				} else {
					p.errprint("Verify lc", fmt.Errorf("lk not ok"))
				}

			default:
				time.Sleep(time.Millisecond * 20)
			}
		}
	}()
}

func (p Points) goKeyPC() {
	//p.log.Debug("Start goKeyPC")
	go func() {
		for {
			select {
			// continue вместо os.exit()
			case data := <-p.LicenseData:

				uuid, err := uuid.GetUUID()
				if err != nil {
					p.errprint("GetUUID", err)
					continue
				}

				//p.log.Debug("Completed create machine id")

				ok, err := sign.VerifySign(uuid, data.Sign, data.PublicKey)
				if err != nil {
					//p.log.Debug("sign data",
					//	slog.String("uuid", uuid),
					//	slog.String("sign", data.Sign),
					//	slog.String("sign_pub", data.PublicKey),
					//)
					p.errprint("VerifySign", err)
					continue
				}

				if ok {
					if p.log != nil {
						p.log.Debug("license verify: ok")
					}

					p.Ok <- Run(next{p.runCore})
				} else {
					p.errprint("VerifySign", fmt.Errorf("sign not ok"))
					p.Ok <- Run(nil)
				}

			default:
				time.Sleep(time.Millisecond * 20)
			}
		}
	}()
}

func (p *Points) errprint(title string, err error) {
	if p.log != nil {
		p.log.Debug("Error in license check", slog.String("title", title), sl.Err(err))
	}
}

type next struct {
	interf coreRunner
}

func (r next) Run() {
	r.interf.Run()
}
