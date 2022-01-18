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
all_region = []

for i in groups:
    for k in i["Region"]:
        if k not in all_region:
            all_region.append(k)

if 'http://' in OAUTH2_REDIRECT_URI:
    os.environ['OAUTHLIB_INSECURE_TRANSPORT'] = 'true'

def token_updater(token):
    session['oauth2_token'] = token


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
    return redirect(url_for('.guilds'))


@app.route('/guilds')
def guilds():
    discord = make_session(token=session.get('oauth2_token'))
    guilds = []
    for i in discord.get(API_BASE_URL + '/users/@me/guilds').json():
        if int(i["permissions"]) & 8208:
            isVtbot = requests.get(API_BASE_URL+f'/guilds/{i["id"]}',headers=bot_headers)
            if isVtbot.status_code == 200:
                guilds.append(i)
    #return jsonify(guilds=guilds)
    return render_template("guilds.html",guilds = guilds)

@app.route('/guilds/<guild_id>/channels')
def channels(guild_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return redirect("/")
    channels = requests.get(API_BASE_URL + f'/guilds/{guild_id}/channels',headers=bot_headers).json()
    channels_info = {
        "agency" : {},
        "guild_channel": [],
        "all_region": all_region,
    }
    
    for i in channels:
        if i["parent_id"] is None:
            continue
        channels_info["guild_channel"].append(i)
        
    channels_info["agency"] =  groups

    #return jsonify(channels_info=channels_info)
    return render_template("channels.html",channels_info=channels_info)

@app.route('/channel/<channel_id>/<agency_id>',methods=['GET'])
def channel(channel_id,agency_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return redirect("/")
    
    channel_data = get_channel_info(agency_id,channel_id)

    region = {}
    for i in groups:
        if int(i["id"]) == int(agency_id):
            if channel_data == []:
                return  jsonify(channel_data={
                "channel_id": channel_id,
                "dynamic": None,
                "indie_notif": None,
                "upcoming": None,
                "lite": None,
                "region": i["Region"],
                "type" : None,
            })            

            for k in i["Region"]:
                region[k] = False
                if k in channel_data[0]["Region"]:
                    region[k] = True

    return jsonify(channel_data={
        "channel_id": channel_data[0]["DiscordChannelID"],
        "dynamic": bool(channel_data[0]["Dynamic"]),
        "indie_notif": bool(channel_data[0]["IndieNotif"]),
        "lite": bool(channel_data[0]["Lite"]),
        "upcoming": bool(channel_data[0]["NewUpcoming"]),
        "region": channel_data[0]["Region"].split(","),
        "type" : channel_data[0]["Type"],
    })



@app.route('/channel/<channel_id>/update',methods=['POST'])
def update_channel(channel_id):
    discord = make_session(token=session.get('oauth2_token'))
    if len(discord.token.keys()) < 1:
        return redirect("/")
    print(channel_id,request.form)
    return "ok"

if __name__ == '__main__':
    app.run()