package tweet

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	list := TweetList{}

	if len(list.tweets) != 0 {
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

	if len(list.tweets) != 2 {
		t.Fatal("Expected 2 tweets to be added")
	}

	if list.tweets[0] != tw1 {
		t.Fatalf("Expected tweet Id %q at index 0", tw1.Id)
	}

	if list.tweets[1] != tw2 {
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
