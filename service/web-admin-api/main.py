import os,time,requests,json
from flask import Flask, g, session, redirect, request, url_for, jsonify,render_template
from requests_oauthlib import OAuth2Session
from db import *


OAUTH2_CLIENT_ID = "719540207552167936"
OAUTH2_CLIENT_SECRET =  os.environ['SECRET']
OAUTH2_REDIRECT_URI =  os.environ['URL']
bot_headers = {"Authorization": "Bot "+ os.environ['BOT'],"Content-Type":"application/json"}

API_BASE_URL = os.environ.get('API_BASE_URL', 'https://discordapp.com/api')
AUTHORIZATION_BASE_URL = API_BASE_URL + '/oauth2/authorize'
TOKEN_URL = API_BASE_URL + '/oauth2/token'

app = Flask(__name__,template_folder='template')
app.debug = True
app.config['SECRET_KEY'] = OAUTH2_CLIENT_SECRET

groups = get_groups()

if 'http://' in OAUTH2_REDIRECT_URI:
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = 'true'

def token_updater(token):
    session['oauth2_token'] = token

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
    return redirect("http://localhost:8080/#/guilds")


@app.route('/guilds')
def guilds():
    discord = make_session(token=session.get('oauth2_token'))
    guilds = []
    for i in discord.get(API_BASE_URL + '/users/@me/guilds').json():
        if int(i["permissions"]) & 8208:
            isVtbot = requests.get(API_BASE_URL+f'/guilds/{i["id"]}',headers=bot_headers)
            if isVtbot.status_code == 200:
                guilds.append(i)            

    response = jsonify({'guilds': guilds})
    response.headers.add('Access-Control-Allow-Origin', 'http://localhost:8080')    
    response.headers.add('Access-Control-Allow-Credentials', 'true')
    return response
    #return render_template("guilds.html",guilds = guilds)

@app.route('/guilds/<guild_id>/channels')
def channels(guild_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return redirect("/")
    channels = requests.get(API_BASE_URL + f'/guilds/{guild_id}/channels',headers=bot_headers).json()
    channels_info = {
        "guild_channel": [],
    }

    for i in channels:
        if i["parent_id"] is None:
            continue
        channels_info["guild_channel"].append(i)

    response = jsonify(channels_info)
    response.headers.add('Access-Control-Allow-Origin', 'http://localhost:8080')    
    response.headers.add('Access-Control-Allow-Credentials', 'true')
    return response

@app.route('/channel/<channel_id>/agency',methods=['GET'])
def channel(channel_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return redirect("/")
    
    agency_list = []
    for k in groups:
        agency_id = k["id"]
        channel_data = get_channel_info(agency_id,channel_id)
        if channel_data == []:
            agency_list.append({
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
            agency_list.append({
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

    response = jsonify(agency_list)
    response.headers.add('Access-Control-Allow-Origin', 'http://localhost:8080')    
    response.headers.add('Access-Control-Allow-Credentials', 'true')                
    return response



@app.route('/channel/<channel_id>/update',methods=['POST','GET','OPTIONS'])
def update_channel(channel_id):
    response = jsonify({"Status":"OK"})
    response.headers.add('Access-Control-Allow-Origin', 'http://localhost:8080')    
    response.headers.add('Access-Control-Allow-Credentials', 'true')     
    if request.method == "OPTIONS":
        response.headers.add('Access-Control-Allow-Headers', 'content-type')  
        response.headers.add('Access-Control-Allow-Methods', 'POST')
        return response

    print(request.json)
    response.headers.add('Access-Control-Allow-Credentials', 'true')           
    return response

if __name__ == '__main__':
    app.run()