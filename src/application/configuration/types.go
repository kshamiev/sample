package configuration // import "application/configuration"

import (
	webserverTypes "application/webserver/types"

	"gopkg.in/webnice/kit.v1/modules/db"
)

// Application Описание основной структуры конфигурационного файла
type Application struct {
	ApplicationName  string `yaml:"ApplicationName"`  // Название приложения
	WorkingDirectory string `yaml:"WorkingDirectory"` // Рабочая директория сервера, сразу после запуска приложение меняет текущую директорию
	PidFile          string `yaml:"PidFile"`          // PID file файл, в который записывается текущий идентификатор процесса
	TempPath         string `yaml:"TempPath"`         // Путь к временной директории сервера
	CachePath        string `yaml:"CachePath"`        // Путь к директории сервера хранения кеша
	LogConfiguration string `yaml:"LogConfiguration"` // Путь и имя файла конфигурации системы логирования (ApplicationLog)
	LogPath          string `yaml:"LogPath"`          // Путь к папке размещения лог файлов
	StateFile        string `yaml:"StateFile"`        // Файл состояния
	SocketFile       string `yaml:"SocketFile"`       // Сокет файл для коммуникации с CLI

	Database   db.Configuration               `yaml:"Database"`   // Database configuration
	WEBServers []webserverTypes.Configuration `yaml:"WEBServers"` // Параметры web серверов
	Storage    string                         `yaml:"Storage"`    // Путь к хранилищу файлов
}

// ApplicationLog Конфигурация системы логирования
type ApplicationLog struct {
	Graylog        string `yaml:"Graylog"` // Описание доступа к graylog2 серверу для отправки логов по протоколу UDP GELF
	GraylogProto   string `yaml:"-"`       // Протокол подключения, tcp или udp, по умолчанию udp
	GraylogAddress string `yaml:"-"`       // Адрес graylog сервера
	GraylogPort    uint16 `yaml:"-"`       // Порт сервера, по умолчанию 12201
	GraylogEnable  bool   `yaml:"-"`       // =true - Включено логирование в graylog2
}
