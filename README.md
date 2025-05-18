# Multipart Form Encoder for Go

A Go package that simplifies encoding structs as multipart/form-data requests. This is particularly useful for API clients and file upload scenarios where you need to convert a Go struct into form fields and file uploads.

## Features

- Encode any Go struct as multipart form-data
- Automatic handling of various field types:
  - String fields as form fields
  - Numeric fields (int, uint, float) as form fields
  - Boolean fields as form fields 
  - []byte fields as file uploads with automatic extension detection
  - Slices as multiple form fields with the same name
  - Nested structs as JSON-encoded form fields
- Customizable field names via struct tags
- Optional custom filenames for file uploads

## Installation

```bash
go get github.com/yourusername/multipart
```

## Usage

### Basic Example

```go
package main

import (
    "fmt"
    "github.com/mdhesari/multipart-encoder"
    "net/http"
)

type UploadRequest struct {
    Username string `form:"username"`
    Email    string
    File     []byte `form:"attachment" filename:"document.pdf"`
    Tags     []string
}

func main() {
    req := UploadRequest{
        Username: "johndoe",
        Email:    "john@example.com",
        File:     []byte("file content"),
        Tags:     []string{"golang", "upload"},
    }

    // Encode the struct as multipart form-data
    buf, contentType, err := multipart.Encode(req)
    if err != nil {
        panic(err)
    }

    // Use the buffer and content type to send an HTTP request
    resp, err := http.Post("https://example.com/upload", contentType, buf)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    fmt.Println("Response status:", resp.Status)
}
```

### Struct Tag Options

- `form:"fieldname"` - Specify the name of the form field (defaults to lowercase field name)
- `form:"-"` - Skip this field when encoding
- `filename:"custom.ext"` - Set a custom filename for []byte fields (defaults to field name + detected extension)

## API Documentation

### EncodeMultipart

```go
func Encode(req any) (*bytes.Buffer, string, error)
```

Converts a struct into multipart form-data format. 

Returns:
- A buffer containing the encoded data
- The content type string (including boundary)
- Any error that occurred during encoding

## File Type Detection

The package includes basic file type detection for common formats:
- PNG (.png)
- JPEG (.jpg)
- GIF (.gif)
- PDF (.pdf)

For unrecognized file types, it defaults to the extension specified in `DefaultFileExtension` (default: `.bin`).

## License

BSD License
