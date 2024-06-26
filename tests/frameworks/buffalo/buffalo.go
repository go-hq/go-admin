package buffalo

import (
	// add buffalo adapter
	_ "github.com/go-hq/go-admin/adapter/buffalo"
	"github.com/go-hq/go-admin/modules/config"
	"github.com/go-hq/go-admin/modules/language"
	"github.com/go-hq/go-admin/plugins/admin/modules/table"
	"github.com/go-hq/themes/adminlte"

	// add mysql driver
	_ "github.com/go-hq/go-admin/modules/db/drivers/mysql"
	// add postgresql driver
	_ "github.com/go-hq/go-admin/modules/db/drivers/postgres"
	// add sqlite driver
	_ "github.com/go-hq/go-admin/modules/db/drivers/sqlite"
	// add mssql driver
	_ "github.com/go-hq/go-admin/modules/db/drivers/mssql"
	// add adminlte ui theme
	_ "github.com/go-hq/themes/adminlte"

	"github.com/go-hq/go-admin/template"
	"github.com/go-hq/go-admin/template/chartjs"

	"net/http"
	"os"

	"github.com/go-hq/go-admin/engine"
	"github.com/go-hq/go-admin/plugins/admin"
	"github.com/go-hq/go-admin/plugins/example"
	"github.com/go-hq/go-admin/tests/tables"
	"github.com/gobuffalo/buffalo"
)

func internalHandler() http.Handler {
	bu := buffalo.New(
		buffalo.Options{
			Env:  "test",
			Addr: "127.0.0.1:9033",
		},
	)

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(tables.Generators)
	adminPlugin.AddGenerator("user", tables.GetUserTable)

	examplePlugin := example.NewExample()

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfigFromJSON(os.Args[len(os.Args)-1]).
		AddPlugins(adminPlugin, examplePlugin).Use(bu); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	bu.ServeFiles("/uploads", http.Dir("./uploads"))

	return bu
}

func NewHandler(dbs config.DatabaseList, gens table.GeneratorList) http.Handler {
	bu := buffalo.New(
		buffalo.Options{
			Env:  "test",
			Addr: "127.0.0.1:9033",
		},
	)

	eng := engine.Default()

	adminPlugin := admin.NewAdmin(gens)

	template.AddComp(chartjs.NewChart())

	if err := eng.AddConfig(
		&config.Config{
			Databases: dbs,
			UrlPrefix: "admin",
			Store: config.Store{
				Path:   "./uploads",
				Prefix: "uploads",
			},
			Language:    language.EN,
			IndexUrl:    "/",
			Debug:       true,
			ColorScheme: adminlte.ColorschemeSkinBlack,
		},
	).
		AddPlugins(adminPlugin).Use(bu); err != nil {
		panic(err)
	}

	eng.HTML("GET", "/admin", tables.GetContent)

	bu.ServeFiles("/uploads", http.Dir("./uploads"))

	return bu
}
