package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/greenblat17/yet-another-messenger/clients/internal/cli/processor"
	"github.com/greenblat17/yet-another-messenger/clients/pkg/args"
)

// Notification содержит информацию о состоянии выполнения команды
type Notification struct {
	WorkerID int
	Message  string
}

type Job struct {
	CommandName string
	Args        []string
}

type CLI struct {
	processor   *processor.CommandProcessor
	commandList []command
	jobs        chan Job
	errors      chan error
	wg          sync.WaitGroup
	mu          sync.Mutex // Мьютекс для синхронизации доступа к cancelFuncs
	running     bool
}

// New создает новый экземпляр CLI.
func New(handler *processor.CommandProcessor) *CLI {
	return &CLI{
		processor:   handler,
		commandList: initCommandList(handler),
		running:     false,
	}
}

// Run запускает ввод команд от пользователя.
func (c *CLI) Run(ctx context.Context) {
	c.jobs = make(chan Job)
	c.errors = make(chan error)

	go c.handleErrors(ctx)
	go c.inputCommands(ctx)
}

func (c *CLI) inputCommands(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			// Завершение горутины при получении сигнала завершения
			return
		default:
			inputCommand, err := reader.ReadString('\n')
			if err != nil {
				c.recordError(fmt.Errorf("%w: %w", ErrCommandReading, err))
				continue
			}

			inputCommand = strings.TrimSpace(inputCommand)
			if inputCommand == "" {
				c.recordError(ErrCommandEmpty)
				continue
			}

			args := strings.Split(inputCommand, " ")
			c.jobs <- Job{CommandName: args[0], Args: args[1:]}
		}
	}
}

// handleCommands направляет команды от пользователя в нужный метод.
func (c *CLI) handleCommands(ctx context.Context, job Job) (any, error) {
	commandName := job.CommandName
	args := job.Args

	switch commandName {
	case helpCommand:
		c.help()
	case exitCommand:
		c.exit()
	default:
		resp, err := c.executeCommands(ctx, commandName, args)
		if err != nil {
			return "", err
		}
		return resp, nil
	}

	return "", nil
}

func (c *CLI) executeCommands(ctx context.Context, commandName string, arguments []string) (any, error) {
	for _, cmd := range c.commandList {
		if cmd.name == commandName {
			mapArgs := args.ConvertArgsSliceToMap(arguments)
			return cmd.call(ctx, mapArgs)
		}
	}

	return "", ErrCommandNotFound
}

func (c *CLI) handleErrors(ctx context.Context) {
	for {
		select {
		case err := <-c.errors:
			fmt.Println("Error: " + err.Error())
		case <-ctx.Done():
			// Завершение горутины при получении сигнала завершения
			return
		}
	}
}

// IsRunning проверяет, запущено ли приложение.
func (c *CLI) IsRunning() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.running
}

// SetRunning устанавливает состояние приложения.
func (c *CLI) SetRunning(state bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.running = state
}

// recordError записывает ошибку в канал ошибок, если приложение запущено.
func (c *CLI) recordError(err error) {
	if c.IsRunning() {
		c.errors <- err
	}
}
