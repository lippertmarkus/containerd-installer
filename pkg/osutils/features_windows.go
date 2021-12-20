package osutils

import (
	"github.com/google/glazier/go/dism"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
	"strings"
	"syscall"
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
		// TODO returns err 1 on Windows Server versions if features already enabled, see https://github.com/google/glazier/pull/460
		if e, ok := err.(syscall.Errno); ok && int(e) == 1 {
			logrus.Warnf("DISM: ignoring error code %d with message \"%s\" for Windows Server systems indicating all required features are installed", int(e), err)
		} else {
			return false, err
		}
	}

	return false, nil
}
