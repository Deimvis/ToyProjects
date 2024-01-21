FROM python:3.11.2

COPY user_service/requirements.txt requirements.txt
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY auth_sdk /auth_sdk
RUN python -m pip install ./auth_sdk

COPY user_service/db.py         /db.py
COPY user_service/exceptions.py /exceptions.py
COPY user_service/main.py       /main.py
COPY user_service/models.py     /models.py
COPY user_service/shortcuts.py  /shortcuts.py

ENTRYPOINT [ "uvicorn", "main:app" ]
