# Zalora Assignment

## Problem Statement
* Build api with CRUD functionality for Ben&Jerry's Icecreams.
* Build a database schema to save icecream data.
* Secure api with authentication mechanism.

## Technologies Used
* Langauge used: Golang, go1.13.1
* Database choice: mysql, 8.0.17
* Dependencies used:
  * https://github.com/gin-gonic/gin to build routing and rest apis.
  * https://github.com/sirupsen/logrus to build structured error logging.
  * https://github.com/dgrijalva/jwt-go to build api authentication using JWT.
  * https://github.com/go-sql-driver/mysql to connect and interact with MySql DB.
  
## Database schema
![Image of DBSchema](https://github.com/shruti-madan09/zalora/blob/master/zalora.png)

## Code Structure & Implementation Details
### vendor
* Directory that contains code for all dependencies (e.g. gin-gonic, logrus, jwt-go, go-sql-driver).
* Path to this folder needs to be set in the $GOPATH for the dependencies to be accessible.

### src
* ***uploader package***: For bulk upload of ice cream data into the DB.
  * Contains a script that reads data from a json file and inserts it into the database.
  * List of icecream data will be iterated over and unique names of sourcing values, ingredients and dietary certifications will be stored in a map of string:boolean.
  * These maps will then be used to insert entries into tables ***sourcingvalue***, ***ingredient***, ***dietarycertification***.
  * ***'name'*** column in these tables has a unique constraint. For this reason ***insert ignore*** query will be run, to ignore errors if value already exists in table.
  * The list of icecream data is iterated over again and for each data
    * An entry will be made in the table ***product***.
    * Ids of sourcing values will be selected from table ***sourcingvalue*** and using them, entries will be made in ***product_sourcingvalue*** table.
    * Similar thing will be done for product ingredients.
  * All of the mysql queries needed in above steps will be executed in a single atomic transaction.
  * How to run
    * Navigate to the directory ***src/uploader***
    * Run the command: go run ***upload.go***
    * File icecream.json should be present in this folder.
* ***bennjeery package***: contains route, controller and model to implement CRUD endpoints.
  * **Create api**: Accepts data of one ice cream product and inserts it into DB.
    * Names of Sourcing values, Ingredients and Dietary Certifications that don't already exist in DB will be inserted.
    * An entry will be inserted in table ***product***.
    * Ids of sourcing values will be selected from table ***sourcingvalue*** and using them entries will be made in ***product_sourcingvalue*** table.
    * Similar thing will be done for product ingredients.
    * All of the mysql queries needed in above steps will be executed in a single atomic transaction.
    * File name: src/bennjerry/controller.go
    * Function name: ***CreateData***
    ```
    Sample Url: 0.0.0.0:8080/bennjerry/
    Request method: POST
    Post form data:
      * Key: "data"
      * Value:
      {
        "productId": "123",
        "name": "Name of Ice Cream",
        "description": "Description of Ice Cream",
        "story": "Story of Ice Cream",
        "image_closed": "Link of closed image",
        "image_open": "Link of open image",
        "sourcing_values": ["List", "of", "sourcing", "values"],
        "ingredients": ["List", "of", "ingredients"],
        "allergy_info": "Allergy related information",
        "dietary_certifications": "Name of dietary certifications"
      }
    Request headers:
      * Key: "JWT-TOKEN"
      * Value: "valid auth token"
    Response data:
      {
        "success": true/false,
        "id": 1/0, //id of inserted record, 0 incase of an error
        "message": "success or failure message"
      }
    ```
  * **Read api**: Accepts product id, fetches from DB and returns, all information corresponding to that product.
    * Information of a product will be returned only if it is not marked as inactive in DB.
    * File name: src/bennjerry/controller.go
    * Function name: ***ReadData***
    ```
    Sample Url: 0.0.0.0:8080/bennjerry/product_id/
    Request method: GET
    Request headers:
      * Key: "JWT-TOKEN"
      * Value: "valid auth token"
    Response data:
    {
      "message": "success or failure message",
      "success": true/false,
      "data": {
        "productId": "123",
        "name": "Name of Ice Cream",
        "description": "Description of Ice Cream",
        "story": "Story of Ice Cream",
        "image_closed": "Link of closed image",
        "image_open": "Link of open image",
        "sourcing_values": ["List", "of", "sourcing", "values"],
        "ingredients": ["List", "of", "ingredients"],
        "allergy_info": "Allergy related information",
        "dietary_certifications": "Name of dietary certifications"
      }
    }
    ```
  
  * **Update api**: Accepts a product_id, name of the fields and ice cream data and updates the provided fields with the provided values from ice cream data.
    * The user may want to update only specific properties of an ice cream product.
    * In such cases, un-marshalling of ice cream data will set default values for the missing properties. These missing properties can be overwritten in DB with default values of their datatypes.
    * Therefore, information regarding which fields to be updated is mandatory in the request.
    * All properties except sourcing values and ingredients will be updated in the **product** table.
    * If sourcing values/ingredients need to be update, the new list(s) will be compared with the data already in DB.
    * For any new sourcing values/ingredients record will be inserted in necessary tables.
    * If there are any sourcing values/ingredients that were already in DB but not present in the new list, such entries will be deleted from the DB.
    * All of the mysql queries needed in above steps will be executed in a single atomic transaction.
    * File name: src/bennjerry/controller.go
    * Function name: ***UpdateData***
    ```
    Sample Url: 0.0.0.0:8080/bennjerry/product_id/
    Request method: PUT
    Post form data:
      * Key: "data"
      * Value:
      {
        "productId": "123",
        "name": "New Name of Ice Cream",
        "description": "Description of Ice Cream",
        "sourcing_values": ["New", "List"],
      }
      * Key: "fields"
      * Value: "name,description,sourcing_values"
    Request headers:
      * Key: "JWT-TOKEN"
      * Value: "valid auth token"  
    Response data:
      {
        "success": true/false,
        "id": 1/0, //id of updated record, 0 incase of an error
        "message": "success/failure message"
      }
    ```
  * **Delete api**: Accepts product id and deletes(temporarily/permanently) all information corresponding to the product.
    * ***Soft Delete***: Product is simply marked as inactive (updating column ***'is_inactive'*** = 1) but not actually deleted from the DB. 
    * ***Permanent delete***: All information corresponding to the requested product_id is deleted from the table.
    * First the references of the product will be deleted from relation tables.
    * Then the actual record will be deleted from the ***product*** table.
    * All of the mysql queries needed in above steps will be executed in a single atomic transaction.
    * Once the above transaction has been successfully executed, any unused sourcing values, ingredients and dietary certifications will be deleted from the tables.
    * File name: src/bennjerry/controller.go
    * Function name: ***DeleteData***
    ```
    Sample Url: 0.0.0.0:8080/bennjerry/product_id/
    Request method: DELETE
    Request headers:
      * Key: "JWT-TOKEN"
      * Value: "valid auth token"
    Response data:
      {
        "success": true/false,
        "id": 1/0, // id of deleted record, 0 incase of an error
        "message": "success/failure message"
      }
    ```
   
* ***authenticator package***: Secures each api endpoint with authentication using JWT.
  * Every request to the application will be passed through a middleware which will look for a token in the request header.
  * The token will be parsed using a JWT signing key (the same that was used to create it) to check its validity.
  * If the token is valid, the remaining logic will be executed, else response with 401 error code will be returned.
  * File name: src/authenticator/authenticate.go
  * Function name: ***IsAuthorized***
  * ***token_generator package***
    * Ideally, a registration/login functionality should be built, which would return a jwt token for a user.
    * For simplicity, a token generator script has been created which generates a token valid for 30 minutes. The same can be used for testing out the apis.
    * How to run
      * Navigate to the package ***src/authenticator/token_generator/***
      * Run the command: go run ***generate.go***

* ***logger package***: To log errors.
  * Path to log file: ***logs/zalora.log***
  * Components of error log:
    * ***Bucket Name***: To identify which package of the code resulted in error, e.g. bennjerry, auth, mysql.
    * ***Identifier***: To identify the exact function which resulted in error.
    * ***Level***: error (for now set to error but can be used later to have other types of logs e.g. info)
    * ***Message***: The actual error message.

* ***Testing***
  * Unit tests for Create endpoint: src/bennjerry/test/create_test.go
    1. Calling api without auth token.
    2. Calling api with empty post form data.
    3. Calling api with invalid structure in post form data.
    4. Calling api with correct request data and request headers.
    5. Calling api with a product_id that already exists in DB.
  
  * Unit tests for Read endpoint: src/bennjerry/test/read_test.go
    1. Calling api without auth token.
    2. Calling api with a product_id that doesn't exist in the DB.
    3. Calling api with the correct product_id and request headers.
    
  * Unit tests for Update endpoint: src/bennjerry/test/update_test.go
    1. Calling api without auth token.
    2. Calling api with empty post form data.
    3. Calling api with invalid structure in post form data.
    4. Calling api with a product_id that doesn't exist in the DB.
    5. Calling api with the correct product_id, request data and request headers.
    
  * Unit tests for Delete endpoint: src/bennjerry/test/delete_test.go
    1. Calling api without auth token.
    2. Calling api with a product_id that doesn't exist in the DB.
    3. Calling api with correct product_id, request data and request headers but without permanent=1 query param.
    4. Calling api with correct product_id, request data, request headers and with permanent=1 query param.
    
* ***constants package***: Some of the information in the code has been kept as constants, to make them configurable.
  * Server related info (File name: ***src/constants/common.go***)
    * ***ServerHost***: Server host ip
    * ***ServerPort***: Server port number
  * MySql related info (File name: ***src/constants/db.go***)
    * ***MySQLDBName***: Database name
    * ***MySQLUserName***: MySql username
    * ***MySQLPassword***: MySql password
    * ***MySQLMaxOpenConnection***: Maximum number of open connections to mysql
    * ***MySQLMaxIdleConnection***: Maximum number of idle connections to mysql
  * Auth related info (File name: ***src/constants/auth.go***)
    * ***JWTSigningKey***: JWT signing key
    * ***JWTTokenKeyNameInHeader***: Key name to be passed in request header for sending auth token
  * Logger related info (File name: ***src/constants/logger.go***)
    * ***LoggerFilePath***: Path to log file
    * All relavant bucket names
  * Apis related info
    * All success/error messages to be sent in response or logs
