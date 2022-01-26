import os,time,requests,json,logging
from flask import Flask, g, session, redirect, request, url_for, jsonify,render_template
from requests_oauthlib import OAuth2Session
from db import *

OAUTH2_CLIENT_ID = "719540207552167936"
OAUTH2_CLIENT_SECRET =  os.environ['SECRET']
OAUTH2_REDIRECT_URI =  os.environ['CALLBACK_URL']
URL = os.environ['URL']
bot_headers = {"Authorization": "Bot "+ os.environ['BOT'],"Content-Type":"application/json"}

API_BASE_URL = os.environ.get('API_BASE_URL', 'https://discordapp.com/api')
AUTHORIZATION_BASE_URL = API_BASE_URL + '/oauth2/authorize'
TOKEN_URL = API_BASE_URL + '/oauth2/token'

app = Flask(__name__,template_folder='template')
app.config['SECRET_KEY'] = OAUTH2_CLIENT_SECRET

logging.basicConfig(level=logging.INFO,format='%(asctime)s - %(message)s', datefmt='%d-%b-%y %H:%M:%S')
groups = get_groups()

class User():
    def __init__(self):
        self.user_id = ''
        self.guild_id = ''
        self.channel_id = ''

    def add_user(self,user_id):
        self.user_id = user_id    

    def add_guild(self,guild_id):
        self.guild_id = guild_id

    def add_channel(self,channel):
        self.channel_id = channel

    def verify_user(self,discord):
        user = discord.get(API_BASE_URL + '/users/@me')
        if user.status_code != 200:
            return user

        user = user.json()
        if user['id'] == self.user_id:
            return True

discord_user = User()

if 'http://' in OAUTH2_REDIRECT_URI:
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = 'true'
        
def token_updater(token):
    session['oauth2_token'] = token

def prepare_response(res_object, status_code):
    response = jsonify(res_object)
    response.headers.add('Access-Control-Allow-Origin', URL)    
    response.headers.add('Access-Control-Allow-Credentials', 'true')
    return response, status_code
    
def channel_type(num):
    if num == 1:
        return "Fanart"    
    elif num == 2:
        return "Livestream"
    elif num == 3:
        return "Fanart & Livestream"
    elif num == 69:
        return  "Lewd"
    else:
        return "Fanart & Lewd"          

def channel_type_rev(num):
    if num == "Fanart":
        return 1    
    elif num == "Livestream":
        return 2
    elif num == "Fanart & Livestream":
        return 3
    elif num == "Lewd":
        return 69 
    else:
        return 70

def make_session(token=None, state=None, scope=None):
    return OAuth2Session(
        client_id=OAUTH2_CLIENT_ID,
        token=token,
        state=state,
        scope=scope,
        redirect_uri=OAUTH2_REDIRECT_URI,
        auto_refresh_kwargs={
            'client_id': OAUTH2_CLIENT_ID,
            'client_secret': OAUTH2_CLIENT_SECRET,
        },
        auto_refresh_url=TOKEN_URL,
        token_updater=token_updater)

@app.route('/')
def index():
    scope = request.args.get(
        'scope',
        'identify guilds guilds.join guilds.members.read')
    discord = make_session(scope=scope.split(' '))
    authorization_url, state = discord.authorization_url(AUTHORIZATION_BASE_URL)
    session['oauth2_state'] = state
    return redirect(authorization_url)


@app.route('/callback')
def callback():
    if request.values.get('error'):
        return request.values['error']
    discord = make_session(state=session.get('oauth2_state'))
    token = discord.fetch_token(
        TOKEN_URL,
        client_secret=OAUTH2_CLIENT_SECRET,
        authorization_response=request.url)
    session['oauth2_token'] = token
    return redirect(URL+"/#/guilds")


@app.route('/guilds')
def guilds():
    discord = make_session(token=session.get('oauth2_token'))
    guilds = []
    discord_guilds = discord.get(API_BASE_URL + '/users/@me/guilds')
    if discord_guilds.status_code != 200:
        return prepare_response({"message":discord_guilds.json(),"error":True},discord_guilds.status_code)
        
    for i in discord_guilds.json():
        if int(i["permissions"]) & 8208:
            isVtbot = requests.get(API_BASE_URL+f'/guilds/{i["id"]}',headers=bot_headers)
            if isVtbot.status_code == 200:
                guilds.append(i)            

    return prepare_response({'guilds': guilds},200)

