# ShopCx FrontEnd Service

A Go-based frontend service for the ShopCx demo application with intentionally vulnerable endpoints. This service handles user authentication, product management, file uploads, and serves as the main web interface for the e-commerce platform.

## Security Note

⚠️ **This is an intentionally vulnerable application for security testing purposes. Do not deploy in production or sensitive environments.**

## Overview

The FrontEnd Service is a Go application built with the Gin web framework that provides the primary web interface and API gateway for the ShopCx platform. It handles user authentication, product operations, file uploads, and proxies requests to other microservices with intentionally vulnerable SQL injection patterns.

## Key Features

- **User Authentication**: JWT-based login and session management
- **Product Management**: Create, search, and manage product listings
- **File Upload**: Handle product images and document uploads
- **Search Functionality**: Product search with dynamic results
- **Admin Operations**: Administrative user and product management
- **Rate Limiting**: Basic request rate limiting and throttling
- **CORS Support**: Cross-origin resource sharing configuration
- **Session Management**: Cookie-based session handling
- **User Profile Management**: Retrieve user information
- **HTML Template Rendering**: Dynamic web page generation

## Technology Stack

- **Go 1.23**: Core programming language
- **Gin Framework**: High-performance HTTP web framework
- **MySQL**: Database for user and product data
- **JWT**: JSON Web Token authentication
- **Gorilla Sessions**: Session management
- **HTML Templates**: Go template engine for web pages
- **CORS Middleware**: Cross-origin request handling
- **Database/SQL**: Standard library database interface
- **File System Operations**: File upload and storage handling
