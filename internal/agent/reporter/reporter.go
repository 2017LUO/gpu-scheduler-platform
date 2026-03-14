package reporter

import "context"

type Reporter interface {
	Report(ctx context.Context, payload any) error
}
