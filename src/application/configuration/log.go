package configuration // import "application/configuration"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var (
	rexGraylog = regexp.MustCompile(`(?i)^((tcp|udp)://)*([^:]+)(:(\d+))*$`)
)

// InitLog Initialization log configuration
func (cnf *impl) InitLog(filename string) (err error) {
	var buf []byte
	var cnk []string
	var prt uint64

	// Read all file data
	if buf, err = ioutil.ReadFile(filename); err != nil {
		err = fmt.Errorf("Can't read configuration from file '%s': %s", filename, err.Error())
		return
	}

	// Unmarshal yaml configuration data to structure
	cnf.appLog = new(ApplicationLog)
	if err = yaml.Unmarshal(buf, cnf.appLog); err != nil {
		err = fmt.Errorf("Can't unmarshal data from yaml file '%s': %s", filename, err.Error())
		return
	}

	cnk = rexGraylog.FindStringSubmatch(cnf.appLog.Graylog)
	if len(cnk) == 6 {
		cnf.appLog.GraylogProto = strings.ToLower(cnk[2])
		cnf.appLog.GraylogAddress = strings.ToLower(cnk[3])
		prt, _ = strconv.ParseUint(cnk[5], 10, 64)
	}

	// Defaults
	switch cnf.appLog.GraylogProto {
	case "tcp", "udp":
	default:
		cnf.appLog.GraylogProto = "udp"
	}
	if prt == 0 {
		prt = 12201
	}
	cnf.appLog.GraylogPort = uint16(prt)
	if cnf.appLog.GraylogAddress != "" {
		cnf.appLog.GraylogEnable = true
	}

	return
}
