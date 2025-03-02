package database

/*
The only purpose of this files is to properly Load data from disk [remote or local]
into memeory for this application to run properly this should
all be handling in this file alone. As of right now this is looking like a MongoDB implementation. just because its quick and easy.

Goes without saying but everything should rely on interfaces and not implementaion details
*/
// This is used to store application state, Each Service should have their own implemntaion of this
type Storage interface {
	Save() error // Save data to disk
	Load() error // Load data from disk into ram
}

// Example Template below
func newStorage() Storage {
	return &StorageSolution{}
}

type StorageSolution struct {
}

func (s *StorageSolution) Save() error { return nil }
func (s *StorageSolution) Load() error { return nil }
