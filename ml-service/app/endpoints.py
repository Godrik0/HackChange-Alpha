import pandas as pd
import logging
from fastapi import APIRouter, HTTPException
from uuid import uuid4

from .schemas import PredictionRequest, PredictionResponse, HealthResponse, ModelInfoResponse
from .inference import explain_features_split
from .loader import is_model_loaded, get_model, load_assets
from .config import settings
from .utils import logger

logger = logging.getLogger(__name__)
router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
async def predict_income(request: PredictionRequest):
    logger.info(f"[PREDICT] Received request: user_id={request.user_id}, features={len(request.features)}")
    
    if not is_model_loaded():
        logger.error("[PREDICT] Model not loaded!")
        raise HTTPException(status_code=503, detail="Model not loaded")

    uid = str(uuid4())
    
    try:
        df = pd.DataFrame([request.features])
        model = get_model()
        prediction_value = model.predict(df)[0]
    except ValueError as ve:
        logger.error(f"[PREDICT] Invalid features: {ve}")
        raise HTTPException(status_code=400, detail=f"Invalid features: {ve}")
    except Exception as e:
        logger.error(f"[PREDICT] Prediction failed: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail="Prediction failed")
    
    try:
        explanation = explain_features_split(df, model)
    except Exception as e:
        logger.error(f"[PREDICT] Explanation failed: {e}", exc_info=True)
        raise HTTPException(status_code=500, detail="Explanation failed")
    
    logger.info(f"[PREDICT] prediction={prediction_value:.2f}, uid={uid}")
    
    return PredictionResponse(
        prediction=float(prediction_value),
        explanation=explanation,
        uid=uid,
    )


@router.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="healthy" if is_model_loaded() else "unhealthy",
        service=settings.APP_NAME,
        is_model_loaded=is_model_loaded(),
        model_version=settings.APP_VERSION,
    )


@router.get("/model-info", response_model=ModelInfoResponse)
async def model_info():
    if not is_model_loaded():
        raise HTTPException(status_code=503, detail="Model not loaded")

    model = get_model()
    features = model.feature_names_in_.tolist() if hasattr(model, "feature_names_in_") else []
    
    return ModelInfoResponse(
        model_type=type(model).__name__,
        model_version=settings.APP_VERSION,
        features=features,
        allowed_categories=[],
        model_loaded=True,
    )


@router.post("/reload-model")
async def reload_model():
    load_assets()
    return {"status": "success"}
