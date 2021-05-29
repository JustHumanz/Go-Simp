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
        x,y = getData(request.State,request.Name,356)
        if x == [] and y == []:
            return prediction_pb2.MessageResponse(Code=1,Prediction=0,Score=0)

        self.Predic.SpellMagic(x,y)
        num = self.Predic.DoBlackMagic(len(x)+request.Limit)
        score = self.Predic.CheckBlackMagic()
        logging.info("state %s name %s limit %s predick %s score %s",request.State,request.Name,request.Limit,num,score)
        return prediction_pb2.MessageResponse(Code=0,Prediction=int(num),Score=score)

    def GetReverseSubscriberPrediction(self,request, context):
        x,y = getData(request.State,request.Name,356)
        if x == [] and y == []:
            return prediction_pb2.MessageResponse(Code=1,Prediction=0,Score=0)

        self.Predic.SpellMagic(y,x)
        num = self.Predic.DoBlackMagic(request.Limit)
        score = self.Predic.CheckBlackMagic()
        fromordinal = datetime.fromordinal(int(num))
        logging.info("state %s name %s limit %s predick %s score %s",request.State,request.Name,request.Limit,fromordinal,score)
        return prediction_pb2.MessageResponse(Code=0,Prediction=int(fromordinal.timestamp()),Score=score)
#Prometheus Query : get_subscriber{vtuber="Parerun",state="Twitter"}
class SubscriberPredic:
    def __init__(self):
        self.lg = LinearRegression()

    def SpellMagic(self,x,y):
        self.x = np.array(x).reshape(-1, 1)
        self.y = np.array(y)

    def DoBlackMagic(self,num):
        self.lg.fit(self.x,self.y)   
        predict_x = self.lg.predict([[num]])
        return predict_x

    def CheckBlackMagic(self):
        accuracy = self.lg.score(self.x,self.y)
        return accuracy

def getData(state,name :str,lmt :int):
    start = datetime.today().replace(hour=23,minute=59,second=59)-timedelta(days=lmt)
    y = []
    x = []
    counter = 1
    try:
        response = requests.get(prometheus+'/api/v1/query_range?query=get_subscriber{state="'+state+'",vtuber="'+name+'"}',params={'start':start.strftime('%s'),'end':datetime.today().strftime('%s'),'step':'86400'})
        logging.info("%s",response.url)
        response.raise_for_status()
        data = response.json()
        if len(data["data"]["result"]) > 0:
            for values in data["data"]["result"][0]['values']:
                x.append(counter)                
                y.append(int(values[1]))
                counter += 1
        elif len(data["data"]["result"]) == 0:
            logging.error("data null,can't be processed")
            return [],[]
    except requests.HTTPError as exception:
        logging.error("%s",exception)
        return [],[]        

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