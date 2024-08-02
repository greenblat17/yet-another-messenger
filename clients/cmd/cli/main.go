package main

import (
	"fmt"
	"log"

	"github.com/greenblat17/yet-another-messenger/clients/internal/app"
	"github.com/greenblat17/yet-another-messenger/clients/internal/cli"
	"github.com/greenblat17/yet-another-messenger/clients/internal/cli/processor"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/auth"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/chat"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/friendship"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/clients/api/proto/user"
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
	authClient := auth.NewAuthServiceClient(authConn)
	userClient := user.NewUserServiceClient(userConn)
	friendshipClient := friendship.NewFriendshipServiceClient(friendshipConn)
	chatClient := chat.NewChatServiceClient(chatConn)

	commandClient := processor.NewCommandClient(authClient, userClient, friendshipClient, chatClient)
	cliApp := cli.New(commandClient)

	app.Run(cliApp)
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
