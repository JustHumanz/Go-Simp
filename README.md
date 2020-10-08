# Vtuber DiscordBot

![alt text](https://raw.githubusercontent.com/JustHumanz/Go-Simp/master/Img/go-simp.png "Go-Simp")  
##### [Original Source](https://twitter.com/any_star_/status/1288184424320790528)
![Inline docs](https://badgen.net/badge/Code%20style/Toxic-Asean/blue?icon=github) ![Inline docs](https://badgen.net/badge/Code%20quality/Better%20than%20yandere%20dev%20code/green?icon=github)
----

## Introduction
A simple VTuber bot to serve notification Fanart and Live stream  


### Current notification support


##### Group
| Group     | Youtube | Live.BiliBili | Space.BiliBili | Twitter | T.bilibili | 
|-----------|---------|---------------|----------------|---------|------------|
| Hanayori  | ✓       | ✓             | ✓              | ✓       |✓           |
| Hololive  | ✓       | ✓             | ✓              | ✓       |✓           |
| Nijisanji | ✓       | ✓             | ✓              | ✓       |✓           |
|Kizuna_AI  | ✓       | ✓             |✓                | ✓       |✖           |
|ENTUM  | ✓       | ✓             |✓               | ✓       |✓       |

##### Member

|   NICKNAME   |    FULL NAME     | REGION |   GROUP   |
|--------------|------------------|--------|-----------|
| Kanochi      | Kano             | JP     | Hanayori  |
| Parerun      | Hanamaru Hareru  | JP     | Hanayori  |
| Kohi         | Kohigashi Hitona | JP     | Hanayori  |
|~~Nononon~~      | ~~Nonomiya Nonono~~  | ~~JP~~     | ~~Hanayori~~  |
| Risu         | Ayunda Risu      | ID     | Hololive  |
| Moona        | Moona Hoshinova  | ID     | Hololive  |
| Lofi         | Airani Iofifteen | ID     | Hololive  |
| Buki         | Shirakami Fubuki | JP     | Hololive  |
| Koronen      | Inugami Korone   | JP     | Hololive  |
| Matsuri      | Natsuiro Matsuri | JP     | Hololive  |
| Aqua         | Minato Aqua      | JP     | Hololive  |
| Pekora       | Usada Pekora     | JP     | Hololive  |
| Senchou      | Houshou Marine   | JP     | Hololive  |
| Haato        | Akai Haato       | JP     | Hololive  |
| Sheep        | Tsunomaki Watame | JP     | Hololive  |
| Kaichou      | Kiryu Coco       | JP     | Hololive  |
| Choco-sensei | Yuzuki Choco     | JP     | Hololive  |
| Konlulu      | Suzuhara Lulu    | JP     | Nijisanji |
| Hanamaki     | Hana Macchia     | ID     | Nijisanji |
| Zea          | Zea Cornelia     | ID     | Nijisanji |
| Miyu         | Miyu Ottavia     | ID     | Nijisanji |
| Takamama     | Taka Radjiman    | ID     | Nijisanji |
| Jukut        | Riksa Dhirendra  | ID     | Nijisanji |
| Cia          | Amicia Michella  | ID     | Nijisanji |
| Rai          | Rai Galilei      | ID     | Nijisanji |
| Zura         | Azura Cecillia   | ID     | Nijisanji |
| Layla        | Layla Alstroemeria | ID   | Nijisanji |
| Nara         | Nara Haramaung   | ID     | Nijisanji |
| Bobon        | Bonnivier Pranaja |ID     | Nijisanji |
| Etna         | Etna Crimson     | ID     | Nijisanji |
| Siska        | Siska Leontyne   | ID     | Nijisanji |
| Alice        | Mononobe Alice   | JP     | Nijisanji |
| Noraneko     | Fumino Tamaki    | JP     | Nijisanji |
| Debiru       | Debidebi Debiru  | JP     | Nijisanji |
| Ai-chan      | Kizuna AI        | JP     | Kizuna_AI |
| hina-chan    | Nekomiya Hinata  | JP     | ENTUM     |
| Nana         | Kagura Nana      | JP     | Independen|
| Ui-mama      | Shigure Ui       | JP     | Independen|
| Kotone       | Tenjin Kotone    | JP     | Independen|
### Command
```vtbot>Enable [art/live/all] [Vtuber Group]```
This command will declare if [Vtuber Group] enable in this channel
Example: `vtbot>enable hanayori so other users can use vtbot>tag me kanochi`  

```vtbot>Update [art/live/all] [Vtuber Group]```
Use this command if you want to change enable state  

```vtbot>Disable [Vtuber Group]```
Just like enable but this disable command :3  

```art>[Group/Member nickname]```  
Show fanart with randomly with their fanart hashtag  
Example: `art>Kanochi or art>hololive`  

```vtbot>Tag me [Group/Member nickname]```  
This command will add you to tags list if any new fanart  
Example: ```art>tag me Kanochi```,then you will get tagged when there is a new fanart of kano  

```vtbot>Del tag [Group/Member nickname]```  
This command will remove you from tags list  

```vtbot>My tags```  
Shows all lists that you are subscribed  

```vtbot>Channel tags```  
Shows what is enable in this channel    

```vtbot>Vtuber data```  
Shows available Vtuber data  

```yt>Upcoming [Vtuber Group/Member]```  
This command will show all Upcoming live streams on Youtube  

```yt>Live [Vtuber Group/Member]```  
This command will show all live streams right now on Youtube  

```yt>Last [Vtuber Group/Member]```  
This command will show all past streams on Youtube [only 5]  

```yt>[Upcoming/Live/Last] [Member name]```  
This command will show all Vtuber member Upcoming/Live/Past streams on Youtube  

~~```bl>Upcoming [Vtuber Group]```  
This command will show all Upcoming live streams on BiliBili~~  

```bl>Live [Vtuber Group/Member]```  
This command will show all live streams right now on BiliBili  

```bl>Last [Vtuber Group/Member]```  
This command will show all past streams on BiliBili [only 5]  

```bl>[Upcoming/Live/Last] [Member name]```  
This command will show all Vtuber member Upcoming/Live/Past streams on BiliBili  

```sp_bl>[Group/Member name]```  
This command will show latest video on bilibili  

```vtbot>Help EN```  
Well,you using it right now  

```vtbot>Help JP```  
Like this but in Japanese  

### Notification&Command 
#### New Upcoming live stream
![alt text](https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/New%20Upcoming.png "New Upcoming live stream")   

#### Reminder  
```Reminder every 1 hours and 30 minutes before stream start```  
![alt text](https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/Reminder.png "Reminder")  

##### Upcoming command
![alt text](https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/Youtube%20Upcoming.png "Upcoming command")  

##### New Fanart
![alt text](https://raw.githubusercontent.com/JustHumanz/Go-simp/master/Img/New%20Fanart.png "New fanart")   

[Invite link](https://discord.com/oauth2/authorize?client_id=721964514018590802&permissions=449536&scope=bot)