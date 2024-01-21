FROM python:3.11.2

COPY order_service/requirements.txt requirements.txt
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY auth_sdk /auth_sdk
RUN python -m pip install ./auth_sdk

COPY order_service/conf.py       /conf.py
COPY order_service/db.py         /db.py
COPY order_service/main.py       /main.py
COPY order_service/models.py     /models.py
COPY order_service/shortcuts.py  /shortcuts.py

ENTRYPOINT [ "uvicorn", "main:app" ]
