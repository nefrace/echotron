| <img src="assets/readme_banner.png" alt="logo" width="800"><br/><br/> [![Language](https://img.shields.io/badge/Language-Go-blue.svg)](https://golang.org/) [![PkgGoDev](https://pkg.go.dev/badge/github.com/NicoNex/echotron/v3)](https://pkg.go.dev/github.com/NicoNex/echotron/v3) [![Go Report Card](https://goreportcard.com/badge/github.com/NicoNex/echotron/v3)](https://goreportcard.com/report/github.com/NicoNex/echotron/v3) [![License](http://img.shields.io/badge/license-LGPL3.0-orange.svg?style=flat)](https://github.com/NicoNex/echotron/blob/master/LICENSE) [![Coverage Status](https://coveralls.io/repos/github/NicoNex/echotron/badge.svg?branch=master)](https://coveralls.io/github/NicoNex/echotron?branch=master) [![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go) [![Telegram](https://img.shields.io/badge/Echotron%20News-blue?logo=telegram&style=flat)](https://t.me/echotronnews) |
| :------: |

**Echotron** is an elegant and concurrent library for the Telegram bot API in Go.

Fetch with

```bash
go get github.com/NicoNex/echotron/v3
```

## Design

**Echotron** is heavily based on concurrency: for example, every call to the `Update` method of each bot is executed on a different goroutine. This makes sure that, even if one instance of the bot is deadlocked, the other ones keep running just fine, making the bot work for other users without any issues and/or slowdowns.

**Echotron** is designed to be as similar to the official [Telegram API](https://core.telegram.org/bots/api) as possible, but there are some things to take into account before starting to work with this library.

- The methods have the exact same name, but with a capital first letter, since in Go methods have to start with a capital letter to be exported.
_Example: `sendMessage` becomes `SendMessage`_
- The order of the parameters in some methods is different than in the official Telegram API, so refer to the [docs](https://pkg.go.dev/github.com/NicoNex/echotron/v3) for the correct one.
- The only `chat_id` (or, in this case, `chatID`) type supported is `int64`, instead of the "Integer or String" requirement of the official API. That's because numeric IDs can't change in any way, which isn't the case with text-based usernames.
- In some methods, you might find a `InputFile` type parameter. [`InputFile`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#InputFile) is a struct with unexported fields, since only three combination of fields are valid, which can be obtained through the methods [`NewInputFileID`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#NewInputFileID), [`NewInputFilePath`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#NewInputFilePath) and [`NewInputFileBytes`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#NewInputFileBytes).
- In some methods, you might find a `MessageIDOptions` type parameter. [`MessageIDOptions`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#MessageIDOptions) is another struct with unexported fields, since only two combination of field are valid, which can be obtained through the methods [`NewMessageID`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#NewMessageID) and [`NewInlineMessageID`](https://pkg.go.dev/github.com/NicoNex/echotron/v3#NewInlineMessageID).
- Optional parameters can be added by passing the correct struct to each method that might request optional parameters. If you don't want to pass any optional parameter, `nil` is more than enough. Refer to the [docs](https://pkg.go.dev/github.com/NicoNex/echotron/v3) to check for each method's optional parameters struct: it's the type of the `opts` parameter.
- Some parameters are hardcoded to avoid putting random stuff which isn't recognized by the Telegram API. Some notable examples are [`ParseMode`](https://github.com/NicoNex/echotron/blob/master/options.go#L21), [`ChatAction`](https://github.com/NicoNex/echotron/blob/master/options.go#L54) and [`InlineQueryType`](https://github.com/NicoNex/echotron/blob/master/inline.go#L27). For a full list of custom hardcoded parameters, refer to the [docs](https://pkg.go.dev/github.com/NicoNex/echotron/v3) for each custom type: by clicking on the type's name, you'll get the source which contains the possible values for that type.

## Usage

### Long Polling

A very simple implementation:

```golang
package main

import (
	"log"

	"github.com/NicoNex/echotron/v3"
)

// Struct useful for managing internal states in your bot, but it could be of
// any type such as `type bot int64` if you only need to store the chatID.
type bot struct {
	chatID int64
	echotron.API
}

const token = "YOUR TELEGRAM TOKEN"

// This function needs to be of type 'echotron.NewBotFn' and is called by
// the echotron dispatcher upon any new message from a chatID that has never
// interacted with the bot before.
// This means that echotron keeps one instance of the echotron.Bot implementation
// for each chat where the bot is used.
func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

// This method is needed to implement the echotron.Bot interface.
func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage("Hello world", b.chatID, nil)
	}
}

func main() {
	// This is the entry point of echotron library.
	dsp := echotron.NewDispatcher(token, newBot)
	log.Println(dsp.Poll())
}
```

Functional example with bot internal states:

```golang
package main

import (
	"log"
	"strings"

	"github.com/NicoNex/echotron/v3"
)

// Recursive type definition of the bot state function.
type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatID int64
	state  stateFn
	name   string
	echotron.API
}

const token = "YOUR TELEGRAM TOKEN"

func newBot(chatID int64) echotron.Bot {
	bot := &bot{
		chatID: chatID,
		API:	echotron.NewAPI(token),
	}
	// We set the default state to the bot.handleMessage method.
	bot.state = bot.handleMessage
	return bot
}

func (b *bot) Update(update *echotron.Update) {
	// Here we execute the current state and set the next one.
	b.state = b.state(update)
}

func (b *bot) handleMessage(update *echotron.Update) stateFn {
	if strings.HasPrefix(update.Message.Text, "/set_name") {
		b.SendMessage("Send me my new name!", b.chatID, nil)
		// Here we return b.handleName since next time we receive a message it
		// will be the new name.
		return b.handleName
	}
	return b.handleMessage
}

func (b *bot) handleName(update *echotron.Update) stateFn {
	b.name = update.Message.Text
	b.SendMessage(fmt.Sprintf("My new name is %q", b.name), b.chatID, nil)
	// Here we return b.handleMessage since the next time we receive a message
	// it will be handled in the default way.
	return b.handleMessage
}

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	log.Println(dsp.Poll())
}
```

Example with self destruction for lower RAM usage:

```golang
package main

import (
	"log"
	"time"

	"github.com/NicoNex/echotron/v3"
)

type bot struct {
	chatID int64
	echotron.API
}

const token = "YOUR TELEGRAM TOKEN"

var dsp *echotron.Dispatcher

func newBot(chatID int64) echotron.Bot {
	bot := &bot{
		chatID,
		echotron.NewAPI(token),
	}
	go bot.selfDestruct(time.After(time.Hour))
	return bot
}

func (b *bot) selfDestruct(timech <-chan time.Time) {
	<-timech
	b.SendMessage("goodbye", b.chatID, nil)
	dsp.DelSession(b.chatID)
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage("Hello world", b.chatID, nil)
	}
}

func main() {
	dsp = echotron.NewDispatcher(token, newBot)
	log.Println(dsp.Poll())
}
```

### Webhook

```golang
package main

import "github.com/NicoNex/echotron/v3"

type bot struct {
	chatID int64
	echotron.API
}

const token = "YOUR TELEGRAM TOKEN"

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage("Hello world", b.chatID, nil)
	}
}

func main() {
	dsp := echotron.NewDispatcher(token, newBot)
	dsp.ListenWebhook("https://example.com:443/my_bot_token")
}
```


### Webhook with a custom http.Server

This is an example for a custom http.Server which handles your own specified routes
and also the webhook route which is specified by ListenWebhook.

```golang
package main

import (
	"github.com/NicoNex/echotron/v3"

	"context"
	"log"
	"net/http"
)

type bot struct {
	chatID int64
	echotron.API
}

const token = "YOUR TELEGRAM TOKEN"

func newBot(chatID int64) echotron.Bot {
	return &bot{
		chatID,
		echotron.NewAPI(token),
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage("Hello world", b.chatID, nil)
	}
}

func main() {
	termChan := make(chan os.Signal, 1) // Channel for terminating the app via os.Interrupt signal
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// Handle user login
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Handle user logout
	})
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		// Tell something about your awesome telegram bot
	})

	// Set custom http.Server
	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		<-termChan
		// Perform some cleanup..
		if err := server.Shutdown(context.Background()); err != nil {
			log.Print(err)
		}
	}()

	// Capture the interrupt signal for app termination handling
	dsp := echotron.NewDispatcher(token, newBot)
	dsp.SetHTTPServer(server)
	// Start your custom http.Server with a registered /my_bot_token handler.
	log.Println(dsp.ListenWebhook("https://example.com/my_bot_token"))
}
```
