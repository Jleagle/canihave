package sites

const BoughtItOnceID = "boughtitonce"

type Boughtitonce struct {
}

func (b Boughtitonce) GetID() string {
	return BoughtItOnceID
}

func (b Boughtitonce) GetLatest() Result {
	panic("implement me")
}
