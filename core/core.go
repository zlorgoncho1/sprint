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
	Headers  map[string]string
	Query    []string
	Body     interface{}
	Key      string
}

type Response struct {
	Content interface{}
}

type EndpointNode struct {
	Method      string
	Function    func(request Request) Response
	DynamicNode *EndpointNode
	NextNodeMap map[string]*EndpointNode
	Level       int
}
