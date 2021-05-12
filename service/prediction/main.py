import numpy as np
import matplotlib.pyplot as plt
import seaborn as sns
import requests,os,time,logging
from concurrent import futures
from sklearn.linear_model import LinearRegression
from sklearn.metrics import accuracy_score
from datetime import timedelta, date, datetime

import grpc

import prediction_pb2
import prediction_pb2_grpc

prometheus = os.environ['PrometheusURL']

class Prediction(prediction_pb2_grpc.PredictionServicer):
    def __init__(self):
        self.Predic = SubscriberPredic()        

    def GetSubscriberPrediction(self,request, context):
        x,y = getData(request.State,request.Name,request.Limit)
        if x == [] and y == []:
            return prediction_pb2.MessageResponse(Code=1,Prediction="0",Score="0")

        self.Predic.addData(x,y)
        num = self.Predic.predic(len(x)+7)
        score = self.Predic.accuracy()
        logging.info("state %s name %s limit %s predick %s score %s",request.State,request.Name,request.Limit,num,score)
        return prediction_pb2.MessageResponse(Code=0,Prediction=str(int(num)),Score=str(int(score *100))+"%")


#Prometheus Query : get_subscriber{vtuber="Parerun",state="Twitter"}
class SubscriberPredic:
    def __init__(self):
        self.lg = LinearRegression()

    def addData(self,x,y):
        self.x = np.array(x).reshape((-1, 1))
        self.y = np.array(y)

    def predic(self,num):
        self.lg.fit(self.x,self.y)
        now = datetime.today()
        nowformat = now-timedelta(days=num)    
        predict_x = self.lg.predict([[nowformat.strftime('%y%m%d')]])
        return predict_x
    def accuracy(self):
        accuracy = self.lg.score(self.x,self.y)
        return accuracy

def getData(state,name :str,lmt :int):
    now = datetime.today()
    x = []
    y = []
    for i in range(lmt,0,-1):
        tmp = now.replace(hour=23,minute=59,second=59)-timedelta(days=i)
        try:
            response = requests.get(prometheus+'/api/v1/query?query=get_subscriber{state="'+state+'",vtuber="'+name+'"}&time='+tmp.strftime('%s'))
            response.raise_for_status()
            data = response.json()
            print(data["data"]["result"],tmp)
            if len(data["data"]["result"]) > 0:
                y.append(int(data["data"]["result"][0]["value"][1]))
                x.append(tmp.strftime('%y%m%d'))
            elif len(data["data"]["result"]) == 0:
                y = []
                x = y
                logging.error("data null,can't be processed")
                break
        except requests.HTTPError as exception:
            logging.error("%s",exception)
            break

    x.reverse()
    return x,y


def serve():
    listen_addr = '[::]:9001'    
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    prediction_pb2_grpc.add_PredictionServicer_to_server(Prediction(), server)
    logging.info("Starting server on %s", listen_addr)
    server.add_insecure_port(listen_addr)
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()