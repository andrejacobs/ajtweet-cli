package app

type Config struct {
	Datastore Datastore
}

type Datastore struct {
	Filepath string
}
