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
	cfg := gomvc.ReadConfig("./config/config.yml")

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

	// initialize controller
	c.Initialize(db, cfg)

	// load template files ... path : /web/templates
	// required : homepagefile & template file
	// see [template names] for details
	c.CreateTemplateCache("home.view.tmpl", "base.layout.html")

	// *** Start registering urls, actions and models ***

	// RegisterAction(url, next, action, model)
	// url = url routing path
	// next = redirect after action complete, use in POST actions if necessary
	// model = database model object for CRUD operations

	// home page : can have two urls "/" and "/home"
	c.RegisterAction(gomvc.ActionRouting{URL: "/"}, gomvc.ActionView, nil)
	c.RegisterAction(gomvc.ActionRouting{URL: "/home"}, gomvc.ActionView, nil)

	// create model for [products] database table
	// use the same model for all action in this example
	pModel := gomvc.Model{DB: db, PKField: "id", TableName: "products"}

	// view products ... / show all products || /products/view/{id} for one product
	c.RegisterAction(gomvc.ActionRouting{URL: "/products"}, gomvc.ActionView, &pModel)
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/view/*"}, gomvc.ActionView, &pModel)

	// build create product action ... this url has two actions
	// #1 View page -> empty product form no redirect url (no next url required)
	// #2 Post form data to create a new record in table [products] -> then redirect to [next] url -> products page
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/create"}, gomvc.ActionView, &pModel)
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/create", NextURL: "/products"}, gomvc.ActionCreate, &pModel)

	// build edit product actions ... this url has two actions
	// #1 View page with the product form -> edit form (no next url required)
	// #2 Post form data to update record in table [products] -> then redirect to [next] url -> products page
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/edit/*"}, gomvc.ActionView, &pModel)
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/edit/*", NextURL: "/products"}, gomvc.ActionUpdate, &pModel)

	// build delete product actions ... this url has two actions
	// #1 View page with the product form -> edit form [locked] to confirm detetion (no next url required)
	// #2 Post form data to delete record in table [products] -> then redirect to [next] url -> products page
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/delete/*"}, gomvc.ActionView, &pModel)
	c.RegisterAction(gomvc.ActionRouting{URL: "/products/delete/*", NextURL: "/products"}, gomvc.ActionDelete, &pModel)

	// build about page ... static page, no table/model, no [next] url
	c.RegisterAction(gomvc.ActionRouting{URL: "/about"}, gomvc.ActionView, nil)

	// build contact page ... static page, no table/model, no [next] url
	c.RegisterAction(gomvc.ActionRouting{URL: "/contact"}, gomvc.ActionView, nil)

	// build contact page POST action ... static page, no table/model, no [next] url
	// Demostrating how to register a custom func to handle the http request/response using your oun code
	// and handle POST data and have access to database through the controller and model object
	c.RegisterCustomAction(gomvc.ActionRouting{URL: "/contact"}, gomvc.HttpPOST, nil, ContactPostForm)
	return c.Router
}

// Custom handler for specific page and action
func ContactPostForm(w http.ResponseWriter, r *http.Request) {

	//test : I have access to products model !!!
	fmt.Print("\n\n")
	fmt.Println("********** ContactPostForm **********")
	fmt.Println("Table Fields : ", c.Models["/products"].Fields)

	//read data from table products (Model->products) even if this is a POST action for contact page
	fmt.Print("\n\n")
	rows, _ := c.Models["/products"].GetRecords([]gomvc.Filter{}, 100)
	fmt.Println("Select Rows Example 1 : ", rows)

	//read data from table products (Model->products) even if this is a POST action for contact page
	fmt.Print("\n\n")
	id, _ := c.Models["/products"].GetLastId()
	fmt.Println("Select Rows Example 1 : ", id)

	//read data from table products (Model->products) with filter (id=1)
	fmt.Print("\n\n")
	var f = make([]gomvc.Filter, 0)
	f = append(f, gomvc.Filter{Field: "id", Operator: "=", Value: 1})
	rows, err := c.Models["/products"].GetRecords(f, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Select Rows Example 2 : ", rows)

	//test : Print Posted Form fields
	fmt.Print("\n\n")
	fmt.Println("Form fields : ", r.Form)

	//test : Set session message
	c.GetSession().Put(r.Context(), "error", "Hello From Session")

	//redirect to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
