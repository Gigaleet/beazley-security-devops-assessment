# EC2 Metadata Fetcher

A simple CLI tool written in Go that retrieves EC2 instance metadata either directly from the Instance Metadata Service (IMDS) or, when provided an EC2 instance ID, via the AWS EC2 API. Outputs all data in tidy JSON format.

---

## Features

* **IMDS Lookup**: Recursively fetches the full metadata tree or a single slash‑delimited key (e.g. `network/interfaces/mac`).
* **API Lookup**: When run locally, accepts an `--instance-id` flag to describe an instance via AWS SDK v2 and returns its metadata structure.
* **Configurable Endpoint**: Override the metadata endpoint (e.g. for local stub servers) with `--endpoint`.
* **JSON Output**: All results are marshaled with `encoding/json` and pretty‑printed.
* **Unit Tested**: Uses `net/http/httptest` to stub IMDS responses and verify parsing logic.

---

## Requirements

* Go 1.18 or newer
* (For API mode) AWS credentials/config in your environment (e.g. `~/.aws/credentials`)

---

## Installation

```bash
git clone <repo-url>
cd <repo-folder>
go build -o ec2meta metadata.go
```

This produces the `ec2meta` binary.

---

## Usage

### IMDS Mode (on EC2 or stub server)

Fetch the entire metadata tree:

```bash
./ec2meta
```

Fetch a single key:

```bash
./ec2meta -key instance-id
```

Override endpoint (for testing against a stub at `localhost:8080`):

```bash
./ec2meta -endpoint http://localhost:8080/latest/meta-data/ -key instance-id
```

### API Mode (local lookup)

Describe an EC2 instance by ID:

```bash
./ec2meta -instance-id i-0123456789abcdef0
```

This mode uses AWS SDK’s `DescribeInstances` call. Ensure your IAM user or role has `ec2:DescribeInstances` permission.

---

## Flags

| Flag            | Description                                                           | Default                                    |
| --------------- | --------------------------------------------------------------------- | ------------------------------------------ |
| `-key string`   | Specific metadata key to fetch (slash‑delimited). Fetch all if empty. | `""`                                       |
| `-endpoint str` | IMDS base URL (for EC2 metadata).                                     | `http://169.254.169.254/latest/meta-data/` |
| `-instance-id`  | EC2 instance ID to describe via AWS API (skips IMDS mode).            | `""`                                       |

---

## Testing

Run all unit tests:

```bash
go test -v ./...
```

Tests spin up an HTTP stub server to simulate IMDS for both folder and leaf responses.

---

## Local Stub Example

```bash
# Start stub server:
cat << 'EOF' > stub.go
package main
import (
  "net/http"
  "log"
)
func main() {
  http.HandleFunc("/latest/meta-data/instance-id", func(w, _ ) { w.Write([]byte("i-STUB123")) })
  log.Fatal(http.ListenAndServe(":8080", nil))
}
go run stub.go

# In another terminal:
./ec2meta -endpoint http://localhost:8080/latest/meta-data/ -key instance-id
```

Outputs:

```json
"i-STUB123"
```
