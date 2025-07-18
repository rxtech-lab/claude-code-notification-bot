# Claude Code Telegram Notification

A Go application that handles webhook events from Claude Code and sends notifications to Telegram. This tool allows you to receive real-time notifications about your Claude Code interactions directly in your Telegram chat.

## Features

- **Multiple Hook Types**: Supports various Claude Code hook events:
  - `user-prompt-submit`: User submits a prompt
  - `tool-call`: Tool execution events
  - `file-write`: File modification events
  - `session-start`: Session initialization
  - `session-end`: Session termination
  - `generic`: Fallback for unknown event types

- **Rich Telegram Messages**: Formatted messages with emojis and structured information
- **Template-Based**: Uses Go templates for customizable message formatting
- **Environment Configuration**: Secure configuration via environment variables
- **Comprehensive Testing**: Unit tests for all components

## Installation

### Prerequisites

- Go 1.19 or higher
- Telegram Bot Token
- Telegram Chat ID

### Building

```bash
# Clone the repository
git clone <repository-url>
cd claude-code-telegram-notification

# Install dependencies
make install

# Build the application
make build
```

### Installing CLI Globally

To install the CLI binary globally so you can use it from anywhere:

```bash
# Install the CLI to /usr/local/bin
make install-cli

# Now you can use it from anywhere
claude-telegram-notifier user-prompt-submit

# To uninstall later
make uninstall-cli
```

## Configuration

### Environment Variables

Set the following environment variables:

```bash
export CLAUDE_CODE_TELEGRAM_BOT_TOKEN="your-bot-token-here"
export CLAUDE_CODE_TELEGRAM_CHAT_ID="your-chat-id-here"
```

### Creating a Telegram Bot

1. Message [@BotFather](https://t.me/BotFather) on Telegram
2. Send `/newbot` and follow the instructions
3. Copy the bot token provided
4. Set your bot token: `export CLAUDE_CODE_TELEGRAM_BOT_TOKEN="your-bot-token"`

### Discovering Your Chat ID

Use the built-in chat client to easily discover your chat ID:

```bash
# Set your bot token first
export CLAUDE_CODE_TELEGRAM_BOT_TOKEN="your-bot-token"

# Run the chat client
make chat-client
# or
make discover-chat-id
# or if installed globally
claude-telegram-notifier chat-client
```

The chat client will:
1. Start listening for messages to your bot
2. Display instructions to send a message
3. Show your chat ID when you send any message to the bot
4. Provide the exact export command to set your chat ID

Example output:
```
🤖 Telegram Chat ID Discovery
===============================
Send any message to your bot to discover your chat ID...
Press Ctrl+C to exit

📨 New message received!
   Chat ID: 123456789
   Chat Type: private
   From: John Doe (@johndoe)
   Message: hello

💡 To use this chat, set your environment variable:
   export CLAUDE_CODE_TELEGRAM_CHAT_ID="123456789"
```

## Usage

### Basic Usage

The application reads JSON hook data from stdin and sends formatted notifications to Telegram.

```bash
# Run with specific hook type
echo '{"event_type":"user-prompt-submit","timestamp":"2024-01-01T12:00:00Z","message":"Hello Claude"}' | ./bin/claude-telegram-notifier user-prompt-submit

# Run with default (generic) handler
echo '{"event_type":"custom","timestamp":"2024-01-01T12:00:00Z"}' | ./bin/claude-telegram-notifier
```

### Hook Types

Specify the hook type as the first command-line argument:

- `user-prompt-submit`: For user prompt submission events
- `tool-call`: For tool execution events
- `file-write`: For file modification events
- `session-start`: For session start events
- `session-end`: For session end events
- `chat-client`: Run chat ID discovery mode
- `generic` (default): For any other event types

### Integration with Claude Code

Configure Claude Code hooks to pipe data to this application:

```bash
# If installed globally
claude-code --hook="user-prompt-submit:echo '%s' | claude-telegram-notifier user-prompt-submit"

# If using local binary
claude-code --hook="user-prompt-submit:echo '%s' | /path/to/claude-telegram-notifier user-prompt-submit"
```

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage
make test-coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run go vet
make vet

# Run all checks
make check
```

### Development Setup

```bash
# Set up development environment
make dev-setup
```

## Examples

### Example Messages

The application generates rich, formatted Telegram messages:

**User Prompt Submit:**
```
🤖 Claude Code Hook: User Prompt Submit

📝 Event: user-prompt-submit
⏰ Timestamp: 2024-01-01 12:00:00
🔧 Tool Used: Read
💬 Message: Please analyze this file
📁 File: /path/to/file.go
```

**Tool Call:**
```
🔨 Claude Code Hook: Tool Call

📝 Event: tool-call
🛠️ Tool: Bash
⏰ Timestamp: 2024-01-01 12:00:00
📁 File: /path/to/script.sh
✅ Status: Success
```

### Testing Examples

Use the Makefile to test different hook types:

```bash
# Test user prompt submit
make example-user-prompt

# Test tool call
make example-tool-call

# Test file write
make example-file-write

# Test session start
make example-session-start

# Test session end
make example-session-end
```

## Project Structure

```
├── cmd/
│   └── main.go                 # Main application entry point
├── internal/
│   ├── handler/               # Hook event handlers
│   │   ├── interface.go       # Handler interface
│   │   ├── user_prompt_submit.go
│   │   ├── tool_call.go
│   │   ├── file_write.go
│   │   ├── session_start.go
│   │   ├── session_end.go
│   │   ├── generic.go
│   │   └── *_test.go         # Unit tests
│   └── telegram/
│       ├── client.go          # Telegram API client
│       └── client_test.go     # Client tests
├── pkg/
│   └── templates/
│       ├── templates.go       # Message templates
│       └── templates_test.go  # Template tests
├── docs/
│   └── DesignDocs.md         # Design documentation
├── Makefile                   # Build and development commands
├── go.mod                     # Go module definition
└── README.md                  # This file
```

## Error Handling

The application handles various error conditions:

- Invalid JSON in hook data
- Missing environment variables
- Telegram API errors
- Template rendering errors

Errors are logged to stderr, and the application exits with a non-zero code on failure.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make check` to ensure code quality
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Troubleshooting

### Common Issues

1. **Missing Environment Variables**
   ```
   Error: TELEGRAM_BOT_TOKEN environment variable is required
   ```
   Solution: Set the required environment variables

2. **Invalid Chat ID**
   ```
   Error: telegram API error: Bad Request: chat not found
   ```
   Solution: Verify your chat ID is correct

3. **Network Issues**
   ```
   Error: failed to send request: dial tcp: lookup api.telegram.org
   ```
   Solution: Check your internet connection and firewall settings

### Debug Mode

For debugging, you can test the application manually:

```bash
# Set environment variables
export TELEGRAM_BOT_TOKEN="your-token"
export TELEGRAM_CHAT_ID="your-chat-id"

# Test with sample data
echo '{"event_type":"test","timestamp":"2024-01-01T12:00:00Z"}' | ./bin/claude-telegram-notifier
```