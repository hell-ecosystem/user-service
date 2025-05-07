package retry

import (
	"context"
	"time"
)

// Retrier объединяет backoff-стратегию и политику фильтрации ошибок.
type Retrier struct {
	Attempts     int // сколько раз пробовать
	Backoff      *ExponentialBackoff
	JitterFactor float64          // например, 0.2 для ±20%
	RetryIf      func(error) bool // кастомный фильтр: true — можно ретраить
}

// Do пытается выполнить fn до тех пор, пока fn не вернёт nil
// или пока не исчерпаются Attempts, пропуская non-retryable ошибки.
func (r *Retrier) Do(ctx context.Context, fn func() error) error {
	var lastErr error

	for i := 1; i <= r.Attempts; i++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		// если явно помечено как non-retryable — выходим сразу
		if IsNonRetryable(err) {
			return err
		}
		// если внешний фильтр запрещает ретраи — тоже выходим
		if r.RetryIf != nil && !r.RetryIf(err) {
			return err
		}

		// если это последняя попытка — выходим с ошибкой
		if i == r.Attempts {
			break
		}

		// считаем задержку: экспонента + джиттер
		delay := r.Backoff.Next(i)
		delay = WithJitter(delay, r.JitterFactor)

		select {
		case <-time.After(delay):
			// следующая попытка
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}
