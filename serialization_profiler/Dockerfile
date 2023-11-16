FROM python:latest

COPY /bin/ /bin/
COPY /src/ /src/
COPY /scripts/ /scripts/
COPY requirements.txt requirements.txt
RUN pip install -r requirements.txt

ENTRYPOINT [ "/scripts/launch_server" ]
