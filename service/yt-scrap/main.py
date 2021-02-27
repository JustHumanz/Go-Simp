from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.firefox.options import Options
import time,re,urllib.request,json,itertools,os,isodate
from db import mydb,dbconn
from datetime import datetime

options = Options()
options.add_argument('--headless')
options.add_argument('--hide-scrollbars')
options.add_argument('--disable-gpu')
api_key = os.environ['YtToken']


class Filter:
    def __init__(self, VideoIdList, VtuberName,VtuberMember_id):
        self.VideoIdList = VideoIdList
        self.VtuberName = VtuberName
        self.VtuberMember_id = VtuberMember_id

    def FilterYt(self):
        url = "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails,contentDetails&fields=items(id,snippet(publishedAt,title,description,channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount),contentDetails(duration))&id=" + ",".join(self.VideoIdList) + "&key=" + api_key
        print(url)
        inp = urllib.request.urlopen(url)
        if inp.getcode() != 200:
            print("can't curl :( ")
            
        resp = json.loads(inp.read())
        tmp = []
        for ytdata in resp['items']:
            if re.search(r'(?m)(cover|song|feat|music|mv|covered|op)', (ytdata['snippet']['title']).lower()):
                YtType = "Covering"
            else:
                YtType = "Streaming"

            print("Video type "+YtType)
            if (ytdata['snippet']['liveBroadcastContent']).lower() == "none":
                ytdata['snippet']['liveBroadcastContent'] = "past"
            ytdata.update({'type': YtType})
            tmp.append(ytdata)

        return tmp

    def InputYt(self,YtData):
        for data in YtData:
            Thumb = "http://i3.ytimg.com/vi/"+data["id"]+"/hqdefault.jpg"
            
            if "scheduledStartTime" in  data['liveStreamingDetails']:
                scheduledStart = changedatetime(data['liveStreamingDetails']['scheduledStartTime'])
            else:
                scheduledStart = "0000-00-00 00:00:00"

            if "actualEndTime" in data['liveStreamingDetails']:
                End = changedatetime(data['liveStreamingDetails']['actualEndTime'])
            else:
                End = "0000-00-00 00:00:00"

            sql = "INSERT INTO Youtube (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,Length,VtuberMember_id) values(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)"
            dat = data["id"],data['type'],data['snippet']['liveBroadcastContent'], \
                data['snippet']['title'],Thumb, data['snippet']['description'], \
                changedatetime(data['snippet']['publishedAt']), scheduledStart, End, data['statistics']['viewCount'],isodate.parse_duration(data['contentDetails']['duration']),self.VtuberMember_id
            dbconn.execute(sql,dat)
            mydb.commit()


def scrapper(Channelid):
    driver = webdriver.Firefox(firefox_options = options)
    url = "https://www.youtube.com/channel/"+Channelid+"/videos"
    driver.get(url)

    #i made this to 70 second because some crazy rabbit and fox made too many video
    for _ in range(70):
        driver.execute_script("return document.body.scrollHeight")
        time.sleep(1)
        driver.find_element_by_tag_name('a').send_keys(Keys.END)

    links = driver.find_elements_by_xpath('//a[@href and @id="thumbnail"]')
    ChannelIDList = []
    for link in links:
        fix = link.get_attribute('href').replace("https://www.youtube.com/watch?v=","")
        if get_video_status(fix):
            ChannelIDList.append(fix)

    driver.quit()
    return list(dict.fromkeys(ChannelIDList))

def get_video_status(VideoID):
    dbconn.execute("SELECT id FROM Youtube where VideoID = %s",(VideoID,) )
    dbconn.fetchall()
    if not dbconn.rowcount:
        print("New Video "+ VideoID)
        return True
    else:
        print("Video already in database "+ VideoID)
        return False

def changedatetime(date):
    return datetime.fromisoformat(date.replace('Z', '+00:00'))


def main():
    print("Get VtuberName_EN,Youtube_ID,id")
    dbconn.execute("SELECT VtuberName_EN,Youtube_ID,id FROM Vtuber.VtuberMember order by id DESC")
    res = dbconn.fetchall()
    for i in res:
        VideoIdList = scrapper(i[1])
        if len(VideoIdList) >= 50:
            n = 50
        else:
            n = 25
        
        print("Done get Video List,total ",len(VideoIdList))
        final = [VideoIdList[i * n:(i + 1) * n] for i in range((len(VideoIdList) + n - 1) // n )]
        for k in range(len(final)):
            #filtervideo(final[k],i[0],i[2])
            Yt = Filter(final[k],i[0],i[2])
            Yt.InputYt(Yt.FilterYt()) 
    os._exit(0)
main()




"""
def filtervideo(VideoIdList,VtuberName,VtuberMember_id):
    try:
        url = "https://www.googleapis.com/youtube/v3/videos?part=statistics,snippet,liveStreamingDetails&fields=items(snippet(publishedAt,title,description,channelTitle,liveBroadcastContent),liveStreamingDetails(scheduledStartTime,concurrentViewers,actualEndTime),statistics(viewCount))&id=" + ",".join(VideoIdList) + "&key=" + api_key
        print(url)
        inp = urllib.request.urlopen(url)
        if inp.getcode() != 200:
            print("can't curl :( ")
    except:
        print("some error")

    resp = json.loads(inp.read())
    for i,j in zip(resp['items'],VideoIdList):
        match = re.search(r'(?m)(cover|song|feat|music|mv)', (i['snippet']['title']).lower())
        if match:
            YtType = "Covering"
        else:
            YtType = "Streaming"

        print("Video type "+YtType)
        if (i['snippet']['liveBroadcastContent']).lower() == "none":
            i['snippet']['liveBroadcastContent'] = "past"
        i.update({'type': YtType})

        print("Input data to database ",VtuberName,VtuberMember_id,j,i['snippet']['title'])
        inputvideo(i,j,VtuberMember_id)
        time.sleep(1)

def inputvideo(Data,VideoID,VtuberMember_id):
    Thumb = "http://i3.ytimg.com/vi/"+VideoID+"/hqdefault.jpg"
    try:
        published = changedatetime(Data['snippet']['publishedAt'])
        scheduledStart = changedatetime(Data['liveStreamingDetails']['scheduledStartTime'])
        End = changedatetime(Data['liveStreamingDetails']['actualEndTime'])
    except:
        published = "0000-00-00 00:00:00"
        scheduledStart = "0000-00-00 00:00:00"
        End = "0000-00-00 00:00:00"

    sql = "INSERT INTO Youtube (VideoID,Type,Status,Title,Thumbnails,Description,PublishedAt,ScheduledStart,EndStream,Viewers,VtuberMember_id) values(%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s)"
    dat = VideoID,Data['type'],Data['snippet']['liveBroadcastContent'], \
        Data['snippet']['title'],Thumb, Data['snippet']['description'], \
        published, scheduledStart, End, Data['statistics']['viewCount'],VtuberMember_id
    dbconn.execute(sql,dat)
    mydb.commit()
    """