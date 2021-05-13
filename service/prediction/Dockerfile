FROM python:3.9-alpine
RUN apk add --update --no-cache git g++
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN pip install scikit-learn numpy scipy
RUN pip install -r requirements.txt
EXPOSE 9001
CMD ["python","main.py"]