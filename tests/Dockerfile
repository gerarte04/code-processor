FROM python:3.9-slim

RUN pip install pytest requests

WORKDIR /app

COPY tests .

CMD [ "pytest", "tests.py" ]
