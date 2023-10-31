FROM python:3.12.0-alpine3.18

RUN pip install requests
WORKDIR /app
COPY ./main.py .

CMD ["python", "main.py"]
