package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/greenblat17/yet-another-messenger/clients/internal/cli/processor"
)

const (
	// auth
	loginCommand    = "login"
	logoutCommand   = "logout"
	registerCommand = "register"

	// chat
	chatHistoryCommand = "chat-history"
	sendMessageCommand = "send-msg"

	// user
	profileCommand       = "profile"
	updateProfileCommand = "update-profile"

	// friendship
	sendFriendRequestCommand   = "send-friend-request"
	acceptFriendRequestCommand = "accept-friend-request"
	rejectFriendRequest        = "reject-friend-request"
	removeFriendCommand        = "remove-friend"
	getFriendsCommand          = "get-friends"

	helpCommand = "help"
	exitCommand = "exit"
)

var (
	ErrCommandNotFound = errors.New("invalid command")
	ErrCommandEmpty    = errors.New("command is not set")
	ErrCommandReading  = errors.New("error reading command")
)

type command struct {
	name        string
	description string
	call        func(context.Context, map[string]string) (any, error)
}

// initCommandList - описание для команд
func initCommandList(processor *processor.CommandProcessor) []command {
	return []command{
		{
			name:        loginCommand,
			description: "Вход",
			call:        processor.LoginUser,
		},
		{
			name:        logoutCommand,
			description: "Выход",
			call:        processor.Logout,
		},
		{
			name:        registerCommand,
			description: "Зарегестрироваться",
			call:        processor.RegisterUser,
		},
		{
			name:        chatHistoryCommand,
			description: "История сообщений",
			call:        processor.GetChatHistory,
		},
		{
			name:        sendMessageCommand,
			description: "Отправить сообщение",
			call:        processor.SendMessage,
		},
		{
			name:        profileCommand,
			description: "Профиль",
			call:        processor.GetUser,
		},
		{
			name:        updateProfileCommand,
			description: "Обновить профиль",
			call:        processor.UpdateUser,
		},
		{
			name:        sendFriendRequestCommand,
			description: "Отправить заявку в друзья",
			call:        processor.SendFriendRequest,
		},
		{
			name:        acceptFriendRequestCommand,
			description: "Принять заявку в друзья",
			call:        processor.AcceptFriendRequest,
		},
		{
			name:        rejectFriendRequest,
			description: "Не принимать заявку в друзья",
			call:        processor.RejectFriendRequest,
		},
		{
			name:        removeFriendCommand,
			description: "Удалить из друзей",
			call:        processor.RemoveFriend,
		},
		{
			name:        getFriendsCommand,
			description: "Список друзей",
			call:        processor.GetFriends,
		},
		{
			name:        helpCommand,
			description: "Получить справку",
		},
		{
			name:        exitCommand,
			description: "Выйти из программы: использование exit или ctrl+c",
		},
	}
}

// help - справка со всеми командами
func (c *CLI) help() {
	for _, cmd := range c.commandList {
		fmt.Println("", cmd.name, cmd.description)
	}
}

// exit - остановка программы
func (c *CLI) exit() {
	c.Close()
	os.Exit(0)
}

// Close закрывает все каналы и ожидает завершения всех горутин.
func (c *CLI) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		close(c.errors)
		close(c.jobs)

		c.wg.Wait()
		c.running = false
	}
}
