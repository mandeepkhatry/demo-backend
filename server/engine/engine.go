package engine

import (
	"demo-backend/server/def"
	"demo-backend/server/io"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/RoaringBitmap/roaring"
)

/*
Design considerations`
---------------------
A typical key consists of following parts:

- db_name [2 bytes] ~ 65k values
- collection_name [4 bytes]
- namespace [4 bytes]
- unique_id [4 bytes]
Total key size for a document will be 14 bytes.
*/

type Engine struct {
	DBName      string
	DBID        []byte
	Namespace   string
	NamespaceID []byte
	Session     map[string][]byte //session is used to check whether given d,c,n creds are correct
	Store       io.Store
}

//ConnectDB initializes engine with DBName, Namespace
func (e *Engine) ConnectDB() error {
	fmt.Println("[[ConnectDB]] inside function")
	dbname := []byte(e.DBName)
	namespace := []byte(e.Namespace)

	dbID, err := e.GetDBIdentifier(dbname)
	if err != nil {
		return err
	}

	e.Session = make(map[string][]byte)

	e.Session[e.DBName] = dbID
	e.DBID = dbID

	namespaceID, err := e.GetNamespaceIdentifier(namespace)
	if err != nil {
		return err
	}

	e.Session[e.Namespace] = namespaceID
	e.NamespaceID = namespaceID

	return nil
}

