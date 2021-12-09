package osutils

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

func IsAdmin() (bool, error) {
	var sid *windows.SID

	logrus.Debug("Checking for admin privileges")

	// https://docs.microsoft.com/en-us/windows/win32/api/securitybaseapi/nf-securitybaseapi-checktokenmembership#examples
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false, errors.Wrap(err, "SID Error")
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)

	isAdmin, err := token.IsMember(sid)
	if err != nil {
		return false, errors.Wrap(err, "Token Membership Error")
	}

	return isAdmin, nil
}
