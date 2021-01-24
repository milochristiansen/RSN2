/*
Copyright 2020-2021 by Milo Christiansen

This software is provided 'as-is', without any express or implied warranty. In
no event will the authors be held liable for any damages arising from the use of
this software.

Permission is granted to anyone to use this software for any purpose, including
commercial applications, and to alter it and redistribute it freely, subject to
the following restrictions:

1. The origin of this software must not be misrepresented; you must not claim
that you wrote the original software. If you use this software in a product, an
acknowledgment in the product documentation would be appreciated but is not
required.

2. Altered source versions must be plainly marked as such, and must not be
misrepresented as being the original software.

3. This notice may not be removed or altered from any source distribution.
*/

package main

import _ "github.com/mattn/go-sqlite3"
import "database/sql"

var DB *sql.DB

var InitCode = `
pragma foreign_keys = on;

create table if not exists Users (
	ID text primary key,
	Email text unique not null,
	Password text not null,
	CanLogin integer
);
create unique index if not exists Emails on Users(Email);

create table if not exists Feeds (
	ID text primary key,

	URL text unique not null
);
create unique index if not exists FeedURLs on Feeds(URL);

create table if not exists Articles (
	ID text primary key,
	Feed text not null,

	Title text collate nocase,
	URL text unique not null,

	Published integer,

	foreign key (Feed) references Feeds(ID) on delete cascade
);
create unique index if not exists ArticleURLs on Articles(URL);

create table if not exists ReadFlags (
	User text not null,
	Article text not null,

	foreign key (User) references Users(ID) on delete cascade,
	foreign key (Article) references Articles(ID) on delete cascade
);

create table if not exists PausedFlags (
	User text not null,
	Feed text not null,

	foreign key (User) references Users(ID) on delete cascade,
	foreign key (Feed) references Feeds(ID) on delete cascade
);

create table if not exists Subscribed (
	User text not null,
	Feed text not null,
	Name text,

	foreign key (User) references Users(ID) on delete cascade,
	foreign key (Feed) references Feeds(ID) on delete cascade
);
`

