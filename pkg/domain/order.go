package domain

type Order struct {
	ID    string
	Lines []*OrderLine
}
