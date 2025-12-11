# Gator - RSS Blog Aggregator

Gator is a high-performance, CLI-based RSS feed aggregator built in Go. It enables users to subscribe to feeds, aggregate posts from various sources into a PostgreSQL database, and browse content directly from the command line.

## Prerequisites
Before running Gator, ensure you have the following installed:
Go (v1.23+): Required if building from source.
PostgreSQL: Gator requires a running Postgres database to store user and feed data.
Docker (Optional): For running Gator in a containerized environment.

## Installation
Method 1: Go Install (Recommended)
To install the gator binary directly to your system's Go bin directory:
go install [github.com/Driadix/GoBlogAggregattor@latest](https://github.com/Driadix/GoBlogAggregattor@latest)

Note: Ensure your $GOPATH/bin is added to your system $PATH so you can run gator from anywhere.

Method 2: Build from Source
git clone [https://github.com/username/repo-name.git](https://github.com/Driadix/GoBlogAggregattor.git)
cd repo-name
go build -o gator .
./gator <command>

Method 3: Docker
You can build a lightweight, statically linked container image:
docker build -t gator .

## Configuration
Gator uses a JSON configuration file to manage database connections and the current user session.
Create a file named .gatorconfig.json in your home directory (~):
touch ~/.gatorconfig.json

Add the following content, replacing the db_url with your Postgres connection string:

{
  "db_url": "postgres://your_user:your_password@localhost:5432/gator_db?sslmode=disable",
  "current_user_name": ""
}

Tip: current_user_name will be populated automatically when you run gator login.

## Usage
### User Management
Register a new user:
gator register <username>

Login as an existing user:
gator login <username>

List all users:
gator users

Reset database (WARNING: Deletes all users/feeds):
gator reset

### Feed Management
Add a new feed:
gator addfeed <feed_name> <url>
Example: gator addfeed "Lane's Blog" "https://wagslane.dev/index.xml"

List all feeds:
gator feeds

Follow a feed:
gator follow <url>


Unfollow a feed:
gator unfollow <url>


List followed feeds:
gator following

### Aggregation & Browsing

Start Aggregator:
This command runs indefinitely, fetching new posts from followed feeds every time interval.
gator agg 1m
(Fetches every 1 minute)

Browse Posts:
View the most recent posts from feeds you follow.
gator browse 5
(Shows the latest 5 posts)

## Running with Docker

To run gator via Docker, you must handle two things: Network Access (to reach your local DB) and State Persistence (to read your config file).

docker run -it \
  -v ~/.gatorconfig.json:/root/.gatorconfig.json \
  --network="host" \
  gator <command>


-v ~/.gatorconfig.json:/root/.gatorconfig.json: Mounts your local config file into the container so Gator knows your DB URL and current user.

--network="host": Allows the container to access your host machine's localhost (Postgres).

## Development

During development, you can run the code directly without compiling a binary:
go run . <command>

However, for production usage, the gator binary or Docker image is recommended.