package retry

import (
	"database/sql"
	"errors"
)

// IsTransientSQLError возвращает true для тех SQL-ошибок, которые можно попробовать ещё раз.
// Тут можно смотреть на конкретные коды Postgres (pq.Error.Code) или `errors.Is(err, sql.ErrConnDone)` и т.п.
func IsTransientSQLError(err error) bool {
	if errors.Is(err, sql.ErrConnDone) {
		return true
	}
	// ... ваша логика по коду ошибки драйвера ...
	return false
}

// HTTPError — пример интерфейса кастомной ошибки HTTP-клиента,
// из которого можно вытащить StatusCode().
type HTTPError interface {
	error
	StatusCode() int
}

// Is5xxHTTPError возвращает true, если это HTTP-ошибка 500–599.
func Is5xxHTTPError(err error) bool {
	if he, ok := err.(HTTPError); ok {
		code := he.StatusCode()
		return code >= 500 && code < 600
	}
	return false
}
