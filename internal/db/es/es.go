package es

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/olivere/elastic/v7/config"
)
import "github.com/olivere/elastic/v7"

func NewESRepository(log log.Logger) (*Repository, error) {
	esConfig, err := config.Parse("http://127.0.0.1:9200?sniff=false&healthcheck=false")
	if err != nil {
		return nil, err
	}
	client, err := elastic.NewClientFromConfig(esConfig)
	if err != nil {
		return nil, err
	}
	return &Repository{
		log:    log,
		client: client,
	}, nil
}

// Repository is used to search and create an index for players
type Repository struct {
	log    log.Logger
	client *elastic.Client
}

// NewIndex will create an index
func (es *Repository) NewIndex(ctx context.Context, index string) (error){
	svc := es.client.CreateIndex("ethan")

	_, err := svc.Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
