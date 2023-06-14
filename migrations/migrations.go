package migrations

import (
	"database/sql"
	"github.com/GuiaBolso/darwin"
)

var items = []darwin.Migration{
	{
		Version:     1,
		Description: `users_tweeter table`,
		Script: "CREATE TABLE `users_tweeter` (" +
			"  `transaction_id` integer NOT NULL PRIMARY KEY AUTOINCREMENT" +
			",  `name` text NOT NULL" +
			",  `password` text NOT NULL" +
			",  `email` text NOT NULL" +
			",  `emailtoken` text DEFAULT NULL" +
			",  `confirmemailtoken` boolen DEFAULT NULL" +
			",  `resetpasswordtoken` text DEFAULT NULL" +
			",  `birthdate` text NOT NULL" +
			",  `nickname` text NOT NULL" +
			",  `bio` text DEFAULT NULL" +
			",  `location` text DEFAULT NULL" +
			",  `logintoken` text DEFAULT NULL" +
			");"},
	{
		Version:     2,
		Description: `tweets`,
		Script: "CREATE TABLE `tweets` (" +
			"  `tweet_id` integer NOT NULL PRIMARY KEY AUTOINCREMENT" +
			",  `user_id` integer NOT NULL" +
			",  `text` text NOT NULL" +
			",  `created_at` timestamp with time zone NOT NULL" +
			",  `public` boolean NOT NULL" +
			",  `only_followers` boolean NOT NULL" +
			",  `only_mutual_followers` boolean NOT NULL" +
			",  `only_me` boolean NOT NULL" +
			")",
	},
	{
		Version:     3,
		Description: `followers_subscriptions`,
		Script: "CREATE TABLE `followers_subscriptions` (" +
			"  `follower_id` integer NOT NULL " +
			",  `subscription_id` integer NOT NULL" +
			")",
	},
}

func Run(db *sql.DB) error {
	return darwin.New(darwin.NewGenericDriver(db, darwin.SqliteDialect{}), items, nil).Migrate()
}
