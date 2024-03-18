package main

import (
	"context"
	"errors"
	"example/internal"
	"fmt"
	"log"
	"net/http"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
)

func main() {
	interceptors := connect.WithInterceptors(internal.NewAuthInterceptor())
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
		interceptors,
	)
	req := connect.NewRequest(&greetv1.GreetRequest{Name: "Joker"})
	req.Header().Set("Acme-Tenant-Id", "1234")
	res, err := client.Greet(
		context.Background(),
		req,
	)
	if err != nil {
		fmt.Println(connect.CodeOf(err))
		if connectErr := new(connect.Error); errors.As(err, &connectErr) {
			fmt.Println(connectErr.Message())
			fmt.Println(connectErr.Details())
		}
	}
	if res == nil {
		return
	}
	log.Println(res.Msg.Greeting)
	log.Println(res.Header().Get("Greet-Version"))
	encoded := res.Header().Get("Greet-Emoji-Bin")
	if emoji, err := connect.DecodeBinaryHeader(encoded); err == nil {
		fmt.Println(string(emoji))
	}
	log.Println(res.Trailer().Get("Greet-Version"))
}
