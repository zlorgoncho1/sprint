package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/zlorgoncho1/sprint/core"
	"github.com/zlorgoncho1/sprint/logger"
	"github.com/zlorgoncho1/sprint/utils"

	"strings"
	"time"
)

// Server struct defines the basic properties of the server including Host, Port, and a route tree for routing.
type Server struct {
	Host        string
	Port        string
	routeTree   core.EndpointNode
	middlewares []core.Middleware
}

// __logger is a global logger instance, initialized to a default logger.
var __logger logger.Logger = logger.Logger{}

// Use method adds a middleware to the server's middleware chain.
func (server *Server) Use(middleware core.Middleware) {
	server.middlewares = append(server.middlewares, middleware)
}

// Start initiates the server to listen on the specified Host and Port.
// It resolves routes, logs server starting, listens for incoming connections, and spawns goroutines to handle each connection.
func (server *Server) Start(mainModule *core.Module) (net.Listener, error) {
	__logger.Log("Starting Sprint Application ...", "ServerCore")

	// Resolve routes from the provided controllers in the mainModule.
	server.routeTree = server.routesResolver(mainModule.Controllers)

	// Record the start time for performance logging.
	startTime := time.Now()

	// Combine Host and Port to form the address and start listening on that TCP address.
	addr := server.Host + ":" + server.Port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		__logger.Error(fmt.Sprintf("Error starting server: %v", err), "ServerCore")
		return nil, err // Return error immediately after logging the failure
	}

	// Log the server startup time.
	endTime := time.Now()
	__logger.Plog("Sprint application successfully started", endTime.Sub(startTime), "ServerCore", "0", "OK")

	// Log the listening address.
	__logger.Log(fmt.Sprintf("Listening on http://%s:%s", server.Host, server.Port), "ServerCore")

	// Accept incoming connections in an infinite loop.
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Log the error if a connection cannot be accepted and continue listening.
			__logger.Error(fmt.Sprintf("Error during connection acceptance: %v", err), "ServerCore")
			continue
		}
		// Handle each connection in a separate goroutine for concurrent processing.
		go server.readBuffer(conn)
	}
}

func (server *Server) routesResolver(controllers []*core.Controller) core.EndpointNode {
	// Initialize the server's route tree.
	server.routeTree = core.EndpointNode{Level: 0, NextNodeMap: make(map[string]*core.EndpointNode)}

	for _, controller := range controllers {
		__logger.Log(fmt.Sprintf("%s | %s", controller.Name, controller.Path), "ControllerResolver")

		for _, route := range controller.Routes {
			startTime := time.Now()

			// Concatenate module, controller, and route paths.
			fullPath := utils.JoinPaths(controller.Path, route.Endpoint)
			route.Endpoint = fullPath

			// Add the route to the server's routing tree.
			server.addEndpoint(&server.routeTree, route)

			endTime := time.Now()
			__logger.Plog(fmt.Sprintf("Mapped %s, {{ %s }}", route.Method, fullPath), endTime.Sub(startTime), "ViewResolver", "0", "OK")
		}
	}

	return server.routeTree
}

func (server *Server) addEndpoint(node *core.EndpointNode, route *core.Route) *core.EndpointNode {
	workingNode := node
	if node.Level == 0 {
		method := string(route.Method)
		_, exists := node.NextNodeMap[method]
		if exists {
			workingNode = node.NextNodeMap[method]
		} else {
			workingNode = &core.EndpointNode{Endpoint: method, Level: node.Level + 1, NextNodeMap: make(map[string]*core.EndpointNode)}
			node.NextNodeMap[method] = workingNode
		}
	}

	routeSplited := strings.Split(route.Endpoint, "/")
	numberOfSubPath := len(routeSplited)
	if numberOfSubPath-workingNode.Level >= 0 {
		path := routeSplited[workingNode.Level-1]
		existingNode, exists := workingNode.NextNodeMap[path]
		if exists {
			if numberOfSubPath == 0 {
				return existingNode
			}
			return server.addEndpoint(existingNode, route)
		} else {
			newNode := &core.EndpointNode{Endpoint: path, Level: workingNode.Level + 1, NextNodeMap: make(map[string]*core.EndpointNode), Function: route.Function}
			if strings.HasPrefix(path, ":") {
				workingNode.DynamicNode = newNode
			} else {
				workingNode.NextNodeMap[path] = newNode
			}
			if numberOfSubPath == 0 {
				return newNode
			}
			return server.addEndpoint(newNode, route)
		}
	}
	return workingNode
}

