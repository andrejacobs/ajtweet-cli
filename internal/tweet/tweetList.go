package tweet

type TweetList struct {
	tweets []Tweet
}

func (list *TweetList) Add(tweet *Tweet) error {
	//TODO: Check if the tweet already exists, if it does return an error
	return nil
}
