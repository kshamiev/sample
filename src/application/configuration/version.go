package configuration // import "application/configuration"

import (
	"log"

	"application/configuration/semver"
)

// NewVersion Создание объекта описания версии приложения на основе строковых переменных: версии и сборки
func NewVersion(version string, build string) (ret *semver.Version) {
	var err error
	var tmp string

	if build != "" {
		tmp = version + "+build." + build
	} else {
		tmp = version
	}

	if ret, err = semver.NewVersion(tmp); err != nil {
		log.Fatalf("Application version error: %s", err.Error())
	}
	return
}