func (server *Server) extractHeadData(head string) (string, string, string, map[string]string, []string, error) {
	headParts := strings.Split(head, "\n")

	// Ensure there is at least one line for the request line
	if len(headParts) == 0 {
		return "", "", "", nil, nil, errors.New("empty HTTP head")
	}

	requestLine := strings.Fields(headParts[0]) // Fields automatically trims spaces and splits
	if len(requestLine) < 3 {
		return "", "", "", nil, nil, errors.New("invalid HTTP request line")
	}

	method := requestLine[0]
	_endpoint := requestLine[1]
	protocol := requestLine[2]

	// Endpoint processing
	endpointParts := strings.SplitN(_endpoint, "?", 2)
	endpoint := strings.TrimPrefix(endpointParts[0], "/")

	var query []string
	if len(endpointParts) > 1 {
		query = strings.Split(endpointParts[1], "&")
	}

	// Headers processing
	headers := make(map[string]string)
	for _, unFormattedHeader := range headParts[1:] {
		headerParts := strings.SplitN(strings.TrimSpace(unFormattedHeader), ":", 2)
		if len(headerParts) != 2 {
			continue // This skips malformed headers
		}
		headers[strings.TrimSpace(headerParts[0])] = strings.TrimSpace(headerParts[1])
	}
	return method, endpoint, protocol, headers, query, nil
}

func (server *Server) extractHTTPBufferData(data string) (core.Request, error) {
	formattedData := strings.ReplaceAll(data, "\x00", "")
	formattedData = strings.ReplaceAll(formattedData, "\r", "")
	parts := strings.SplitN(formattedData, "\n\n", 2) // Use SplitN to ensure only one split at the first occurrence

	if len(parts) < 1 {
		return core.Request{}, errors.New("invalid buffer format: no header found")
	}
	head := parts[0]
	var body interface{}
	if len(parts) > 1 {
		body = parts[1]
	}

	method, endpoint, protocol, headers, query, err := server.extractHeadData(head)
	if err != nil {
		return core.Request{}, err
	}
	contentType, hasContentType := headers["Content-Type"]
	if hasContentType {
		if strings.HasPrefix(contentType, string(core.PLAINTEXT)) || strings.HasPrefix(contentType, string(core.HTML)) {
		} else if strings.HasPrefix(contentType, string(core.JSON)) {
			var jsonObj interface{}
			err := json.Unmarshal([]byte(body.(string)), &jsonObj)
			if err != nil {
				return core.Request{}, err
			}
			body = jsonObj
		} else {
			return core.Request{}, errors.New("ContentTypeException")
		}
	}
	return core.Request{Method: method, Endpoint: endpoint, Protocol: protocol, Headers: headers, Query: query, Body: body, Params: make(map[string]string)}, nil
}

func (server *Server) applyMiddlewares(req core.Request, finalHandler func(core.Request) core.Response) core.Response {
	var execMiddleware func(int, core.Request) core.Response
	execMiddleware = func(index int, req core.Request) core.Response {
		if index < len(server.middlewares) {
			return server.middlewares[index](req, func(req core.Request) core.Response {
				return execMiddleware(index+1, req)
			})
		}
		return finalHandler(req)
	}
	return execMiddleware(0, req)
}

func (server *Server) readBuffer(conn net.Conn) {
	startTime := time.Now()
	defer conn.Close()
	var req []byte
	keep_loop := true
	for keep_loop {
		buffer := make([]byte, 1024)
		conn.Read(buffer)
		if buffer[len(buffer)-1] == 0 {
			keep_loop = false
		}
		req = append(req, buffer...)
	}
	msg := ""
	for _, ch := range req {
		msg += string(ch)
	}
	request, err := server.extractHTTPBufferData(msg)
	if err != nil {
		__logger.Error(string(err.Error()), "ServerCore")
	}

	// Appliquer les middlewares
	response := server.applyMiddlewares(request, func(req core.Request) core.Response {
		return server.handleRequest(&server.routeTree, req)
	})

	server.handleResponse(&conn, request.Headers["Accept"], request.Protocol, &response)
	endTime := time.Now()
	responseMessage := fmt.Sprintf("%s ==> %s - {{ %s }}", conn.RemoteAddr().String(), request.Method, request.Endpoint)
	__logger.Plog(responseMessage, endTime.Sub(startTime), "RequestHandler", "2", "OK")
}

