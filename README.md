# WTR<sup>2</sup> - Webex Token Retriever / Refresher

<!-- [![Status](https://img.shields.io/badge/status-wip-yellow)](https://github.com/darrenparkinson/wtr) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/darrenparkinson/wtr)](https://goreportcard.com/report/github.com/darrenparkinson/wtr)  [![License: MIT](https://img.shields.io/badge/License-MIT-lightgreen.svg)](https://opensource.org/licenses/MIT) [![published](https://static.production.devnetcloud.com/codeexchange/assets/images/devnet-published.svg)](https://developer.cisco.com/codeexchange/github/repo/darrenparkinson/wtr)
<!-- [![GoDoc](https://godoc.org/github.com/darrenparkinson/wtr?status.svg)](https://godoc.org/github.com/darrenparkinson/wtr)  -->
<!-- ![GitHub All Releases](https://img.shields.io/github/downloads/darrenparkinson/wtr/total) -->


A small utility CLI for retrieving and subsequently refreshing an integration token from Webex.  Primarily useful for when you just want to get a token to test with.

See the Webex documentation for further information on [Webex Integrations and Authorization](https://developer.webex.com/docs/integrations)

## Simplified Sequence Diagram

### Retrieve

This command is used to obtain the initial token.

* `wtr retrieve`

![Retrieve Diagram](images/retrieve-sequence-diagram.png "Retrieve Sequence Diagram")

### Refresh

This command is used to refresh a token before it expires using the details provided by Webex.  The details used to refresh the token are obtained by the initial token retrieval.

* `wtr refresh`

![Refresh Diagram](images/refresh-sequence-diagram.png "Refresh Sequence Diagram")

## Installation

You can download an executable for your platform from the [releases](https://github.com/darrenparkinson/wtr/releases) page.

On MacOS, you will need to jump though a few hoops to get it to execute.  You can do the following:

* Run `chmod +x wtr_0.0.0_darwin_amd64` to make it executable (changing the filename as appropriate)
* In Finder, right-click the executable and use "Open With" to open it in terminal and accept the security warning.

You should then be able to use the wtr utility.

As a final note, you might also like to rename it directly to `wtr`.

## Configuration

Information required to begin the process includes:

| Item          | Description                                                                                                  | json         | Env Var            |
|---------------|--------------------------------------------------------------------------------------------------------------|--------------|--------------------|
| Client ID     | Webex Client ID for integration                                                                              | clientid        | WEBEX_CLIENTID        |
| Secret        | Secret Associated to the Client ID                                                                           | secret       | WEBEX_SECRET       |
| Scopes        | Required Scopes as configured for the Client ID.  Space separated list of scopes.  Ensure you add `spark:kms` | scopes       | WEBEX_SCOPES       |
| Redirect Port | Local port to retrieve responses                                                                             | redirectPort | WEBEX_REDIRECTPORT |

There are a couple of default values so these items are optional:

| Item          | Default             |
|---------------|---------------------|
| Scopes        | spark:kms spark:all |
| Redirect Port | 6855                |

Note that the scopes would still need to be configured in the webex developer portal for your integration, along with a Redirect URI.

If you've got this far, you probably already know what you need, but as a recap, to configure your integration, 
* Go to [https://developer.webex.com/my-apps](https://developer.webex.com/my-apps);
* Select "Create a New App";
* Select "Integration" and fill in the required details, ensuring that you:
  * select the required **Scopes** you need 
  * and add a **Redirect URI** of `http://localhost:6855` (*replace the port number with the one you want to use if you're not using the default one*)

At this point, you will want to copy the Client ID and Client Secret and add them to your configuration file:

`.wtr-cli.json`:
```json
{
    "clientid": "<YOUR CLIENT ID HERE>",
    "secret": "<YOUR CLIENT SECRET HERE>",
    "scopes": "spark:kms spark:all"
}
```

*Ensure you put in the scopes you specified, remembering to add `spark:kms`.*

Alternatively you can use environment variables for `clientid`, `secret`, `scopes` and `redirectPort` too.  The token details will still be written to this configuration file though.

## Usage

There are essentially two commands:

* `wtr retrieve` - Retrieve an initial token
* `wtr refresh` - Refresh an existing token that you already retrieved

Each has some options:

* `-d` or `--debug` - see verbose log output as to what the command is doing
* `-o` or `--output` - output the token to the console as text
* `-j` or `--json` - output the token to the console as json 
* `-t` or `--timeout` - how many seconds to wait for a response during the initial retrieval process

Using the `--help` option on either command will show all the options available.

