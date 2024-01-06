package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.41

import (
	"context"
	"db"
	"fmt"
	"gqlgen/graph/model"
)

// Bots is the resolver for the bots field.
func (r *queryResolver) Bots(ctx context.Context) ([]*model.Bot, error) {
	var ss db.TelegramBot
	var s []*model.Bot
	bts := ss.ListBots()
	for i := range bts {
		s = append(s, &model.Bot{
			ID:   bts[i].Id,
			Name: bts[i].Name,
		})
	}
	return s, nil
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *botResolver) ID(ctx context.Context, obj db.Bot) (string, error) {
	panic(fmt.Errorf("not implemented: ID - id"))
}
func (r *botResolver) Name(ctx context.Context, obj db.Bot) (string, error) {
	panic(fmt.Errorf("not implemented: Name - Name"))
}

type botResolver struct{ *Resolver }