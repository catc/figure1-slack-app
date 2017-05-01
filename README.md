# Figure 1 slackbot oembed

Provides Figure 1 case oembed on slack.

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

## TODO
- add tests
- change name of app + repo
- update readme.md to contain info how to set up correct slack permissions and steps
- todo
	- handle errors better, http.Error(...) doesn't provide revelant data to slack
	- make a wrapper interface? takes:
		- res
		- public message
		- custom msg
		- err
		- other content to add
- todo
	- finish remaining handlers	
		- collection
	- change colors to constants
	- change urls to constants?


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

**description:** Display collection preview

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
