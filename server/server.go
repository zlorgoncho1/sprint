package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	core "github.com/zlorgoncho1/sprint/core"
	logger "github.com/zlorgoncho1/sprint/logger"
	"github.com/zlorgoncho1/sprint/utils"

	"strings"
	"time"
)

type Server struct {
	Host      string
	Port      string
	routeTree core.EndpointNode
}

var __logger logger.Logger = logger.Logger{}

func (server Server) Start(mainModule core.Module) (net.Listener, error) {
	__logger.Log("Starting Sprint Application ...", "ServerCore")
	server.routeTree = server.routesResolver(mainModule.Controllers)
	startTime := time.Now()
	addr := server.Host + ":" + server.Port
	listener, err1 := net.Listen("tcp", addr)
	if err1 != nil {
		fmt.Println("Server never started")
	}
	endTime := time.Now()
	__logger.Plog("Sprint application successfully started", endTime.Sub(startTime), "ServerCore", "0", "OK")
	__logger.Log(fmt.Sprintf("Listenning on http://%s:%s", server.Host, server.Port), "ServerCore")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur lors de l'acceptation d'une connexion :", err)
			continue
		}
		go server.readBuffer(conn)
	}
}

func (server Server) routesResolver(controllers []core.Controller) core.EndpointNode {
	server.routeTree = core.EndpointNode{Level: 0, NextNodeMap: make(map[string]*core.EndpointNode)}
	for i := 0; i < len(controllers); i++ {
		controller := (controllers)[i]
		controllerName, controllerPath := controller.Init()
		controllerPath = strings.TrimPrefix(controllerPath, "/")
		__logger.Log(fmt.Sprintf("%s | %s", controllerName, controllerPath), "ControllerResolver")
		routes := controller.Routes()
		// REVOIT LES SLASH
		for o := 0; o < len(routes); o++ {
			startTime := time.Now()
			route := routes[o]
			route.Endpoint = strings.TrimPrefix(route.Endpoint, "/")
			route.Endpoint = fmt.Sprintf("%s/%s", controllerPath, route.Endpoint)
			server.addEndpoint(&server.routeTree, route)
			endTime := time.Now()
			__logger.Plog(fmt.Sprintf("Mapped %s, {{ %s }}", route.Method, route.Endpoint), endTime.Sub(startTime), "ViewResolver", "0", "OK")
		}
	}
	return server.routeTree
}

func (server Server) addEndpoint(node *core.EndpointNode, route core.Route) *core.EndpointNode {
	workingNode := node
	if node.Level == 0 {
		_, exists := node.NextNodeMap[route.Method]
		if exists {
			workingNode = node.NextNodeMap[route.Method]
		} else {
			workingNode = &core.EndpointNode{Endpoint: route.Method, Level: node.Level + 1, NextNodeMap: make(map[string]*core.EndpointNode)}
			node.NextNodeMap[route.Method] = workingNode
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

func (server Server) extractHeadData(head string) (string, string, string, map[string]string, []string, error) {
	headParts := strings.Split(head, "\n")

	requestLine := strings.Split(headParts[0], " ")
	if len(requestLine) < 3 {
		return "", "", "", nil, nil, errors.New("EN TETE HTTP INCORRECT")
	}

	method := requestLine[0]
	_endpoint := requestLine[1]
	protocol := requestLine[2]

	endpointParts := strings.SplitN(_endpoint, "?", 2)
	endpoint := strings.TrimPrefix(endpointParts[0], "/")

	var query []string
	if len(endpointParts) > 1 {
		query = strings.Split(endpointParts[1], "&")
	}

	headers := make(map[string]string)
	for _, unFormattedHeader := range headParts[1:] {
		headerParts := strings.SplitN(unFormattedHeader, ": ", 2)
		if len(headerParts) < 2 {
			return "", "", "", nil, nil, errors.New("header format is incorrect")
		}
		headers[headerParts[0]] = headerParts[1]
	}
	return method, endpoint, protocol, headers, query, nil
}

func (server Server) extractHTTPBufferData(data string) (core.Request, error) {
	formatedData := strings.ReplaceAll(data, "\x00", "")
	formatedData = strings.ReplaceAll(formatedData, "\r", "")
	parts := strings.Split(formatedData, "\n\n")
	if len(parts) == 0 {
		return core.Request{}, errors.New("FORMAT DU BUFFER INVALIDE")
	}
	head := parts[0]
	var body interface{}
	if len(parts) == 2 {
		body = parts[1]
	}
	method, endpoint, protocol, headers, query, err := server.extractHeadData(head)
	if err != nil {
		return core.Request{}, err
	}
	contentType, keyExists := headers["Content-Type"]
	if keyExists {
		if strings.HasPrefix(contentType, "text/plain") {
		} else if strings.HasPrefix(contentType, "application/json") {
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

func (server Server) readBuffer(conn net.Conn) {
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
	response := server.handleRequest(&server.routeTree, request)
	server.handleResponse(&conn, request.Headers["Accept"], request.Protocol, &response)
	endTime := time.Now()
	responseMessage := fmt.Sprintf("%s ==> %s - {{ %s }}", conn.RemoteAddr().String(), request.Method, request.Endpoint)
	__logger.Plog(responseMessage, endTime.Sub(startTime), "RequestHandler", "2", "OK")
}

func (server Server) handleRequest(node *core.EndpointNode, request core.Request) core.Response {
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

func (server Server) handleResponse(conn *net.Conn, acceptHeader string, protocol string, response *core.Response) {
	if response.ContentType == "text/html" && (acceptHeader == "text/html" || acceptHeader == "*/*") {
		utils.HandleHTML(response)
	} else if response.ContentType == "application/json" && (acceptHeader == "application/json" || acceptHeader == "*/*") {
		utils.HandleJSON(response)
	} else if response.ContentType == "text/plain" && (acceptHeader == "text/plain" || acceptHeader == "*/*") {
		utils.HandlePlainText(response)
	} else {
		response.ContentType = "text/plain"
	}
	if response.StatusCode == 0 {
		response.StatusCode = 200
	}
	if response.StatusText == "" {
		response.StatusText = "OK"
	}
	if len(response.Headers) == 0 {
		response.Headers = utils.GetDefaultHeader((response.Content).(string), response.ContentType)
	}
	responseStatus := utils.FormatStatusResponse(response.StatusCode, response.StatusText, protocol)
	headers := utils.DictToHTTPHeadersResponse(response.Headers)
	(*conn).Write(utils.FormatHTTPResponse(responseStatus, headers, (response.Content).(string)))
}
