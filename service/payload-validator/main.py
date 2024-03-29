import requests,json,os,time,logging

logging.basicConfig(level=logging.INFO,format='%(asctime)s - %(message)s', datefmt='%d-%b-%y %H:%M:%S')
DATA_DIR = "service/migrate/json/"
API = os.environ["APIURL"]

vtuber_json = os.listdir(DATA_DIR)
list_agency = requests.get(API+"groups/").json()
variables_checker = {
    "Nickname": [],
    "Twitter_Hashtag" : [],
    "Twitter_Lewd_Hashtag": [],
    "Twitter_Username": [],
    "BiliBili_Hashtag": [],
    "BiliBili_ID": [],
    "BiliBili_Room_ID": [],
    "Youtube_Channel": [],
}

error_report = []

def get_member(agency_id):
    members_raw = requests.get(API+f"members/?groupid={agency_id}")
    member_id = []
    for i in members_raw.json():
        member_id.append(str(i["ID"]))

    member = requests.get(API+f"members/{','.join(member_id)}")
    return member.json()

for i in vtuber_json:
    f = open(DATA_DIR+"/"+i)
    agency_json = json.load(f)

    if i == "Independent.json":
        agency_json["GroupName"] = "independent"
        agency_json["ID"] = "10"
        agency_json["GroupChannel"] = {
            "Youtube": None,
            "BiliBili": None
        }

    for agency in list_agency:
        if agency["GroupName"] == agency_json["GroupName"]:
            if agency_json["GroupChannel"]["Youtube"] != None:
                for i in agency_json["GroupChannel"]["Youtube"]:
                    if i["Region"] == "":
                        error_report.append(f"Detect null region in agency youtube channel {i['ChannelID']}")

            logging.info(f"[{agency['GroupName']}] Get member data")

            members = get_member(agency["ID"])
            time.sleep(10)
            
            for member_json in agency_json["Members"]:
                for member in members:
                    if member["NickName"] == member_json["Name"]:

                        logging.info(f"[{agency['GroupName']}] check member [{member_json['Name']}]")

                        if member_json["Name"] not in variables_checker["Nickname"]:
                            variables_checker["Nickname"].append(member_json["Name"])

                        else:
                            error_report.append(f"Detect duplicate nickname vtuber {member_json['Name']}")                        

                        if member["EnName"] != "" and member["EnName"] != member_json["EN_Name"]:
                            print(f"[{member_json['Name']}] EN name was updated, old {member['EnName']} new {member_json['EN_Name']}")

                        if member["JpName"] != "" and member["JpName"] != member_json["JP_Name"]:
                            print(f"[{member_json['Name']}] JP name was updated, old {member['JpName']} new {member_json['JP_Name']}")

                        if member["Region"] != "" and member["Region"] != member_json["Region"]:
                            print(f"[{member_json['Name']}] Region was updated, old {member_json['Name']['Region']} new {member_json['Region']}")                            

                        if member["Fanbase"] != "" and member["Fanbase"] != member_json["Fanbase"]:
                            print(f"[{member_json['Name']}] Fanbase was updated, old {member['Fanbase']} new {member_json['Fanbase']}")                                                        

                        if member["Status"] != "" and member["Status"] != member_json["Status"]:
                            print(f"[{member_json['Name']}] Status was updated, old {member['Status']} new {member_json['Status']}")                                                                                    

                        if member["Status"] != "" and member["Status"] != member_json["Status"]:
                            print(f"[{member_json['Name']}] Status was updated, old {member['Status']} new {member_json['Status']}")                             

                        #################BiliBili Stuff#################
                        if member["BiliBili"] == None != None and member_json["BiliBili"]:
                            print(f"[{member_json['Name']}] Bilibili was added, old None new {member_json['BiliBili']}")

                        if member["BiliBili"] != None and member_json["BiliBili"] == None:
                            print(f"[{member_json['Name']}] Bilibili was deleted, old {member['BiliBili']} new None")

                        if member["BiliBili"] != None and member_json["BiliBili"] != None:

                            if member_json["BiliBili"]["BiliBili_Fanart"] != "":
                                if member_json["BiliBili"]["BiliBili_Fanart"] not in variables_checker["BiliBili_Hashtag"]:
                                    variables_checker["BiliBili_Hashtag"].append(member_json["BiliBili"]["BiliBili_Fanart"])
                                else:
                                    error_report.append(f"Detect duplicate bilibili hashtag {member_json['BiliBili']['BiliBili_Fanart']} vtuber {member_json['Name']}")

                            if member_json["BiliBili"]["BiliRoom_ID"] != 0:    
                                if member_json["BiliBili"]["BiliRoom_ID"] not in variables_checker["BiliBili_Room_ID"]:
                                    variables_checker["BiliBili_Room_ID"].append(member_json["BiliBili"]["BiliRoom_ID"])
                                else:
                                    error_report.append(f"Detect duplicate bilibili room {member_json['BiliBili']['BiliRoom_ID']} vtuber {member_json['Name']}")

                            if member_json["BiliBili"]["BiliBili_ID"] != 0:    
                                if member_json["BiliBili"]["BiliBili_ID"] not in variables_checker["BiliBili_ID"]:
                                    variables_checker["BiliBili_Room_ID"].append(member_json["BiliBili"]["BiliRoom_ID"])
                                else:
                                    error_report.append(f"Detect duplicate bilibili id {member_json['BiliBili']['BiliBili_ID']} vtuber {member_json['Name']}")
                        
                        #################Twitter Stuff#################
                        if member["Twitter"] == None and member_json["Twitter"]:
                            print(f"[{member_json['Name']}] Twitter was added, old None new {member_json['Twitter']}")

                            if member_json["Twitter"]["Twitter_Fanart"] == "#" or member_json["Twitter"]["Twitter_Lewd"] == "#" or member_json["Twitter"]["Twitter_Lewd"].startswith('#') == False:
                                error_report.append(f"Invalid twitter hashtag {member_json['Twitter']['Twitter_Fanart']} vtuber {member_json['Name']}")

                        if member["Twitter"] != None and member_json["Twitter"] == None:
                            print(f"[{member_json['Name']}] Twitter was deleted, {member['Twitter']} %s new None")

                        if member["Twitter"] != None and member_json["Twitter"] != None:
                            if member_json["Twitter"]["Twitter_Fanart"] != "":
                                if member_json["Twitter"]["Twitter_Fanart"] not in variables_checker["Twitter_Hashtag"]:
                                    variables_checker["Twitter_Hashtag"].append(member_json["Twitter"]["Twitter_Fanart"])
                                else:
                                    error_report.append(f"Detect duplicate twitter hashtag {member_json['Twitter']['Twitter_Fanart']} vtuber {member_json['Name']}")

                            if member_json["Twitter"]["Twitter_Username"] != "":    
                                if member_json["Twitter"]["Twitter_Username"] not in variables_checker["Twitter_Username"]:
                                    variables_checker["Twitter_Username"].append(member_json["Twitter"]["Twitter_Username"])
                                else:
                                    error_report.append(f"Detect duplicate Twitter username {member_json['Twitter']['Twitter_Username']} vtuber {member_json['Name']}")

                            if member_json["Twitter"]["Twitter_Lewd"] != "":    
                                if member_json["Twitter"]["Twitter_Lewd"] not in variables_checker["Twitter_Lewd_Hashtag"]:
                                    variables_checker["Twitter_Lewd_Hashtag"].append(member_json["Twitter"]["Twitter_Lewd"])
                                else:
                                    error_report.append(f"Detect duplicate twitter lewd hashtag {member_json['Twitter']['Twitter_Lewd']} vtuber {member_json['Name']}")

                        #################Youtube Stuff#################
                        if member["Youtube"] == None != None and member_json["Youtube"]:
                            print(f"[{member_json['Name']}] Twitter was added, old None new {member_json['Twitter']}")

                        if member["Youtube"] != None and member_json["Youtube"] == None:
                            print(f"[{member_json['Name']}] Twitter was deleted, old {member['Twitter']} new None")

                        if member["Youtube"] != None and member_json["Youtube"] != None:
                            if member_json["Youtube"]["Yt_ID"] != "":
                                if member_json["Youtube"]["Yt_ID"] not in variables_checker["Youtube_Channel"]:
                                    variables_checker["Youtube_Channel"].append(member_json["Youtube"]["Yt_ID"])
                                else:
                                    error_report.append(f"Detect duplicate youtube channel id {member_json['Youtube']['Yt_ID']} vtuber {member_json['Name']}")
                                                                                                                                        
    f.close()

if error_report != []:
    for i in error_report:
        print(i)
    os._exit(1)                        

print("all payload ok")