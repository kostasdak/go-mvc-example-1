# Example 01 - gomvc package

Example 01 for MVC (Model View Controller) implementation with Golang using MySql database

## Overview
Web app with 5 pages :</br>
    - Home (static)</br>
    - Products -> View, Edit, Create, Delete product</br>
    - About (static)</br>
    - Contact (static)</br>

Database :</br>
`/database/example_1.sql`</br>

Steps :</br>
* Edit config file `configs/config.yml`
* Setup MySql database `example_1.sql`
* Load config file `configs/config.yaml`</br>
* Connect to MySql database</br>
* Write code to initialize your Models and Controllers</br>
* Write your standard text/Template files (Views)</br>
* Start your server</br>
* Enjoy</br>


### Edit configuration file

```
#UseCache true/false 
#Read files for every request, use this option for debug and development, set to true on production server
UseCache: false

#EnableInfoLog true/false
#Enable information log in console window, set to false in production server
EnableInfoLog: true

#InfoFile "path.to.filename"
#Set info filename, direct info log to file instead of console window
InfoFile: ""

#ShowStackOnError true/false
#Set to true to see the stack error trace in web page error report, set to false in production server
ShowStackOnError: true

#ErrorFile "path.to.filename"
#Set error filename, direct error log to file instead of web page, set this file name in production server
ErrorFile: ""

#Server Settings
server:
  #Listening port
  port: 8080

  #Session timeout in hours 
  sessionTimeout: 24

  #Use secure session, set to tru in production server
  sessionSecure: true

#Database settings
database:
  #Database name
  dbname: "golang"

  #Database server/ip address
  server: "localhost"

  #Database user
  dbuser: "root"

  #Database password
  dbpass: ""
```

