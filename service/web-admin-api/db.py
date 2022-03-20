import mysql.connector,os,logging

def get_mysql_conn():
  VTDB = mysql.connector.connect(
    host=os.environ['DB_HOST'],
    user=os.environ['DB_USER'],
    password= os.environ['DB_PASS'],
    database="Vtuber"
  )
  VTcursor = VTDB.cursor(dictionary=True)
  return VTDB,VTcursor


logging.basicConfig(level=logging.INFO,format='%(asctime)s - %(message)s', datefmt='%d-%b-%y %H:%M:%S')

def get_channel_info(group_id,channel_id):
    _,VTcursor = get_mysql_conn()

    VTcursor.execute("SELECT Channel.*,VtuberGroup.VtuberGroupName,VtuberGroupIcon FROM Channel INNER JOIn VtuberGroup on VtuberGroup.id=Channel.VtuberGroup_id WHERE VtuberGroup_id= %s AND DiscordChannelID= %s", (group_id,channel_id,))
    return VTcursor.fetchall()


def get_groups():
    _,VTcursor = get_mysql_conn()

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

def update_channel_db(data):
  VTDB,VTcursor = get_mysql_conn()
  if get_channel_info(data["agency_id"],data["channel_id"]) != []:
    logging.info('Update channel %s agency %d',data["channel_id"],data["agency_id"])
    query = "UPDATE Channel SET Type = %s,NewUpcoming = %s,Dynamic = %s,Region = %s,Lite= %s,IndieNotif= %s WHERE DiscordChannelID = %s and VtuberGroup_id = %s"
    var = (
      data["channel_type"],
      data["upcoming"],
      data["dynamic"],
      data["region"],
      data["lite"],
      data["indie"],
      data["channel_id"],
      data["agency_id"]
    )
    VTcursor.execute(query,var)
    VTDB.commit()
  else:
    logging.info('Add new channel %s "agency" %d',data["channel_id"],data["agency_id"])
    query = "INSERT INTO Channel (Type,NewUpcoming,Dynamic,Region,Lite,IndieNotif,DiscordChannelID,VtuberGroup_id) VALUES (%s,%s,%s,%s,%s,%s,%s,%s)"
    var = (
      data["channel_type"],
      data["upcoming"],
      data["dynamic"],
      data["region"],
      data["lite"],
      data["indie"],
      data["channel_id"],
      data["agency_id"]      
    )
    VTcursor.execute(query,var)
    VTDB.commit()    

def delete_channel_db(data):
    VTDB,VTcursor = get_mysql_conn()
    logging.info('Add new channel %s "agency" %d',data["channel_id"],data["agency_id"])
    query = "DELETE FROM Channel where DiscordChannelID = %s and VtuberGroup_id = %s"
    var = (
      data["channel_id"],
      data["agency_id"]      
    )
    VTcursor.execute(query,var)
    VTDB.commit()    
