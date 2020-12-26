# Create your views here.
from django.http import HttpResponse
from django.shortcuts import render,redirect
from backend.engine import *

git = GitGood(os.environ['GITKEY'])
LOGINURL = "https://discord.com/api/oauth2/authorize?client_id=719540207552167936&permissions=522304&redirect_uri=http%3A%2F%2Flocalhost%3A8000%2FDiscord%2Flanding&response_type=code&scope=bot%20guilds%20identify"
Discord = Discortttt()
Vtubers = GetVtubers()
Groups = Vtubers.GetGroups()
regex = r"^(?:(http|https):\/\/[a-zA-Z-]*\.{0,1}[a-zA-Z-]{3,}\.[a-z]{2,})\/channel\/([a-zA-Z0-9_]{3,})$"

def GetChannelID(url):
    matches = re.finditer(regex, url, re.MULTILINE)
    for matchNum, match in enumerate(matches, start=1):
        return match.group(2)

def go_simps_index(request):
    return render(request, 'index.html', {'Groups':Groups})

def go_simps_group(request, GroupID):
    Members = Vtubers.GetMemberGroups(GroupID)
    Region = Vtubers.GetRegList
    return render(request, 'group.html',{'Members':Members,'Region':Region,'Add':False})

def go_simps_members(request):
    Vtubers.ResizeImg("128")
    Members = Vtubers.Members
    RegList = Vtubers.GetRegList

    return render(request, 'group.html',{'Members':Members,'Region':RegList,'Add':True})

def go_simps_command(request):
    return render(request,'exec.html')

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
            Payload = {"State":"Error","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}
            return render(request,'done.html',Payload)
        else:
            POSTData["Youtube"] = ChannelID

        if issuenum is None:
            issuenum = git.PushNewIssues(POSTData,Title)
            Payload = {"State":"New","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}
            return render(request,'done.html',Payload)
        else:
            git.UpdateIssues(POSTData,issuenum,Title)
            Payload = {"State":"Duplicate","URL":"https://github.com/JustHumanz/Go-Simp/issues/"+str(issuenum)}
            return render(request,'done.html',Payload)

    else:
        Payload = {'Groups':Vtubers.GetGroups()}
        return render(request,'add.html',Payload)

def go_simps_member(request,MemberID):
    Vtubers.ResizeImg("512")
    Member,Subs = Vtubers.GetMemberSubs(MemberID)
    return render(request, 'member.html',{'Member': Member,'SubsInfo':Subs})


def go_simps_support(request,Type):
    Payload = ""
    if Type == "hug":
        Payload = 'https://i.ibb.co/rMt2Cqz/hug2.gif'
    elif Type == "airforceone":
        Payload = "https://img-comment-fun.9cache.com/media/a2NgKoZ/azaXgVx4_700w_0.jpg"
    else:
        Payload = "https://cdn.human-z.tech/404.jpg"    

    return render(request, 'support.html',{'Data': Payload})

def go_simps_guide(request):
    return render(request,'guide.html')   

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

    return render(request,'Pilot/guild.html',{'Guilds':Guilds})


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
    return render(request,'Pilot/channel.html',{'ChannelName':ChannelInfo,'Groups':Groups})