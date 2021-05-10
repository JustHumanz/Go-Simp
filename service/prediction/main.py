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
        return prediction_pb2.MessageResponse(Code=0,Prediction=str(num),Score=str(score))


#Prometheus Query : get_subscriber{vtuber="Parerun",state="Twitter"}
class SubscriberPredic:
    def __init__(self):
        self.lg = LinearRegression()

    def addData(self,x,y):
        self.x = np.array(x).reshape((-1, 1))
        self.y = np.array(y)

    def predic(self,num):
        self.lg.fit(self.x,self.y)
        predict_x = self.lg.predict([[num]])
        return predict_x
    def accuracy(self):
        accuracy = self.lg.score(self.x,self.y)
        return accuracy

def getData(state,name :str,lmt :int):
    now = datetime.today()
    x = []
    y = []
    for i in range(lmt,0,-1):
        tmp = now-timedelta(days=i)
        try:
            response = requests.get(prometheus+'/api/v1/query?query=get_subscriber{state="'+state+'",vtuber="'+name+'"}&time='+tmp.strftime('%s'))
            response.raise_for_status()
            data = response.json()
            if len(data["data"]["result"]) > 0:
                y.append(int(data["data"]["result"][0]["value"][1]))
                x.append(i)
        except requests.HTTPError as exception:
            print(exception)
            break

    x.reverse()
    return x,y


def serve():
    listen_addr = '[::]:9000'    
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    prediction_pb2_grpc.add_PredictionServicer_to_server(Prediction(), server)
    logging.info("Starting server on %s", listen_addr)
    server.add_insecure_port(listen_addr)
    server.start()
    server.wait_for_termination()


if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    serve()

"""
def main():
    
    x,y = getVtuberData("Twitter","Pekora")
    if x == [] and y == []:
        print("null data")

    Predic.addData(x,y)
    print("ID","Followers data")
    for i,j in zip(x,y):
        print(i,j)
    hasil = Predic.predic(len(x)+7)
    accuracy = Predic.accuracy()
    print("Prediksi",len(x)+7,"Hari kedepan",int(hasil),"presentase akurasi",accuracy)
main()
"""