#!/usr/bin/env python3
import os
import sys
from fastapi import FastAPI
import uvicorn
import psycopg2

app = FastAPI(title="Care4U Auth & Healthcare Core API", version="1.0.0")

DB_URL = os.getenv("DATABASE_URL", "postgresql://infortts_admin:InforttsSecureDB2026!@127.0.0.1:5432/care4u_db")

@app.get("/")
def root():
    return {
        "status": "online",
        "service": "Care4U Healthcare Suite",
        "node": "Oracle Cloud ap-mumbai-1 (130.210.24.48)",
        "db": "care4u_db (PostgreSQL 16 Shared)"
    }

@app.get("/health")
def health():
    return {"status": "healthy", "service": "care4u-api"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8002)
