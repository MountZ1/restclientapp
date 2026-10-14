package types

type KeyValue struct {
	Key   string
	Value string
}

type RequestTab struct {
	ID        string
	Method    string
	Title     string
	Path      string
	URL       string
	Section   string
	Params    []KeyValue
	Headers   []KeyValue
	Body      string
	AuthType  string
	AuthToken string
	AuthUser  string
	AuthPass  string
}
