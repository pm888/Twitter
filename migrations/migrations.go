package migrations

import (
	"database/sql"

	"github.com/GuiaBolso/darwin"
)

var items = []darwin.Migration{
	{
		Version:     1,
		Description: `users_tweeter table`,
		Script: `CREATE TABLE users_tweeter (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		email TEXT NOT NULL,
		emailtoken TEXT DEFAULT NULL,
		confirmemailtoken BOOLEAN DEFAULT NULL,
		resetpasswordtoken TEXT DEFAULT NULL,
		birthdate DATE NOT NULL,
		nickname TEXT NOT NULL,
		bio TEXT DEFAULT NULL,
		location TEXT DEFAULT NULL
	)`,
	},
	{
		Version:     2,
		Description: `tweets table`,
		Script: `CREATE TABLE tweets (
			tweet_id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			text TEXT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			parent_tweet_id INTEGER,
			public BOOLEAN NOT NULL,
			only_followers BOOLEAN NOT NULL,
			only_mutual_followers BOOLEAN NOT NULL,
			only_me BOOLEAN NOT NULL,
			retweet INTEGER NOT NULL
		)`,
	},
	{
		Version:     3,
		Description: `followers_subscriptions table`,
		Script: `CREATE TABLE followers_subscriptions (
			id SERIAL PRIMARY KEY,
			follower_id INTEGER NOT NULL,
			subscription_id INTEGER NOT NULL
		)`,
	},
	{
		Version:     4,
		Description: `likes table`,
		Script: `CREATE TABLE likes (
			id SERIAL PRIMARY KEY,
			tweet_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL
		)`,
	},
	{
		Version:     5,
		Description: `retweets table`,
		Script: `CREATE TABLE retweets (
			id SERIAL PRIMARY KEY,
			tweet_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL
		)`,
	},
	{
		Version:     6,
		Description: `users_session table`,
		Script: `CREATE TABLE user_session (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL,
			login_token TEXT DEFAULT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL
		)`,
	},
	{
		Version:     7,
		Description: `reply`,
		Script: `CREATE TABLE reply (
			reply_id SERIAL PRIMARY KEY,
			text_reply TEXT NOT NULL,
			parent_tweet_id INTEGER DEFAULT NULL,
			tweet_id INTEGER NOT NULL,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL
		)`,
	},
}

func Run(db *sql.DB) error {
	return darwin.New(darwin.NewGenericDriver(db, darwin.PostgresDialect{}), items, nil).Migrate()
}
