package core

// Module represents a core module in Sprint, which can contain other modules, controllers, and routes.
type Module struct {
	Name        string        // Unique identifier for the module.
	Imports     []*Module     // Other modules that this module depends on.
	Exports     []*Module     // Sub-modules that this module provides to the outside world.
	Controllers []*Controller // Controllers associated with this module.
}

type Middleware func(Request, func(Request) Response) Response

// Controller handles incoming HTTP requests and routes them to their respective handler functions.
type Controller struct {
	Name   string   // Name of the controller.
	Path   string   // Base path to which this controller's routes will be appended.
	Routes []*Route // Routes defined for this controller.
}

// AddRoute is a method to add new routes to a Controller.
func (controller *Controller) AddRoute(method HttpMethod, endpoint string, handler func(request Request) Response) {
	if controller.Routes == nil {
		controller.Routes = []*Route{}
	}
	// Adds a new Route to the Controller's Routes slice.
	controller.Routes = append(controller.Routes, &Route{Endpoint: endpoint, Method: method, Function: handler})
}

// Route defines a single route, its method, endpoint, and the handler function.
type Route struct {
	Method   HttpMethod                     // HTTP method (GET, POST, etc.)
	Endpoint string                         // Endpoint path for the route.
	Function func(request Request) Response // Handler function to execute when the route is accessed.
}

// Request represents the HTTP request data received by the server.
type Request struct {
	Method   string            // HTTP method used for the request.
	Endpoint string            // Target endpoint of the request.
	Protocol string            // Protocol used for the request, e.g., HTTP, HTTPS.
	Params   map[string]string // URL parameters.
	Headers  map[string]string // HTTP headers.
	Query    []string          // Query parameters.
	Body     interface{}       // Request body.
	// TODO : ADD METADATA
}

// Response represents the structure of the HTTP response to be sent back to the client.
type Response struct {
	Content     interface{}       // Content of the response (could be a string, JSON, etc.)
	ContentType ContentType       // Type of the content, e.g., application/json.
	StatusCode  int               // HTTP status code, e.g., 200 (OK), 404 (Not Found), etc.
	StatusText  string            // Textual representation of the status code.
	Headers     map[string]string // Response headers.
}

// EndpointNode is a structure used in Sprint's internal routing mechanism to map
// endpoint strings to their corresponding handler functions.
type EndpointNode struct {
	Endpoint    string                         // Endpoint path.
	Function    func(request Request) Response // Handler function for the endpoint.
	DynamicNode *EndpointNode                  // Pointer to a node representing a dynamic segment in the route.
	NextNodeMap map[string]*EndpointNode       // Map of next possible nodes in the route tree.
	Level       int                            // Depth level of the node in the route tree.
}

// HttpMethod represents the type for various HTTP methods used in web requests.
type HttpMethod string

// Enumeration of HttpMethod. These constants define standard HTTP methods
// and provide a clear, type-safe way of using them throughout the code.
const (
	GET    HttpMethod = "GET"    // GET method for HTTP requests, typically used for retrieving data.
	POST   HttpMethod = "POST"   // POST method for HTTP requests, commonly used for submitting data.
	PUT    HttpMethod = "PUT"    // PUT method for HTTP requests, often used for updating or replacing resources.
	DELETE HttpMethod = "DELETE" // DELETE method for HTTP requests, used for deleting resources.
	PATCH  HttpMethod = "PATCH"  // PATCH method for HTTP requests, applied for partially updating resources.
)

// ContentType defines the MIME type of the content being sent or received in HTTP transactions.
type ContentType string

// Constants for various ContentType. These values are used to set the 'Content-Type'
// header in HTTP responses and to interpret the content type of HTTP requests.
const (
	HTML      ContentType = "text/html"        // HTML content type, used for sending HTML-formatted data.
	JSON      ContentType = "application/json" // JSON content type, used for sending JSON-formatted data.
	PLAINTEXT ContentType = "text/plain"       // PlainText content type, used for sending plain text data.
)
