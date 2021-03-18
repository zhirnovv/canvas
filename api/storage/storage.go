package storage

// Storage is an interface describing a basic storage solution.
type Storage interface {
	Create() 
	Read() 
	Update() 
	Delete() 
}
