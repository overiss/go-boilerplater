package scaffold

import "path/filepath"

var dirs = []string{
	"internal/app",
	"internal/behavior",
	"internal/config",
	"internal/model/dto",
	"internal/model/dao",
	"internal/provider",
	"internal/repository/domain",
	"internal/repository/integrations",
	"internal/server/http/handler",
	"internal/server/http/middleware",
	"internal/service",
	"internal/vars",
	"pkg/utils",
	"deploy",
	"docs",
}

const (
	appGoPath             = "internal/app/app.go"
	creatorGoPath         = "internal/app/creator.go"
	routingGoPath         = "internal/app/routing.go"
	behaviorGoPath        = "internal/behavior/behavior.go"
	configGoPath          = "internal/config/config.go"
	requestGoPath         = "internal/model/request.go"
	responseGoPath        = "internal/model/response.go"
	serverContainerGoPath = "internal/server/container.go"
	httpServerGoPath      = "internal/server/http/server.go"
	handlerContainerPath  = "internal/server/http/handler/container.go"
	serviceContainerPath  = "internal/service/container.go"
	varsGoPath            = "internal/vars/vars.go"
	utilsGoPath           = "pkg/utils/utils.go"
)

func buildFiles(serviceName, moduleName string) map[string]string {
	return map[string]string{
		filepath.Join("cmd", serviceName, "main.go"): mainGoTemplate(moduleName),
		appGoPath:             appTemplate(moduleName),
		creatorGoPath:         creatorTemplate,
		routingGoPath:         routingTemplate,
		behaviorGoPath:        behaviorTemplate,
		configGoPath:          configTemplate(serviceName),
		requestGoPath:         modelTemplate,
		responseGoPath:        modelTemplate,
		serverContainerGoPath: serverContainerTemplate(moduleName),
		httpServerGoPath:      httpServerTemplate(moduleName),
		handlerContainerPath:  handlerContainerTemplate,
		serviceContainerPath:  serviceContainerTemplate,
		varsGoPath:            varsTemplate,
		utilsGoPath:           utilsTemplate,
	}
}
