package cli

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/greenblat17/yet-another-messenger/clients/internal/cli/processor"
)

const (
	loginCommand         = "login"
	logoutCommand        = "logout"
	registerCommand      = "register"
	chatHistoryCommand   = "chat-history"
	sendMessageCommand   = "send-msg"
	profileCommand       = "profile"
	updateProfileCommand = "update-profile"
	helpCommand          = "help"
	exitCommand          = "exit"
	workersCommand       = "workers"
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
			name:        helpCommand,
			description: "Получить справку",
		},
		{
			name:        exitCommand,
			description: "Выйти из программы: использование exit или ctrl+c",
		},
		{
			name:        workersCommand,
			description: "Изменить количество горутин",
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
