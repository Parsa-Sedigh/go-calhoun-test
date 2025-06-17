package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var (
	templates struct {
		Orders struct {
			New    *template.Template
			Review *template.Template
		}
		Campaigns struct {
			Show *template.Template
		}
	}
)

const (
	formTemplateHTML = `
		<div class="w-full mb-6">
<label class="block uppercase tracking wide text-grey-darker text-xs font-bold mb-2" for="{{.Name}}">
{{. Label}}
</label>
<input class="bg-grey-lighter appearance-none border-2
border-grey-lighter hover: border-orange rounded w-full py-2 px-4
OUTLINE
text-grey-darker leading-tight" name="{{.Name}}" type="{{.Type}}"
`
	stripeSecretKey = "sk_test_..."
	stripePublicKey = "pk_test_..."
)

func init() {
	formTemplate := template.Must(template.New("").Parse(formTemplateHTML))

	templates.Orders.New = template.Must(template.New("new_order.gohtml")).Funcs(template.FuncMap{
		"form_for": func(strct interface{}) (template.HTML, error) {
			return form.HTML(formTemplate, strct)
		},
	}).ParseFiles("./templates/new_order.gohtml")

	templates.Orders.Review = template.Must(template.ParseFiles("./templates/review_order.gohtml"))
	templates.Campaigns.Show = template.Must(template.ParseFiles("./templates/show_campaign.gohtml"))
}

func main() {
	defer db.DB.Close()

	db.CreateCampaign(time.Now(), time.Now().Add(time.Hour), 1200)

	mux := http.NewServeMux()
	resourceMux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets/"))

	// NOTE: The html folder is not directly accessible in the fileserver
	mux.Handle("/img/", fs)
	mux.Handle("/css/", fs)
	mux.Handle("/favicon.ico", http.FileServer(http.Dir("./assets/img/")))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = cleanPath(r.URL.Path)
		resourceMux.ServeHTTP(w, r)
	})
	resourceMux.HandleFunc("/", showActiveCampaign)
	resourceMux.Handle("/campaigns/", http.StripPrefix("/campaigns/", campaignsMux()))
	resourceMux.Handle("/orders/", http.StripPrefix("/orders", ordersMux()))

	port := os.Getenv("SWAG_PORT")
	if port == "" {
		port = "3000"
	}

	addr := fmt.Sprintf(":%s", port)

	// NOTE: `resourceMux` is called by `mux` in handler of "/" above.
	log.Fatal(http.ListenAndServe(addr, mux))
}

func ordersMux() http.Handler {
	// The order mux expects the order to be set in the context and the ID to be trimmed from the path.
	ordMux := http.NewServeMux()
	ordMux.HandleFunc("/confirm", confirmOrder)
	ordMux.HandleFunc("/", showOrder)

	// trim the id from the path, set the campaign in the ctx, and call the cmpMux.
	// This acts like a middleware.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payCusID, path := splitPath(r.URL.Path)

		order, err := db.GetOrderViaPayCus(payCusID)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		ctx := context.WithValue(r.Context(), "order", order)

		r = r.WithContext(ctx)
		r.URL.Path = path
		ordMux.ServeHTTP(w, r)
	})
}

func campaignsMux() http.Handler {
	// Paths like /campaigns/:id/orders/new are handled here, but most of the path - the /campaigns/:id/orders part - is stripped and
	// processed beforehand.
	cmpOrderMux := http.NewServeMux()
	cmpOrderMux.HandleFunc("/new", newOrder)
	cmpOrderMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createOrder(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// the campaign mux expects the campaign to be set in the context and the ID to be trimmed from the path.
	cmpMux := http.NewServeMux()
	cmpMux.Handle("/orders/", http.StripPrefix("/orders", cmpOrderMux))

	// trim the id from the path, set the campaign in the ctx, and call the cmpMux.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr, path := splitPath(r.URL.Path)

		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		campaign, err := db.GetCampaign(id)
		if err != nil {
			http.NotFound(w, r)

			return
		}

		ctx := context.WithValue(r.Context(), "campaign", campaign)
		r = r.WithContext(ctx)
		r.URL.Path = path
		cmpMux.ServeHTTP(w, r)
	})
}

func cleanPath(pth string) string {
	pth = path.Clean("/" + pth)
	if pth[len(pth)-1] != '/' {
		pth += "/"
	}

	return pth
}

func splitPath(pth string) (head, tail string) {
	pth = cleanPath(pth)
	parts := strings.SplitN(pth[1:], "/", 2)

	if len(parts) < 2 {

	}
}
