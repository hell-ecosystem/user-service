package retry

import "time"

type Option func(*Retrier)

// New создаёт Retrier с дефолтными значениями и применяет все Option.
func New(opts ...Option) *Retrier {
	r := &Retrier{
		Attempts:     1,
		Backoff:      &ExponentialBackoff{Initial: 100 * time.Millisecond, Max: 2 * time.Second, Multiplier: 2.0},
		JitterFactor: 0,
	}
	for _, o := range opts {
		o(r)
	}
	return r
}

// WithMaxAttempts задаёт, сколько всего попыток сделать.
func WithMaxAttempts(n int) Option {
	return func(r *Retrier) {
		r.Attempts = n
	}
}

// WithBackoffExponential задаёт параметры экспоненциального бэкоффа.
func WithBackoffExponential(initial time.Duration, multiplier float64) Option {
	return func(r *Retrier) {
		r.Backoff = &ExponentialBackoff{
			Initial:    initial,
			Multiplier: multiplier,
			Max:        r.Backoff.Max, // можно тоже сделать опцией, если нужно
		}
	}
}

// WithJitter задаёт коэффициент джиттера ±factor (0.1 → ±10%).
func WithJitter(factor float64) Option {
	return func(r *Retrier) {
		r.JitterFactor = factor
	}
}

// RetryIf задаёт пользовательский фильтр — возвращает true, если по этой ошибке стоит повторить.
func RetryIf(fn func(error) bool) Option {
	return func(r *Retrier) {
		r.RetryIf = fn
	}
}