var Queries = map[string]*queryHolder{
	// Background updater
	"GetAllFeeds": &queryHolder{`
		select URL, ID from Feeds;
	`, nil},
	"ArticleExistsByURL": &queryHolder{`
		select ID from Articles where URL = ?1 union select "" order by 1 desc limit 1;
	`, nil},
	"ArticleAdd": &queryHolder{`
		insert into Articles (ID, Feed, Title, URL, Published) values (?1, ?2, ?3, ?4, ?5);
	`, nil},
	"FeedListSubs": &queryHolder{`
		select User from Subscribed where Feed = ?1;
	`, nil},

	// /api/user/confirm-email
	"ConfirmEmail": &queryHolder{`
		update Users set CanLogin = 1 where ID = ?1;
	`, nil},
	"GetEmail": &queryHolder{`
		select Email from Users where ID = ?1;
	`, nil},
	// /api/user/delete-email
	"DeleteEmail": &queryHolder{`
		delete from Users where ID = ?1 and CanLogin = 0;
	`, nil},
	// /api/user/login (one row)
	"UserLogin": &queryHolder{`
		select ID, Password, CanLogin from Users where Email = ?1;
	`, nil},
	// /api/user/new
	"UserNew": &queryHolder{`
		insert into Users (ID, Email, Password, CanLogin) values (?1, ?2, ?3, 0);
	`, nil},
	"UserEmailExists": &queryHolder{`
		select exists(select 1 from Users where Email = ?1);
	`, nil},
	// /api/user/new-pass
	"UserNewPass": &queryHolder{`
		update Users set Password = ?2 where ID = ?1;
	`, nil},
	"UserGetPass": &queryHolder{`
		select Password from Users where ID = ?1;
	`, nil},
	// /api/user/new-name
	"UserNewName": &queryHolder{`
		update Users set Email = ?2, CanLogin = 0 where ID = ?1;
	`, nil},

	// /api/feed/list
	"FeedList": &queryHolder{`
		select ID, (
				select Name from Subscribed where Feed = Feeds.ID and User = ?1
			), URL, (
			ID in (select Feed from PausedFlags where User = ?1)
		) from Feeds where (
			ID in (select Feed from Subscribed where User = ?1)
		);
	`, nil},
	// /api/feed/details (one row)
	"FeedDetails": &queryHolder{`
		select ID, (
			select Name from Subscribed where Feed = ?2 and User = ?1
		), URL, (
			ID in (select Feed from PausedFlags where User = ?1)
		) from Feeds where (
			ID = ?2 and
			ID in (select Feed from Subscribed where User = ?1)
		);
	`, nil},
	// /api/feed/articles
	"FeedArticles": &queryHolder{`
		select ID, Title, URL, Published, (
			ID in (select Article from ReadFlags where User = ?1)
		) from Articles where (
			Feed = ?2 and
			Feed in (select Feed from Subscribed where User = ?1)
		) order by Published;
	`, nil},
	// /api/feed/subscribe
	"FeedExistsByURL": &queryHolder{`
		select ID from Feeds where URL = ?1 union select "" order by 1 desc limit 1;
	`, nil},
	"FeedAdd": &queryHolder{`
		insert into Feeds (ID, URL) values (?1, ?2);
	`, nil},
	"FeedAlreadySubscibed": &queryHolder{`
		select exists(select 1 from Subscribed where User = ?1 and Feed = ?2);
	`, nil},
	"FeedSubscibe": &queryHolder{`
		insert into Subscribed (User, Feed, Name) values (?1, ?2, ?3);
	`, nil},
	// /api/feed/unsubscribe
	"FeedUnsub1": &queryHolder{`
		delete from Subscribed where User = ?1 and Feed = ?2;
	`, nil},
	"FeedUnsub2": &queryHolder{`
		delete from PausedFlags where User = ?1 and Feed = ?2;
	`, nil},
	"FeedHasSubs": &queryHolder{`
		select exists(select 1 from Subscribed where Feed = ?1 limit 1);
	`, nil},
	"FeedDelete": &queryHolder{`
		delete from Feeds where ID = ?1;
	`, nil},
	// /api/feed/pause
	"FeedExists": &queryHolder{`
		select exists(select 1 from Feeds where ID = ?1 limit 1);
	`, nil},
	"FeedPause": &queryHolder{`
		insert into PausedFlags (User, Feed) values (?1, ?2);
	`, nil},
	// //api/feed/unpause
	"FeedUnpause": &queryHolder{`
		delete from PausedFlags where User = ?1 and Feed = ?2;
	`, nil},

	// /api/article/read
	"ArticleRead": &queryHolder{`
		insert into ReadFlags (User, Article) values (?1, ?2);
	`, nil},
	// /api/article/unread
	"ArticleUnread": &queryHolder{`
		delete from ReadFlags where User = ?1 and Article = ?2;
	`, nil},
	// /api/article/feed
	"GetUnread": &queryHolder{`
		select a.ID, a.Title, a.URL, fn.Name, a.Published from Articles a
		left join Subscribed fn on fn.Feed = a.Feed and fn.User = ?1 where (
			not a.ID in (select Article from ReadFlags where User = ?1) and
			not a.Feed in (select Feed from PausedFlags where User = ?1)
		) order by Published;
	`, nil},
}

func init() {
	var err error
	DB, err = sql.Open("sqlite3", "file:feeds.db")
	if err != nil {
		panic(err)
	}

	_, err = DB.Exec(InitCode)
	if err != nil {
		panic("Error loading DB init code:\n" + err.Error())
	}

	for _, v := range Queries {
		err := v.Init()
		if err != nil {
			panic("Error loading query: " + v.Code + "\n\n" + err.Error())
		}
	}
}

type queryHolder struct {
	Code   string
	Preped *sql.Stmt
}

func (q *queryHolder) Init() error {
	var err error
	q.Preped, err = DB.Prepare(q.Code)
	return err
}
