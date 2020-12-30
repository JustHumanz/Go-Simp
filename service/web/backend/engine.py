from github import Github
import requests,os,json,re

GroupURL = "http://rest_api:2525/Groups/"
MemberURL = "http://rest_api:2525/Members/"
SubscriberURL = "http://rest_api:2525/Subscribe/"
ChannelURL = "http://rest_api:2525/channel/"
Youtube = "http://rest_api:2525/Youtube/"
API_ENDPOINT = 'https://discord.com/api/v6'

class GetVtubers:
    def __init__(self):
        response = requests.get(MemberURL)
        self.Members = response.json()

    def GetGroups(self):
        return requests.get(GroupURL).json()
        
    def GetMemberSubs(self,ID):
        for Member in self.Members:
            if int(Member["ID"]) == int(ID):
                SubInfo = requests.get(SubscriberURL+"Member/"+str(Member["ID"]))    
                return Member,SubInfo.json()[0]
        return None        

    def GetMemberGroups(self,GroupID):
        GroupMember = []
        LiveInfo = CheckLive(GroupID)
        for MemberData in self.Members:
            if int(MemberData["GroupID"]) == int(GroupID):
                MemberData["YtLive"] = False
                if LiveInfo is not None:
                    for Live in LiveInfo:
                        if int(Live["MemberID"]) == int(MemberData["ID"]):
                            MemberData["YtLive"] = True
                GroupMember.append(MemberData)

        return GroupMember

    def ResizeImg(self,size):
        for i in range(len(self.Members)):
            self.Members[i]["YoutubeAvatar"] = self.Members[i]["YoutubeAvatar"].replace("s800","s"+size)    

def GetRegList(Members):
    Region = []
    for Member in Members:
        if Member['Region'] not in Region:
            Region.append(Member['Region'])
    return Region    

def CheckLive(GroupID):
    response = requests.get(Youtube+"Group/"+GroupID+"/Live")
    if response.ok:
        return response.json()
    else:
        return None    




class GitGood:
    def __init__(self, Token):    
        self.g = Github(Token)
        self.repo = self.g.get_repo("JustHumanz/Go-Simp")

    def CheckIssues(self,title):
        IssueState = {"open","closed"}
        for state in IssueState:
            open_issues = self.repo.get_issues(state=state)
            for issue in open_issues:
                if issue.title == title:
                    return issue.number

    def PushNewIssues(self,Payload,Title):
        del Payload["csrfmiddlewaretoken"]
        del Payload["Group"]
        PayloadStr = json.dumps(Payload,indent = 1,ensure_ascii=False)
        label = self.repo.get_label("enhancement")
        issue = self.repo.create_issue(title=Title, body=PayloadStr,labels=[label],assignee="JustHumanz")
        return issue.number

    def UpdateIssues(self,Payload,Number,Title):
        del Payload["csrfmiddlewaretoken"]
        del Payload["Group"]
        PayloadStr = json.dumps(Payload,indent = 1,ensure_ascii=False)
        issue = self.repo.get_issue(Number)
        issue.edit(title=Title,body=PayloadStr,assignee="JustHumanz")

class Discortttt:
    def __init__(self):
        self.CLIENT_ID = os.environ['CLIENT_ID']
        self.CLIENT_SECRET = os.environ['CLIENT_SECRET']
        self.URLI = 'http://localhost:8000/Discord/landing'
        self.DisocrdBot = os.environ["DISCORD_BOT"]

    def GetAccessToken(self,code):    
        r = requests.post('%s/oauth2/token' % API_ENDPOINT, data={
            'client_id': self.CLIENT_ID,
            'client_secret': self.CLIENT_SECRET,
            'grant_type': 'authorization_code',
            'code': code,
            'redirect_uri': self.URLI,
            'scope': 'identify guild'
        }, headers={
            'Content-Type': 'application/x-www-form-urlencoded'
        })
        r.raise_for_status()
        return r.json()        

    def GetUserGuild(self,token):
        ResultUser = requests.get(API_ENDPOINT+"/users/@me/guilds",headers={
            'Authorization': 'Bearer '+token
        })

        ResultBot = requests.get(API_ENDPOINT+"/users/@me/guilds",headers={
            'Authorization': 'Bot '+self.DisocrdBot
        })
        GuildList = []
        for UserGuild in ResultUser.json():
            for BotGuild in ResultBot.json():
                if UserGuild["id"] == BotGuild["id"]:
                    GuildList.append(UserGuild)               
        return GuildList

    def GetChannels(self,GuildID):
        self.GuildID = GuildID
        Result = requests.get(API_ENDPOINT+"/guilds/%s/channels" % GuildID,headers={
            'Authorization': 'Bot '+self.DisocrdBot
        })
        Channels = []
        for Channel in Result.json():
            if Channel['type'] == 0:
                Channels.append(Channel)
        return Channels

    def GetChannelInfo(self,ChannelID):        
        Result = requests.get(API_ENDPOINT+"/channels/%s" % ChannelID,headers={
            'Authorization': 'Bot '+self.DisocrdBot
        })

        return Result.json(),requests.get(ChannelURL+ChannelID).json()

    def GetGuildRols(self):
        Result = requests.get(API_ENDPOINT+"/guilds/%s/roles" % self.GuildID,headers={
            'Authorization': 'Bot '+self.DisocrdBot
        })
        return Result.json()