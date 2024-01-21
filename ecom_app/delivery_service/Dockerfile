FROM python:3.11.2

COPY delivery_service/requirements.txt requirements.txt
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY auth_sdk /auth_sdk
RUN python -m pip install ./auth_sdk

COPY delivery_service/db.py         /db.py
COPY delivery_service/main.py       /main.py
COPY delivery_service/models.py     /models.py
COPY delivery_service/shortcuts.py  /shortcuts.py

ENTRYPOINT [ "uvicorn", "main:app" ]