@app.route('/guilds/<guild_id>/channels')
def channels(guild_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return prepare_response({"message":"Missing token"},301)

    user = discord.get(API_BASE_URL + '/users/@me')
    if user.status_code != 200:
        return prepare_response({"message":user.json(),"error":True},user.status_code)

    isAdmin = False
    #Check role permissions
    user_guild = discord.get(API_BASE_URL+ f'/users/@me/guilds/{guild_id}/member')    
    if user_guild.status_code != 200:
        return prepare_response({"message":user_guild.json(),"error":True},user_guild.status_code)

    else:
        user_role = user_guild.json()['roles']
        role = requests.get(API_BASE_URL + f'/guilds/{guild_id}/roles',headers=bot_headers)
        if role.status_code != 200:
            return prepare_response({"message":role.json(),"error":True},role.status_code)
        else:
            for ii in role.json():
                if ii['id'] in user_role:
                    if int(ii["permissions"]) & 8208:
                        isAdmin = True
                        break    
    if isAdmin == False:
        return prepare_response({"message":"Your don't have access to manage guilds channel"}, 301)

    discord_user.add_guild(guild_id)
    discord_user.add_user(user.json()['id'])

    channels = requests.get(API_BASE_URL + f'/guilds/{guild_id}/channels',headers=bot_headers).json()
    channels_info = {
        "guild_channel": [],
    }

    for i in channels:
        if i["parent_id"] is None:
            continue
        channels_info["guild_channel"].append(i)

    return prepare_response(channels_info,200)
    

@app.route('/channel/<channel_id>/agency',methods=['GET'])
def channel(channel_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return "Missing token",301
    
    discord_channel = requests.get(API_BASE_URL+f'/channels/{channel_id}',headers=bot_headers)
    if discord_channel.status_code != 200:
        return prepare_response({"message":discord_channel.json(),"error":True},discord_channel.status_code)

    discord_channel = discord_channel.json()
    discord_user.add_channel(channel_id)

    agency_list = {
        "agency_list": [],
        "channel_id": channel_id,
        "channel_name": discord_channel["name"],
        "is_nsfw": discord_channel["nsfw"]
    }
    for k in groups:
        agency_id = k["id"]
        channel_data = get_channel_info(agency_id,channel_id)
        if channel_data == []:
            agency_list['agency_list'].append({
                "agency_name": k["VtuberGroupName"],
                "agency_id": agency_id,
                "agency_icon": k["VtuberGroupIcon"],
                "agency_region": k["Region"],
                "channel_region": [],
                "channel_id": channel_id,
                "dynamic": False,
                "indie_notif": False,
                "lite": False,
                "upcoming": False,
                "type" : None})

        else:
            agency_list['agency_list'].append({
                "agency_name": k["VtuberGroupName"],
                "agency_id": agency_id,                
                "agency_icon": k["VtuberGroupIcon"],
                "channel_id": channel_data[0]["DiscordChannelID"],
                "dynamic": bool(channel_data[0]["Dynamic"]),
                "indie_notif": bool(channel_data[0]["IndieNotif"]),
                "lite": bool(channel_data[0]["Lite"]),
                "upcoming": bool(channel_data[0]["NewUpcoming"]),
                "agency_region": k["Region"],
                "channel_region": channel_data[0]["Region"].split(","),                
                "type" : channel_type(channel_data[0]["Type"]),})

    return prepare_response(agency_list,200)



@app.route('/channel/<channel_id>/update',methods=['POST','GET','OPTIONS'])
def update_channel(channel_id):
    if request.method == "OPTIONS":
        response = jsonify({"message":"OK"})
        response.headers.add('Access-Control-Allow-Headers', 'content-type')  
        response.headers.add('Access-Control-Allow-Methods', 'POST')
        response.headers.add('Access-Control-Allow-Origin', URL)    
        response.headers.add('Access-Control-Allow-Credentials', 'true')        
        return response

    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return prepare_response({"message":"Missing token"},401)

    logging.info('%s %s %s',"update channel",discord_user.verify_user(discord),channel_id,discord_user.channel_id)
    if discord_user.verify_user(discord) and channel_id == discord_user.channel_id:
        form = request.json
        if form["agency_id"] != 10:
            form["indie_notif"] = False

        if form["channel_region"] == []:
            form["channel_region"] = form["agency_region"]

        data = {
            "channel_type": channel_type_rev(form["type"]),
            "upcoming": bool(form["upcoming"]),
            "dynamic": bool(form["dynamic"]),
            "region": str(','.join(form["channel_region"])),
            "lite": bool(form["lite"]),
            "indie": bool(form["indie_notif"]),
            "channel_id": channel_id,
            "agency_id": form["agency_id"]
        }

        logging.info(data)

        update_channel_db(data)
    else:
        return prepare_response({"message":"Missmatch channel id"},401)
    
    return prepare_response({"message":"ok"},200)
    
@app.route('/channel/<channel_id>/delete',methods=['POST','OPTIONS'])
def delete_channel(channel_id):
    if request.method == "OPTIONS":
        response = jsonify({"message":"OK"})
        response.headers.add('Access-Control-Allow-Headers', 'content-type')  
        response.headers.add('Access-Control-Allow-Methods', 'POST')
        response.headers.add('Access-Control-Allow-Origin', URL)    
        response.headers.add('Access-Control-Allow-Credentials', 'true')        
        return response

    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return prepare_response({"message":"Missing token"},401)

    logging.info('%s %s %s',"delete channel",discord_user.verify_user(discord),channel_id,discord_user.channel_id)
    if discord_user.verify_user(discord) and channel_id == discord_user.channel_id:
        form = request.json
        if form["agency_id"] != 10:
            form["indie_notif"] = False

        data = {
            "channel_type": channel_type_rev(form["type"]),
            "upcoming": bool(form["upcoming"]),
            "dynamic": bool(form["dynamic"]),
            "region": str(','.join(form["channel_region"])),
            "lite": bool(form["lite"]),
            "indie": bool(form["indie_notif"]),
            "channel_id": channel_id,
            "agency_id": form["agency_id"]
        }
        
        logging.info(data)

        delete_channel_db(data)
    else:
        return prepare_response({"message":"Missmatch channel id"},401)
    
    return prepare_response({"message":"ok"},200)
    


if __name__ == '__main__':
    app.run()