# goandzer-load-balancer

A simple load balancer in Go built with [Echo](https://echo.labstack.com) and the standard library’s `net/http/httputil` reverse proxy. Ideal for prototyping and learning.

## Features

- Round-robin request distribution  
- Periodic health checks to keep backends up-to-date  
- Customizable path-based routing  
- Easy integration with Echo middleware
- Admin URL to check services status  (base_url/admin/status)

## Getting Started

### Prerequisites

- Go 1.18 or newer  
- `git`  

### Installation

```bash
git clone https://github.com/Juanmagc99/goandzer.git
cd goandzer
go mod tidy
```

### HOW to septup
- Create an ADMIN_API_KEY env variable
- Create a config.yaml file inside config folder following example_config.yaml structure
