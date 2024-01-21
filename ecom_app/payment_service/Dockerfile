FROM python:3.11.2

COPY payment_service/requirements.txt requirements.txt
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY auth_sdk /auth_sdk
RUN python -m pip install ./auth_sdk

COPY payment_service/db.py         /db.py
COPY payment_service/main.py       /main.py
COPY payment_service/models.py     /models.py
COPY payment_service/shortcuts.py  /shortcuts.py

ENTRYPOINT [ "uvicorn", "main:app" ]
