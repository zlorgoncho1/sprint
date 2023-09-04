package core

type Module struct {
	Name        string
	Imports     []Module
	Exports     []Module
	Controllers []Controller
}

type Controller interface {
	Routes() []Route
	Init() (string, string)
}

type Route struct {
	Method   string
	Endpoint string
	Function func(request Request) Response
}

type Request struct {
	Method   string
	Endpoint string
	Protocol string
	Params   map[string]string
	Headers  map[string]string
	Query    []string
	Body     interface{}
}

type Response struct {
	Content     interface{}
	ContentType string
	StatusCode  int
	StatusText  string
	Headers     map[string]string
}

type EndpointNode struct {
	Endpoint    string
	Function    func(request Request) Response
	DynamicNode *EndpointNode
	NextNodeMap map[string]*EndpointNode
	Level       int
}

const (
	ContentTypeHTML  = "text/html"
	ContentTypeJSON  = "application/json"
	ContentTypePlain = "text/plain"
)
