<p align="center"><img src = "https://github.com/nicewook/gopt/assets/6977358/6da29951-b25d-492f-967d-4f7be8bddbe9" width="50%" height="50%"></p>

# GOPT - GPT CLI made of Golang

GOPT is a command line interface that uses OpenAI's `gpt-3.5-turbo` API to generate responses to user input.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Issues](#issues)
4. [Contact](#contact)
5. [License](#license)

## Installation

### Prerequisites
Before you begin, ensure you have met the following requirements:

- You have installed the latest version of Go (1.20 or higher as of this writing).

### Go install

You can install the `gopt` binary directly, by using `go install`
``` bash
go install github.com/nicewook/gopt@latest
```
After running this command, Go downloads the source code, compiles the package and installs 
the binary to `$GOPATH/bin` (or `$GOBIN` if set).

Or, you can `git clone` the whole repository and build like this
``` bash
git clone https://github.com/nicewook/gopt.git
cd gopt
go build
  ```

### Configuration
The app need `OPENAI_API_KEY` to use OpenAI API.
- Get your API Key here: https://platform.openai.com/account/api-keys

Your OpenAI API key. Required to make requests to the OpenAI API.
- You can set it as a environment variable. 
  ```
  export OPENAI_API_KEY=YOUR_KEY_HERE
  ```
- OR, GOPT will ask you to enter API key. the key will be saved in `$HOME/.gopt/config.json` as a plaintext.

## Usage
- Write a message and press enter to send it. 
- The app will send the message to OpenAI's API and return a response. 
- The response will be displayed below your message.
- Keep entering messages to continue the conversation.
  - The conversation will be saved and sent along with the next message.
  - And conversation will be deleted from the oldest one due to token limitation.
- Enter `exit` at any time to exit the app.

### Commands on gpt shell mode
- `help` - Displays command usage.
- `config` - Displays configuration information. 
- `context` - Displays the conversation context which reserved at the moment.
- `reset` - Reset all the conversation context.
- `clear` - Clear terminal.
- `exit` or `q` - Exits the app.

## Issues 
Please open an issue on GitHub to report any problems or make feature requests. 

## Contact
If you want to contact me you can reach me at <nicewook@hotmail.com>.

## License
This project is licensed under the terms of the MIT license. See the [LICENSE](LICENSE) file for details.
