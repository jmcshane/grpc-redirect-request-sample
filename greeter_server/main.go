/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/redirect-sample/helloworld"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 50051, "The server port")
	ip   = flag.String("ip", "127.0.0.1", "the ip address")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	if dec, ok := localMatch(in.Ip); !ok {

		conn, err := grpc.Dial(fmt.Sprintf("%s:50051", dec), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("did not connect: %v", err)
			return nil, err
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()
		log.Printf("Delegating request to %s", dec)
		return c.SayHello(ctx, in)
	}
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

// falls back to processing locally so we don't get stuck in a loop
func localMatch(reqIp string) (string, bool) {
	if reqIp == "" {
		return "", true
	}
	sDec, err := base64.StdEncoding.DecodeString(reqIp)
	if err != nil {
		return "", true
	}
	if string(sDec) == "127.0.0.1" {
		return "", true
	}
	fmt.Printf("My ip is %s and req ip is %s", *ip, string(sDec))
	return string(sDec), string(sDec) == *ip
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
