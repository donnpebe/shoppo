package domain

import "time"

type Promotion struct {
	StartDate time.Time
	EndDate   time.Time
	Condition PromotionCondition
}
