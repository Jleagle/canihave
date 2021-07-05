package sites

type Site interface {
	GetID() string
	GetLatest() Result
}

var Sites = []Site{
	&Boughtitonce{},
	&Canopy{},
	&Dattwenty{},
	&Shityoucanafford{},
	&Wannaspend{},
}

type Result struct {
}