func (server *Server) handleRequest(node *core.EndpointNode, request core.Request) core.Response {
	workingNode := node
	var exists bool
	if workingNode.Level == 0 {
		workingNode, exists = workingNode.NextNodeMap[request.Method]
		if !exists {
			return core.Response{}
		}
	}
	routeSplited := strings.Split(request.Endpoint, "/")
	numberOfSubPath := len(routeSplited)
	if numberOfSubPath-workingNode.Level >= 0 {
		path := routeSplited[workingNode.Level-1]
		existingNode, exists := workingNode.NextNodeMap[path]
		if !exists {
			if workingNode.DynamicNode == nil {
				return core.Response{}
			}
			existingNode = workingNode.DynamicNode
			request.Params[strings.TrimPrefix(existingNode.Endpoint, ":")] = path
		}
		if numberOfSubPath-workingNode.Level == 0 {
			return existingNode.Function(request)
		}
		return server.handleRequest(existingNode, request)
	}
	return core.Response{}
}

func (server *Server) handleResponse(conn *net.Conn, acceptHeader string, protocol string, response *core.Response) {
	acceptTypes := strings.Split(acceptHeader, ",")
	// Assuming "*/*" or matching ContentType is acceptable
	isAcceptableType := func(content core.ContentType) bool {
		for _, t := range acceptTypes {
			if t == string(content) || t == "*/*" {
				return true
			}
		}
		return false
	}

	switch {
	case response.ContentType == core.HTML && isAcceptableType(core.HTML):
		utils.HandleHTML(response)
	case response.ContentType == core.JSON && isAcceptableType(core.JSON):
		utils.HandleJSON(response)
	case response.ContentType == core.PLAINTEXT && isAcceptableType(core.PLAINTEXT):
		utils.HandlePlainText(response)
	default:
		response.ContentType = core.PLAINTEXT
		// Consider setting a more appropriate status code and message for unsupported content types.
	}

	if response.StatusCode == 0 {
		response.StatusCode = 200
	}
	if response.StatusText == "" {
		response.StatusText = "OK"
	}
	if len(response.Headers) == 0 {
		contentString := server.FormatContentString(response.Content)
		response.Headers = utils.GetDefaultHeader(contentString, response.ContentType)
	}
	responseStatus := utils.FormatStatusResponse(response.StatusCode, response.StatusText, protocol)
	headers := utils.DictToHTTPHeadersResponse(response.Headers)

	if _, err := (*conn).Write(utils.FormatHTTPResponse(responseStatus, headers, server.FormatContentString(response.Content))); err != nil {
		// Log or handle the error based on your application's requirements
		__logger.Error(fmt.Sprintf("Error writing response: %s", err), "ServerCore")
	}
}

func (server *Server) FormatContentString(content interface{}) string {
	var contentString string
	var err error
	switch v := content.(type) {
	case string:
		contentString = v
	default:
		var jsonData []byte
		jsonData, err = json.Marshal(v)
		if err != nil {
			log.Fatalf("Erreur lors de la sérialisation en JSON : %v", err)
		}
		contentString = string(jsonData)
	}
	if err != nil {
		panic(err)
	}
	return contentString
}

func CORSMiddleware(req core.Request, next func(core.Request) core.Response) core.Response {
	req.Headers["Access-Control-Allow-Origin"] = "*"
	req.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
	req.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization"

	if req.Method == "OPTIONS" {
		return core.Response{
			StatusCode:  200,
			ContentType: core.PLAINTEXT,
			Content:     "",
			StatusText:  "OK",
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		}
	}

	response := next(req)

	fmt.Println("Response:", response, req.Body)

	// Initialiser la map Headers si elle est nil
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers["Access-Control-Allow-Origin"] = "*"
	response.Headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE, OPTIONS"
	response.Headers["Access-Control-Allow-Headers"] = "Content-Type, Authorization"

	return response
}
