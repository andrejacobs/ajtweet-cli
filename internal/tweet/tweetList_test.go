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
