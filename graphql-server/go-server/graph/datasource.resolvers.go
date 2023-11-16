package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"strings"

	"github.com/kubeagi/arcadia/graphql-server/go-server/graph/model"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/auth"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/client"
	"github.com/kubeagi/arcadia/graphql-server/go-server/pkg/datasource"
)

// CreateDatasource is the resolver for the createDatasource field.
func (r *mutationResolver) CreateDatasource(ctx context.Context, input model.CreateDatasourceInput) (*model.Datasource, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}

	url, authSecret, bucket, displayname := "", "", "", ""
	var insecure bool
	if input.Endpointinput != nil {
		if input.Endpointinput.URL != nil {
			url = *input.Endpointinput.URL
		}
		insecure = *input.Endpointinput.Insecure
		if input.Endpointinput.AuthSecret.Name != "" {
			authSecret = input.Endpointinput.AuthSecret.Name
		}
	}
	if input.Ossinput != nil && input.Ossinput.Bucket != nil {
		bucket = *input.Ossinput.Bucket
	}
	if input.DisplayName != "" {
		displayname = input.DisplayName
	}
	return datasource.CreateDatasource(ctx, c, input.Name, input.Namespace, url, authSecret, bucket, displayname, insecure)
}

// UpdateDatasource is the resolver for the updateDatasource field.
func (r *mutationResolver) UpdateDatasource(ctx context.Context, input *model.UpdateDatasourceInput) (*model.Datasource, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name, displayname := "", ""
	if input.DisplayName != "" {
		displayname = input.DisplayName
	}
	if input.Name != "" {
		name = input.Name

	}
	return datasource.UpdateDatasource(ctx, c, name, input.Namespace, displayname)
}

// DeleteDatasource is the resolver for the deleteDatasource field.
func (r *mutationResolver) DeleteDatasource(ctx context.Context, input *model.DeleteDatasourceInput) (*string, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name := ""
	labelSelector, fieldSelector := "", ""
	if input.Name != nil {
		name = *input.Name
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	return datasource.DeleteDatasource(ctx, c, name, input.Namespace, labelSelector, fieldSelector)
}

// Datasource is the resolver for the datasource field.
func (r *queryResolver) Datasource(ctx context.Context, name string, namespace string) (*model.Datasource, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	return datasource.ReadDatasource(ctx, c, name, namespace)
}

// DatasourcesPaged is the resolver for the datasourcesPaged field.
func (r *queryResolver) DatasourcesPaged(ctx context.Context, input model.ListDatasourceInput) (*model.PaginatedDatasource, error) {
	token := auth.ForOIDCToken(ctx)
	c, err := client.GetClient(token)
	if err != nil {
		return nil, err
	}
	name, displayName, labelSelector, fieldSelector := "", "", "", ""
	page, pageSize := 1, 10
	if input.Name != nil {
		name = *input.Name
	}
	if input.DisplayName != nil {
		displayName = *input.DisplayName
	}
	if input.FieldSelector != nil {
		fieldSelector = *input.FieldSelector
	}
	if input.LabelSelector != nil {
		labelSelector = *input.LabelSelector
	}
	if input.Page != nil && *input.Page > 0 {
		page = *input.Page
	}
	if input.PageSize != nil && *input.PageSize > 0 {
		pageSize = *input.PageSize
	}
	result, err := datasource.ListDatasources(ctx, c, input.Namespace, labelSelector, fieldSelector)
	if err != nil {
		return nil, err
	}
	var filteredResult []*model.Datasource
	for idx, u := range result {
		if (name == "" || strings.Contains(u.Name, name)) && (displayName == "" || strings.Contains(u.DisplayName, displayName)) {
			filteredResult = append(filteredResult, result[idx])
		}
	}

	totalCount := len(filteredResult)
	end := page * pageSize
	if end > totalCount {
		end = totalCount
	}
	return &model.PaginatedDatasource{
		TotalCount:  totalCount,
		HasNextPage: end < totalCount,
		Nodes:       filteredResult[(page-1)*pageSize : end],
	}, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }