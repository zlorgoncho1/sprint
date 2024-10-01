# Sprint Web Framework

Sprint is a modular, high-performance web framework for Go, inspired by the architecture of NestJS. Designed to streamline backend development, Sprint empowers developers to build scalable and maintainable applications with ease. It emphasizes a modular structure, where functionality is divided into controllers and modules, allowing for clear organization and flexibility in development.

## Table of Contents

- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Creating Your First Sprint Application](#creating-your-first-sprint-application)
- [Key Features](#key-features)
- [Main Components](#main-components)
- [Design Decisions and Highlights](#design-decisions-and-highlights)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Installation

Before you start using Sprint, ensure that you have Go installed on your machine. If not, you can follow Go's [official documentation](https://golang.org/doc/install) for installation guidance.

Once you have Go installed, you can install Sprint as a Go module using the following command:

```bash
go get github.com/zlorgoncho1/sprint/server
```

### Creating Your First Sprint Application

#### Step 1: Creating the Project

Start by setting up a new Go project:

```bash
mkdir app && cd app
go mod init app
```

#### Step 2: Project Structure

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

#### Step 3: Writing Code

##### `main.go`

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

##### `module.go`

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

##### `controllers.go`

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

#### Step 4: Running Your Application

Run your application:

```bash
go run main.go
```

Your Sprint application should now be running on `localhost:8000`. Visiting `localhost:8000/hello/HTML/Sprint` in your browser or through a tool like `curl` should return a "Hello World" message.

## Key Features

- **Modular Architecture**: Like NestJS, Sprint organizes code into modules, making it easier to manage and scale large applications.
- **Intuitive Controller Setup**: Simplified creation of controllers to handle various routes and requests.
- **Simplicity and Performance**: Leverages Go's efficiency and simplicity, providing a framework that is both easy to use and high-performing.

## Main Components

### Core (`core.go`)

The core of our framework lies within `core.go`, which defines the primary structs and interfaces used throughout the application. This includes definitions for:

- **Modules**: Conceptually similar to controllers in MVC frameworks, they're used to organize the handling of different routes and their associated functionalities.
- **EndpointNode**: A fundamental part of the routing mechanism, these nodes represent individual endpoints in a tree-like structure for efficient route resolution.
- **Response and Request**: Custom types to handle incoming requests and outgoing responses, providing an abstraction over Go's native HTTP handling facilities.

### Utils (`utils.go`)

This utility package (`utils.go`) provides a set of handy functions for converting dictionaries to JSON strings, formatting HTTP responses, and generating standard headers. It abstracts some repetitive tasks, such as:

- Converting map objects to JSON.
- Formatting status responses for HTTP.
- Converting a dictionary of headers into a proper HTTP header format.
- Handling different types of responses (HTML, JSON, Plain Text) based on content type.

### Logger (`logger.go`)

Logging is a critical aspect of any application. Our `logger.go` file describes a custom logger capable of handling various logging levels (e.g., Info, Error) and providing a way to output logs with different severities. This is crucial for both debugging and runtime monitoring of the application.

### Server (`server.go`)

The `server.go` file is the heart of the framework where the HTTP server is configured and launched. Key functionalities include:

- **Starting the Server**: It listens on a specified host and port, manages incoming connections, and processes incoming data.
- **Routing**: It resolves the routes from the provided modules and efficiently manages the mapping of HTTP requests to their respective handlers.
- **Request and Response Handling**: It handles the parsing of requests, including headers and body, and ensures that responses are correctly formatted and sent back to the client.

## Design Decisions and Highlights

- **Routing Mechanism**: A tree-like structure for routing (as seen in `server.go`) allows efficient resolution of endpoints. This is particularly advantageous for applications with a large number of routes.
- **Dynamic Endpoints**: The framework supports dynamic routing (e.g., `/user/:id`), allowing for more flexible endpoint definitions.
- **Custom Error Handling**: Each component, especially in request parsing and response generation in `server.go`, includes comprehensive error handling, enhancing the robustness of the application.
- **Performance Logging**: Performance for critical operations, like route resolving and request handling, is monitored and logged, which can help in optimization and debugging.

## Contributing

We welcome contributions to Sprint! If you'd like to get involved, here are some areas where you can make a difference:

- **Code Improvement and Review**: Continuous improvement, optimization, and refactoring of the codebase.
- **Middleware Support**: Introducing middleware for pre-processing requests and post-processing responses.
- **Enhanced JSON and XML Parsing**: More utilities for handling different types of request and response payloads.
- **Authentication and Security**: Adding modules for handling authentication, authorization, and security.
- **HTTP/2 Support**: Evolving the framework to support HTTP/2 for better performance.
- **Template Rendering**: Integrating a template engine for serving HTML.
- **Command Line Interface (CLI) Tooling**: Developing a CLI tool to automate routine tasks.

For more details, please refer to our [Contribution Guidelines](CONTRIBUTING.md) or explore the [Wiki](https://github.com/zlorgoncho1/sprint/wiki).

## License

Sprint is licensed under the [MIT License](LICENSE).
