package swag

import (
	"fmt"
	"git.moumentei.com/plat_go/golib/tools/qbtool/cmd/base"
	util "git.moumentei.com/plat_go/golib/tools/qbtool/cmd/internal"
	"github.com/swaggo/swag"
	"github.com/swaggo/swag/gen"
	"log"
	"os"
	"path"
)

var Cmd = &base.Command{
	CustomFlags: true,
	Run:         runSwag,
	UsageLine:   base.ToolName + " swag command [command options] [arguments...]",
	Short:       "swag - Automatically generate RESTful API documentation with Swagger 2.0 for Go.",
	Long: `
VERSION:
   ` + swag.Version + `

COMMANDS:
   init  swagger.json & swagger.yml & js/css/icon
COMMAND OPTIONS
   -g  Go file path in which 'swagger general API Info' is written, default is main.go
   -d  Directory you want to parse, default is .(current directory) 
   -p  Property Naming Strategy like snakecase,camelcase,pascalcase, default is camelcase 
   -o  Output directory for all the generated files(swagger.json, swagger.yaml), default is swagger 
  -md  Parse folder containing markdown files to use as description, disabled by default 
 -ven  Parse go files in 'vendor' folder, disabled by default
 -dep  Parse go files in outside dependency folder, disabled by default
`,
}

var (
	g, dir, p, output, md string
	vendor, dependency    bool
)

func init() {
	Cmd.Flag.StringVar(&g, "g", "main.go", "Go file path in which 'swagger general API Info' is written")
	Cmd.Flag.StringVar(&dir, "d", ".", "Directory you want to parse")
	Cmd.Flag.StringVar(&p, "p", "camelcase", "Property Naming Strategy like snakecase,camelcase,pascalcase")
	Cmd.Flag.StringVar(&output, "o", "swagger", "Output directory for all the generated files(swagger.json, swagger.yaml and doc.go)")
	Cmd.Flag.StringVar(&md, "md", "", "Parse folder containing markdown files to use as description, disabled by default")
	Cmd.Flag.BoolVar(&vendor, "ven", false, "Parse go files in 'vendor' folder, disabled by default")
	Cmd.Flag.BoolVar(&dependency, "dep", false, "Parse go files in outside dependency folder, disabled by default")
}

func runSwag(cmd *base.Command, args []string) {
	if len(args) >= 2 {
		cmd.Flag.Parse(args[2:])
	} else {
		fmt.Printf("please see %s help swag\n", base.ToolName)
		return
	}
	if args[1] != "init" {
		fmt.Printf("please see %s help swag\n", base.ToolName)
		return
	}

	switch p {
	case swag.CamelCase, swag.SnakeCase, swag.PascalCase:
	default:
		fmt.Printf("please see %s help swag\n", base.ToolName)
		return
	}

	prepareSwagger(output)
	// swag v1.6.5版本
	gen.New().Build(&gen.Config{
		SearchDir:          dir,
		MainAPIFile:        g,
		PropNamingStrategy: p,
		OutputDir:          output,
		ParseVendor:        vendor,
		ParseDependency:    dependency,
		MarkdownFilesDir:   md,
	})

}

func prepareSwagger(output string) {
	log.Println("Prepare swagger, output:", output)
	if _, err := os.Stat(path.Join(output, "index.html")); err != nil {
		if os.IsNotExist(err) {
			if err := util.DownloadFromURL(util.SwaggerLink, "swagger.zip"); err == nil {
				util.UnzipAndDelete(output, "swagger.zip")
			}
		}
	}
}
