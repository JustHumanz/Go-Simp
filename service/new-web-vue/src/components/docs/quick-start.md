  # Quick Start

  ### Setup Live Streaming
  ```slash
  /setup channel-type livestream slash{channel-name: channel{#hololive}} slash{vtuber-group: hololive}
  ```
  router-link{« Read More}(/docs/configuration#setup-live-streaming)

  ### Update Stage
  ```slash
  /channel-update slash{channel-name: channel{#hololive}}
  ```
  router-link{« Read More}(/docs/configuration#update-stage)

  ### Tag Roles
  ```slash
  /tag-role slash{role-name: role{@Holo Simps}} slash{vtuber-group: hololive}
  ```
  this means that the bot will mention role{@Holo Simps} when any new hololive fan arts or live streams are uploaded. 
  router-link{« Read More}(/docs/roles-and-taging#tag-roles)

  ### Get Vtuber Group and Vtuber Name
  center{![get-group-name](/src/assets/docs/get-group.png)}
  Your can find **vtuber-group** by click **group menu** inside vtuber list and get **vtuber-group** in right side of Group Name.
  center{![get-vtuber-name1](/src/assets/docs/get-member-card.png)}
  For **vtuber-name**, your can find below **Name Charachers** inside card list.
  center{![get-vtuber-name2](/src/assets/docs/get-member-detail.png)}
  Or in **Member Detail**, **vtuber-name** it's on the right of the **Name Charachers** or in **Nickname**.