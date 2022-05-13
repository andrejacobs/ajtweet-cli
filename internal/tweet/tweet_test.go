package tweet

import (
	"testing"
	"time"
)

func TestUniqueId(t *testing.T) {
	message := "The quick brown fox jumped over the lazy dog!"
	scheduledTime := time.Now()
	tw1 := New(message, scheduledTime)
	tw2 := New(message, scheduledTime)

	if tw1.Id == tw2.Id {
		t.Fatalf("Tweets should have unique identifiers. Result Id = %q", &tw1.Id)
	}
}

func TestSendNow(t *testing.T) {
	testCases := []struct {
		name          string
		scheduledTime time.Time
		expected      bool
	}{
		{"Send tomorrow", time.Now().AddDate(0, 0, 1), false},
		{"Send later", time.Now().Add(5 * time.Minute), false},
		{"Send now", time.Now(), true},
		{"Really send now!", time.Now().AddDate(-1, 0, 0), true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tw := New(tc.name, tc.scheduledTime)
			if result := tw.SendNow(); result != tc.expected {
				t.Fatalf("Expected: %t, Result: %t, Tweet: %s", tc.expected, result, tw)
			}
		})
	}
}
