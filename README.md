# admincheckapi microservice

## Purpose

The API responsibility is to classify a JWT token in order to determine 
if the sender/owner of the token is a member of Azure admin group.

The main businnes logic entity is the table **CLIENT_ADMIN_GROUP**. It links
the **CLIENT** with a **ADMIN GROUP**. The server receices a request carrying
a JWT token. The question is:

- Does the token belong to an admin user?

The token is checked and decoded and a set of group ids is
extractd from the claims. The id may be stored in the caching DB. If yes
it means that the owner of the token is an admin group member. If not
then the effective group info is to be loaded from MS Azure graph
for each of the id loaded from JWT token. If some of them decodes
into well known group name then this token is an admin token
and the group id can be associated in the local DB with the client
to be used as cache in next requests.

## Entities

- **CLIENT_ADMIN_GROUPS**: maps admin group ID to client

## API methods and resources

- **POST:/client/{client}/admin/token** with Body:Token -> True, False
- **GET:/client/{client}/group/{group}/admin** -> True, False
- **GET:/client/{client}/admin/group** -> Read
- **POST:/client/{client}/admin/group/{group}** -> Create
- **DELETE:/client/{client}/admin/group/{group}** -> Delete
- **POST:/client/{client}/admin/auth/{method}** with Body:Claims -> Token

Remark:

(1) The JWT token used to access MS graph is the token of the tenent
of the application. It is not token of the client. The client has to enable access
to their groups. 

(2) Auth method is used to get the JWT token of the application uses secret
stored in the config.yaml file. It can be overloaded with the value from env
variable where the real secret value is provided. Later this value may
be obtained from other source of secrets like AWS key store.

(3) The first level cache stores the group id mappings in memeory. The are
materialized in 2nd level cache - relational DB. Both of the caches
have the main source of the information: Azure graph. A scheme of cache
ivalidation may be considered as additional feature.

Additional technical methods are to be added like:

- **GET:/system/health**
- **GET:/system/alive**
- **GET:/system/stat**
- **GET:/system/version**

## Functional components

- cmd: main function location
- api/secretstore: keeps JWT token of the application
- api/token: handles group id extraction from JWT token
- api/graph: method used to acces Azure graph to decode group id to group name
- api/config: loading config from local yaml file and oeverwriting values with env variables
- api/backend: handles basic relational database access
- api/controller: it links routes with repository handles processing input and output JSON structures
- api/model: main business entity defined here is CLIENT_ADMIN_GROUP with GORM injection
- api/repository: handles GORM access to backend relational database, Postgres now
- api/resource: defines API interface structures in JSON
- api/router: creates routes for resources and defines methods
- api/server: HTTP and HTTPS request routing
- db/postgres: Postgres db configuration and maintenance scripts

## Building

The build process uses Makefile and Dockerfile. The 1st one is used to uild
the image. The 2nd one is used by the CI inorder to build the Docker image.

The basic build commands are:

- make build: it build executable in the local directory
- make test: it starts unit tests
- make integrtest: runs integraton tests from test/integr folder
- make image: builds Docker image
- make clean: cleans old artefacts
- etc.

## Testing

### Unit tests

It is enough to run make test (or make vtest for verbose output) in order to start the unit test.
The connection to Postgres DB must exist as it it tested. It is fixed as:

- test/test @ argonadmindb @ localhost

It is used in the unit tests only if the env variable POSTGRES or MYSQL are defined. Otherwise
only mock DB is used for unit testing.

The unit tests using dependendent components from the environment are started with command:

```
make envtest
```

### Integration tests

The curl scripts contain popular scenarios of events like create, read, delete admin
memebership for a group or tocken check with positive or negative result. The scripts
are stored n the folder:

- test/integr

The integration tests are started with command:

```
make integrtest
```

The Postgres DB instance shall be accesible as in case of unit tests.

The executable admincheckapi can be started with commandline like: ./admincheckapi
in a separate sesion so that the log can be printed on the stdout. In other session
the scripts can be started. The log shows the status of each request, URL path used
and timing statistics of the request serverd by the API server.

## Configuration

The configuration is stores in the confi.yaml file. This files load yaml definitions of env variables.
They are used later in some ponts of the process to make various decisions. An example is:

```
loggers:
- log:
  kind: log
  env:
    logrus: Info  
    httplog: True
    gorm: True
providers:
- msad:
  kind: msad
  env:
    tenant_id: XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX
    client_id: YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY
    client_secret: ZZZZZ~ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ
    admin_group_name: NeonAdmin
servers:
- http:
  kind: http
  env:
    port: 1234
    address: localhost
sqloptions:
- sql:
  kind: sql
  env:
    Max_Idle_Conns: 10
    Max_Open_Conns: 100
    Max_Lifetime: 1
backends:
- postgres:
  kind: postgres
  env:
    user: test
    pass: test
    dbname: argonadmindb
    host: localhost
    port: 5432
```

The variable names are created with kind field value concatenated with '_' and value env list.
The values of the env variables are loaded from each yaml value.

For example the Postgres credentials are defined as variables:

- POSTGRES_USER
- POSTGRES_PASS
- POSTGRES_DBNAME
- POSTGRES_HOST
- POSTGRES_PORT

The corresponding MYSQL variables are:

- MYSQL_USER
- MYSQL_PASS
- MYSQL_DBNAME
- MYSQL_HOST
- MYSQL_PORT

They are created from the yaml as listed above. Later, the DB connect stringer module uses them
to build the appropriate Connect String used to open DB connection.

If the variable exists in the env then its value is not loaded from yaml but it is left intact
and the module using them will used the value from the environment. This is the way
how some sectet credentials cvan be passed from K8S environment and bot from yaml config file.

### Loggers

This section defines the loggin levels for 3 main components:

- **logrus**: this is general logging component used all over the service. It can be used with levels of logging like Trace, Debug, Info, etc.

- **httplog**: logging event of each request with time and basic info. There is timing logged there as well. It may be True or False.

- **gorm**: this part is to trigger logging of SQL DB queries when used with ther than Silent level. It may be True or False

### Providers

This section defines parameters necessary to connect to identity provides like Miscrosoft Active Directory.

### Servers

As HTTP and HTTPS servers are in scope, this section defines necessary parameters like host address
or port numbers. IP and DNS addresses are both allowed. 

### Backends

Several backends like Postgres or MYSQL are very easy to be used with GORM so this section triggers
usage of one of them.
