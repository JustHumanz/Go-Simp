# Create your views here.
from django.shortcuts import render
import requests,asyncio



class GetVtubers:
    def GetGroups(self):
        response = requests.get('https://api.human-z.tech/vtbot/group')
        return response.json()

    def GetMembers(self,GroupName):
        response = requests.get('https://api.human-z.tech/vtbot/member/'+GroupName)        
        return response.json()  

    def GetSubs(self,MemberName):
        SubsInfo = requests.get('https://api.human-z.tech/vtbot/subscriber/'+MemberName)    
        return SubsInfo.json()

Vtubers = GetVtubers()

def go_simps_index(request):
    Payload = {'Groups':Vtubers.GetGroups()}
    return render(request, 'index.html', Payload)

def go_simps_group(request, GroupName):
    Members = Vtubers.GetMembers(GroupName)
    Region = []
    i = 0
    for Member in Members:
        if Member['Region'] not in Region:
            Region.append(Member['Region'])
        Members[i]["Youtube_Avatar"] = Members[i]["Youtube_Avatar"].replace("s800","s100")    

    Payload = {'Members':Members,'Region':Region}
    return render(request, 'group.html',Payload)

def go_simps_member(request,MemberName):
    Member = Vtubers.GetMembers(MemberName)
    Member[0]["Youtube_Avatar"] = Member[0]["Youtube_Avatar"].replace("s800","s200")    
    Subs = Vtubers.GetSubs(MemberName)

    Payload = {'Member': Member,'Subs': Subs}
    return render(request, 'member.html',Payload)
