package help

import (
	"bufio"
	"fmt"
	"github.com/hero1s/golib/tools/qbtool/cmd/base"
	"io"
	"strings"
	"text/template"
	"unicode"
	"unicode/utf8"
)

var Cmd = &base.Command{
	UsageLine: "",
	Short:     "",
	Long:      ``,
}

// Help implements the 'help' command.
func Help(w io.Writer, args []string) {
	cmd := base.Cmd
Args:
	for i, arg := range args {
		for _, sub := range cmd.Commands {
			if sub.Name() == arg {
				cmd = sub
				continue Args
			}
		}

		// helpSuccess is the help command using as many args as possible that would succeed.
		helpSuccess := base.ToolName + " help"
		if i > 0 {
			helpSuccess += " " + strings.Join(args[:i], " ")
		}
		fmt.Fprintf(w, base.ToolName+" help %s: unknown help topic. Run '%s'.\n", strings.Join(args, " "), helpSuccess)
		base.SetExitStatus(2) // failed at 'go help cmd'
		base.Exit()
	}

	if len(cmd.Commands) > 0 {
		PrintUsage(w, cmd)
	} else {
		tmpl(w, helpTemplate, cmd)
	}
	// not exit 2: succeeded at 'go help cmd'.
	return
}

var usageTemplate = `{{.Long | trim}}

Usage:

	{{.UsageLine}} <command> [arguments]

The commands are:
{{range .Commands}}{{if or (.Runnable) .Commands}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "go help{{with .LongName}} {{.}}{{end}} <command>" for more information about a command.
{{if eq (.UsageLine) "go"}}
Additional help topics:
{{range .Commands}}{{if and (not .Runnable) (not .Commands)}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}{{end}}

Use "go help{{with .LongName}} {{.}}{{end}} <topic>" for more information about that topic.
{{end}}
`

var helpTemplate = `{{if .Runnable}}usage: {{.UsageLine}}

{{end}}{{.Long | trim}}
`

// An errWriter wraps a writer, recording whether a write error occurred.
type errWriter struct {
	w   io.Writer
	err error
}

func (w *errWriter) Write(b []byte) (int, error) {
	n, err := w.w.Write(b)
	if err != nil {
		w.err = err
	}
	return n, err
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "capitalize": capitalize})
	template.Must(t.Parse(text))
	ew := &errWriter{w: w}
	err := t.Execute(ew, data)
	if ew.err != nil {
		// I/O error writing. Ignore write on closed pipe.
		if strings.Contains(ew.err.Error(), "pipe") {
			base.SetExitStatus(1)
			base.Exit()
		}
		base.Fatalf("writing output: %v", ew.err)
	}
	if err != nil {
		panic(err)
	}
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToTitle(r)) + s[n:]
}

func PrintUsage(w io.Writer, cmd *base.Command) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, cmd)
	bw.Flush()
}
