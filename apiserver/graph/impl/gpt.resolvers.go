package impl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/kubeagi/arcadia/apiserver/graph/generated"
	"github.com/kubeagi/arcadia/apiserver/pkg/gpt"
)

// GetGpt is the resolver for the getGPT field.
func (r *gPTQueryResolver) GetGpt(ctx context.Context, obj *generated.GPTQuery, name string) (*generated.Gpt, error) {
	c, err := getAdminClient()
	if err != nil {
		return nil, err
	}
	return gpt.GetGPT(ctx, c, name)
}

// ListGpt is the resolver for the listGPT field.
func (r *gPTQueryResolver) ListGpt(ctx context.Context, obj *generated.GPTQuery, input generated.ListGPTInput) (*generated.PaginatedResult, error) {
	c, err := getAdminClient()
	if err != nil {
		return nil, err
	}
	return gpt.ListGPT(ctx, c, input)
}

// Gpt is the resolver for the GPT field.
func (r *queryResolver) Gpt(ctx context.Context) (*generated.GPTQuery, error) {
	return &generated.GPTQuery{}, nil
}

// GPTQuery returns generated.GPTQueryResolver implementation.
func (r *Resolver) GPTQuery() generated.GPTQueryResolver { return &gPTQueryResolver{r} }

type gPTQueryResolver struct{ *Resolver }