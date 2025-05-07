package retry

import (
	"math"
	"math/rand"
	"time"
)

// ExponentialBackoff хранит параметры для классического экспоненциального роста.
type ExponentialBackoff struct {
	Initial    time.Duration // начальная задержка, например 100ms
	Max        time.Duration // максимум, например 2s
	Multiplier float64       // обычно 2.0
}

// Next возвращает задержку для given попытки (номер attempt, начиная с 1),
// с учётом экспоненты и ограничения Max.
func (b *ExponentialBackoff) Next(attempt int) time.Duration {
	// d = Initial * Multiplier^(attempt-1)
	d := float64(b.Initial) * math.Pow(b.Multiplier, float64(attempt-1))
	if time.Duration(d) > b.Max {
		return b.Max
	}
	return time.Duration(d)
}

// WithJitter добавляет ±jitterFactor (например 0.2) разброс к базовой задержке.
func WithJitter(base time.Duration, jitterFactor float64) time.Duration {
	// случайно в диапазоне [1-jitterFactor, 1+jitterFactor]
	delta := (rand.Float64()*2 - 1) * jitterFactor
	return time.Duration(float64(base) * (1 + delta))
}
