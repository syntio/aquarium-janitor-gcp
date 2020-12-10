// Copyright Syntio d.o.o.
// All Rights Reserved
//
// Package database is used a an interface between the REST Server and the underlying database.
//
// It holds general information about the database infrastructure and is the mediator between the REST Server and whatever
// database technology is used underneath.
//
// It currently doesn't serve its full purpose but will soon be a valid interface.
package database

import (
	"context"

	. "github.com/syntio/schema-registry/model"
	"github.com/syntio/schema-registry/model/dto"
)

//
// DBConnector is an interface between the REST Server and the underlying database.
//
type DBExecutor interface {
	CreateSchema(ctx context.Context, dto *dto.SchemaDTO) (*InsertInfo, bool, error)
	GetSchemaByIdAndVersion(ctx context.Context, id string, version int32) (*Schema, bool)
	UpdateSchemaById(ctx context.Context, id string, schema []byte, autogenerated bool) (*InsertInfo, bool, error)
	GetSchemaVersions(ctx context.Context, id string) (*[]*SchemaDetails, error)
	DeleteById(ctx context.Context, id string) error
}