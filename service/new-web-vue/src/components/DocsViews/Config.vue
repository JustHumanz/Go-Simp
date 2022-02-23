<script setup>
import "./style.scss"
</script>

<template>
  <h3 class="content-title">Configuration</h3>
  <h4 class="content-subtitle">Setup (Old Version)</h4>
  <p class="content-code">vtbot>Setup</p>
  <p class="content-text">
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
    <b>Role permission required: Manage Channel or Higher</b>
  </p>
  <h4 class="content-subtitle">Setup (Slash Command)</h4>
  <p class="content-code">
    /setup channel-type <span>livestream/fanart/lewd</span>
  </p>
  <p class="content-text">
    Same like above, but you can set any schedules/fanart on on a different
    channel in <b>channel-name</b>.
    <i
      >When <b>vtuber-group</b> can't use group name, using
      <b>Group ID</b> instead. For list, check below!</i
    >
    <br />
    <br />
    <b>Vtuber Groups:</b>
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
    regular content like cover or other video) (livestream stage
    only)<br />
    <span class="code-span">newupcoming</span>: Show upcoming live
    streaming (livestream stage only)<br />
    <span class="code-span">dynamic</span>: Deleting schedule after ending live
    streaming (livestream stage only)<br />
    <span class="code-span">lite-mode</span>: Show only schedule live
    streaming without notification (livestream stage only)<br />
    <span class="code-span">indie-notif</span>:When get a notification from indie vtuber, recommend for set <b><i>True</i></b> <span class="code-span">vtuber-group</span> to <b>independent</b> or <b>10</b> (livestream stage only)<br />
    <span class="code-span">fanart</span>: Add fanart post about same group in same channel (optional, livestream stage only)<br /><br />
    <b>Example 1:</b>
  </p>
  <p class="content-code">
    /setup channel-type livestream <span>channel-name: #hololive</span>
    <span>vtuber-group: hololive</span><span>liveonly: False</span>
    <span>newupcoming: False</span><span>dynamic: False</span>
    <span>lite-mode: False</span><span>indie-notif: False</span>
  </p>
  <p>
    <b>Example 2:</b>
  </p>
  <p class="content-code">
    /setup channel-type fanart <span>channel-name: #niji-art</span>
    <span>vtuber-group: 6</span>
  </p>
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
