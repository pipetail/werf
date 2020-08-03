package helm_v3

import (
	"context"
	"time"
)

type InstallOptions struct {
	Timeout time.Duration
}

func Install(ctx context.Context, chartPath, releaseName, namespace string, values []string, secretValues []map[string]interface{}, set, setString []string, opts InstallOptions) error {
	return nil
}
