# Utilities

## Prediction Followers

```slash
/prediction slash{platform} slash{vtuber-name}
```

small{**Role permission required: No permission needed**}
%br%
Show member subscriber or followers prediction from **YouTube**, **Twitch**, **BiliBili**, and **Twitter** in next 7 days.

### Example

```slash
/prediction slash{platform: twitch} slash{vtuber-name: iofi}
```

## Get Random Fan Art

```slash
/fanart slash{vtuber-group} slash{vtuber-name (Optional)}
```

small{**Role permission required: No permission needed**}
%br%
Show randomly selected fan art for a vtuber member in **twitter** or **pixiv**. And because slash{vtuber-group} required, make sure it's the same between slash{vtuber-group} and slash{vtuber-name}).

### Example 1 (Get Random Fan Art for random members in group)

```slash
/fanart slash{vtuber-group: nijisanji}
```

### Example 2 (Get Random Fan Art for selected member)

```slash
/fanart slash{vtuber-group: independent} slash{vtuber-name: Ui-mama}
```

## Get Random Fan Art (NSFW)

```slash
/lewd vtuber-group slash{vtuber-name (Optional)}
```

small{**Role permission required: No permission needed**}
%br%
Same like random fan art command, but for NSFW/Lewd from lewd hastag or #R-18 in pixiv.

### Example 1 (Get Random Fan Art for random members in group)

```slash
/lewd slash{vtuber-group: yume_reality}
```

### Example 2 (Get Random Fan Art for selected member)

```slash
/lewd slash{vtuber-group: hololive} slash{vtuber-name: konfauna}
```

## Get Detail Member

```slash
/info slash{vtuber-name}
```

small{**Role permission required: No permission needed**}
%br%
Show detail and total follower any platform from single member

### Example

```slash
/info slash{vtuber-name: gura}
```

## Livesteam Utilites

```slash
/livesteam slash{state} slash{status} slash{vtuber-group} slash{vtuber-name (optional)} slash{region (optional)}
```

Show any livestream from members or group in **YouTube**, **Twitch**, and **BiliBili**. When you find by groups, live streaming info is limited to 5 only (exept Live Status).

### Available States

| Platform | Upcoming        | Live        | Past            |
| -------- | --------------- | ----------- | --------------- |
| YouTube  | ✓ Max 5 Members | ✓ Available | ✓ Max 5 Members |
| Twitch   | ✕ Not Available | ✓ Available | ✓ Max 5 Members |
| BiliBili | ✕ Not Available | ✓ Available | ✓ Max 5 Members |

### Example 1 (Show Upcoming Group Livestreams)

```slash
/livesteam slash{state: YouTube} slash{status: Upcoming} slash{vtuber-group: nijisanji}
```

### Example 2 (Show Group Live Today in a Single Region)

```slash
/livesteam slash{state: Twitch} slash{status: Live} slash{vtuber-group: Hololive} slash{region: JP}
```

### Example 3 (Show Past Livestream from Single Member)

```slash
/livesteam slash{state: YouTube} slash{status: Past} slash{vtuber-group: Re:Memories} slash{vtuber-name: reynald}
```

## Bot Information

```slash
/help EN
```

small{**Role permission required: No permission needed**}
%br%
Shows the basic information's about the bot.
