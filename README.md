# sakila-inventory

[![Travis status]][Travis]

a sample REST API for [sakila](https://dev.mysql.com/doc/sakila/en/ "main page") database that uses the [go-chi](https://github.com/go-chi/chi "chi") library with the subpackages [middleware](https://github.com/go-chi/chi/blob/master/middleware "middleware") and [render](https://github.com/go-chi/render "render").

### Prerequisites

* Latest Go's version

* MySQL 8.0.18 or higher

* Chi and subpackages, MySQL driver and go-cmp. If you want to get these in the go's package-style, in you need to run in your system command prompt:
```
go get -u github.com/go-chi/chi github.com/go-chi/render github.com/go-sql-driver/mysql github.com/google/go-cmp
```

### Installation, run and testing

If you have set the GOPATH and GOBIN environment variables, you can run an installation by placing the source code in the /src folder and running

```Go
> go install
```

Otherwise, main.go can be executed to build a temporal file by running

```Go
> go run main.go (+args)
```

To test if the router is working correctly, you can and run:

```Go
> go test -v (+args)
```

#### Arguments

The default arguments in the main file are used for my own pc settings, so if you're here and by some reason you want to run the code, you need to replace them or by adding them in the command prompt.

| Argument | Description | type | example(default) |
| :---: | ------ | ------ | ------ |
| user | Database user | string | jean |
| pass | Database password | string | 64fa4632 |
| schema | Database schema | string | sakila |
| protocol | Database protocol connection | string | tcp |
| ip | Database IP | string | localhost |
| port | Database port | unsigned integer | 3306 |
| server-port | port where API will run | unsigned integer | 80 |


# Access routes

The access routes for the API are:
1.  /API/rest/inventory/
2.  /API/rest/inventory/search/...
3.  /API/rest/inventory/film/...
1.  /API/rest/categories/
1.  /API/rest/actors/

In the inventory section, there are different URL query options that are divided by filter and pagination:

#### filter

| parameter | description | type |
| :---: | ------ | ------ |
| q | search string, either in the title or description of the film | string |
| year | release year | unsigned integer |
| lang | film language | string |
| rdur | rental duration | unsigned decimal |
| price | rental price | unsigned decimal |
| len | film duration | unsigned integer |
| replcos | replacement cost | unsigned decimal |
| rating | film age ratings | set(G, PG, PG-13, R, NC-17) |
| spec | additional features | array(Trailers, Commentaries, Deleted Scenes, Behind the Scenes) |

NOTE:
The additional features can be separated by a "+" sign and the blanks are filled with "_"

#### pagination

| parameter | description | type |
| :---: | ------ | ------ |
| page | Access the requested page number. If this option is not included, the first page is returned | unsigned integer |
| lim | number of results per page. by default 20 | unsigned integer |
| orderby | contains a filter listed above to apply before paging | filter param |
| ord | order `asc` (ascending) or `desc` (descending), by default adc | string |

### Usage examples

The main section is the inventory of film, so you can pretty much do what you want when add, read, modify or delete films. Others sections  (categories and actors) are minimized just for show how the data would mapped in a JSON file, so not all functions and stored procedures in the sql script were implemented.

```
to get all films:                           localhost/API/rest/inventory
to get some films by filtering and paging:  localhost/API/rest/inventory/search/q=ace&price=2.99/page=2&ord=desc
to get the last five films added:           localhost/API/rest/inventory/search/q=/lim=5&orderby=id&ord=desc
to get a specific film:                     localhost/API/rest/inventory/film/2/ace_goldfinger
to get all categories:                      localhost/API/rest/categories
to get all actors:                          localhost/API/rest/actors
```

## HTTP status codes
These status codes vary according to the http method requested and the response obtained

### GET
| code | description |
| :---: | --- |
| 200 | request has succeeded |
| 404 | can't find the requested resource |
| 500 | Database connection can't be done correctly |
| 503 | Database is in maintenance, that means the connection can't be done correctly or something in the schemas, views or stored procedures has changed |
### POST
| code | description |
| :---: | --- |
| 201 | resources are created successfully. no response body returned |
| 400 | Request data invalid or expressed without a valid JSON format |
| 403 | data insertion denied |
| 422 | Wrong input(s) |
| 500 | Database connection can't be done correctly |
### PUT
| code | description |
| :---: | --- |
| 202 | updated successfully. no response body returned |
| 400 | Request data invalid or expressed without a valid JSON format |
| 404 | can't find the requested resource |
| 422 | Wrong input(s) |
| 500 | Database connection can't be done correctly |
### DELETE
| code | description |
| :---: | --- |
| 204 | deleted successfully. no response body returned |
| 404 | can't find the requested resource |
| 500 | Database connection can't be done correctly |

[Travis]: https://travis-ci.com/rohan5564/sakila-inventory
[Travis status]: https://travis-ci.com/rohan5564/sakila-inventory.svg?token=KdHNqeoecbGixP8VXSSy&branch=master