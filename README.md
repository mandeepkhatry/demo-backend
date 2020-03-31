
## Description
-----
### A working POC for Collateral Intelligence (backend).

## Run Locally
-----

### Linux

```1. chmod +x ./command.sh```

```2. ./command.sh```

## Current Deployment
-----

```
Digital Ocean
IP Address: 139.59.66.121

```

## API Endpoints
-----
```/form```  
``` 
Methods: POST

Eg. 

Request:
        {
            "table" : "nica",
            "schema" : {
                "title" : "account",
                "type" : "object",
                "required" : ["name"],
                "properties":{
                    "name":{
                        "type": "string",
                        "description" : "name",
                        "minLength" : 2,
                        "maxLength"  :100
                        
                    }
                }
                
            }
        }

Response:

        {
            "status": "config added successfully"
        }
```
```/form```  
``` 
Methods: GET



Response:

        {
            "status": "forms fetched successfully",
            "results": [
                "nabil"
            ]
        }
```
```/form/{table}```  
``` 
Methods: GET
		
Eg: /form/nica

Response:

        {
            "status": "schema fetched successfully",
            "results": [
                {
                    "schema": {
                        "properties": {
                            "name": {
                                "description": "name",
                                "maxLength": 100,
                                "minLength": 2,
                                "type": "string"
                            }
                        },
                        "required": [
                            "name"
                        ],
                        "title": "account",
                        "type": "object"
                    },
                    "table": "nib"
                }
            ]
        }
```


```/api/table/{table}```  
``` 
Methods: POST

Eg : /api/table/nica

Request:
        {
            "name" : "Mandeep",
            "age" : 50,
            "salary" : 50000,
            "address" : "Dhobighat"
        }

Response:
        
        {
            "status": "data added successfully"
        }

Note : Validation error results respons as follows
{
    "validation_error": {
        "name": [
            "name is required"
        ]
    }
}

{
    "validation_error": {
        "name": [
            "String length must be greater than or equal to 2"
        ]
    }
}

400 status code
```


```/api/table/{table}```  
``` 
Methods: GET

with params : Eg: age (?age=50)

Eg: /api/table/nica?age=50

Response:
        
        {
            "status": "query successful",
            "results": [
                {
                    "_id": 5,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                },
                {
                    "_id": 7,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                }
            ]
        }

without params

Eg : /api/table/nica (get all data of the table)

Response:
        
        {
            "status": "query successful",
            "results": [
                {
                    "_id": 5,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                },
                {
                    "_id": 7,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                },
                {
                    "_id": 9,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                },
                {
                    "_id": 10,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                }
            ]
        }

```

```/api/table/{table}/{id}```  
``` 
Methods: GET

Eg : /api/table/nica/5

Response:
        
       {
            "status": "query successful",
            "results": [
                {
                    "_id": 5,
                    "address": "Dhobighat",
                    "age": 20,
                    "name": "Bishal",
                    "salary": 50000
                }
            ]
    }



```

```/api/table/{table}/{id}```  
``` 
Methods: PATCH

Eg: /api/table/nica/2

Request :
        {
            "name" : "updatedName"
        }

Response:
        
        {
            "status": "data updated successfully"
        }
```

```/api/table/{table}/{id}```  
``` 
Methods: DELETE

Eg : /api/table/nica/2

Response:
        
        {
            "status": "data deleted successfully"
        }
```

```/api/query```  
``` 
Methods: POST

Request :
        {
            "query" : "@nica address=\"Dhobighat\" AND age<=50"
        }


Response:
        
        {
            "status": "query successful",
            "results": [
                {
                    "_id": 5,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                },
                {
                    "_id": 7,
                    "address": "Dhobighat",
                    "age": 50,
                    "name": "newTest",
                    "salary": 50000
                }
            ]
        }
```

