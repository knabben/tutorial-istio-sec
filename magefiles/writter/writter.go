package writter

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"github.com/muesli/termenv"
	"path"
)

var (
	Istioctl = RunCmd("istioctl")
	Kubectl  = RunCmd("kubectl")
	Kind     = RunCmd("kind")
)

// RunCmd uses Exec underneath, so see those docs for more details.
func RunCmd(cmd string, args ...string) func(args ...string) error {
	return func(args2 ...string) error {
		result := append(args, args2...)
		p := termenv.ColorProfile()
		fmt.Println(
			"\n",
			termenv.String(cmd).Foreground(p.Color("#71BEF2")),
			termenv.String(result...).Bold(),
			"\n",
		)
		return sh.Run(cmd, result...)
	}
}

func AppendFolder(specFolder, p string) string {
	return path.Join(specFolder, p)
}

func Output(out string) {
	fmt.Println()
	fmt.Println(termenv.String(out).Bold())
	fmt.Println()
}
