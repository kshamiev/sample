package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Заголовок помощи коммандной строки
func initCLIHelp(name string) {
	var tmp []string
	if name != "" {
		kingpin.CommandLine.Help = fmt.Sprintf("--- ### %s ### ---", name)
		return
	}
	if tmp = strings.Split(os.Args[0], string(os.PathSeparator)); len(tmp) > 0 {
		kingpin.CommandLine.Help = fmt.Sprintf("--- ### %s ### ---", tmp[len(tmp)-1])
		return
	}
	kingpin.CommandLine.Help = "--- ####### ---"
}

// Получение аргументов коммандной строки
func initCLIArgs(cnf *impl) {
	const appConfiguration = `APPLICATION_CONFIGURATION`

	initCLIHelp("")
	cnf.args.Version = kingpin.Command("version", `Print version and exit`)
	cnf.args.Daemon = kingpin.Command("daemon", `Run server`)
	cnf.args.Cli = kingpin.Command("cli", `Command line interface (default)`).Default()

	kingpin.Flag("debug", `Debug application mode`).
		Short('d').
		BoolVar(&cnf.debug)
	kingpin.Flag("conf", `Custom path to configuration file. ENV: `+appConfiguration).
		Short('c').
		PlaceHolder("/etc/service/configuration.file.yml").
		StringVar(&cnf.args.FilePath)
	if os.Getenv(appConfiguration) != "" {
		cnf.args.FilePath = os.Getenv(appConfiguration)
	}

	// Custom args
	//	cnf.args.Cli.Command("start", i18n.T("eng", "Start streaming", "Cообщение помощь для CLI - Запуск потока"))
	//	cnf.args.Cli.Command("stop", i18n.T("eng", "Stop streaming", "Cообщение помощь для CLI - Остановка потока"))
	//	cnf.args.Cli.Command("restart", i18n.T("eng", "Restart streaming", "Cообщение помощь для CLI - Перезапуск потока"))
	//	cnf.args.Cli.Command("status", i18n.T("eng", "Status of streaming", "Cообщение помощь для CLI - Состояние"))
}
