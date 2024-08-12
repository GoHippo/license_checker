package uuid

import (
	"github.com/denisbrodbeck/machineid"
	"github.com/jaypipes/ghw"
)

func GetUUID() (string, error) {
	id, err := machineid.ID()
	if err != nil {
		return "", err
	}

	info, err := ghw.CPU()
	if err != nil {
		return "", err
	}

	info2, err := ghw.GPU()
	if err != nil {
		return "", err
	}

	id, err = machineid.ProtectedID(id + info.String() + info2.String())
	if err != nil {
		return "", err
	}

	return id, err
}
