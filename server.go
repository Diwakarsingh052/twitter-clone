package main

import (
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"twitter-clone/controller"
	"twitter-clone/graph"
	"twitter-clone/graph/generated"
	"twitter-clone/graph/model"
	"twitter-clone/middleware"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println(os.Getenv("NAME"))
	us, err := model.NewUserService(os.Getenv("NAME")+":"+os.Getenv("PASSWORD")+"@/"+os.Getenv("DATABASE")+"?charset=utf8&parseTime=True&loc=Local")
	must(err)
	defer us.Close()
	us.AutoMigrate()
	usersC := controller.NewUsers(us)

	const defaultPort = "8080"
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	requireUserMw := middleware.RequireUser{Us: us}
	requireUser := requireUserMw.ApplyFn(usersC.SayHello)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Db: us.DB, MW: requireUserMw}}))
	requireUser = requireUserMw.Apply(srv)

	http.Handle("/query", requireUser)

	http.HandleFunc("/signup", usersC.Create)

	http.HandleFunc("/login", usersC.Login)



	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
