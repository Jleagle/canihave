package models

type Relation struct {
	ID          string
	RelatedID   string
	DateCreated string
	Type        string
}

const (
	RELATION_TYPE_SIMILAR   string = "similar"
	RELATION_TYPE_SAME_NODE string = "node"
)
