import mysql.connector
import os

VTDB = mysql.connector.connect(
  host="localhost",
  user="root",
  password= os.environ['DB_PASS'],
  database="Vtuber"
)

VTcursor = VTDB.cursor(dictionary=True)

def get_channel_info(group_id,channel_id):
    VTcursor.execute("SELECT Channel.*,VtuberGroup.VtuberGroupName,VtuberGroupIcon FROM Channel INNER JOIn VtuberGroup on VtuberGroup.id=Channel.VtuberGroup_id WHERE VtuberGroup_id= %s AND DiscordChannelID= %s", (group_id,channel_id,))
    return VTcursor.fetchall()


def get_groups():
    VTcursor.execute("SELECT * FROM Vtuber.VtuberGroup")
    agency = []
    for i in VTcursor.fetchall():
      VTcursor.execute(f"SELECT Region FROM Vtuber.VtuberMember WHERE VtuberGroup_id={i['id']} group by Region")
      region = []
      for k in VTcursor.fetchall():
          region.append(k["Region"]) 
      
      i["Region"] = region

      agency.append(i)

    return agency
