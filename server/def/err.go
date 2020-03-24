package def

/*
Package def defines constants, error messages and their status codes
*/

import (
	"errors"
)

var (
	DbNameEmpty               error = errors.New("dbname empty")
	DbDoesNotExist            error = errors.New("database doesn't exist")
	DbIdentifierEmpty         error = errors.New("db identifier empty")
	CollectionNameEmpty       error = errors.New("collection name empty")
	CollectionIdentifierEmpty error = errors.New("collection identifier empty")

	NamespaceIdentifierEmpty error = errors.New("namespace identifier empty")
	NamesCannotBeEmpty       error = errors.New("database/collection/namespace names can't be empty")
	NamespaceDoesNotExist    error = errors.New("namespace doesn't exist")
	KeyEmpty                 error = errors.New("key is empty")
	EmptyKeyCannotBeDeleted  error = errors.New("can't delete empty key")
	StartOrEndKeyEmpty       error = errors.New("start or end key is empty")
	StartKeyUnknown          error = errors.New("can't scan from last without knowing startKey")
	IdentifierNotFound       error = errors.New("id not found for given db/collection/namespace")

	ConnectionCouldNotBeEstablished error = errors.New("connectionstore to database could not be established")
	ResultsNotFound                 error = errors.New("results not found")
)
