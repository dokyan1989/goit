package grpc

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/dokyan1989/goit/app/hellogrpc/pb"
	"google.golang.org/grpc"
)

var (
	// tls        = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	// certFile   = flag.String("cert_file", "", "The TLS cert file")
	// keyFile    = flag.String("key_file", "", "The TLS key file")
	// jsonDBFile = flag.String("json_db_file", "", "A json file containing a list of features")
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	// services
	todoStore TodoStore
}

func NewServer(todoStore TodoStore) *server {
	return &server{todoStore: todoStore}
}

func (s *server) Serve() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = data.Path("x509/server_cert.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = data.Path("x509/server_key.pem")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials: %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTodoManagementServer(grpcServer, &todoSvc{todoStore: s.todoStore})
	grpcServer.Serve(lis)
}
