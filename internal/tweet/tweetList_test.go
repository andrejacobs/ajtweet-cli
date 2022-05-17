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
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	list := TweetList{}

	if len(list.Tweets) != 0 {
		t.Fatal("TweetList must be empty on initialization")
	}

	tw1 := New("Tweet1", time.Now())
	tw2 := New("Tweet2", time.Now())

	if err := list.Add(tw1); err != nil {
		t.Fatal(err)
	}
	if err := list.Add(tw2); err != nil {
		t.Fatal(err)
	}

	if len(list.Tweets) != 2 {
		t.Fatal("Expected 2 tweets to be added")
	}

	if list.Tweets[0] != tw1 {
		t.Fatalf("Expected tweet Id %q at index 0", tw1.Id)
	}

	if list.Tweets[1] != tw2 {
		t.Fatalf("Expected tweet Id %q at index 1", tw2.Id)
	}

}

func TestAddingTwiceIsInvalid(t *testing.T) {
	list := TweetList{}

	tw := New("Tweet1", time.Now())

	if err := list.Add(tw); err != nil {
		t.Fatal(err)
	}

	expectedErr := fmt.Errorf("%w: %q", ErrExists, tw.Id)
	if err := list.Add(tw); errors.Is(err, expectedErr) {
		t.Fatalf("Expected error: %q, Result: %q", expectedErr, err)
	}
}

func TestFind(t *testing.T) {
	list := TweetList{}
	tw1 := New("Tweet1", time.Now())
	tw2 := New("Tweet2", time.Now())

	found, index := list.Find(tw1.Id)
	if found != false || index != -1 {
		t.Fatal("Tweet should not have been found in the list")
	}

	list.Add(tw1)
	list.Add(tw2)

	found, index = list.Find(tw1.Id)
	if found != true || index != 0 {
		t.Fatalf("Tweet should have been found in the list at index 0. Result: %t, %d", found, index)
	}

	found, index = list.Find(tw2.Id)
	if found != true || index != 1 {
		t.Fatalf("Tweet should have been found in the list at index 1. Result: %t, %d", found, index)
	}

}

func TestListIsSorted(t *testing.T) {
	list := TweetList{}
	tw1 := New("Tweet1", time.Now())
	tw2 := New("Tweet2", time.Now().Add(-5*time.Second))
	tw3 := New("Tweet3", time.Now().Add(5*time.Second))
	expectedOrder := []Tweet{tw2, tw1, tw3}

	list.Add(tw1)
	list.Add(tw2)
	list.Add(tw3)

	result := list.List()
	if count := len(result); count < 3 {
		t.Fatalf("Expected %d tweets. Result: %d", 3, count)
	}
	// Check the original list has not been modified
	if list.Tweets[0].Id != tw1.Id {
		t.Fatalf("Expected the original list to remain unchanged")
	}

	// Check expected order
	for i, tw := range result {
		if tw.Id != expectedOrder[i].Id {
			t.Errorf("Expected %q. Result %q", expectedOrder[i].Message, tw.Message)
		}
	}
}

func TestDelete(t *testing.T) {
	list := TweetList{}
	tw1 := New("Tweet1", time.Now())
	tw2 := New("Tweet2", time.Now().Add(-5*time.Second))
	tw3 := New("Tweet3", time.Now().Add(5*time.Second))

	list.Add(tw1)
	list.Add(tw2)
	list.Add(tw3)

	if count := len(list.Tweets); count < 3 {
		t.Fatalf("Expected %d tweets. Result: %d", 3, count)
	}

	if err := list.Delete(tw1.Id); err != nil {
		t.Fatal(err)
	}
	if err := list.Delete(tw3.Id); err != nil {
		t.Fatal(err)
	}
	if count := len(list.Tweets); count < 1 {
		t.Fatalf("Expected %d tweet. Result: %d", 1, count)
	}

	if err := list.Delete(tw2.Id); err != nil {
		t.Fatal(err)
	}
	if count := len(list.Tweets); count != 0 {
		t.Fatalf("Expected %d tweet. Result: %d", 0, count)
	}

	expectedErr := fmt.Errorf("%w: %q", ErrNotExists, tw1.Id)
	if err := list.Delete(tw1.Id); errors.Is(err, expectedErr) {
		t.Fatalf("Expected an error of type ErrNotExists, instead result is: %s", err)
	}
}

func TestDeleteAll(t *testing.T) {
	list := TweetList{}
	tw1 := New("Tweet1", time.Now())
	tw2 := New("Tweet2", time.Now().Add(-5*time.Second))
	tw3 := New("Tweet3", time.Now().Add(5*time.Second))

	list.Add(tw1)
	list.Add(tw2)
	list.Add(tw3)

	if err := list.DeleteAll(); err != nil {
		t.Fatal(err)
	}

	if len(list.Tweets) != 0 {
		t.Fatal("Expected all tweets to have been deleted")
	}
}

func TestSaveLoad(t *testing.T) {
	l1 := TweetList{}
	l2 := TweetList{}

	l1.Add(New("Tweet1", time.Now()))
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)

	}
	defer os.Remove(tempFile.Name())

	if err := l1.Save(tempFile.Name()); err != nil {
		t.Fatalf("Error saving list to file: %s", err)
	}

	if err := l2.Load(tempFile.Name()); err != nil {
		t.Fatalf("Error getting list from file: %s", err)
	}

	if l1.Tweets[0].Id != l2.Tweets[0].Id {
		t.Errorf("Tweet %q should match %q.", l1.Tweets[0], l2.Tweets[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}

	if err := os.Remove(tempFile.Name()); err != nil {
		t.Fatalf("Error deleting temp file: %s", err)
	}

	list := TweetList{}
	if err := list.Load(tempFile.Name()); err != nil {
		t.Errorf("Not expecting an error. Result: %q", err)
	}
}
