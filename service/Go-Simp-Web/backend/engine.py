from github import Github
import requests,os,json,re

GroupURL = "http://rest_api:2525/group"
MemberURL = "http://rest_api:2525/member/"
SubscriberURL = "http://rest_api:2525/subscriber/"
API_ENDPOINT = 'https://discord.com/api/v6'

class GetVtubers:
    def __init__(self, InputData):
        self.InputData = InputData
        self.Members = ""

    def GetGroups(self):
        response = requests.get(GroupURL)
        return response.json()

    def GetMembers(self):
        response = requests.get(MemberURL+self.InputData)        
        self.Members = response.json()  
        return response.json()  

    def GetSubs(self):
        SubsInfo = requests.get(SubscriberURL+self.InputData)    
        return SubsInfo.json()

    def GetRegList(self):
        Region = []

        for Member in self.Members:
            if Member['Region'] not in Region:
                Region.append(Member['Region'])
        return Region    

    def ResizeImg(self,size):
        Members = self.Members

        for i in range(len(Members)):
            Members[i]["Youtube_Avatar"] = Members[i]["Youtube_Avatar"].replace("s800",size)    

        return Members   

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

    def GetAccessToken(self,code):    
        data = {
            'client_id': self.CLIENT_ID,
            'client_secret': self.CLIENT_SECRET,
            'grant_type': 'authorization_code',
            'code': code,
            'redirect_uri': self.URLI,
            'scope': 'identify guild'
        }
        headers = {
            'Content-Type': 'application/x-www-form-urlencoded'
        }
        r = requests.post('%s/oauth2/token' % API_ENDPOINT, data=data, headers=headers)
        r.raise_for_status()
        return r.json()        
