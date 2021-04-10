from github import Github
import requests,os,json,re

GroupsURL = "http://"+os.environ['REST_API']+":2525/groups/"
MembersURL = "http://"+os.environ['REST_API']+":2525/members/?groupid="
Youtube = "http://"+os.environ['REST_API']+":2525/livestream/youtube/member/"
API_ENDPOINT = 'https://discord.com/api/v6'

class GetVtubers:
    def __init__(self):        
        Members = []
        Groups = []
        response = requests.get(GroupsURL)
        for Data in response.json():
            r = requests.get(MembersURL+str(Data["ID"]))
            for Member in r.json():
                Members.append(Member)

            Data["Members"] = r.json()    
            Groups.append(Data)                
        
        self.Groups = Groups
        self.Members = Members

    def GetMemberInfo(self,ID):
        for Member in self.Members:
            if int(ID) == int(Member["ID"]):
                return Member

    def GetGroups(self):
        return self.Groups
           

    def GetMemberGroups(self,GroupID):
        GroupMember = []
        for Group in self.Groups:
            if int(Group["ID"]) == int(GroupID):
                for Member in Group["Members"]:
                    if MemberCheckLive(Member["ID"]) is not None:
                        Member["YtLive"] = True
                    else:
                        Member["YtLive"] = False    
                """
                MemberData["YtLive"] = False
                if LiveInfo is not None:
                    for Live in LiveInfo:
                        if int(Live["MemberID"]) == int(MemberData["ID"]):
                            MemberData["YtLive"] = True
                """
                GroupMember = Group["Members"]            
                #GroupMember.append(MemberData)

        return GroupMember

    def ResizeImg(self,size):
        for i in range(len(self.Members)):
            if self.Members[i]["Youtube"] is not None:
                self.Members[i]["Youtube"]["Avatar"] = self.Members[i]["Youtube"]["Avatar"].replace("s800","s"+size)    

def GetRegList(Members):
    Region = []
    for Member in Members:
        if Member['Region'] not in Region:
            Region.append(Member['Region'])
    return Region    

def MemberCheckLive(MemberID):
    response = requests.get(Youtube+str(MemberID)+"/Live")
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

"""
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

"""        