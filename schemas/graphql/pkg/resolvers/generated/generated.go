package generated

// THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

import (
	"context"

	"github.com/jay-manday/perspex/schemas/graphql/pkg/graphql/model"
	"github.com/jay-manday/perspex/schemas/graphql/pkg/graphql/source"
)

type Resolver struct{}

// // foo
func (r *mutationResolver) NodeMutation(ctx context.Context) (model.Node, error) {
	panic("not implemented")
}

// // foo
func (r *queryResolver) NodeQuery(ctx context.Context) (model.Node, error) {
	panic("not implemented")
}

// // foo
func (r *subscriptionResolver) NodeSubscription(ctx context.Context) (<-chan model.Node, error) {
	panic("not implemented")
}

// Mutation returns source.MutationResolver implementation.
func (r *Resolver) Mutation() source.MutationResolver { return &mutationResolver{r} }

// Query returns source.QueryResolver implementation.
func (r *Resolver) Query() source.QueryResolver { return &queryResolver{r} }

// Subscription returns source.SubscriptionResolver implementation.
func (r *Resolver) Subscription() source.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
