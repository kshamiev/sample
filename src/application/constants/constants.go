package constants // import "application/constants"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "time"

const (
	// EmailConfirmationCodeLength Длинна кода подтверждения отправляемоего на электронный адрес
	EmailConfirmationCodeLength = uint8(12)

	// PhoneConfirmationCodeLength Длинна кода подтверждения отправляемоего на мобильный телефон
	PhoneConfirmationCodeLength = uint8(6)

	// AccountInactiveLifetimeDefault Время жизни аккаунта при отсутствии активности по умолчанию
	AccountInactiveLifetimeDefault = TypeConstantUint64(time.Hour * 24 * 366 * 1)

	// TimeoutFilestoreMarkDeleted Таймаут пометки временных файлов файлового хранилища как удалённых
	TimeoutFilestoreMarkDeleted = TypeConstantUint64(time.Hour)

	// TimeoutFilestoreCleanup Таймаут очистки файлового хранилища от временных файлов
	TimeoutFilestoreCleanup = TypeConstantUint64(time.Hour * 24)
)
