package configuration // import "application/configuration"

import "gopkg.in/webnice/web.v1"
import "gopkg.in/webnice/kit.v1/modules/db"

// WEBServerConfiguration WEB server configuration structure
type WEBServerConfiguration struct {
	Server            web.Configuration `yaml:"Server"`            // Конфигурация WEB сервера
	DocumentRoot      string            `yaml:"DocumentRoot"`      // Корень http сервера
	Pages             string            `yaml:"Pages"`             // Расположение html шаблонов страниц, чей код генерируется на стороне сервера
	ErrorCodeTemplate map[int]string    `yaml:"ErrorCodeTemplate"` // Шаблоны для Content-type text/html соответствующие кодам http ответа, файл шаблона ищется в Template
}

// Application Описание основной структуры конфигурационного файла
type Application struct {
	ApplicationName  string `yaml:"ApplicationName"`  // Название приложения
	WorkingDirectory string `yaml:"WorkingDirectory"` // Рабочая директория сервера, сразу после запуска приложение меняет текущую директорию
	PidFile          string `yaml:"PidFile"`          // PID file файл, в который записывается текущий идентификатор процесса
	TempPath         string `yaml:"TempPath"`         // Путь к временной директории сервера
	LogConfiguration string `yaml:"LogConfiguration"` // Путь и имя файла конфигурации системы логирования (ApplicationLog)
	LogPath          string `yaml:"LogPath"`          // Путь к папке размещения лог файлов
	StateFile        string `yaml:"StateFile"`        // Файл состояния
	SocketFile       string `yaml:"SocketFile"`       // Сокет файл для коммуникации с CLI

	Database   db.Configuration         `yaml:"Database"`   // Database configuration
	WEBServers []WEBServerConfiguration `yaml:"WEBServers"` // Параметры web серверов
}

// ApplicationLog Конфигурация системы логирования
type ApplicationLog struct {
	Graylog        string `yaml:"Graylog"` // Описание доступа к graylog2 серверу для отправки логов по протоколу UDP GELF
	GraylogProto   string `yaml:"-"`       // Протокол подключения, tcp или udp, по умолчанию udp
	GraylogAddress string `yaml:"-"`       // Адрес graylog сервера
	GraylogPort    uint16 `yaml:"-"`       // Порт сервера, по умолчанию 12201
	GraylogEnable  bool   `yaml:"-"`       // =true - Включено логирование в graylog2
}
