package main

import (
	"context"
	"fmt"
	"log"
	"logger/data"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/petarnenov/myutils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoUrl = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	fmt.Println("Logger Service")

	//connect to mongo
	mongoClient, err := connectToMongo()
	myutils.CheckNillError(err)
	client = mongoClient

	//create context in order to disconnect from mongo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	defer func() {
		err := client.Disconnect(ctx)
		myutils.CheckNillError(err)
	}()

	app := Config{
		Models: data.New(client),
	}

	//Register RPC server
	err = rpc.Register(&RPCServer{})
	myutils.CheckNillError(err)
	//start the RPC server
	go app.rpcListener()

	//start the gRPC server
	go app.gRPCListen()

	//start the web server
	app.serve()
}

func (app *Config) serve() {
	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	myutils.CheckNillError(err)
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoUrl)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	connection, err := mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Println("error connecting to mongo: ", err)
		return nil, err
	}

	return connection, nil
}

func (app *Config) rpcListener() error {
	log.Println("Starting RPC server on port: ", rpcPort)

	listen, err := net.Listen("tcp", "0.0.0.0:"+rpcPort)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("error accepting connection: ", err)
			continue
		}

		go rpc.ServeConn(conn)
	}
}
