from fastapi import FastAPI
from routes import config_routes

app = FastAPI(
    title="Secure Configuration Management API",
    description="API for managing system settings and critical maintenance operations.",
    version="1.0.0"
)

# Include routers
app.include_router(config_routes.router)

@app.get("/", tags=["root"])
def read_root():
    return {"message": "Welcome to the Secure System Configuration API. Access /docs for API documentation."}

# To run this application:
# uvicorn main:app --reload
