package main

import (
	_ "context"
	_ "log"

	_ "cloud.google.com/go/datastore"
	"github.com/xiahongze/pricetracker/server"
	_ "github.com/xiahongze/pricetracker/types"
)

func main() {
	server.Run()
	// entity := types.Enity{}

	// ctx := context.Background()
	// var (
	// 	alertType = "change"
	// 	email     = "x@x.com"
	// )
	// entity := types.Entity{
	// 	AlertType: &alertType,
	// 	Email:     &email,
	// }

	// dsClient, err := datastore.NewClient(ctx, "project-order-management")
	// if err != nil {
	// 	// Handle error.
	// 	log.Fatal("failed to new a dsClient")
	// }
	// entity.Save(ctx, "price-tracks", *dsClient)

	// t := test{Str: "", Number: 1, B: true}
	// b, e := json.Marshal(t)
	// if e == nil {
	// 	fmt.Printf("%s\n", b)
	// }
}
