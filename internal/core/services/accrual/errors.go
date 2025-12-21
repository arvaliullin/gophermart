package accrual

import (
	"fmt"
	"time"
)

// RetryAfterError возвращается при превышении лимита запросов.
type RetryAfterError struct {
	Duration time.Duration
}

func (e *RetryAfterError) Error() string {
	return fmt.Sprintf("превышен лимит запросов, повторить через %v", e.Duration)
}

