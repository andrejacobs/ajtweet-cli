/*
Copyright © 2022 André Jacobs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

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
