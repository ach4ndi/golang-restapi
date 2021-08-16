
For this test im using go 1.1.6 windows version with GO111MODULE=off, VSCODIUM as code editor, Markdown Monster as markdown editor, PostMan as Request Tool, and Chrome as Web Browser. Thanks to <a href="https://github.com/victorsteven/Go-JWT-Postgres-Mysql-Restful-API" target="_blank">this tutorial</a> for make me understand about using golang rest api.

#### Running The Application

    go run main.go

Before run the application, you should install go first from <a href="https://golang.org/doc/install" target="_blank">Here</a>. Make sure GOPATH and GO111MODULE=off is already defined on OS ENVIROMENT.

After setting OS ENVIROMENT, you must installing some modules for this application, you can run the following commands:

    go get github.com/badoux/checkmail
    go get github.com/jinzhu/gorm
    go get golang.org/x/crypto/bcrypt
    go get github.com/dgrijalva/jwt-go
    go get github.com/gorilla/mux
    go get github.com/joho/godotenv
    go get gopkg.in/go-playground/assert.v1
    go get github.com/google/uuid
 
If using mysql:

    go get github.com/jinzhu/gorm/dialects/mysql

If using postgres:

    go get github.com/jinzhu/gorm/dialects/postgres
    
If using sqlite3:

    go get github.com/mattn/go-sqlite3
    
After Installing modules, set some configuration on ENV :

    WEB_PORT = :8780
    IMG_DIR = F:\a\majootestapi // directory to store image
    IMG_DIR_PRODUCT = \images\products 
    IMG_LIMIT = 10
    SEED_LOAD = 0 

    #Postgres (postgres), mysql (mysql), sqlite (sqlite3)
    DB_HOST=127.0.0.1
    DB_DRIVER=postgres
    API_SECRET=98hbun98h #Used for creating a JWT. Can be anything 
    DB_USER=steven
    DB_PASSWORD=
    DB_NAME=majoo_minipos_test
    DB_PORT=5432

> If you want change DB Drive, just change DB_DRIVER and its configuration on .env file (Since is using GORM as ORM). Also you can change limitation of file upload (in MB), directory img root, web port (started with : ) and directory product IMG. Seed Load is for regenerate sample data, turn off by set to 0 for not always refresh database and remove recent changes on it.

#### End-Point

##### Users 

JWT Token

> you need defined header "Authorization: 'bearer {token}'"

[GET] /users 

> for listing all users, no need token to access this end point

Example Result :

    [
        {
            "id": 1,
            "nickname": "Hyna Tester",
            "email": "hyna1234@gmail.com",
            "password": "$2a$10$anSBUMbgB38ebQWlhZ.q1.dNFcJ/hT9bvMNVI1Z9qgljB5F1cxeMa",
            "created_at": "2021-08-13T08:49:43.5176436+07:00",
            "updated_at": "2021-08-13T08:49:43.5176436+07:00"
        },
        {
            "id": 2,
            "nickname": "Martin Brunch",
            "email": "martin@gmail.com",
            "password": "$2a$10$.Kz3Q3yzhCAS02vmn7ao5OgcfHwlm69i.6esMpXRf5FvL3mM40waG",
            "created_at": "2021-08-13T08:49:43.6546376+07:00",
            "updated_at": "2021-08-13T08:49:43.6546376+07:00"
        }
    ]

[POST] /login 

> for login, no need token to access this end point

    {
    	"email":"temp1@gmail.com",
    	"password":"password"
    }

Example Result :

    "token"

[POST] /users 

> for creating new user, no need token to access this end point

    {
    	"nickname": "Template User 1",
    	"email":"temp1@gmail.com",
    	"password":"password"
    }

Example Result :

> it will return user detail

[PUT] /users/{id}

> for update selected user, need token to access this end point

    {
    	"nickname": "Template User 1",
    	"email":"temp1@gmail.com",
    	"password":"password"
    }

Example Result :

> it will return user detail

[DELETE] /users/{id}

> for delete selected user, need token to access this end point

    {
    	"nickname": "Template User 1",
    	"email":"temp1@gmail.com",
    	"password":"password"
    }

Example Result :

> if get respond 204, is mean deletion is success (no data to return), if fail it will return error response

[GET] /products

> for listing all products, no need token to access this end point

Example Result :

    [
        {
            "id": 1,
            "code": "BA097094101",
            "name": "TestGambar",
            "user": {
                "id": 1,
                "nickname": "Hyna Tester",
                "email": "hyna1234@gmail.com",
                "password": "$2a$10$anSBUMbgB38ebQWlhZ.q1.dNFcJ/hT9bvMNVI1Z9qgljB5F1cxeMa",
                "created_at": "2021-08-13T08:49:43.5176436+07:00",
                "updated_at": "2021-08-13T08:49:43.5176436+07:00"
            },
            "user_id": 1,
            "description": "tidak ada",
            "default_price": 0,
            "pic_name": "product_4812b6911a504043a30e7392e97527e8.png",
            "created_at": "2021-08-13T08:54:42.1071007+07:00",
            "updated_at": "2021-08-13T08:54:42.1071007+07:00"
        }
    ]

[POST] /products 

> for creating new products, need token to access this end point
> 
> is need form-data request [code, name, description, default_price, myImage (File)]

Example Result :

    {
        "id": 1,
        "code": "BA097094101",
        "name": "TestGambar",
        "user": {
            "id": 1,
            "nickname": "Hyna Tester",
            "email": "hyna1234@gmail.com",
            "password": "$2a$10$anSBUMbgB38ebQWlhZ.q1.dNFcJ/hT9bvMNVI1Z9qgljB5F1cxeMa",
            "created_at": "2021-08-13T08:49:43.5176436+07:00",
            "updated_at": "2021-08-13T08:49:43.5176436+07:00"
        },
        "user_id": 1,
        "description": "tidak ada",
        "default_price": 0,
        "pic_name": "product_4812b6911a504043a30e7392e97527e8.png",
        "created_at": "2021-08-13T08:54:42.1071007+07:00",
        "updated_at": "2021-08-13T08:54:42.1071007+07:00"
    }
    
[PUT] /products/{id}

> for update selected product, need token to access this end point and same userId
> 
> is need form-data request [code, name, description, default_price, update_image, myImage (File)]
> update_image flag is used when image is changed or not. if true then it process changing image path and upload new image or delete it

Example Result :

    {
        "id": 1,
        "code": "BA097094101",
        "name": "TestGambar",
        "user": {
            "id": 1,
            "nickname": "Hyna Tester",
            "email": "hyna1234@gmail.com",
            "password": "$2a$10$anSBUMbgB38ebQWlhZ.q1.dNFcJ/hT9bvMNVI1Z9qgljB5F1cxeMa",
            "created_at": "2021-08-13T08:49:43.5176436+07:00",
            "updated_at": "2021-08-13T08:49:43.5176436+07:00"
        },
        "user_id": 1,
        "description": "tidak ada",
        "default_price": 0,
        "pic_name": "product_4812b6911a504043a30e7392e97527e8.png",
        "created_at": "2021-08-13T08:54:42.1071007+07:00",
        "updated_at": "2021-08-13T08:54:42.1071007+07:00"
    }

[DELETE] /products/{id}

> for delete selected product, need token to access this end point and same userId

Example Result :

> if get respond 204, is mean deletion is success (no data to return), if fail it will return error response