package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/kostasdak/gomvc"
)

var c gomvc.Controller

func main() {

	// Load Configuration file
	cfg := gomvc.LoadConfig("./config/config.yml")

	// Connect to database
	db, err := gomvc.ConnectDatabase(cfg.Database.Dbuser, cfg.Database.Dbpass, cfg.Database.Dbname)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()

	//Start Server
	srv := &http.Server{
		Addr:    ":" + strconv.FormatInt(int64(cfg.Server.Port), 10),
		Handler: AppHandler(db, cfg),
	}

	fmt.Println("Web app starting at port : ", cfg.Server.Port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// App handler ... Builds the structure of the app !!!
func AppHandler(db *sql.DB, cfg *gomvc.AppConfig) http.Handler {

	// initialize
	c.Initialize(db, cfg)

	// load template files ... path : /web/templates
	c.CreateTemplateCache("home.view.tmpl", "base.layout.html")

	// *** Start registering urls, actions and models ***
	// home page
	c.RegisterAction("/", "", gomvc.ActionView, nil)
	c.RegisterAction("/home", "", gomvc.ActionView, nil)

	// create model for [products] table
	pModel := gomvc.Model{DB: db, IdField: "id", TableName: "products"}

	// view products ... /products for all records || /products/view/{id} for one product
	c.RegisterAction("/products", "", gomvc.ActionView, &pModel)
	c.RegisterAction("/products/view/*", "", gomvc.ActionView, &pModel)

	// create product actions ... this url has two actions
	// #1 View page -> empty form (no next url required)
	// #2 Post form data to create a new record in table [products] -> then redirect to [next] url
	c.RegisterAction("/products/create", "", gomvc.ActionView, &pModel)
	c.RegisterAction("/products/create", "/products", gomvc.ActionCreate, &pModel)

	// create edit product actions ... this url has two actions
	// #1 View page with product data -> edit form (no next url required)
	// #2 Post form data to update record in table [products] -> then redirect to [next] url
	c.RegisterAction("/products/edit/*", "", gomvc.ActionView, &pModel)
	c.RegisterAction("/products/edit/*", "/products", gomvc.ActionUpdate, &pModel)

	// create delete product actions ... this url has two actions
	// #1 View page with product data -> edit form [locked] to confirm detetion (no next url required)
	// #2 Post form data to delete record in table [products] -> then redirect to [next] url
	c.RegisterAction("/products/delete/*", "", gomvc.ActionView, &pModel)
	c.RegisterAction("/products/delete/*", "/products", gomvc.ActionDelete, &pModel)

	// create about page ... static page, no table/model, no [next] url
	c.RegisterAction("/about", "", gomvc.ActionView, nil)

	// contact page ... static page, no table/model, no [next] url
	c.RegisterAction("/contact", "", gomvc.ActionView, nil)

	// contact page POST action ... static page, no table/model, no [next] url
	// Register a custom func to handle the request/response using your oun code
	c.RegisterCustomAction("/contact", "", gomvc.HttpPOST, nil, ContactPostForm)
	return c.Router
}

// Custom handler for specific page and action
func ContactPostForm(w http.ResponseWriter, r *http.Request) {

	//test if I have access to products Model
	fmt.Print("Table Fields : ")
	fmt.Println(c.Models["products"].Fields)

	//read all records from table products
	rows, _ := c.Models["products"].GetAllRecords(100)
	fmt.Print("Select Rows : ")
	fmt.Println(rows)

	//print form fields
	fmt.Print("Form Fields : ")
	fmt.Println(r.Form)

	//test session -> send hello message
	c.GetSession().Put(r.Context(), "error", "Hello From Session")

	//redirect to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
