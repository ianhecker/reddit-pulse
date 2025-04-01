# Reddit Pulse

!!! Warning - This program does not abide by Reddit's guidelines for user
credentials, and is a proof of concept for a coding assignment !!!

> "According to Reddit’s API guidelines, applications should not require users
to directly enter their Reddit username and password. Instead, Reddit strongly
prefers that developers implement the OAuth authorization flow, which allows
users to securely grant access to their account without exposing their
credentials."

## Setup

Create an .env file. You will need:
- Client ID
- Client Secret
- Reddit Subreddit
- Reddit Account Password
- Reddit Account Username

```bash
cp .env.example .env
```

### Client Credentials

You will need to create an App at the following link: https://www.reddit.com/prefs/apps/

1. You will be able to name the app anything
2. Make sure to select the "Script" option
3. Other options do not matter

### User Credentials

Why do i need user credentials?
- As a proof of concept, it was quicker to require user credentials in the .env
file, rather than figure out Reddit's full OAuth process for granting an app
permissions

When you enter in user credentials, you gain increased access to Reddit's API.
You also get rate limit headers that tell you:
- How many requests you have left
- How many requests you have sent
- How many seconds you have remaining in the rate limit cycle

Using a "read only" client does not supply the rate limits via header, and has
slower rate limits than authenticated clients

## Run

You can use the Makefile, and `make run` will compile for your system

```bash
make run
```

Or, use Golang directly

```bash
go run main.go
```

Cancel the program with `Ctrl-C`

## Output

This program finds:
- The top posts for a subreddit
- The top authors for a subreddit

> Note, it is somewhat rare to find authors with mutliple posts, even in the top
100 posts for a subreddit

This progam will log to sdtout, and output prettyified JSON to `output.json`

You can open `output.json` in your file browser, and watch the updates

The JSON file will persist after running
