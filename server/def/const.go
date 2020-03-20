package def

/*
Package def defines constants, error messages and their status codes
*/

const (
	MetaDbidentifier                 = "meta:dbidentifier"
	MetaCollectionidentifier         = "meta:collectionidentifier"
	MetaNamespaceidentifier          = "meta:namespaceidentifier"
	MetaDbid                         = "meta:dbid:"
	MetaCollectionid                 = "meta:collectionid:"
	MetaNamespaceid                  = "meta:namespaceid:"
	MetaDb                           = "meta:db:"
	MetaCollection                   = "meta:collection:"
	MetaNamespace                    = "meta:namespace:"
	IndexKey                         = "_index:"
	UniqueId                         = "_uniqueid:"
	UniqueIdInitialcount             = uint32(1)
	DbidentifierInitialcount         = uint16(1)
	CollectionidentifierInitialcount = uint32(1)
	NamespaceidentifierInitialcount  = uint32(1)
)
