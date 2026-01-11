# 🐊 Gator Blog Aggregator
Gator is a CLI blog aggregator written in Golang. Gator uses PostgreSQL.
Gator supports multiple users on a single device via Config file.

## Dependencies
- PostgresSQL
- Go

## Installation
Run the following in your command line:
```bash
go install https://github.com/RussNavas/Gator
```
## Setup
Create a .gatorconfig.json file in the root directory
#### Linux/maxOS
`~/.gatorconfig.json`
#### Windows:
`C:\\Users\<YourName>\.gatorconfig.json`

Inside that .gatorconfig.json we want to map the below struct to a json object:
```go
type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}
```
Where `db_url` is your Postgres connection string and `current_user_name` is your CLI username. So .gatorconfig.json should contain the following:
``` bash
 {"db_url":"YourConnectionString", "current_user_name":"YourUserName"}
```
Save that and you will be read to start using Gator.
## Installed Usage
To register a user:
`gator register <name>`

To login:
`gator login <name>`

To list registered users:
`gator users`

To begin aggregating RSS feeds, in a seprate terminal:
`gator agg`

To add a feed for the currently logged in user to follow:
`gator addfeed <URL>`

To list the feeds the currently logged in user follows:
`gator following`

To unfollow a feed for the currently logged in user:
`gator unfollow <URL>`

To browse a posts, note that limit is optional with a default of 2 posts:
`gator browse <limit>`

---

## Development Usage
To register a user:
`go run . register <name>`

To login:
`go run . login <name>`

To list registered users:
`go run . users`

To begin aggregating RSS feeds, in a seprate terminal:
`go run . agg`

To add a feed for the currently logged in user to follow:
`go run . addfeed <URL>`

To list the feeds the currently logged in user follows:
`go run . following`

To unfollow a feed for the currently logged in user:
`go run . unfollow <URL>`

To browse a posts, note that limit is optional with a default of 2 posts:
`go run . browse <limit>`
