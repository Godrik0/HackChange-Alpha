from pydantic import BaseModel, Field
from typing import Optional, List, Dict
from datetime import datetime, timezone


class PredictionRequest(BaseModel):
    model_version: str
    pipeline_version: str
    features: Dict[str, float | int | str]
    user_id: str = "anonymous"


class PredictionResponse(BaseModel):
    prediction: float
    explanation: Optional[Dict[str, float]] = None
    uid: str


class HealthResponse(BaseModel):
    status: str
    service: str
    is_model_loaded: bool
    model_version: str
    timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))


class ModelInfoResponse(BaseModel):
    model_type: str
    model_version: str
    features: List[str]
    allowed_categories: List[str]
    model_loaded: bool
