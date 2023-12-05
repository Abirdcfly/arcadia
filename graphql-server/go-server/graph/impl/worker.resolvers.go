package impl

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"

	"github.com/kubeagi/arcadia/graphql-server/go-server/graph/generated"
	md "github.com/kubeagi/arcadia/graphql-server/go-server/pkg/worker"
)

// Worker is the resolver for the Worker field.
func (r *mutationResolver) Worker(ctx context.Context) (*generated.WorkerMutation, error) {
	return &generated.WorkerMutation{}, nil
}

// Worker is the resolver for the Worker field.
func (r *queryResolver) Worker(ctx context.Context) (*generated.WorkerQuery, error) {
	return &generated.WorkerQuery{}, nil
}

// CreateWorker is the resolver for the createWorker field.
func (r *workerMutationResolver) CreateWorker(ctx context.Context, obj *generated.WorkerMutation, input generated.CreateWorkerInput) (*generated.Worker, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return md.CreateWorker(ctx, c, input)
}

// UpdateWorker is the resolver for the updateWorker field.
func (r *workerMutationResolver) UpdateWorker(ctx context.Context, obj *generated.WorkerMutation, input *generated.UpdateWorkerInput) (*generated.Worker, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return md.UpdateWorker(ctx, c, input)
}

// DeleteWorkers is the resolver for the deleteWorkers field.
func (r *workerMutationResolver) DeleteWorkers(ctx context.Context, obj *generated.WorkerMutation, input *generated.DeleteCommonInput) (*string, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return md.DeleteWorkers(ctx, c, input)
}

// GetWorker is the resolver for the getWorker field.
func (r *workerQueryResolver) GetWorker(ctx context.Context, obj *generated.WorkerQuery, name string, namespace string) (*generated.Worker, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}
	return md.ReadWorker(ctx, c, name, namespace)
}

// ListWorkers is the resolver for the listWorkers field.
func (r *workerQueryResolver) ListWorkers(ctx context.Context, obj *generated.WorkerQuery, input generated.ListCommonInput) (*generated.PaginatedResult, error) {
	c, err := getClientFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	return md.ListWorkers(ctx, c, input)
}

// WorkerMutation returns generated.WorkerMutationResolver implementation.
func (r *Resolver) WorkerMutation() generated.WorkerMutationResolver {
	return &workerMutationResolver{r}
}

// WorkerQuery returns generated.WorkerQueryResolver implementation.
func (r *Resolver) WorkerQuery() generated.WorkerQueryResolver { return &workerQueryResolver{r} }

type workerMutationResolver struct{ *Resolver }
type workerQueryResolver struct{ *Resolver }