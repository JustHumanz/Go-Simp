<script setup>
import "./style.scss"
</script>

<template>
  <h3 class="content-title">Configuration</h3>
  <h4 class="content-subtitle">Setup</h4>
  <p class="content-code">vtbot>Setup</p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    This command will bring you in setup mode and
    <i
      >all schedules/fanart will be displayed on the channel that has been
      commanded</i
    >
    <img
      src="https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/RemakeSetup.png"
      alt=""
      class="content-image"
    />
    You can setup <b>Slash Command</b>. Unlike legacy command, you can set any schedules/fanart on on a different
    channel in <span class="code-span">channel-name</span>. <br><br> <b>Setup Livestreaming:</b>
  </p>
  <p class="content-code">
    /setup channel-type livestream <span>​channel-name</span> <span>vtuber-group</span> <span>liveonly</span> <span>newupcoming</span> <span>dynamic</span> <span>lite-mode</span> <span>indie-notif</span> <span>fanart (optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    <b>Setup Fanart:</b>
  </p>
  <p class="content-code">
    /setup channel-type fanart <span>​channel-name</span> <span>vtuber-group</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    <b>Setup NSFW/Lewd Fanart:</b>
  </p>
  <p class="content-code">
    /setup channel-type lewd <span>​channel-name</span> <span>vtuber-group</span> 
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
  <!-- Make a table -->
  <!-- <div class="flex items-start">
    <table class="content-table" v-for="group in groups" :key="group">
      <tr>
        <th>Group ID</th>
        <th>Group Name</th>
      </tr>
      <tr v-for="g in group" :key="g.ID">
        <td>
          {{ g.ID }}
        </td>
        <td>
          {{ g.GroupName }}
        </td>
      </tr>
    </table>
  </div> -->
    <b>Options:</b> <br />
    <span class="code-span">channel-name</span> Available <b>channels</b> on
    your discord server <br />
    <span class="code-span">vtuber-group</span> Any Name Group or ID Group<br />
    <span class="code-span">liveonly</span> Show only live streaming (without
    regular content like cover or other video) (livestream stage only)<br />
    <span class="code-span">newupcoming</span> Show upcoming live streaming
    (livestream stage only)<br />
    <span class="code-span">dynamic</span> Show schedule and deleted after past
    live streaming (livestream stage only)<br />
    <span class="code-span">lite-mode</span> Disabling ping user/role function
    (livestream stage only)<br />
    <span class="code-span">indie-notif</span> When get a notification from
    indie vtuber, when <span class="code-span">vtuber-group</span> is <b>independent</b> (livestream stage only)<br />
    <span class="code-span">fanart</span> Additional post fanart in same channel (optional, livestream stage only)<br /><br />
    <b>Example 1:</b>
  </p>
  <p class="content-code">
    /setup channel-type livestream <span>channel-name: #hololive</span>
    <span>vtuber-group: hololive</span><span>liveonly: False</span>
    <span>newupcoming: False</span><span>dynamic: False</span>
    <span>lite-mode: False</span><span>indie-notif: False</span>
  </p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">
    /setup channel-type fanart <span>channel-name: #niji-art</span>
    <span>vtuber-group: 6</span>
  </p>
  <h4 class="content-subtitle">Checking Stage (Slash Command)</h4>
  <p class="content-code">/channel-state <span>channel-name</span></p>
  <p class="content-text">
      <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Checking existing vtuber group on
    <span class="code-span">channel-name</span><br /><br /><b>Example:</b>
  </p>
  <p class="content-code">
    /channel-state <span>channel-name: #re-memories</span>
  </p>
  <h4 class="content-subtitle">Change Stage</h4>
  <p class="content-code">vtbot>Update</p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Change any existing vtuber group on channel your command, like adding
    another region/removing a region or changing from live stream to fan art.<br />
    For aternative, use <b>Slash Command</b> for update in different channel.
    <br /><br /><b>Slash Commands:</b>
  </p>
  <p class="content-code">/channel-update <span>channel-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    >
    <br /><br /><b>Example:</b>
  </p>
  <p class="content-code">
    /channel-update <span>channel-name: #kizuna-ai</span>
  </p>
  <h4 class="content-subtitle">Disable/Remove Stage</h4>
  <p class="content-text">
    <b>Disable Partial Group (Legacy only):</b>
  </p>
  <p class="content-code">vtbot>Disable <span>group vtuber</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Disable/remove single or several vtuber groups in channel your command.
    <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Disable <span>VOMS</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>Disable <span>hololive,holostars</span></p>
  <p class="content-text">
    <b>Disable All group (Slash Command):</b>
  </p>
  <p class="content-code">/channel-update <span>channel-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Disable/remove any vtuber groups in channel inside <span class="code-span">channel-name</span>.
    <br /><br /><b>Example :</b>
  </p>
  <p class="content-code">/channel-update <span>channel-name: #vshojo</span></p>
  <br />
</template>

<script>
import axios from "axios"
import Config from "../../config.json"

export default {
  data() {
    return {
      groups: [],
    }
  },
  async created() {
    // Get all groups
    const group_data = await axios
      .get(Config.REST_API + "/groups/")
      .then((response) => response.data)

    const array1 = []
    const array2 = []

    group_data.forEach((group) => {
      if (array1.length < group_data.length / 2) array1.push(group)
      else array2.push(group)
    })

    this.groups = [array1, array2]
  },
}
</script>
