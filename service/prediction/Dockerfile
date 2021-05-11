FROM python:3.9-alpine
RUN apk add --update --no-cache git
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN pip install -r requirements.txt
CMD ["python","main.py"]