v3.10.1
- Change pixiv fanart workflow

v3.10.0
- Fix prediction
- Change date to counter num in prediction
- Fix state icon on fanart command
- Add twitch followers & viwers

v3.9.2
- Add nijisanji EN
- Change error to panic when drop cache
- Update new vtuber message
- Change twitter fanart workflow
- Change twitter fanart interval
- Add `LiveURL` in rest-api

v3.9.1
- Update error payload when nil payload
- Swap X to Y in prediction
- Remove days in api prediction

v3.9.0
- Add Vtuber followers/subscriber prediction
- Change youtube channel checker interval
- Remove usertags when no one user tagged in upcoming yt
-  Add Tsunderia 

v3.8.0
- Change Youtube workflow
- Add Calli lewd hashtag
- Add Metric
- Add Prometheus

v3.7.0
- Add dark mode on go-simp web
- Add search bar on go-simp web

v3.6.1-beta
- Update rest-api
- Add rest-api doc
- disable new upcoming livestream if livestream less than 6h

v3.6.0
- Add Group/Agency in youtube/bilibili
- Fix pilot memory leak
- Update donate message
- Remove request payload every 2h
- Fix image not showing in web
- Update Vshojo&arknet
- Update Rest-API

v3.6.0-beta
- Add Group/Agency in youtube/bilibili
- Fix pilot memory leak
- Update donate message

v3.5.6
- Change default views to "???"
- Update fanart notif&Add sync in pixiv fanart

v3.5.5
- Fix pixiv not send any fanart
- Fix rest-api

v3.5.4
- Blacklist loli tag in NSFW
- Remove `-` in fanart notif
- Fix youtube viwers

v3.5.3
- Fix `update` command
- Add Pixiv session outdate alert

v3.5.2
- Update rest api to V2
- Update Go simp web
- Update utility
- Update dependencies
- Change twitch auth

v3.5.1
- Fix Pixiv 401
- Change data.bilibili to live.bilibili

v3.5.0
- Update pixiv fanart scrape
- Add Twitch `live`/`past` command
- Fix live bilibili not send notif
- Fix space bilibili vtuber name

v3.5.0-beta
- Update pixiv fanart scrape
- Add Twitch `live`/`past` command

v3.4.3
- Fix bilibili space not send any new video
- Fix live bilibili not update status
- Rework logic livestream service for sending notif

v3.4.2
- Add ViViD vtuber group
- Add Kuri from Independen #150
- Update go-simp web exec

v3.4.1
- Add NSFW checker

v3.4.0
- Remove Danbooru lewd
- Add pixiv R18
- Add `lewd` command

v3.3.5
- Fix name not showed in upcoming livestream
- Fix livestream not sended

v3.3.4
- Fix bug on `setup`
- Add color on pixiv image
- Fix livestream conter not running

v3.3.3
- Fix swaped youtube thub and avatar
- Use const for paramater

v3.3.2
- Fix Youtube upcoming/live/past command
- Fix Bilibili fanart nill pointer
- Add twitch embed inline 

v3.3.1
- Add pixiv as fanart
- Rebase fanart & livestream struct
- Remove call sql

v3.2.0
- Add `upvote` command
- Add fanbase name
- Fix lewd
- Rework Fanart logic
- Remove fanart checker

v3.1.5
- Change some grammar (web)
- Update go simp web
- Lewd disable multiple pic
- Select Group by name
- Fix worng vtuber name
- Fanart fix filter

v3.1.4
- Change Helsinki/Finland to Europe/Helsinki

v3.1.3
- Fix danbooru not send lewd
- Update dependency 
- Fix Independen notif in update menu
- Fix twitch live status
- Change some grammar
- Change bot desc
- Add bot status
- Add VShojo group
- Fix migrate can't send message

v3.1.2
- Remove get color in danbooru

v3.1.1
- Add caching in static files
- send fanart from Pixvi use discord Embed
- Show version info in help menu

v3.0.1
- Add `Change Livestream state menu`
- Remove parent in danbooru

## v3.0.0

v2.9.8
- Add Mea
- Change hard config to variable config
- Upgrade dependency
- Fix nill channel & nill upcoming
- Add danbooru fanart

v2.9.8-alpha
- Add danbooru fanart

v2.8.8
- Remove add,update command
- Fix guild invite log
- Change webhook to channelID
- Add `LiteMode`
- Add `IndieNotif`
- Update change state

v2.4.7
- Change fe wait workflow
- Enhance `updatev2` to `update`
- Add Support region
- Move time config to config package

v2.4.7-beta
- Add `updatev2`
- Add Support region

v2.3.7
- Change HeartBeat workflow
- Add register emoji when a new vtuber added

v2.2.7
- Add gRPC for for service communicate
- Update 774inc
- Add Ameka&Lili
- Separate invite log with pilot log

v2.2.7-beta
- Add gRPC for for service communicate

v2.1.7
- Fix guild channel notif
- Change Tbilibili fanart workflow
- Add LowResources options

v2.0.6
- Change 774.inc to 774inc
- Update 774inc
- Add k8s mainfest
- Fix go-simp-web debug env

v2.0.5
- Add Twitch support
- Add Ririsya #112
- Add Yayanehi #111
- Change Lofi to Iofi
- Fix 404.jpg

v2.0.4
- Update redis
- Update cronjob to v.3
- Force retrun when failed send message
- Add Discord session Close 
- Fix reminder

v2.0.3
- Youtube checker livestream by count time
- Fix User not taged when vtuber got a milestone subs/followers
- Specially for Independent livestream will not be send if there no one user/role tagged
- Add Higuchi Kaede
- Disable youtube checker send notif
- Fix add in go-simp-web

v2.0.2
- Fix duplicate notif
- Fix users not pinged
- Add context to get userlist
- Fix typo in help
- Add async to get userlist
- Add wait when get userlist

v2.0.1
- Update help command
- Fix redis error handling
- Fix go-simp django

## v2.0.0
- Add Register user via reacting
- Add Dynamic mode on bilibili

----------------------------------------------------------------------------------
v1.3.14
- Add Dynamic notif #83
- Fix twitter fanart high cpu usage
- Add DiscordChannel struct

v1.2.13
- Fix subs milistone not send
- Fix member cache
- Fix reminder time

v1.2.12
- Add Top.gg server count
- Remove async from bilibili fanart
- Fix youtube notif

v1.2.11
- Bring back module info
- Change reminder interval
- Add Youtube live bilibili checker 

v1.2.10
- Update discordgo
- Change `EnName` to `Name` in fanart scraper
- Change fanart scraper logic
- Add config format
- Update twitter-scraper
- Update youtube send notif format

v1.2.9
- Change Viwers/Online/Followers number to human readable
- Change `subscribe` to `info`

v1.1.9
- Fix redis malfunction
- Change bilibili-fanart logic
- Add youtube changer status in user handler
- Remove twitter with Quoted and Reply

v1.1.7
- Change backend to micro service
- Fix & change redis TTL
- Update reminder logic

v1.1.6
- Fix youtube not send notif
- Fix migrate token out of limit

v1.1.5
- Add redis for cache
- Set yt state to cache
- Set fanart checker to cache

v1.1.4
- Change wg.wait to sleep
- Change cronjob
- Add kano bilibili hashtag

v1.1.3
- Add wg.wait() every 10 members
- Change cron job

v1.1.2
- Change Go-Simp-Web to Web
- Fix twitter search query
- Change "/tmp" to const
- Change Youtube scraper to offical API
- Upgrade twitter-scraper

v1.1.1
- Add donation message
- Move hardcode config to toml file
- Change tmp dir

## v1.0.0
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
