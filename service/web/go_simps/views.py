# Create your views here.
from django.http import HttpResponse
from django.shortcuts import render,redirect
from backend.engine import *

git = GitGood(os.environ['GITKEY'])
DOMAIN = os.environ['DOMAIN']

Vtubers = GetVtubers()
Vtubers.ResizeImg("512")
Groups = Vtubers.GetGroups()
regex = r"^(?:(http|https):\/\/[a-zA-Z-]*\.{0,1}[a-zA-Z-]{3,}\.[a-z]{2,})\/channel\/([a-zA-Z0-9_]{3,})$"

def GetChannelID(url):
    matches = re.finditer(regex, url, re.MULTILINE)
    for _, match in enumerate(matches, start=1):
        return match.group(2)

def go_simps_index(request):
    return render(request, 'index.html', {'Groups':Groups,'Domain':DOMAIN})

def go_simps_group(request, GroupID):
    Members = Vtubers.GetMemberGroups(GroupID)
    return render(request, 'group.html',{'Members':Members,'Region':GetRegList(Members),'Add':False,'Domain':DOMAIN})

def go_simps_members(request):
    Members = Vtubers.Members
    return render(request, 'group.html',{'Members':Members,'Region':GetRegList(Members),'Add':True,'Domain':DOMAIN})

def go_simps_command(request):
    return render(request,'exec.html',{'Domain':DOMAIN})

def go_simps_add(request):
    if request.method == "POST":
        if request.POST["Nickname"] == "" or request.POST["Region"] == "" :
            return HttpResponse("????")
        elif request.POST["Youtube"] == "" and request.POST["BiliBili"] == "":
            return HttpResponse("????")
        Title = "Add "+request.POST["Nickname"]+" from "+" "+request.POST["Group"]    
        issuenum = git.CheckIssues(Title)
        ChannelID = GetChannelID(request.POST["Youtube"])
        POSTData = request.POST.copy()

        if ChannelID is None:
            Payload = {"State":"Error","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum),'Domain':DOMAIN}
            return render(request,'done.html',Payload)
        else:
            POSTData["Youtube"] = ChannelID

        if issuenum is None:
            issuenum = git.PushNewIssues(POSTData,Title)
            Payload = {"State":"New","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum),'Domain':DOMAIN}
            return render(request,'done.html',Payload)
        else:
            git.UpdateIssues(POSTData,issuenum,Title)
            Payload = {"State":"Duplicate","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum),'Domain':DOMAIN}
            return render(request,'done.html',Payload)

    else:
        Payload = {'Groups':Vtubers.GetGroups(),'Domain':DOMAIN}
        return render(request,'add.html',Payload)

def go_simps_member(request,MemberID):
    Member = Vtubers.GetMemberInfo(MemberID)
    return render(request, 'member.html',{'Member': Member,'Domain':DOMAIN})


def go_simps_support(request,Type):
    Payload = ""
    if Type == "hug":
        Payload = 'https://i.ibb.co/rMt2Cqz/hug2.gif'
    elif Type == "airforceone":
        Payload = "https://img-comment-fun.9cache.com/media/a2NgKoZ/azaXgVx4_700w_0.jpg"
    else:
        Payload = f"https://cdn.{DOMAIN}/404.jpg"    

    return render(request, 'support.html',{'Data': Payload,'Domain':DOMAIN})

def go_simps_guide(request):
    return render(request,'guide.html',{'Domain':DOMAIN})   

"""
def go_simps_discord_login(request):
    return redirect(LOGINURL)

def go_simps_discord_landing(request):
    code = request.GET["code"]
    access_token = Discord.GetAccessToken(code)

    response = HttpResponse("Cookie Set")  
    response.set_cookie('oauth2', access_token['access_token'],max_age = access_token['expires_in'])      
    response = redirect('/Discord/cp')
    return response

        
def go_simps_discord_cp(request):
    cokkie = ""
    try:
        cokkie = request.COOKIES['oauth2']
    except:
        return redirect('/Discord/login')
    Guilds = Discord.GetUserGuild(cokkie)    

    for i in range(len(Guilds)):
        Guilds[i]['Channels'] = Discord.GetChannels(Guilds[i]['id'])

    return render(request,'Pilot/guild.html',{'Guilds':Guilds,'Domain':DOMAIN})


def go_simps_discord_channel(request,ChannelID):
    ChannelInfo,ChannelSupport = Discord.GetChannelInfo(ChannelID)
    Groups = Vtubers.GetGroups()
    Roles = Discord.GetGuildRols()
    try:
        for i in range(len(Groups)):
            del Groups[i]["VtuberGroupIcon"]
            for Data in ChannelSupport["ChannelData"]:
                if Data["GroupName"] == Groups[i]["VtuberGroupName"]:
                    Groups[i]["Enable"] = True
                    Groups[i]["ChannelData"] = Data
                else:
                    Groups[i]["Enable"] = False                
    except:
        print("Not enable any groups")
    return render(request,'Pilot/channel.html',{'ChannelName':ChannelInfo,'Groups':Groups,'Domain':DOMAIN})
"""    