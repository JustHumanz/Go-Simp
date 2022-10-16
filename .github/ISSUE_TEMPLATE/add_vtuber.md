---
name: Add vtuber request (Legacy)
about: Manually request a vtuber in GitHub Issues
title: "Add [Vtuber Nickname] from [Group/Agency]"
labels: enhancement
assignees: 'JustHumanz'

---
```json
{
	"GroupName": "",
	"GroupIcon": "",
	"GroupChannel": {
        "Youtube": [
            {
              "ChannelID": "",
              "Region": ""
            }
          ],
          "BiliBili": [
            {
              "BiliBili_ID": 0,
              "BiliRoom_ID": 0,
              "Region": ""
            }
          ]
	},
	"Members": [
        {
            "Name": "",
            "EN_Name": "",
            "JP_Name": "",
            "Twitter": {
              "Twitter_Fanart": "",
              "Twitter_Lewd": "",
              "Twitter_Username": ""
            },
            "Youtube": {
              "Yt_ID": ""
            },
            "BiliBili": {
              "BiliBili_Fanart": "",
              "BiliBili_ID": 0,
              "BiliRoom_ID": 0
            },
            "Twitch": null,
            "Region": "",
            "Fanbase": "",
            "Status": ""
        },
	]
}
```

### Example :

```json
{
	"GroupName": "MKLNtic",
	"GroupIcon": "https://cdn.humanz.moe/MKLNtic.png",
	"GroupChannel": {
        "Youtube": [
            {
              "ChannelID": "UCN3M-Nwa-eaZuMhPb5c3Pww",
              "Region": "JP"
            }
          ],
          "BiliBili": null 
	},
	"Members": [
        {
            "Name": "Kanochi",
            "EN_Name": "Mahoro Kano",
            "JP_Name": "鹿乃まほろ",
            "Twitter": {
              "Twitter_Fanart": "#まほろ絵",
              "Twitter_Lewd": "",
              "Twitter_Username": "kanomahoro"
            },
            "Youtube": {
              "Yt_ID": "UCfuz6xYbYFGsWWBi3SpJI1w"
            },
            "BiliBili": {
              "BiliBili_Fanart": "鹿乃絵",
              "BiliBili_ID": 316381099,
              "BiliRoom_ID": 15152878
            },
            "Twitch": null,
            "Region": "JP",
            "Fanbase": "",
            "Status": "Active"
        },
		{
			"Name": "Lon",
			"EN_Name": "Lon",
			"JP_Name": "ろん",
			"Twitter": {
				"Twitter_Fanart": "#LONART",
				"Twitter_Lewd": "",
				"Twitter_Username": "LOOON4"
			},
			"Youtube": {
				"Yt_ID": "UCSuVf5hodNJftAv0Zio-vXA"
			},
			"BiliBili": {
                "BiliBili_Fanart": "",
                "BiliBili_ID": 1253935371,
                "BiliRoom_ID": 25600815
            },
			"Twitch": null,
			"Region": "JP",
			"Fanbase": "",
			"Status": "Active"
		}
	]
}
```