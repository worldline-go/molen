package docs

import (
	"github.com/worldline-go/auth"
	"github.com/worldline-go/swagger"
)

func Info(basePath, version string, provider auth.InfProvider) error {
	return swagger.SetInfo( //nolint:wrapcheck // no need
		swagger.WithTitle("molen"),
		swagger.WithVersion(version),
		swagger.WithCustom(map[string]interface{}{
			"tokenUrl": provider.GetTokenURLExternal(),
			"authUrl":  provider.GetAuthURLExternal(),
		}),
		swagger.WithBasePath(basePath),
	)
}
