# Figure 1 slackbot oembed
Add Figure 1 content functionality to slack.

## Configuration
Create `conf.json` file with structure:

```
{
	"email": "EMAIL",
	"password": "PASSWORD",
	"oauth_access_token": "OAUTH_ACCESS_TOKEN",
	"verification_token": "VERIFICATION_TOKEN"
}
```

## Adding app to slack
Name: **Figure 1 Slack App**


### Required permissions
In the **Permissions** section, add the following permissions and re-authorize the app:
- `commands`
- `chat:write:bot`
- `chat:write:user`

### Required tokens/secrets
Add the following tokens to `conf.json`
- in **Permissions** -> **OAuth & Permissions**, grab the `OAuth Access Token`
- in **App Credentials**, grab `Verification Token`


### Supported commands
In the slash command sections, add the following commands:

#### case
**command:** `/case`

**url:** `https://catc-services.com/fig1-slack/case`

**description:** Display Figure 1 case info

**usage hint:** [case url or case id]

#### user
**command:** `/user`

**url:** `https://catc-services.com/fig1-slack/user`

**description:** Get Figure 1 user

**usage hint:** [username or profile url]

#### collection
**command:** `/collection`

**url:** `https://catc-services.com/fig1-slack/collection`

**description:** Display a collection preview

**usage hint:** [collection id or collection url]

### nginx config
All requests to `https://catc-services.com/fig1-slack/*` redirects to `localhost:3400/*`

```nginx
upstream slack_app {
        server localhost:3400;
}

server {
	# ...
	location ~ ^/fig1-slack/(.*)$ {
		proxy_pass http://slack_app/$1;
	}
}
```
