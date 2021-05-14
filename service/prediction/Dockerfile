FROM python:3.9
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN pip install scikit-learn numpy scipy
RUN pip install -r requirements.txt
EXPOSE 9001
CMD ["python","main.py"]