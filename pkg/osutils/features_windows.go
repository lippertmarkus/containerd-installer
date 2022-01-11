package osutils

import (
	"github.com/google/glazier/go/dism"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"strings"
)

func EnableFeatures(features []string) (bool, error) {
	logrus.Debugf("Enabling Windows features: %s", strings.Join(features, ", "))

	dismSession, err := dism.OpenSession(dism.DISM_ONLINE_IMAGE, "", "", dism.DismLogErrorsWarningsInfo, "", "")
	if err != nil {
		return false, err
	}
	defer dismSession.Close()

	if err := dismSession.EnableFeature(strings.Join(features, ";"), "", nil, true, nil, nil); err != nil {
		if errors.Is(err, windows.ERROR_SUCCESS_REBOOT_REQUIRED) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}
