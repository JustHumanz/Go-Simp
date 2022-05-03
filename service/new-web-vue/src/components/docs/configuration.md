# Configuration

## Setup

### Setup Live Streaming

```slash
/setup channel-type livestream slash{channel-name} slash{vtuber-group} slash{liveonly} slash{newupcoming} slash{dynamic} slash{lite-mode} slash{indie-notif} slash{fanart (optional)}
```

small{**Role permission required: Manage Channel or Higher**}

### Setup Fanart

```slash
/setup channel-type fanart ​slash{channel-name} slash{vtuber-group}
```

small{**Role permission required: Manage Channel or Higher**}

### Setup NSFW/Lewd Fanart

```slash
/setup channel-type lewd ​slash{channel-name} slash{vtuber-group}
```

small{**Role permission required: Manage Channel or Higher**}

### Options

- slash{**channel-name**} your channel name in **Discord server**.
- slash{**vtuber-group**} available router-link{vtuber group name}(/docs/get-data-groups).
- slash{**liveonly**} only show live streaming (**without regular content**) (**_livestream stage only_**).
- slash{**newupcoming**} Show upcoming live streaming (**_livestream stage only_**).
- slash{**dynamic**} Show schedule and deleted after past live streaming (**_livestream stage only_**).
- slash{**lite-mode**} Disabling ping user/role function (**_livestream stage only_**).
- slash{**indie-notif**} Show all vtuber **Independent** notification when set slash{**vtuber-group**} to **_independent_** (**_livestream stage only_**).
- slash{**fanart**} Additional show random fanart in same channel (**_optional_**) (**_livestream stage only_**).

### Example 1 (Set Livestream)

```slash
/setup channel-type livestream slash{channel-name: channel{#hololive}} slash{vtuber-group: hololive} slash{liveonly: False} slash{newupcoming: False} slash{dynamic: False} slash{lite-mode: False} slash{indie-notif: False}
```

### Example 2 (Add independent vtuber and set fanart)

```slash
/setup channel-type livestream slash{channel-name: channel{#inde-notif}} slash{vtuber-group: independent} slash{liveonly: False} slash{newupcoming: True} slash{dynamic: True} slash{lite-mode: False} slash{indie-notif: True} slash{fanart: True}
```

### Example 3 (Set random Fan Art Vtuber)

```slash
/setup channel-type fanart slash{channel-name: channel{#niji-fanart}} slash{vtuber-group: nijisanji}
```

%br%
Or you can set on website [Here](https://web-admin.humanz.moe/login)

## Checking Stage

```slash
/channel-state slash{channel-name}
```

small{**Role permission required: No permission needed**}
%br%
Checking existing vtuber group on slash{channel-name}

### Example

```slash
/channel-state slash{channel-name: channel{#re:memories}}
```

## Update Stage

```slash
/channel-update slash{channel-name}
```

small{**Role permission required: Manage Channel or Higher**}
%br%
Change any existing vtuber group on slash{channel-name}, like adding another region/removing a region or changing/add fan art from livesteam.

### Example

```slash
/channel-update slash{channel-name: channel{#kizuna-ai}}
```

## Disable/Remove Stage

```slash
/channel-delete slash{channel-delete}
```

small{**Role permission required: Manage Channel or Higher**}
%br%
Disable/remove any vtuber groups in channel inside slash{channel-delete}.

### Example

```slash
/channel-delete slash{channel-name: #vshojo}
```
