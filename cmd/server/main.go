package main

import (
	"context"
	"example/internal"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "example/gen/greet/v1"        // generated by protoc-gen-go
	"example/gen/greet/v1/greetv1connect" // generated by protoc-gen-connect-go
)

type GreetServer struct{}

func (s *GreetServer) Greet(
	ctx context.Context,
	req *connect.Request[greetv1.GreetRequest],
) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	log.Println(req.Header().Get("Acme-Tenant-Id"))

	if err := req.Msg.ValidateAll(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if err := ctx.Err(); err != nil {
		return nil, err
	}
	greeting, err := doGreetWork(ctx, req.Msg)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnknown, err)
	}
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: greeting,
	})
	res.Header().Set("Greet-Version", "header:v1")
	res.Header().Set(
		"Greet-Emoji-Bin",
		connect.EncodeBinaryHeader([]byte("👋")),
	)
	res.Trailer().Set("Greet-Version", "trailer:v1")
	return res, nil
}

func main() {
	mux := http.NewServeMux()
	interceptors := connect.WithInterceptors(internal.NewAuthInterceptor())
	mux.Handle(greetv1connect.NewGreetServiceHandler(
		&GreetServer{},
		interceptors,
	))
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c, so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

func doGreetWork(ctx context.Context, req *greetv1.GreetRequest) (string, error) {
	log.Println(ctx)
	return fmt.Sprintf("Hello, %s!", req.Name), nil
}