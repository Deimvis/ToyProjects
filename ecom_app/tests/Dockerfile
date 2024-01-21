FROM python:3.11.2

WORKDIR /tests/

COPY requirements.txt requirements.txt
RUN pip install --upgrade pip && pip install -r requirements.txt

COPY tests /tests
COPY pyproject.toml pyproject.toml

ENTRYPOINT [ "pytest" ]
