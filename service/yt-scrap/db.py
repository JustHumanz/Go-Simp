import mysql.connector,os

mydb = mysql.connector.connect(
  host=os.environ['DBHOST'],
  user=os.environ['SQLUSER'],
  password=os.environ['SQLPASS'],
  database="Vtuber"
)
dbconn = mydb.cursor()