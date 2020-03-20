package kvstore

import (
	"demo-backend/server/io"
	"demo-backend/server/kvstore/badgerdb"
)

//NewBadgerFactory returns badgerdb storeclient as io.Store
//pdAddr => Placement Driver
//dbDir => Db Directory
func NewBadgerFactory(pdAddr []string, dbDIR string) io.Store {
	badger := &badgerdb.StoreClient{}
	err := badger.NewClient(pdAddr, dbDIR)
	if err != nil {
		panic(err)
	}
	return badger
}
