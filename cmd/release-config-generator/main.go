//nolint:lll,funlen
package main

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type goreleaserTempldata struct {
	Binaries          map[string]string
	PrimaryBinaries   []string
	SecondaryBinaries []string
	ToolBinaries      []string
}

type systemdTempldata struct {
	Service string
}

type scriptTempldata struct {
	Bins []string
}

func main() {
	var rootpath string

	var ignorestr string

	var primarybinaries []string

	var secondarybinaries []string

	var toolbinaries []string

	ignoredbinaries := make(map[string]struct{})
	requiredbinaries := make(map[string]string)
	filledtemplate := new(bytes.Buffer)

	workingdir := getWorkingDirectory()

	flag.StringVar(&rootpath, "path", workingdir, "binary path of your application")
	flag.StringVar(&ignorestr, "ignore", "", "comma separated list of ignored binaries")
	flag.Parse()

	getIgnoredBinaries(ignorestr, ignoredbinaries)

	primarybinaries, secondarybinaries, toolbinaries = getProductiveBinaries(rootpath, ignoredbinaries, primarybinaries, secondarybinaries, toolbinaries, requiredbinaries)

	// goreleaser file
	templatefilename := ".goreleaser.yml.tmpl"
	goreleaserData := goreleaserTempldata{
		Binaries:          requiredbinaries,
		PrimaryBinaries:   primarybinaries,
		SecondaryBinaries: secondarybinaries,
		ToolBinaries:      toolbinaries,
	}

	readFileAndExecuteTemplate(rootpath, templatefilename, filledtemplate, goreleaserData)

	outputfilename := ".goreleaser.yml"

	writeFile(rootpath, outputfilename, filledtemplate)

	// systemd files
	templatefilename = "init/systemd.service.tmpl"

	for binary := range requiredbinaries {
		filledtemplate.Reset()

		outputfilename = filepath.Join("init", binary, binary+".service")
		systemdData := systemdTempldata{Service: binary}

		readFileAndExecuteTemplate(rootpath, templatefilename, filledtemplate, systemdData)

		writeFile(rootpath, outputfilename, filledtemplate)
	}

	// package scripts
	scriptnames := []string{"preinstall", "postinstall", "preremove", "postremove"}
	scriptbasepath := "scripts/package"

	// primary package
	for _, name := range scriptnames {
		filledtemplate.Reset()

		templatefilename = filepath.Join(scriptbasepath, name+".sh.tmpl")
		outputfilename = filepath.Join(scriptbasepath, "primary", name+".sh")
		scriptData := scriptTempldata{
			Bins: primarybinaries,
		}

		readFileAndExecuteTemplate(rootpath, templatefilename, filledtemplate, scriptData)

		writeFile(rootpath, outputfilename, filledtemplate)
	}

	// secondary scripts
	for _, name := range scriptnames {
		filledtemplate.Reset()

		templatefilename = filepath.Join(scriptbasepath, name+".sh.tmpl")
		outputfilename = filepath.Join(scriptbasepath, "secondary", name+".sh")
		scriptData := scriptTempldata{
			Bins: secondarybinaries,
		}

		readFileAndExecuteTemplate(rootpath, templatefilename, filledtemplate, scriptData)

		writeFile(rootpath, outputfilename, filledtemplate)
	}
}

func writeFile(rootpath string, outputfilename string, filledtemplate *bytes.Buffer) {
	direrr := os.MkdirAll(filepath.Dir(filepath.Join(rootpath, outputfilename)), 0o777)
	if direrr != nil {
		log.Print(direrr)
	}

	fwriteerr := os.WriteFile(filepath.Join(rootpath, outputfilename), filledtemplate.Bytes(), 0o600)
	if fwriteerr != nil {
		log.Fatal(fwriteerr)
	}
}

//nolint:cyclop
func getProductiveBinaries(rootpath string, ignoredbinaries map[string]struct{}, primarybinaries, secondarybinaries, toolbinaries []string, requiredbinaries map[string]string) ([]string, []string, []string) {
	binarylookuperr := filepath.Walk(filepath.Join(rootpath, "cmd"), func(path string, _ os.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" && filepath.Base(path) == "main.go" {
			binname := filepath.Base(filepath.Dir(path))

			if _, ignored := ignoredbinaries[binname]; !ignored {
				switch {
				case strings.HasPrefix(binname, "."):
					ignoredbinaries[binname] = struct{}{}

					return nil
				case strings.HasPrefix(binname, "pdns-secondary"):
					secondarybinaries = append(secondarybinaries, binname)
				case strings.HasPrefix(binname, "pdns-tool") || strings.HasSuffix(binname, "tool"):
					toolbinaries = append(toolbinaries, binname)
				default:
					primarybinaries = append(primarybinaries, binname)
				}

				requiredbinaries[binname], err = filepath.Rel(rootpath, path)
				if err != nil {
					return err
				}
			}
		}

		return err
	})
	if binarylookuperr != nil {
		log.Fatal(binarylookuperr)
	}

	return primarybinaries, secondarybinaries, toolbinaries
}

func getIgnoredBinaries(ignorestr string, ignoredbinaries map[string]struct{}) {
	for _, bin := range strings.Split(ignorestr, ",") {
		ignoredbinaries[bin] = struct{}{}
	}
}

func getWorkingDirectory() string {
	workingdir, wdgeterr := os.Getwd()
	if wdgeterr != nil {
		log.Fatal(wdgeterr)
	}

	return workingdir
}

func readFileAndExecuteTemplate(rootpath, templatefilename string, filledtemplate io.Writer, templatedata interface{}) {
	templ, tplerr := template.ParseFiles(filepath.Join(rootpath, templatefilename))
	if tplerr != nil {
		log.Fatal(tplerr)
	}

	tplexecerr := templ.Execute(filledtemplate, templatedata)
	if tplexecerr != nil {
		log.Fatal(tplexecerr)
	}
}
