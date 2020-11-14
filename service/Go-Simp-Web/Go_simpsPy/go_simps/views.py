# Create your views here.
from django.http import HttpResponse
from django.shortcuts import render,redirect
from github import Github
import requests,os,json


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
DiscordLink = os.environ['GUILD']
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
        if issuenum is None:
            issuenum = git.PushNewIssues(request.POST,Title)
            Payload = {"State":"New","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum),"Guild":DiscordLink}
            #return redirect("https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum))
            return render(request,'done.html',Payload)
        else:
            git.UpdateIssues(request.POST,issuenum,Title)
            Payload = {"State":"Duplicate","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum),"Guild":DiscordLink}
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