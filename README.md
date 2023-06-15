<p align="center"><img src = "https://github.com/nicewook/gopt/assets/6977358/6da29951-b25d-492f-967d-4f7be8bddbe9" width="50%" height="50%"></p>

# GOPT(GPT CLI made of Golang)
GOPT is a command line interface that uses OpenAI's `gpt-3.5-turbo` API to generate responses to user input.

## Configuration
The app uses the following environment variables:

- `OPENAI_API_KEY` - Your OpenAI API key. Required to make requests to the OpenAI API.
  - You can set it as a environment variable. or GOPT will ask you to enter API key. the key will be saved as a JSON file.
- `RUN_MODE` - Set to `dev` to enable debug logging, any other value to disable. Optional, defaults to disabled.

## Usage
To use GOPT:

1. Install Go.
2. Clone the repository:

```bash
$ git clone https://github.com/nicewook/gopt.git
```
3. Obtain an OpenAI API key and set the OPENAI_API_KEY environment variable:
```bash
export OPENAI_API_KEY=YOUR_KEY_HERE
```
4. Build and run the app:
```bash
$ go build 
$ ./gopt
```
- Write a message and press enter to send it. 
- The app will send the message to OpenAI's API and return a response. 
- The response will be displayed below your message.
- Keep entering messages to continue the conversation.
  - The conversation will be saved and sent along with the next message.
  - And conversation will be deleted from the oldest one due to token limitation.
- Enter `exit` at any time to exit the app.

## Commands on gpt shell mode
- `help` - Displays command usage.
- `config` - Displays configuration information. 
- `context` - Displays the conversation context which reserved at the moment.
- `reset` - Reset all the conversation context.
- `exit` or `q` - Exits the app.

## Issues and Contributing
Please open an issue on GitHub to report any problems or make feature requests. Pull requests are welcome!

