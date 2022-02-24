<script>
import "./style.scss"
</script>

<template>
  <h3 class="content-title">Command Utilites</h3>
  <h4 class="content-subtitle">Prediction Followers</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">
    vtbot>prediction <span>platform</span> <span>Vtuber Name</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /prediction <span>platform</span> <span>vtuber-name</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show subscriber or followers prediction in next 7 days for a vtuber member
    <br /><br />
    <b>Platform (For Legacy Commands):</b> <br />
    <span class="code-span">-yt</span> YouTube <br />
    <span class="code-span">-tw</span> Twitter <br />
    <span class="code-span">-bl</span> Bilibili <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">
    vtbot>prediction <span>-yt</span> <span>iofi</span>
  </p>
  <p class="content-text">
    <b>Example 2 (Slash Command):</b>
  </p>
  <p class="content-code">
    /prediction <span>platform: twitch</span> <span>vtuber-name: miko-chi</span>
  </p>
  <h4 class="content-subtitle">Get Random Fan Art</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>fanart <span>Group/Member</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /fanart <span>vtuber-group</span> <span>vtuber-name (Optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show randomly selected fan art for a vtuber member in twitter or pixiv (for
    <b>Slash Command</b>, because
    <span class="code-span">vtuber-group</span> required, make sure it's the
    same between <span class="code-span">vtuber-group</span> and
    <span class="code-span">vtuber-name</span>). <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>fanart <span>Kanochi</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>fanart <span>Voms</span></p>
  <p class="content-text">
    <b>Example 3 (Slash Command):</b>
  </p>
  <p class="content-code">
    /fanart <span>vtuber-group: independent</span>
    <span>vtuber-name: Ui-mama</span>
  </p>
  <h4 class="content-subtitle">Get Random NSFW/Lewd Fan Art</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>lewd <span>Group/Member</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">
    /lewd <span>vtuber-group</span> <span>vtuber-name (Optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Same like random fan art commands, but for NSFW/Lewd from lewd hastag or
    #R-18 in pixiv.
    <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>lewd <span>Kanochi</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">vtbot>lewd <span>yume_reality</span></p>
  <p class="content-text">
    <b>Example 3 (Slash Command):</b>
  </p>
  <p class="content-code">
    /lewd <span>vtuber-group: nijisanji</span>
    <span>vtuber-name: Kaederon</span>
  </p>
  <h4 class="content-subtitle">Show Vtubers Members from Groups</h4>
  <p class="content-code">
    vtbot>Vtuber data <span>Group/Member</span> -region <span>region</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Shows available Vtuber data from groups. You can filter by region using tag
    <span class="code-span">-region</span>. <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Vtuber data <span>774inc</span></p>
  <p class="content-text">
    <b>Example 2:</b>
  </p>
  <p class="content-code">
    vtbot>Vtuber data <span>hololive,nijisanji</span> -region <span>jp</span>
  </p>
  <h4 class="content-subtitle">Get A Member Data</h4>
  <p class="content-text">
    <b>Legacy Commands:</b>
  </p>
  <p class="content-code">vtbot>Info <span>Vtuber Name</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Commands:</b>
  </p>
  <p class="content-code">/info <span>vtuber-name</span></p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show detail and follower any platform from single member.
    <br /><br /><b>Example 1:</b>
  </p>
  <p class="content-code">vtbot>Info <span>Sora</span></p>
  <p class="content-text">
    <b>Example 2 (Slash Command):</b>
  </p>
  <p class="content-code">/info <span>vtuber-name: Mito-chan</span></p>
  <h4 class="content-subtitle">Livesteam Utilites</h4>
  <p class="content-text">
    <b>Youtube:</b>
  </p>
  <p class="content-code">
    yt><span>Upcoming/Live/Past</span> <span>Group/Member</span> -region
    <span>region</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Twitch:</b>
  </p>
  <p class="content-code">
    tw><span>Live/Past</span> <span>Group/Member</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Bilibili:</b>
  </p>
  <p class="content-code">
    bl><span>Live/Past</span> <span>Group/Member</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    <b>Slash Command:</b>
  </p>
  <p class="content-code">
    /livesteam
    <span>state</span
    >&nbsp;<span>status</span>&nbsp;<span>vtuber-group</span>&nbsp;<span
      >vtuber-name (optional)</span
    >&nbsp;<span>region (optional)</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show any livestream from members or group. When you find by groups, live
    streaming info is limited to 3 only (exept <b>Live</b> Status).
    <br /><br /><b>Example 1 (Upcoming a single member in YouTube):</b>
  </p>
  <p class="content-code">yt>Upcoming <span>Pikamee</span></p>
  <p class="content-text">
    <b>Example 2 (Check members live from a Group in English Region):</b>
  </p>
  <p class="content-code">
    yt>Live <span>nijisanji</span> -region <span>en</span>
  </p>
  <p class="content-text">
    <b>Example 3 (Check passing live streams on Twitch):</b>
  </p>
  <p class="content-code">tw>Past <span>froot</span></p>
  <p class="content-text">
    <b>Example 4 (Check live streams on Bilibili):</b>
  </p>
  <p class="content-code">bl>Live <span>chiharu</span></p>
  <p class="content-text">
    <b
      >Example 5 (Check passing live streams on YouTube using Slash Command):</b
    >
  </p>
  <p class="content-code">
    /livesteam
    <span>state: YouTube</span>&nbsp;<span>status: Past</span>&nbsp;<span
      >vtuber-group: Re:Memories</span
    >&nbsp;<span>vtuber-name: reynald</span>
  </p>
  <h4 class="content-subtitle">Get Latest Video in BiliBili</h4>
  <p class="content-code">
    sp_bl><span>Group/Member</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Show latest video vtuber from BiliBili (only 3 if use Vtuber Group).
    <br /><br /><b>Example:</b>
</p>
  <p class="content-code">sp_bl><span>Kana-Kana</span></p>
  <h2 class="content-subtitle">Bot Information</h2>
  <p class="content-code">
    vtbot>Help <span>EN/JP</span>
  </p>
  <p class="content-text">
    <small><b>Role permission required: No permission needed</b></small
    ><br /><br />
    Shows the basic information's about the bot in English or Japanese.
    </p>
    <br>
</template>
