# WTR<sup>2</sup> - Webex Token Retriever / Refresher

[![Status](https://img.shields.io/badge/status-wip-yellow)](https://github.com/darrenparkinson/wtr)

A small utility CLI for retrieving and subsequently refreshing an integration token from Webex.  

See the Webex documentation for further information on Webex Integrations and Authorization](https://developer.webex.com/docs/integrations)

## Simplified Sequence Diagram

### Retrieve

This command is used to obtain the initial token.

* `wtr retrieve`


```mermaid
sequenceDiagram
    participant User
    participant Webex
    User->>+Webex: webexapis.com/v1/authorize
    Note right of Webex: client_id<br/>response_type=code<br/>redirect_uri<br/>scope<br/>state
    Webex-->>-User: Browser Login Prompt
    User->>+Webex: Username/Password
    Webex-->>-User: Redirect with Code
    User->>+Webex: Exchange Code for token<br/>https://api.ciscospark.com/v1/access_token
    Note right of Webex: grant_type=authorization_code<br/>client_id<br/>client_secret<br/>code<br/>redirect_uri<br/>state
    Webex-->>-User: Access Token Response
```

### Refresh

This command is used to refresh a token before it expires using the details provided by Webex.  The details used to refresh the token are obtained by the initial token retrieval.

* `wtr refresh`

```mermaid
sequenceDiagram
    participant User
    participant Webex
    User->>+Webex: webexapis.com/v1/access_token
    Note right of Webex: grant_type=refresh_token<br/>client_id<br/>client_secret<br/>refresh_token<br/>scope<br/>state
    Webex-->>-User: New Token
```

## Configuration

Information required to begin the process includes:

| Item          | Description                                  | YAML         | Env Var            |
|---------------|----------------------------------------------|--------------|--------------------|
| AppID         | Webex App ID for integration                 | appid        | WEBEX_APPID        |
| Secret        | Secret Associated to the App ID              | secret       | WEBEX_SECRET       |
| Scopes        | Required Scopes as configured for the App ID | scopes       | WEBEX_SCOPES       |
| Redirect Port | Local port to retrieve responses             | redirectPort | WEBEX_REDIRECTPORT |

There are a couple of default values:

| Item          | Default                                                                                                                                                      |
|---------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Scopes        | meeting:recordings_read spark:kms meeting:schedules_read meeting:preferences_write meeting:recordings_write meeting:preferences_read meeting:schedules_write |
| Redirect Port | 6855                                                                                                                                                         |

Note that the scopes would still need to be configured in the webex developer portal for your integration.