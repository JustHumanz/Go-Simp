v2.6.21
- Update redis
- Update cronjob to v.3
- Force retrun when failed send message
- Add Discord session Close 

v2.6.20
- Youtube checker livestream by count time
- Fix User not taged when vtuber got a milestone subs/followers
- Specially for Independent livestream will not be send if there no one user/role tagged
- Add Higuchi Kaede
- Disable youtube checker send notif
- Fix add in go-simp-web

v2.5.19
- Fix duplicate notif
- Fix users not pinged
- Add context to get userlist
- Fix typo in help
- Add async to get userlist
- Add wait when get userlist

v2.5.18
- Update help command
- Fix redis error handling
- Fix go-simp django

v2.5.17
- Add Register user via reacting
- Add Dynamic mode on bilibili

----------------------------------------------------------------------------------
v1.5.16
- Add Dynamic notif #83
- Fix twitter fanart high cpu usage
- Add DiscordChannel struct

v1.4.15
- Fix subs milistone not send
- Fix member cache
- Fix reminder time

v1.4.12
- Add Top.gg server count
- Remove async from bilibili fanart
- Fix youtube notif

v1.4.11
- Bring back module info
- Change reminder interval
- Add Youtube live bilibili checker 

v1.4.10
- Update discordgo
- Change `EnName` to `Name` in fanart scraper
- Change fanart scraper logic
- Add config format
- Update twitter-scraper
- Update youtube send notif format

v1.4.9
- Change Viwers/Online/Followers number to human readable
- Change `subscribe` to `info`

v1.3.8
- Fix redis malfunction
- Change bilibili-fanart logic
- Add youtube changer status in user handler
- Remove twitter with Quoted and Reply

v1.3.7
- Change backend to micro service
- Fix & change redis TTL
- Update reminder logic

v1.3.6
- Fix youtube not send notif
- Fix migrate token out of limit

v1.3.5
- Add redis for cache
- Set yt state to cache
- Set fanart checker to cache

v1.2.5
- Change wg.wait to sleep
- Change cronjob
- Add kano bilibili hashtag

v1.2.4
- Add wg.wait() every 10 members
- Change cron job

v1.2.3
- Change Go-Simp-Web to Web
- Fix twitter search query
- Change "/tmp" to const
- Change Youtube scraper to offical API
- Upgrade twitter-scraper

v1.2.2
- Add donation message
- Move hardcode config to toml file
- Change tmp dir

v1.1.1
- Add modul info
- Change tor node
- Fix youtube avatar scraper
- Add bot version info in help command

----------------------------------------------------------
v0.5.18
- Add flag to backend service
- Disable hololive from bilibili

v0.5.17
- Upgrade twitter-scraper
- Change fanart scraping logic

v0.5.16
- Replicate n0madic/twitter-scraper to JustHumanz/twitter-scraper
- Fix BiliBili ghost notif
- Add sleep for twitter scraping

v0.5.15
- Fix twitter scraping bug
- Update n0madic/twitter-scraper
- Change twitter filter
- Add image filter
- Update exec page
- Remove dirty func&struct

v0.5.14
- Change cron twitter scraping

v0.5.13
- Change logic twitter scraping
- Add donation message
- Change logic delete channel
- Add #53 #54 #55 #56

v0.5.12
- Add Nijisanji ID gen 5

v0.5.11
- Update exec doc
- Add channel remover on db migrate

v0.5.10
- Add disable `reminder time` command
- Add HoloID gen 2

v0.4.10
- Fix fanart #51
- Change some log format
- Add http proxy in .toml config
- Change struct name


v0.4.9
- Change twitter avatar scraping    
- Add discord id submitter

v0.4.8
- Migrate Guild handler to database
- Create `network` module
- Back to twitter-scraper
- Update multitor

v0.4.7
- Fix `Update` command out of array
- Change 404.jpg Thumbnail to null

v0.4.6
- Change Guide URL

v0.4.5
- Add reminder for roles state #35

v0.3.5
- Fix `Enable` command #34

v0.3.4
- Specially for Independent fan art will not be sent if there no one user/role tagged
- Add -liveonly -newupcoming and -rm_liveonly -rm_newupcoming see at [here](https://go-simp.human-z.tech/Exec/) 
- Add ArkNET  
- Remove `@here`

v0.2.4
- Change context timeout from 30s to 120s 
- Remove all user when a channel was disable one or more vtuber groups

v0.2.3
- Change Multiple Region in yt to single Region
- Change `Roles tags` to `Roles info`

v0.2.2
- Add customize reminder

v0.1.2
- Fix `Del me` command

v0.1.1
- Add usertag when Vtuber get milestone
- Change reminder from 1H and 30M to 30M
- Rename command from `channel tags` to `channel state`

v0.0.1
- Move holoCN to Independen
- Add BiliBili fanart
- Change Twitterscraper to official api
