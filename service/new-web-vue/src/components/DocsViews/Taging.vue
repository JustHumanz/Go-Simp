<script setup>
import "./style.scss"
</script>
<template>
  <h3 class="content-title">Roles and Taging</h3>
  <h4 class="content-subtitle">Check Yourself</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>My tags</p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">/mytags</p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show following group or member vtuber in your account (including your
    following roles)
  </p>
  <h4 class="content-subtitle">Check Roles</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>Roles <span>@role-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">/role-info <span>role-name</span></p>
  <p>
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show following group or member vtuber in your roles <br /><br /><b
      >Example 1:</b
    >
  </p>
  <p class="content-code">vtbot>Roles <span>@Hololive</span></p>
  <p class="content-text">
    <b>Example 2 (Slash Command):</b>
  </p>
  <p class="content-code">/role-info <span>role-name: @nijisanji</span></p>
  <h4 class="content-subtitle">Taging Yourself</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>Tag me <span>group/member</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /tag-me <span>vtuber-group</span><span>vtuber-name</span
    ><span>reminder (per-minutes, optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small>
    <br /><br />
    Tag yourself for get any notification schedules and fanarts from one or many
    group or member vtuber. For <b>Slash Command</b>, you can enter
    <span class="code-span">vtuber-group</span> or
    <span class="code-span">vtuber-name</span>, or maybe both. But you can't
    empty both.<br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Tag me <span>kizuna_ai</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>Tag me <span>moona,pekora</span></p>
  <p class="content-text">
    <b>Example 3 (Slash Command):</b>
  </p>
  <p class="content-code">/tag-me <span>vtuber-group: brave</span></p>
  <p class="content-text">
    <b>Example 4 (Slash Command):</b>
  </p>
  <p class="content-code">/tag-me <span>vtuber-name: evelyn</span></p>
  <h4 class="content-subtitle">Taging Roles</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">
    vtbot>Tag roles <span>@role-name</span> <span>group/member</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /tag-roles <span>role-name</span><span>vtuber-group</span
    ><span>vtuber-name</span><span>reminder (per-minutes, optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Tag any roles for get any notification schedules and fanarts from one or
    many group or member vtuber. You can tag 2 or more roles to members or
    groups separated by space (taging multi role only support legacy command).
    Same like <b>Taging Yourself</b>, in <b>Slash Command</b> version you can
    enter <span class="code-span">vtuber-group</span>,
    <span class="code-span">vtuber-name</span>, or both. But you can't empty
    both. <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">
    vtbot>Tag roles <span>@maha5</span> <span>maha5</span>
  </p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">
    vtbot>Tag roles <span>@hololive</span> <span>@fbk-shimps</span>
    <span>fbk</span>
  </p>
  <p class="content-text">
    <b>Example 3:</b>
  </p>
  <p class="content-code">
    vtbot>Tag roles <span>@micomet</span> <span>miko-chi,sui-chan</span>
  </p>
  <p class="content-text">
    <b>Example 4 (Slash Command):</b>
  </p>
  <p class="content-code">
    /tag-roles <span>role-name: @maha5</span><span>vtuber-group: maha5</span>
  </p>
  <p class="content-text">
    <b>Example 5 (Slash Command):</b>
  </p>
  <p class="content-code">
    /tag-roles <span>role-name: @Watame</span><span>vtuber-name: sheep</span>
  </p>
  <h4 class="content-subtitle">Remove tag for yourself</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>Del tag <span>Group/Member</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /del-tag <span>vtuber-group</span> <span>vtuber-name</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Remove tag for group or member vtuber for yourself. And in
    <b>Slash Command</b>, you must fill in one or both like previous slash
    command. <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Del tag <span>arknet</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>Del tag <span>haato</span></p>
  <p class="content-text">
    <b>Example 3 (Slash Command):</b>
  </p>
  <p class="content-code">/del-tag <span>vtuber-group: kamitsubaki</span></p>
  <p class="content-text">
    <b>Example 4 (Slash Command):</b>
  </p>
  <p class="content-code">/del-tag <span>vtuber-name: gura</span></p>
  <h4 class="content-subtitle">Remove tag for roles</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">
    vtbot>Del roles <span>@role-name</span> <span>Group/Member</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /del-roles <span>role-name</span><span>vtuber-group</span
    ><span>vtuber-name</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    Same like previous command, but now you can remove tag for roles.
    <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">
    vtbot>Del roles <span>@kizuna-ai</span> <span>kizuna_ai</span>
  </p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">
    vtbot>Del roles <span>@chloe</span> <span>chloe</span>
  </p>
  <p class="content-text">
    <b>Example 3 (Slash Command):</b>
  </p>
  <p class="content-code">
    /del-roles <span>role-name: @voms</span><span>vtuber-group: voms</span>
  </p>
  <p class="content-text">
    <b>Example 4 (Slash Command):</b>
  </p>
  <p class="content-code">
    /del-roles <span>role-name: @doog</span><span>vtuber-name: Koronen</span>
  </p>
  <h4 class="content-subtitle">Reminders</h4>
  <p class="content-text">
    <b>For Yourself (Legacy Commands):</b>
  </p>
  <p class="content-code">
    vtbot>Set reminder <span>Group/Member</span> <span>Minutes</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>For Roles (Legacy Commands):</b>
  </p>
  <p class="content-code">
    vtbot>Roles reminder <span>@role-name</span> <span>Group/Member</span>
    <span>Minutes</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: Manage Channel or Higher</b></small
    ><br /><br />
    This command will set or update reminder for group or member vtuber. When
    you remove reminder time,you can set minutes to 0.
    <br /><br /><b>Example 1 (Reminder for yourself):</b>
  </p>
  <p class="content-code">
    vtbot>Set reminder <span>hololive</span> <span>20</span>
  </p>
  <p class="content-text">
    <b>Example 2 (Reminder for Roles):</b>
  </p>
  <p class="content-code">
    vtbot>Roles reminder <span>@holox</span> <span>orca</span> <span>30</span>
  </p>
  <p class="content-text">
    <b>Example 3 (Remove Reminder):</b>
  </p>
  <p class="content-code">
    vtbot>Set reminder <span>adelia</span> <span>0</span>
  </p>
  <p class="content-text">
    You can using <b>Slash Command</b> to set reminder, when your taging.
    <br /><br /><b>Example 1 (Reminder and taging yourself):</b>
  </p>
  <p class="content-code">
    /tag-me <span>vtuber-group: 774inc</span> <span>reminder: 20</span>
  </p>
  <p class="content-text">
    <b>Example 2 (Reminder and taging Roles):</b>
  </p>
  <p class="content-code">
    /tag-roles <span>role-name: @Tenshi</span> <span>vtuber-name: Tenshi</span>
    <span>reminder: 30</span>
  </p>
  <br />
</template>
