# Nested Object Key Lookup

A small Go package that provides a utility function to extract values from arbitrarily nested `map[string]interface{}` objects using slash-delimited keys (e.g. `"a/b/c"`). Simplifies deep lookups without boilerplate, with built-in error handling and unit tests.

---

## Features

* **Slash-Delimited Paths**: Lookup nested values with a single string, e.g.:

  ```go
  GetNestedValue(obj, "a/b/c")
  ```
* **Error Handling**:

  * Empty key produces an error.
  * Missing path segment yields a clear "key not found" error.
  * Unexpected types during traversal produce descriptive type errors.
* **Unit Tested**: Comprehensive tests using Go’s `testing` package to cover success and failure scenarios.

---

## Requirements

* Go 1.16 or newer

---

## Installation

```bash
go get github.com/Gigaleet/beazley-security-devops-assessment/cmd/nested
```

Or clone this repository and import locally:

```bash
git clone <repo-url>
cd <repo-folder>
```

---

## Usage

Import the package:

```go
import "github.com/Gigaleet/beazley-security-devops-assessment/cmd/nested"
```

Call the helper:

```go
obj := map[string]interface{}{
  "a": map[string]interface{}{
    "b": map[string]interface{}{ "c": "d" },
  },
}

value, err := nested.GetNestedValue(obj, "a/b/c")
if err != nil {
  // handle error
}
fmt.Println(value) // prints: d
```

---

## API

```go
func GetNestedValue(obj map[string]interface{}, key string) (interface{}, error)
```

* **`obj`**: the nested `map[string]interface{}` structure.
* **`key`**: a non-empty, slash-delimited path. Each segment must correspond to a map key.
* **Returns**: the value at the given path, or an error if the key is empty, missing, or the traversal encounters a non-map.

---

## Testing

Run unit tests:

```bash
go test ./cmd/nested -v
```

Sample tests include:

* Successful retrieval for valid nested paths.
* Error for empty key.
* Error for missing keys.
* Error for wrong types during traversal.
