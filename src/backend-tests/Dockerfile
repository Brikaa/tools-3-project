FROM python:3.12.0-alpine3.18

RUN adduser \
  --no-create-home \
  --system \
  --shell /usr/sbin/nologin \
  tester

RUN pip install requests websocket-client
WORKDIR /app
COPY ./main.py .

USER tester
CMD ["python", "main.py"]
