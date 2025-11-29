from pydantic import BaseModel, Field
from typing import Optional, List, Dict, Any
from datetime import datetime, timezone


class PredictionRequest(BaseModel):
    client_data: dict
    return_explanation: bool = False
    user_id: Optional[str] = None


class PredictionResponse(BaseModel):
    prediction: float
    model_version: str
    status: str
    timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    explanation: Optional[List[Dict[str, float]]] = None
    confidence: Optional[float] = None
    request_id: Optional[str] = None
    message: Optional[str] = None


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
