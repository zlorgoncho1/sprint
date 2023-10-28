### Overview

#### Introduction

Sprint is a modular, high-performance web framework for Go, inspired by the architecture of NestJS. Designed to streamline backend development, Sprint empowers developers to build scalable and maintainable applications with ease. It emphasizes a modular structure, where functionality is divided into controllers and modules, allowing for clear organization and flexibility in development.

#### Key Features

- **Modular Architecture**: Like NestJS, Sprint organizes code into modules, making it easier to manage and scale large applications.
- **Intuitive Controller Setup**: Simplified creation of controllers to handle various routes and requests.
- **Simplicity and Performance**: Leverages Go's efficiency and simplicity, providing a framework that is both easy to use and high-performing.

### Getting Started

#### Setting Up Your First Sprint Application

This guide will help you create a basic "Hello World" API using Sprint.

##### Prerequisites

- Go installed on your machine (see Go's [official documentation](https://golang.org/doc/install) for installation guidance).

##### Step 1: Creating the Project

Start by setting up a new Go project:

```bash
mkdir app && cd app
go mod init app
```

##### Step 2: Install Sprint

Assuming Sprint is available as a Go module, install it using:

```bash
go get github.com/zlorgoncho1/sprint
```

##### Step 3: Project Structure

Organize your project with the following structure:

```
📂 app
┣━━ 📂 hello
┃   ┣━━ 📄 controllers.go
┃   ┗━━ 📄 module.go
┣━━ 📄 go.mod
┣━━ 📄 go.sum
┣━━ 📄 main.go
```

##### Step 4: Writing Code

###### main.go

Set up the main entry point of your application:

```go
package main

import (
	"app/hello"

	"github.com/zlorgoncho1/sprint/server"
)

func main() {
	server := &server.Server{Host: "localhost", Port: "8000"}
	server.Start(hello.HelloModule())
}
```

###### module.go

Define your application module:

```go
package hello

import (
	"github.com/zlorgoncho1/sprint/core"
)

// Module Definition
func HelloModule() *core.Module {
	return &core.Module{
		Name: "HelloModule",
		Controllers: []*core.Controller{
			HelloController(),
		},
	}

}
```

###### controllers.go

Create a basic controller:

```go
package hello

import (
	"github.com/zlorgoncho1/sprint/core"
)

func HelloController() *core.Controller {
	var HelloController = &core.Controller{Name: "HelloController", Path: "hello"}
	HelloController.AddRoute(core.GET, ":name", hello)
	HelloController.AddRoute(core.GET, "JSON/:name", helloJSON)
	HelloController.AddRoute(core.GET, "HTML/:name", helloHTML)
	return HelloController
}

func hello(request core.Request) core.Response {
	return core.Response{Content: "Bonjour " + request.Params["name"]}
}

func helloHTML(request core.Request) core.Response {
	name := request.Params["name"]
	return core.Response{Content: "<h1>Bonjour " + name + "</h1>", ContentType: core.HTML}
}

func helloJSON(request core.Request) core.Response {
	type greetingObj struct {
		Nom     string `json:"nom"`
		Message string `json:"message"`
	}
	name := request.Params["name"]
	return core.Response{Content: greetingObj{name, "Bonjour"}, ContentType: core.JSON}
}
```

##### Step 5: Running Your Application

Run your application:

```bash
go run main.go
```

Your Sprint application should now be running on `localhost:8000`. Visiting `localhost:8000/hello/HTML/Sprint` in your browser or through a tool like `curl` should return a "Hello World" message.

#### Conclusion

This guide covers the basic setup and a simple "Hello World" example in Sprint. As Sprint evolves, more advanced features and documentation will be made available, catering to more complex application needs and empowering developers to fully leverage the Go language in web development.

### For Contributors

Our project is a custom-built web server framework, somewhat reminiscent of popular frameworks like Flask in Python or Express in Node.js, but tailored for use in the Go language environment. This framework aims to provide an easy-to-use, yet efficient way to handle HTTP requests and responses, deal with routing, manage logging, and perform common utility tasks.

### Main Components

#### 1. Core (`core.go`)

The core of our framework lies within `core.go`, which defines the primary structs and interfaces used throughout the application. This includes definitions for:

- **Modules**: Conceptually similar to controllers in MVC frameworks, they're used to organize the handling of different routes and their associated functionalities.
- **EndpointNode**: A fundamental part of the routing mechanism, these nodes represent individual endpoints in a tree-like structure for efficient route resolution.
- **Response and Request**: Custom types to handle incoming requests and outgoing responses, providing an abstraction over Go's native HTTP handling facilities.

#### 2. Utils (`utils.go`)

This utility package (`utils.go`) provides a set of handy functions for converting dictionaries to JSON strings, formatting HTTP responses, and generating standard headers. It abstracts some repetitive tasks, such as:

- Converting map objects to JSON.
- Formatting status responses for HTTP.
- Converting a dictionary of headers into a proper HTTP header format.
- Handling different types of responses (HTML, JSON, Plain Text) based on content type.

#### 3. Logger (`logger.go`)

Logging is a critical aspect of any application. Our `logger.go` file describes a custom logger capable of handling various logging levels (e.g., Info, Error) and providing a way to output logs with different severities. This is crucial for both debugging and runtime monitoring of the application.

#### 4. Server (`server.go`)

The `server.go` file is the heart of the framework where the HTTP server is configured and launched. Key functionalities include:

- **Starting the Server**: It listens on a specified host and port, manages incoming connections, and processes incoming data.
- **Routing**: It resolves the routes from the provided modules and efficiently manages the mapping of HTTP requests to their respective handlers.
- **Request and Response Handling**: It handles the parsing of requests, including headers and body, and ensures that responses are correctly formatted and sent back to the client.

### Design Decisions and Highlights

- **Routing Mechanism**: A tree-like structure for routing (as seen in `server.go`) allows efficient resolution of endpoints. This is particularly advantageous for applications with a large number of routes.
- **Dynamic Endpoints**: The framework supports dynamic routing (e.g., `/user/:id`), allowing for more flexible endpoint definitions.
- **Custom Error Handling**: Each component, especially in request parsing and response generation in `server.go`, includes comprehensive error handling, enhancing the robustness of the application.
- **Performance Logging**: Performance for critical operations, like route resolving and request handling, is monitored and logged, which can help in optimization and debugging.

### Conclusion and Potential Contributions

Contributors can consider various aspects for extending this framework:

- **Code Improvement and Review**:As the framework evolves, continuous improvement and refactoring of the codebase are essential. This involves optimizing existing code for performance, readability, and maintainability. Regular code reviews can ensure adherence to coding standards and best practices.
- **Middleware Support**: Introducing middleware support for pre-processing requests and post-processing responses.
- **Enhanced JSON and XML Parsing**: Further utilities for handling different types of request and response payloads.
- **Authentication and Security**: Adding layers or modules for handling authentication, authorization, and security checks.
- **HTTP/2 Support**: While the current setup is based on HTTP/1.1, evolving it to support HTTP/2 could significantly enhance performance.
- **Template Rendering**: For applications serving HTML, integrating a template engine would be beneficial.
- **Command Line Interface (CLI) Tooling**: Developing a CLI tool for the framework can significantly improve the developer's experience. This tool could automate routine tasks like starting the server, generating boilerplate code for new modules or routes, and managing configurations.
- **Comprehensive Documentation and Examples**: Creating detailed documentation, including usage guides, API references, and example projects, would be invaluable. It would help new users to understand and effectively utilize the framework.
- **Community Engagement and Support**: Establishing a community around the framework through forums, chat groups, or regular meetups can foster collaboration, user support, and attract more contributors.

These contributions can help evolve the framework into a more feature-rich and user-friendly tool, encouraging wider adoption and fostering a robust user and contributor community.
The current framework lays down a robust foundation and provides ample opportunity for enhancements and scaling, making it an excellent starting point for building a full-fledged web framework in Go.
