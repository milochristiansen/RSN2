#!/bin/sh

touch feeds.db
docker run -d -p 2053:443 -v $PWD/feeds.db:/app/feeds.db --env-file project.key rsn2
