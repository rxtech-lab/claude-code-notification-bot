base on this docs: https://docs.anthropic.com/en/docs/claude-code/hooks and python example, create a go program to handle the list of commands sent from claude code and then send notification base on the hook type
(user will config hook type through cli args).

Your job is

1. create list of hook handlers in `internal/handler` function that implements `Handler` interface and each handler should also takes a telegram client in their constructor method

2. Send different message in different hook handler base on the event. Message should be built using go template.

3. Write unit tests for each go file

4. Create a makefile for handling different commands

5. A readme about the project

6. Read telegram API Key and target chat room id through `env` and print error if not se
