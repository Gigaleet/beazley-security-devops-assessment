package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

const defaultMetaBase = "http://169.254.169.254/latest/meta-data/"

// MetadataClient wraps an HTTP client and base URL for IMDS.
type MetadataClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewMetadataClient constructs a client with timeout and given endpoint.
func NewMetadataClient(endpoint string) *MetadataClient {
	return &MetadataClient{
		httpClient: &http.Client{Timeout: 2 * time.Second},
		baseURL:    endpoint,
	}
}

// fetch retrieves raw bytes from IMDS path.
func (c *MetadataClient) fetch(path string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.baseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Treat 404 as “key not found”
	if resp.StatusCode == http.StatusNotFound {
	    return nil, fmt.Errorf("key %q not found", path)
	}
	// Any other non-2xx is an error
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
	    return nil, fmt.Errorf("unexpected HTTP status: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// parseEntry handles folders (newline-delimited) or leaves.
func (c *MetadataClient) parseEntry(path string) (interface{}, error) {
	data, err := c.fetch(path)
	if err != nil {
		return nil, err
	}
	str := string(data)
	if strings.Contains(str, "\n") {
		m := make(map[string]interface{})
		for _, line := range strings.Split(strings.TrimSpace(str), "\n") {
			key := strings.TrimSuffix(line, "/")
			val, err := c.parseEntry(path + line)
			if err != nil {
				return nil, err
			}
			m[key] = val
		}
		return m, nil
	}
	return str, nil
}

// GetAll fetches all metadata recursively.
func (c *MetadataClient) GetAll() (map[string]interface{}, error) {
	v, err := c.parseEntry("")
	if err != nil {
		return nil, err
	}
	return v.(map[string]interface{}), nil
}

// Get fetches a single metadata key.
func (c *MetadataClient) Get(key string) (interface{}, error) {
	return c.parseEntry(strings.TrimPrefix(key, "/"))
}

// fetchEC2Metadata uses AWS SDK to describe an instance by ID.
func fetchEC2Metadata(ctx context.Context, instanceID string) (interface{}, error) {
	// load AWS config
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}
	client := ec2.NewFromConfig(cfg)
	out, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instances error: %w", err)
	}
	return out, nil
}

func main() {
	key := flag.String("key", "", "metadata key to fetch (fetch all if empty)")
	endpoint := flag.String("endpoint", defaultMetaBase, "IMDS base URL")
	instanceID := flag.String("instance-id", "", "EC2 instance ID to fetch via AWS API")
	flag.Parse()

	var result interface{}
	var err error
	ctx := context.Background()

	if *instanceID != "" {
		result, err = fetchEC2Metadata(ctx, *instanceID)
	} else {
		client := NewMetadataClient(*endpoint)
		if *key == "" {
			result, err = client.GetAll()
		} else {
			result, err = client.Get(*key)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	out, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "json marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(out))
}
