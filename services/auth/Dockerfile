FROM python:3.11-alpine
WORKDIR /app
RUN pip install --no-cache-dir fastapi uvicorn psycopg2-binary
COPY server.py /app/server.py
EXPOSE 8002
CMD ["uvicorn", "server:app", "--host", "0.0.0.0", "--port", "8002"]
