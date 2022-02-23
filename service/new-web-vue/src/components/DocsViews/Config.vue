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
  </p>
  <h4 class="content-subtitle">Setup (Slash Command)</h4>
  <p class="content-code">
    /setup channel-type <span>livestream/fanart/lewd</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Same like above, but you can set any schedules/fanart on on a different
    channel in <span class="code-span">channel-name</span>.
    <i
      >When <span class="code-span">vtuber-group</span> can't use <b>group name</b>, using
      <b>Group ID</b> instead. For list, check below!</i
    >
    <br />
    <br />
    <b for="vtuber-groups">Vtuber Groups:</b>
  </p>
  <!-- Make a table -->
  <div class="flex items-start">
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
  </div>
  <p class="content-text">
    <br />
    <b>Context command:</b> <br />
    <span class="code-span">channel-name</span>: Available <b>channels</b> on
    your discord server <br />
    <span class="code-span">vtuber-group</span>: Available Group Name or Group
    ID in above tables <br />
    <span class="code-span">liveonly</span>: Show only live streaming (without
    regular content like cover or other video) (livestream stage only)<br />
    <span class="code-span">newupcoming</span>: Show upcoming live streaming
    (livestream stage only)<br />
    <span class="code-span">dynamic</span>: Show schedule and deleted after past
    live streaming (livestream stage only)<br />
    <span class="code-span">lite-mode</span>: Disabling ping user/role function
    (livestream stage only)<br />
    <span class="code-span">indie-notif</span>:When get a notification from
    indie vtuber, recommend for set <b><i>True</i></b
    >&nbsp;<span class="code-span">vtuber-group</span> to <b>independent</b> or
    <b>10</b> (livestream stage only)<br />
    <span class="code-span">fanart</span>: Add fanart post about same group in
    same channel (optional, livestream stage only)<br /><br />
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
    Checking any existing stage group on
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
    another region/removing a region or changing from live stream to fan art.
  </p>
  <h4 class="content-subtitle">Change Stage (Slash Command)</h4>
  <p class="content-code">/channel-update <span>channel-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Like previous command, but you can changing separated channel using
    <span class="code-span">channel-name</span>
    <br /><br /><b>Example:</b>
  </p>
  <p class="content-code">
    /channel-update <span>channel-name: #kizuna-ai</span>
  </p>
  <h4 class="content-subtitle">Disable</h4>
  <p class="content-code">vtbot>Disable <span>group vtuber</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Disable/remove single or several groups in channel your command.
    <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Disable <span>VOMS</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>Disable <span>hololive,holostars</span></p>
  <h4 class="content-subtitle">Disable (Slash Command)</h4>
  <p class="content-code">/channel-update <span>channel-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Disable/remove any groups in channel inside <span class="code-span">channel-name</span>
    <br /><br /><b>Example :</b>
  </p>
  <p class="content-code">/channel-update <span>channel-name: #vshojo</span></p>
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
