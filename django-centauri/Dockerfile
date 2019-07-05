FROM python:3

WORKDIR /app

COPY . .
#RUN ls
#RUN pip install --no-cache-dir -r requirements.txt
RUN pip3 install --no-cache-dir django
RUN pip3 install --no-cache-dir djangorestframework
#COPY . .

#CMD [ "python", "./your-daemon-or-script.py" ]
EXPOSE 8000

CMD [ "python", "-u", "./manage.py", "runserver", "--noreload", "0.0.0.0:8000"]

