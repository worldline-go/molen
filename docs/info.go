package docs

import (
	"github.com/worldline-go/swagger"
)

func Info(basePath, version string) error {
	return swagger.SetInfo( //nolint:wrapcheck // no need
		swagger.WithTitle("molen"),
		swagger.WithVersion(version),
		swagger.WithBasePath(basePath),
	)
}
