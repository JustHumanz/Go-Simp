FROM python:3.9-alpine
RUN apk add --update --no-cache git g++
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN pip install scikit-learn 
RUN pip install -r requirements.txt
CMD ["python","main.py"]