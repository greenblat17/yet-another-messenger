package main

import (
	"fmt"
	"log"

	"github.com/greenblat17/yet-another-messenger/clients/internal/app"
	"github.com/greenblat17/yet-another-messenger/clients/internal/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// connections
	authConn := createConnection(8080)
	defer func() { _ = authConn.Close() }()

	userConn := createConnection(8081)
	defer func() { _ = userConn.Close() }()

	friendshipConn := createConnection(8082)
	defer func() { _ = friendshipConn.Close() }()

	chatConn := createConnection(8083)
	defer func() { _ = chatConn.Close() }()

	// clients
	auth.UnimplementedAuthServiceServer{}
	//authClient := auth.NewAuthServiceClient(authConn)
	//userClient := user.NewUserServiceClient(userConn)
	//friendshipClient := friendship.NewFriendshipServiceClient(friendshipConn)
	//chatClient := chat.NewChatServiceClient(chatConn)

	client := cli.NewCommandClient(authClient, userClient, friendshipClient, chatClient)
	commands := cli.New(client)

	app.Run(commands)
}

func createConnection(port int) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%d", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
