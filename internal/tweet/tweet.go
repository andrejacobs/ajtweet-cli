package tweet

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	Id            uuid.UUID
	Message       string
	ScheduledTime time.Time
}

func New(message string, scheduledTime time.Time) Tweet {
	tweet := Tweet{
		Id:            uuid.New(),
		Message:       message,
		ScheduledTime: scheduledTime,
	}
	return tweet
}

func (tweet Tweet) SendNow() bool {
	return tweet.ScheduledTime.Before(time.Now())
}

// Stringer implementation
func (tweet Tweet) String() string {
	return fmt.Sprintf("id: %s, time: %s, tweet: %s", tweet.Id, tweet.ScheduledTime, tweet.Message)
}
