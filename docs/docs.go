// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "Show the start page of the API",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Show start page",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/SearchBooks": {
            "get": {
                "description": "Returns an array of books that are similar in name or description to the request",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Outputs an array of books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Looking for a similar book",
                        "name": "search",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination (default: 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of books per page (default: 10)",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns a paginated and sorted list of books",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/addBook": {
            "post": {
                "description": "JWT authentication via cookie only for admin.\nThe JWT token should be stored in a cookie named \"jwt\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Add a new book",
                "parameters": [
                    {
                        "description": "Book Data",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.AddBookRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/deleteBook": {
            "delete": {
                "description": "deletes the book from the library\nJWT authentication via cookie only for admin.\nThe JWT token should be stored in a cookie named \"jwt\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Delete the book",
                "parameters": [
                    {
                        "description": "Book Data",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.DeleteBookRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/getBook": {
            "get": {
                "description": "Get detailed information about a single book by ID\nJWT authentication via cookie.\nThe JWT token should be stored in a cookie named \"jwt\".",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Get one book",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Book ID",
                        "name": "bookId",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Book"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/getBooks": {
            "get": {
                "description": "Retrieve all books, optionally sorted by a specific field",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Get list of books",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Field to sort by (e.g., 'title', 'author', 'published_year')",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Page number for pagination (default: 1)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "Number of books per page (default: 10)",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Returns a paginated and sorted list of books",
                        "schema": {
                            "type": "object",
                            "additionalProperties": true
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/logOut": {
            "post": {
                "description": "Log user from the api",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Log out user",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/login": {
            "post": {
                "description": "Logs in an existing user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Performs user login",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/modifyingBook": {
            "post": {
                "description": "JWT authentication via cookie only for admin.\nThe JWT token should be stored in a cookie named \"jwt\".\nJWT Bearer authentcation only admin",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "book"
                ],
                "summary": "Modifying book",
                "parameters": [
                    {
                        "description": "Book Data",
                        "name": "book",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.ModifyingBookRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/register": {
            "post": {
                "description": "Add a new library User",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Add a new User",
                "parameters": [
                    {
                        "description": "User Data",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/handlers.RegisterUserRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/subMailing": {
            "get": {
                "description": "Subscribes a user to mailing lists",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Subscribe mailing",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        },
        "/unsubMailing": {
            "get": {
                "description": "Describes the user from the mailing list",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Unsubscribe mailing",
                "responses": {
                    "200": {
                        "description": "OK"
                    }
                }
            }
        }
    },
    "definitions": {
        "handlers.AddBookRequest": {
            "type": "object",
            "required": [
                "author",
                "genre",
                "published_year",
                "title"
            ],
            "properties": {
                "author": {
                    "description": "Автор",
                    "type": "string",
                    "example": "John Doe"
                },
                "description": {
                    "description": "Описание книги",
                    "type": "string",
                    "example": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."
                },
                "genre": {
                    "description": "Жанра",
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Учебная литература"
                    ]
                },
                "published_year": {
                    "description": "Год публикации",
                    "type": "string",
                    "example": "2024"
                },
                "title": {
                    "description": "Название книги",
                    "type": "string",
                    "example": "Golang Basics"
                }
            }
        },
        "handlers.DeleteBookRequest": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "id": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "handlers.LoginRequest": {
            "type": "object",
            "required": [
                "email",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string"
                },
                "password": {
                    "type": "string",
                    "example": "123456"
                }
            }
        },
        "handlers.ModifyingBookRequest": {
            "type": "object",
            "required": [
                "id"
            ],
            "properties": {
                "author": {
                    "type": "string",
                    "example": "Jeff Bezos"
                },
                "description": {
                    "type": "string",
                    "example": "Explore the ultimate question: Why is there something rather than nothing? This thought-provoking journey through philosophy, science, and metaphysics challenges readers to ponder existence itself, blending deep inquiry with accessible insight. A must-read for curious minds."
                },
                "genre": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Детектив"
                    ]
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "published_year": {
                    "type": "string",
                    "example": "2021"
                },
                "title": {
                    "type": "string",
                    "example": "Why Does the World Exist?"
                }
            }
        },
        "handlers.RegisterUserRequest": {
            "type": "object",
            "required": [
                "email",
                "mailing",
                "name",
                "password"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "Laminano@mail.ru"
                },
                "mailing": {
                    "type": "boolean",
                    "example": true
                },
                "name": {
                    "type": "string",
                    "example": "Vladislav"
                },
                "password": {
                    "type": "string",
                    "minLength": 6,
                    "example": "123456"
                }
            }
        },
        "models.Book": {
            "type": "object",
            "properties": {
                "author": {
                    "type": "string"
                },
                "description": {
                    "type": "string"
                },
                "genres": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Genre"
                    }
                },
                "published_year": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "models.Genre": {
            "type": "object",
            "properties": {
                "books": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Book"
                    }
                },
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Library API",
	Description:      "This is a sample library server",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
