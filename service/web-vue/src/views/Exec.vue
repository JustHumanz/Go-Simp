<template>
  <div class="Exec container">
         <h3><a id="Command_0"></a>Command of VTbot</h3>
         <h4>Info</h4><br>
         <p><code>vtbot&gt;Help EN</code><br>
            Shows the basic information's about the bot.
         </p>
         <p><code>vtbot&gt;Help JP</code><br>
            Shows the basic information's about the bot but in Japanese.
         </p>
         <p><code>[]</code> =&gt; supports multiple [Names/Groups] (<em>separated by commas</em>)<br>
            <code>{}</code> =&gt; only support single VtuberName/VtuberGroup <br>
            Utility command don't need enable state,you can execute that command anywhere
         </p>
         <p>
            <b>Important rule</b><br>
            Independent groups have a strict rule, if no [Users/roles] are tagged than [live/fan art/lewd art] won't send anything. <br>
            recommended that you enable "independent notif" so notifications will be send even if no [users/roles] are tagged.
         </p>
         <h4>Main</h4><br>
         <p><code>vtbot&gt;Channel state</code><br>
            This command shows the state of the channel, what is enabled or disabled in this channel. <br>
            Role permission required: No permission needed
         </p>
         <p><code>vtbot&gt;Setup</code><br>
            This command will bring you in setup mode <br>
            Example: <code>vtbot&gt;Setup</code> <br>
            <img src="https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/RemakeSetup.png" alt=""><br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p> 
         <p><code>vtbot&gt;Update</code><br>
            This command will update your discord channel state, like adding another region/removing  a region or changing from live stream to fan art.<br>
            Example: <code>vtbot&gt;Update</code> <br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p> 
         <p><code>vtbot&gt;Disable [Vtuber Group]</code><br>
            This command will disable/remove a [Group].<br>
            Example: <code>vtbot&gt;Disable nijisanji</code> <br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p>
         <p>
            Options: <br>
            <code>dynamic</code> Set channel to dynamic notification,Remove past livestream notification from channel<br>
            <code>liveonly</code> Set livestreams in strict mode(ignoring covering or regular video) notification<br>
            <code>newupcoming</code> Enable new livestreams notification(by default will enable this) <br>
            <code>litemode</code> Disabling ping user/role function <br>
            <code>lewd</code> Send a NSFW pic from pixivi r18 or twitter with lewd hashtag <br>
         </p>
         <hr>
         <h4>Roles & Tags</h4>
         <p><code>vtbot&gt;Roles info [roles]</code><br>
            Show all roles info <br>
            Role permission required: No permission needed
         </p>
         <p><code>vtbot&gt;Tag me [Group/Member]</code><br>
            This command will add you to the tag list if any new fan art or live stream is uploaded. <br>
            Example: <code>vtbot&gt;tag me Kanochi</code>,then you will tagged when there is a new fanart or livestream(depends by channel state) of Kanochi <br>
            deprecated <br>
            <s>in other case if you want to remind before livestream started you can add <code>-setreminder {Minutes}</code> in <code>tag me</code><br>
            Example: <code>vtbot&gt;tag me hanayori -setreminder 30</code> so bot will remind you 30 Minutes before livestream started</s><br>
            Role permission required: No permission needed
         </p>
         <p><code>vtbot&gt;Del tag [Group/Member]</code><br>
            This command will remove you from tags list <br>
            Example: vtbot>del tag hololive <br>
            Role permission required: No permission needed
         </p>
         <p><code>vtbot&gt;Tag roles [roles] [Group/Member]</code><br>
            Same like <code>vtbot&gt;tag me</code> command,but this will tag roles (<em>every role separated by space</em>)<br>
            Example: <code>vtbot&gt;tag roles @kanochi-simps hanayori</code> or <code>vtbot&gt;tag roles @fbk-simps @hololive-simps fbk</code> <br>
            deprecated <br>
            <s>In other case if you want to remind before livestream started you can add <code>-setreminder {Minutes}</code><br>
            Example: <code>vtbot&gt;tag roles @kanochi-simps hanayori -setreminder 40</code></s><br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p>
         <p><code>vtbot&gt;Del roles [roles] [Group/Member]</code><br>
            Removes roles from being tagged. <br>
            Example: <code>vtbot&gt;Del roles @kanochi-simps kanochi</code> or <code>vtbot&gt;Del roles @hololive-simps fbk</code><br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p>
         <h4>Reminders</h4>
         <p><code>vtbot&gt;Set reminder</code><br>
            Update or Set your reminder time<br>
            Example: <code>vtbot&gt;Set reminder hololive 20</code> <br>
            if you want remove reminder time,you can set reminder to 0 <br>
            Example: <code>vtbot&gt;Set reminder hololive 0</code> <br>
            Role permission required: No permission needed
         </p>
         <p><code>vtbot&gt;Roles reminder [roles] [Group/Member] {Minutes}</code><br>
            Update or Set Roles reminder time<br>
            Example: <code>vtbot&gt;Roles reminder @kanochi-simps Hanayori 30</code><br>
            PS: if you want remove reminder time,you can set reminder to 0 <br>
            Example: <code>vtbot&gt;Roles reminder @kanochi-simps Hanayori 0</code><br>
            Role permission required: <b>Manage Channel or Higher</b> 
         </p>
         <hr>
         <h4>Utility command</h4>
         <p><code>vtbot&gt;My tags</code><br>
            Shows all your info
         </p>
         <p><code>art&gt;{Group/Member}</code><br>
            Show fanart with randomly with their fanart hashtag<br>
            Example: <code>art&gt;Kanochi</code> or <code>art&gt;hololive</code>
         </p>
         <p><code>lewd&gt;{Group/Member}</code><br>
            Show lewd fanart with randomly from pixivi r-18 or twitter with their lewd hashtag<br>
            Example: <code>lewd&gt;Kanochi</code> or <code>lewd&gt;hololive</code>
         </p>
         <p><code>vtbot&gt;Vtuber data [Group/Member] -region {region}</code><br>
            Shows available Vtuber data<br>
            Example: <code>vtbot&gt;vtuber data nijisanji</code> or <code>vtbot&gt;vtuber data nijisanji,hololive -region id</code>
         </p>
         <p><code>vtbot&gt;Info {Member}</code><br>
            Shows Vtuber Info(subscriber,followers,video count) data(twitter,Bilibili,Youtube)<br>
            Example: <code>vtbot&gt;Info sora</code>
         </p>
         <p><code>yt&gt;Upcoming [Vtuber Group/Member] -region {region}</code><br>
            This command will show Upcoming live streams on Youtube (<em>only 3 if use Vtuber Group</em>)<br>
            Example: <code>yt&gt;Upcoming nijisanji,hololive</code><br>
            with Region: <code>yt&gt;Upcoming nijisanji,hololive -region en</code>
         </p>
         <p><code>yt&gt;Live [Vtuber Group/Member] -region {region}</code><br>
            Show all livestreams right now on Youtube
         </p>
         <p><code>yt&gt;Past [Vtuber Group/Member] -region {region}</code><br>
            Show past streams on Youtube (<em>only 3 if use Vtuber Group</em>)
         </p>
         <p><s><code>bl&gt;Upcoming [Vtuber Group]</code><br>
            Show Upcoming live streams on Bilibili</s>
         </p>
         <p><code>bl&gt;Live [Vtuber Group/Member]</code><br>
            Same like <code>yt&gt;upcoming</code> but this for bilibili
         </p>
         <p><code>bl&gt;Past [Vtuber Group/Member]</code><br>
            Same like <code>yt&gt;past</code> but this for bilibili
         </p>
         <p><code>sp_bl&gt;[Group/Member]</code><br>
            Show latest video on Bilibili (<em>only 3 if use Vtuber Group</em>)
         </p>
  </div>
</template>
