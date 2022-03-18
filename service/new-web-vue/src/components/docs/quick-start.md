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
  
  router-link{« Read More}(/docs/roles-and-taging#taging-roles)

  ### About independent group
  Independent groups have a strict rule, if no **users/roles** are tagged than **live/fan art/lewd** art won't send anything. Recommended set slash{indie-notif} to **True** in router-link{setup}(/docs/configuration#setup-live-streaming) and notifications will be send even if no **users/roles** are tagged. 

  