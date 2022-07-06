package graphdb

import (
	"context"

	_ "github.com/mattn/go-sqlite3"
	"github.com/projectdiscovery/katana/ent"
	"github.com/projectdiscovery/katana/ent/endpoint"
	"github.com/projectdiscovery/katana/ent/predicate"
)

type GraphDB struct {
	entClient *ent.Client
}

func New() (*GraphDB, error) {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		return nil, err
	}

	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return &GraphDB{entClient: client}, nil
}

func (graphDB *GraphDB) Close() error {
	return graphDB.entClient.Close()
}

func (graphDB *GraphDB) AddEndpoint(ctx context.Context, newEndpoint *ent.Endpoint) (*ent.Endpoint, error) {
	return graphDB.entClient.Endpoint.
		Create().
		SetBody(newEndpoint.Body).
		SetSource(newEndpoint.Source).
		SetURL(newEndpoint.URL).
		SetHeaders(newEndpoint.Headers).
		SetMethod(newEndpoint.Method).
		Save(ctx)
}

func (graphDB *GraphDB) ConnectEndpoints(ctx context.Context, source *ent.Endpoint, destinations ...*ent.Endpoint) (*ent.Endpoint, error) {
	return source.Update().AddLinks(destinations...).Save(ctx)
}

func (graphDB *GraphDB) QueryConnections(ctx context.Context, e *ent.Endpoint) ([]*ent.Endpoint, error) {
	return e.QueryLinks().All(ctx)
}

func (graphDB *GraphDB) QueryEndpoints(ctx context.Context) ([]*ent.Endpoint, error) {
	return graphDB.entClient.Endpoint.Query().All(ctx)
}

func (graphDB *GraphDB) QueryEndpoint(ctx context.Context, e *ent.Endpoint) (*ent.Endpoint, error) {
	predicates := []predicate.Endpoint{
		endpoint.URL(e.URL),
		endpoint.Body(e.Body),
		endpoint.Source(e.Source),
		endpoint.Method(e.Method),
	}
	return graphDB.entClient.Endpoint.Query().Where(predicates...).Only(ctx)
}

func (graphDB *GraphDB) GetOrCreate(ctx context.Context, e *ent.Endpoint) (*ent.Endpoint, error) {
	endpoint, err := graphDB.QueryEndpoint(ctx, e)
	if ent.IsNotFound(err) {
		return graphDB.AddEndpoint(ctx, e)
	}
	return endpoint, err
}
