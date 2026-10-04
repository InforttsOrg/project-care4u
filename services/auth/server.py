#!/usr/bin/env python3
import os
import sys
from fastapi import FastAPI
import uvicorn
import psycopg2

app = FastAPI(title="Care4U Auth & Healthcare Core API", version="1.0.0")

# DATABASE_URL must be supplied by the environment. There is deliberately no
# fallback default: this file used to embed a live production DSN, which leaked
# a real credential into git history. Rotate that credential, then pass the
# replacement via the environment only.
DB_URL = os.getenv("DATABASE_URL", "")

@app.get("/")
def root():
    return {
        "status": "online",
        "service": "Care4U Healthcare Suite",
        "db": "care4u_db (PostgreSQL 16 Shared)"
    }

@app.get("/health")
def health():
    return {"status": "healthy", "service": "care4u-api"}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8002)
