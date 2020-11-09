# Create your views here.
from django.http import HttpResponse
from django.shortcuts import render
import requests,asyncio


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
    #return render(request,'exec.html')
    return HttpResponse("Still Dev~")

def go_simps_add(request):
    if request.method == "POST":
        print(request.POST)
        return HttpResponse("Still Dev~")
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
        Payload = "https://raw.githubusercontent.com/JustHumanz/Go-Simp/master/Img/404.jpg"    

    return render(request, 'support.html',{'Data': Payload})