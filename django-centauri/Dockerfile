FROM python:3

WORKDIR /app

COPY . .
RUN pip3 install --no-cache-dir django
RUN pip3 install --no-cache-dir djangorestframework

RUN mkdir /app/db
RUN python3 manage.py makemigrations
RUN python3 manage.py migrate

VOLUME /app/db

EXPOSE 8000

CMD [ "python", "-u", "./manage.py", "runserver", "--noreload", "0.0.0.0:8000"]

