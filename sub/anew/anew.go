package anew

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

func New(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		args = []string{"myapp"}
	}

	mod := args[0]
	modpath := strings.Split(mod, "/")
	modtail := modpath[len(modpath)-1]

	fmt.Println("New a project...")
	err := os.Mkdir(modtail, 0755)
	if err != nil {
		log.Panic("can not create a directory", err)
	}

	makefile(modtail, mod)
	gitignore(modtail)
	license(modtail)
	maingo(modtail)
	gomod(modtail, mod)
	dockerfile(modtail, mod)
}

//go:embed gomod.template
var gomodTemplate string

//go:embed makefile.template
var makefileTemplate string

//go:embed gitignore.template
var gitignoreTemplate string

//go:embed license.template
var licenseTemplate string

//go:embed main.template
var mainTemplate string

//go:embed Dockerfile.template
var dockerFileTemplate string

func gomod(modtail, mod string) {
	f, err := os.Create(modtail + "/go.mod")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	ver := strings.TrimPrefix(runtime.Version(), "go")
	{
		if strings.Count(ver, ".") > 1 {
			lastDot := strings.LastIndex(ver, ".")
			ver = ver[:lastDot]
		}
	}

	t := template.Must(template.New("mod").Parse(gomodTemplate))
	err = t.Execute(f, struct {
		Module    string
		GoVersion string
	}{
		Module:    mod,
		GoVersion: ver,
	})

	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}

func makefile(modtail string, mod string) {
	f, err := os.Create(modtail + "/Makefile")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	t := template.Must(template.New("makefile").Parse(makefileTemplate))
	err = t.Execute(f, struct {
		Module string
	}{
		Module: mod,
	})
	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}

func gitignore(modtail string) {
	f, err := os.Create(modtail + "/.gitignore")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	_, err = f.WriteString(gitignoreTemplate)
	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}

func license(modtail string) {
	f, err := os.Create(modtail + "/LICENSE")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	_, err = f.WriteString(licenseTemplate)
	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}

func maingo(modtail string) {
	f, err := os.Create(modtail + "/main.go")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	_, err = f.WriteString(mainTemplate)
	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}

func dockerfile(modtail string, mod string) {
	f, err := os.Create(modtail + "/Dockerfile")
	if err != nil {
		log.Panic("can not create file", err)
	}
	defer f.Close()

	t := template.Must(template.New("dockerFile").Parse(dockerFileTemplate))
	err = t.Execute(f, struct {
		Module string
	}{
		Module: mod,
	})
	if err != nil {
		log.Panic("can not write file", err)
	}
	f.Sync()
}