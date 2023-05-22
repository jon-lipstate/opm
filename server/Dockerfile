! NOT_USED_RIGHT_NOW !
FROM python:3.9-slim-buster
WORKDIR /app
LABEL maintainer="jon@lipstate.com" \
      version="1.0"
COPY requirements.txt /app
RUN pip install --no-cache-dir -r requirements.txt
COPY . /app
EXPOSE 5000
CMD ["python", "app.py"]
