package auth

import "context"

type Subject struct {
	SubjectID   string
	Name        string
	TenantID    string
	Roles       []string
	Permissions []string
	TokenID     string
	Issuer      string
	IsSystem    bool
}

type subjectContextKey struct{}

func WithSubject(ctx context.Context, sub *Subject) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, subjectContextKey{}, sub)
}

func SubjectFromContext(ctx context.Context) (*Subject, bool) {
	if ctx == nil {
		return nil, false
	}
	v, ok := ctx.Value(subjectContextKey{}).(*Subject)
	return v, ok && v != nil
}

func AnonymousSubject() *Subject {
	return &Subject{
		SubjectID: "anonymous",
		Name:      "anonymous",
		Roles:     []string{"anonymous"},
	}
}

func SystemSubject() *Subject {
	return &Subject{
		SubjectID: "system",
		Name:      "system",
		Roles:     []string{"admin"},
		IsSystem:  true,
	}
}
