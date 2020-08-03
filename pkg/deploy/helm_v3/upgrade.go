package helm_v3

import (
	"context"
	"time"
)

type UpgradeOptions struct {
	Timeout time.Duration
}

func Upgrade(ctx context.Context, chartPath, releaseName, namespace string, values []string, secretValues []map[string]interface{}, set, setString []string, opts UpgradeOptions) error {

}
