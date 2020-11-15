# Create your views here.
from django.http import HttpResponse
from django.shortcuts import render,redirect
from github import Github
import requests,os,json,re


GroupURL = "https://api.human-z.tech/vtbot/group"
MemberURL = "https://api.human-z.tech/vtbot/member/"
SubscriberURL = "https://api.human-z.tech/vtbot/subscriber/"

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
        self.Token = Token
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
        Payload = Payload.copy()
        del Payload["csrfmiddlewaretoken"]
        del Payload["Group"]
        Payloadstr = json.dumps(Payload,indent = 1,ensure_ascii=False)
        label = self.repo.get_label("enhancement")
        issue = self.repo.create_issue(title=Title, body=str(Payloadstr),labels=[label],assignee="JustHumanz")
        return issue.number

    def UpdateIssues(self,Payload,Number,Title):
        Payload = Payload.copy()
        del Payload["csrfmiddlewaretoken"]
        del Payload["Group"]
        Payloadstr = json.dumps(Payload,indent = 1,ensure_ascii=False)
        issue = self.repo.get_issue(Number)
        issue.edit(title=Title,body=str(Payloadstr),assignee="JustHumanz")




git = GitGood(os.environ['GITKEY'])

def GetChannelID(url):
    regex = r"^(?:(http|https):\/\/[a-zA-Z-]*\.{0,1}[a-zA-Z-]{3,}\.[a-z]{2,})\/channel\/([a-zA-Z0-9_]{3,})$"
    matches = re.finditer(regex, url, re.MULTILINE)
    for matchNum, match in enumerate(matches, start=1):
        print ("Match {matchNum} was found at {start}-{end}: {match}".format(matchNum = matchNum, start = match.start(), end = match.end(), match = match.group()))
        for groupNum in range(0, len(match.groups())):
            groupNum = groupNum + 1
            return match.group(groupNum)
            #print ("Group {groupNum} found at {start}-{end}: {group}".format(groupNum = groupNum, start = match.start(groupNum), end = match.end(groupNum), group = match.group(groupNum)))


def go_simps_index(request):
    Vtubers = GetVtubers("")
    Payload = {'Groups':Vtubers.GetGroups()}
    return render(request, 'index.html', Payload)

def go_simps_group(request, GroupName):
    Vtubers = GetVtubers(GroupName)
    Members = Vtubers.GetMembers()
    Region = Vtubers.GetRegList

    Payload = {'Members':Members,'Region':Region,'Add':False}
    return render(request, 'group.html',Payload)

def go_simps_members(request):
    Vtubers = GetVtubers("")
    Vtubers.GetMembers()

    Payload = {'Members':Vtubers.ResizeImg("s100"),'Region':Vtubers.GetRegList(),'Add':True}
    return render(request, 'group.html',Payload)

def go_simps_command(request):
    return render(request,'exec.html')
    #return HttpResponse("Still Dev~")

def go_simps_add(request):
    if request.method == "POST":
        if request.POST["Nickname"] == "" or request.POST["Region"] == "" :
            return HttpResponse("????")
        elif request.POST["Youtube"] == "" and request.POST["BiliBili"] == "":
            return HttpResponse("????")
        Title = "Add "+request.POST["Nickname"]+" from "+" "+request.POST["Group"]    
        issuenum = git.CheckIssues(Title)
        ChannelID = GetChannelID(request.POST["Youtube"])
        if ChannelID is None:
            Payload = {"State":"Error","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}

        if issuenum is None:
            #issuenum = git.PushNewIssues(request.POST,Title)
            Payload = {"State":"New","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}
            #return redirect("https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum))
            return render(request,'done.html',Payload)
        else:
            #git.UpdateIssues(request.POST,issuenum,Title)
            Payload = {"State":"Duplicate","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}
            return render(request,'done.html',Payload)
            #return redirect("https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum))


    else:
        Vtubers = GetVtubers("")
        Payload = {'Groups':Vtubers.GetGroups()}
        return render(request,'add.html',Payload)

def go_simps_member(request,MemberName):
    Vtubers = GetVtubers(MemberName)
    Vtubers.GetMembers()

    Payload = {'Member': Vtubers.ResizeImg("s300"),'Subs': Vtubers.GetSubs()}
    return render(request, 'member.html',Payload)


def go_simps_support(request,Type):
    Payload = ""
    if Type == "hug":
        Payload = 'https://i.ibb.co/rMt2Cqz/hug2.gif'
    elif Type == "airforceone":
        Payload = "https://img-comment-fun.9cache.com/media/a2NgKoZ/azaXgVx4_700w_0.jpg"
    else:
        Payload = "https://cdn.human-z.tech/404.jpg"    

    return render(request, 'support.html',{'Data': Payload})