package types

type QueryRequest struct {
	Key   string
	Value any
}

type RequestURL struct {
	Raw   string
	Host  []string
	Path  []string
	Query []QueryRequest
}

type RequestAuth struct {
	Type     string
	Token    string
	Username string
	Password string
}

type RequestBody struct {
	Type    string
	Content string
}

type RequestHttp struct {
	Method string
	Header []string
	URL    RequestURL
	Auth   RequestAuth
	Body   RequestBody
}

type Request struct {
	Name     string
	Request  RequestHttp
	Response []string
}