//GenerateDBIdentifier return db identifier value and increase identifier by 1
func (e *Engine) GenerateDBIdentifier(dbname []byte) ([]byte, error) {
	val, err := e.Store.Get([]byte(def.MetaDbidentifier))
	if err != nil {
		return []byte{}, err
	}
	//if there is no id present, generate a new one
	if len(val) == 0 {
		identifier := make([]byte, 2)

		binary.BigEndian.PutUint16(identifier, def.DbidentifierInitialcount)
		err := e.Store.Put([]byte(def.MetaDbidentifier), identifier)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint16(val)
		binary.BigEndian.PutUint16(val, uint16(identifier+1))

		err := e.Store.Put([]byte(def.MetaDbidentifier), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetDBIdentifier returns identifier for given db
func (e *Engine) GetDBIdentifier(dbname []byte) ([]byte, error) {
	if len(dbname) == 0 {
		return []byte{}, def.DbNameEmpty
	}
	val, err := e.Store.Get([]byte(def.MetaDb + string(dbname)))
	if err != nil {
		return []byte{}, err
	}

	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateDBIdentifier(dbname)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:db:dname = identifier
		err = e.Store.Put([]byte(def.MetaDb+string(dbname)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:dbid:id=name
		err = e.Store.Put(append([]byte(def.MetaDbid), identifier...), dbname)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		return val, nil
	}
}

//GetDBName returns database name for given db identifier
func (e *Engine) GetDBName(dbIdentifier []byte) (string, error) {
	if len(dbIdentifier) == 0 {
		return "", def.DbIdentifierEmpty
	}
	val, err := e.Store.Get(append([]byte(def.MetaDbid), dbIdentifier...))
	if err != nil {
		return "", err
	}
	if val == nil {
		return "", nil
	}
	return string(val), nil
}

//GenerateCollectionIdentifier generate collection identifier and increases identifier by 1
func (e *Engine) GenerateCollectionIdentifier(collectionname []byte) ([]byte, error) {
	val, err := e.Store.Get([]byte(def.MetaCollectionidentifier))
	if err != nil {
		return []byte{}, err
	}
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.BigEndian.PutUint32(identifier, def.CollectionidentifierInitialcount)
		err := e.Store.Put([]byte(def.MetaCollectionidentifier), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint32(val)
		binary.BigEndian.PutUint32(val, uint32(identifier+1))
		err := e.Store.Put([]byte(def.MetaCollectionidentifier), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetCollectionIdentifier returns identifier for given collection
func (e *Engine) GetCollectionIdentifier(collection []byte) ([]byte, error) {
	if len(collection) == 0 {
		return []byte{}, def.CollectionNameEmpty
	}
	val, err := e.Store.Get([]byte(def.MetaCollection + string(collection)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateCollectionIdentifier(collection)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collection:collectionname = identifier
		err = e.Store.Put([]byte(def.MetaCollection+string(collection)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:collectionid:id=name
		err = e.Store.Put(append([]byte(def.MetaCollectionid), identifier...), collection)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else {
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//GetCollectionName returns collection name for given collection identifier
func (e *Engine) GetCollectionName(collectionIdentifier []byte) (string, error) {
	if len(collectionIdentifier) == 0 {
		return "", def.CollectionIdentifierEmpty
	}
	val, err := e.Store.Get([]byte(def.MetaCollectionid + string(collectionIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), nil
}

//GenerateNamespaceIdentifier generates namespace identifier value and increases identifier by 1
func (e *Engine) GenerateNamespaceIdentifier(namespace []byte) ([]byte, error) {
	val, err := e.Store.Get([]byte(def.MetaNamespaceidentifier))
	if err != nil {
		return []byte{}, err
	}
	//if there is no namespace id, generate a new one
	//TODO: move this logic to separate init file for performance
	if len(val) == 0 {
		identifier := make([]byte, 4)
		binary.BigEndian.PutUint32(identifier, def.NamespaceidentifierInitialcount)
		err := e.Store.Put([]byte(def.MetaNamespaceidentifier), identifier)
		if err != nil {
			return []byte{}, err
		}
		return identifier, nil
	} else {
		identifier := binary.BigEndian.Uint32(val)
		binary.BigEndian.PutUint32(val, uint32(identifier+1))
		err := e.Store.Put([]byte(def.MetaNamespaceidentifier), val)
		if err != nil {
			return []byte{}, err
		}
		return val, nil
	}
}

//GetNamespaceIdentifier returns identifier for given namespace
func (e *Engine) GetNamespaceIdentifier(namespace []byte) ([]byte, error) {
	if len(namespace) == 0 {
		return []byte{}, def.CollectionNameEmpty
	}
	val, err := e.Store.Get([]byte(def.MetaNamespace + string(namespace)))
	if err != nil {
		return []byte{}, err
	}
	//if len(val) is zero, generate a new identifier
	if len(val) == 0 {
		identifier, err := e.GenerateNamespaceIdentifier(namespace)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespace:namespace = identifier
		err = e.Store.Put([]byte(def.MetaNamespace+string(namespace)), identifier)
		if err != nil {
			return []byte{}, err
		}

		//insert meta:namespaceid:id=name
		err = e.Store.Put(append([]byte(def.MetaNamespaceid), identifier...), namespace)
		if err != nil {
			return []byte{}, err
		}

		return identifier, nil
	} else { //else send the value read from db
		//identifier := binary.LittleEndian.Uint32(val)
		return val, nil
	}
}

//GetNamespaceName return namespace name with given namespace identifier
func (e *Engine) GetNamespaceName(namespaceIdentifier []byte) (string, error) {
	if len(namespaceIdentifier) == 0 {
		return "", def.NamespaceIdentifierEmpty
	}
	val, err := e.Store.Get([]byte(def.MetaNamespaceid + string(namespaceIdentifier)))
	if err != nil {
		return "", err
	}
	return string(val), err
}

//GenerateUniqueID generates a 4 byte unique_id
func (e *Engine) GenerateUniqueID(collectionID []byte) ([]byte, error) {
	//unixTimeStamp := getUnixTimestamp() //returns 4 byte UNIX timestamp
	//macAddr := getMACAddress()          //returns 3 byte MAC Address
	//processID := getProcessID()         //returns 2 byte ProcessID
	//counter := generateRandomCount()    //returns 4 byte RANDOM count
	//uniqueID := append(unixTimeStamp, macAddr...)
	//uniqueID = append(uniqueID, processID...)
	//uniqueID = append(uniqueID, counter...)

	//key format: _uniqueid:dbid:colid:namespaceid=idcounter
	idKey := []byte(def.UniqueId + string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID))
	idCounterInBytes, err := e.Store.Get(idKey)
	if err != nil {
		return []byte{}, err
	}
	//if there is no id counter, generate a new one
	if len(idCounterInBytes) == 0 {
		//create 4 byte identifier for each document
		counterByte := make([]byte, 4)
		binary.BigEndian.PutUint32(counterByte, def.UniqueIdInitialcount)
		err := e.Store.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	} else {
		currentCount := binary.BigEndian.Uint32(idCounterInBytes)
		counterByte := make([]byte, 4)
		//increase count by 1 and write to db
		binary.BigEndian.PutUint32(counterByte, (currentCount + 1))
		//insert
		err := e.Store.Put(idKey, counterByte)
		if err != nil {
			return []byte{}, err
		}
		return counterByte, nil
	}
}

//GetIdentifiers returns database, collection and namespace identifiers for respective names given
//and generate new ones if they do not exist
func (e *Engine) GetIdentifiers(database string, collection string,
	namespace string) ([]byte, []byte, []byte, error) {

	dbID, err := e.GetDBIdentifier([]byte(database))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	collectionID, err := e.GetCollectionIdentifier([]byte(collection))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	namespaceID, err := e.GetNamespaceIdentifier([]byte(namespace))
	if err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return dbID, collectionID, namespaceID, nil
}

//SearchIdentifiers retrieves db,collection,namespace identifiers only if they exist
func (e *Engine) SearchIdentifiers(dbname string, collection string,
	namespace string) ([]byte, []byte, []byte, error) {

	dbID, err := e.Store.Get([]byte(def.MetaDb + string(dbname)))
	if len(dbID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	collectionID, err := e.Store.Get([]byte(def.MetaCollection + string(collection)))
	if len(collectionID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	namespaceID, err := e.Store.Get([]byte(def.MetaNamespace + string(namespace)))
	if len(namespaceID) == 0 || err != nil {
		return []byte{}, []byte{}, []byte{}, err
	}

	return dbID, collectionID, namespaceID, nil
}

//InsertDocument retrieves identifiers and inserts document to database
func (e *Engine) InsertDocument(collection string,
	data map[string][]byte, indices []string) error {

	fmt.Println("HERE")
	if len(e.DBName) == 0 || len(collection) == 0 || len(e.Namespace) == 0 {
		return def.NamesCannotBeEmpty
	}
	fmt.Println("HERE")

	//KV pair to insert in batch
	keyCache := make([][]byte, 0)
	valueCache := make([][]byte, 0)

	if _, ok := e.Session[e.DBName]; !ok {
		return def.DbDoesNotExist
	}

	if _, ok := e.Session[e.Namespace]; !ok {
		return def.NamespaceDoesNotExist
	}

	collectionID, err := e.GetCollectionIdentifier([]byte(collection))

	if err != nil {
		return err
	}

	//generate unique_id
	uniqueID, err := e.GenerateUniqueID(collectionID)
	if err != nil {
		return err
	}

	// indexer
	indexKey, indexValue, err := e.IndexDocument(collectionID, uniqueID, data, indices)
	if err != nil {
		return err
	}

	keyCache = append(keyCache, indexKey...)
	valueCache = append(valueCache, indexValue...)

	key := []byte(string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + string(uniqueID))

	indicesBytes, _ := json.Marshal(indices)
	data["_indices"] = indicesBytes
	dataInBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	keyCache = append(keyCache, key)
	valueCache = append(valueCache, dataInBytes)

	//insert in batch
	err = e.Store.PutBatch(keyCache, valueCache)
	if err != nil {
		return err
	}

	return nil
}

//SearchDocument queries document for given query params
func (e *Engine) SearchDocument(collection string,
	query []string) ([][]byte, error) {

	if len(e.DBName) == 0 || len(collection) == 0 || len(e.Namespace) == 0 {
		return [][]byte{}, def.NamesCannotBeEmpty
	}

	if _, ok := e.Session[e.DBName]; !ok {
		return [][]byte{}, def.DbDoesNotExist
	}

	if _, ok := e.Session[e.Namespace]; !ok {
		return [][]byte{}, def.NamespaceDoesNotExist
	}

	//here if collection doesn't exist, do not create new one
	collectionID, err := e.Store.Get([]byte(def.MetaCollection + collection))
	if err != nil {
		return [][]byte{}, err
	}

	//collectionID check is required here
	if len(e.DBID) == 0 || len(collectionID) == 0 || len(e.DBID) == 0 {
		return [][]byte{}, def.IdentifierNotFound
	}

	fmt.Println("[[engine.go]] evaluate postfix expression")
	rb, err := e.EvaluatePostFix(query, collectionID)
	if err != nil {
		return [][]byte{}, err
	}

	resultRoaring := rb.(roaring.Bitmap)

	//retrieve document keys for search
	fmt.Println("result roaring size -->", resultRoaring.GetSerializedSizeInBytes())
	fmt.Println("length of roaring -->", len(resultRoaring.ToArray()))
	searchKeys := make([][]byte, 0)
	searchKeyLength := len(resultRoaring.ToArray())
	uniqueIDArr := resultRoaring.ToArray() //get all IDs
	//get all documents keys
	for i := 0; i < searchKeyLength; i++ {
		uniqueIDByte := make([]byte, 4)

		binary.BigEndian.PutUint32(uniqueIDByte, uniqueIDArr[i])
		documentKeys := []byte(string(e.DBID) + ":" + string(collectionID) + ":" + string(e.NamespaceID) + ":" + string(uniqueIDByte))
		searchKeys = append(searchKeys, documentKeys)
	}
	resultArr, err := e.Store.GetBatch(searchKeys)
	if err != nil {
		return [][]byte{}, err
	}
	return resultArr, nil
}
