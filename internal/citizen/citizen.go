package citizen

type User struct {
	ID           string
	Email        string
	HashPassword string
}

type Address struct {
	Log, Lat float32
}

type Citizen struct {
	*User
	Address *Address
}
