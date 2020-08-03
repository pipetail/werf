package helm_v3

import "time"

type RollbackOptions struct {
	Timeout time.Duration
}

func Rollback(releaseName string, version int64, opts RollbackOptions) error {
	return nil
}
