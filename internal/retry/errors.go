package retry

import (
	"errors"
	"fmt"
)

// ErrNonRetryable — «сигнальная» ошибка, означающая: не пытаться снова.
var ErrNonRetryable = errors.New("retry: non retryable error")

// MarkNonRetryable оборачивает вашу ошибку, чтобы ретраи сразу её прервал.
func MarkNonRetryable(err error) error {
	return fmt.Errorf("%w: %s", ErrNonRetryable, err.Error())
}

// IsNonRetryable проверяет, помечена ли ошибка как «non-retryable».
func IsNonRetryable(err error) bool {
	return errors.Is(err, ErrNonRetryable)
}
