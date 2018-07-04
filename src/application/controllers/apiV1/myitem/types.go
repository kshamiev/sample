package myitem

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"net/http"
	"time"
)

const (
	retryAfter = time.Second * 30
)

// Interface is an interface of package
type Interface interface {
	// Status Получение состояние данных в БД
	Status(wr http.ResponseWriter, rq *http.Request)

	// Create Сохранение в БД данных
	Create(wr http.ResponseWriter, rq *http.Request)

	// Load Получение данных из БД по ID
	Load(wr http.ResponseWriter, rq *http.Request)

	// Delete Удаление данных в БД по ID
	Delete(wr http.ResponseWriter, rq *http.Request)
}

// impl is an implementation of package
type impl struct {
}

// StatusResponse Исходящие данные метода Status
type StatusResponse struct {
	Сount uint64   `json:"count"   cpy:"name=Size"` // Количество записей в БД
	Ids   []uint64 `json:"ids"`                     // Массив идентификаторов сущностей базы данных
}

// CreateRequest Входящие данные метода Create
type CreateRequest struct {
	Date   time.Time `json:"date"      validate:"required"`               // Любая дата и время
	Number int64     `json:"number"    validate:"required,ne=0,max=1024"` // Любое положительное или отрицательное число кроме 0
	Text   string    `json:"text"      validate:"required"`               // Любая не пустая строка
}

// CreateResponse Исходящие данные метода Create
type CreateResponse struct {
	ID uint64 `json:"id"` // Уникальный идентификатор сущности
}

// LoadResponse Исходящие данные метода Load
type LoadResponse struct {
	ID     uint64    `json:"id"`     // Уникальный идентификатор сущности
	Date   time.Time `json:"date"`   // Любая дата и время
	Number int64     `json:"number"` // Любое число
	Text   string    `json:"text"`   // Любая строка
}
