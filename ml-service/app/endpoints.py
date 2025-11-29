import pandas as pd
import logging
from fastapi import APIRouter, HTTPException
from uuid import uuid4

from .schemas import PredictionRequest, PredictionResponse, HealthResponse, ModelInfoResponse
from .inference import predict, explain_features_split, clean_features
from .loader import IS_LOADED, MODEL, load_assets
from .config import settings
from .utils import logger

logger = logging.getLogger(__name__)
router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
async def predict_income(request: PredictionRequest):
    if not IS_LOADED:
        raise HTTPException(status_code=503, detail="Model not loaded")


    # Создаем uid для ответа и логов
    uid = str(uuid4())
    logger.info(f"Prediction request uid={uid}, user_id={request.user_id}")

    try:
        df = pd.DataFrame([request.features])
    except Exception as e:
        logger.exception("Failed to convert features to DataFrame")
        raise HTTPException(status_code=400, detail=f"Invalid features: {e}")

    # Предсказание
    try:
        prediction_value = MODEL.predict(df)[0]  # pipeline уже включает препроцессинг
    except Exception as e:
        logger.exception("Prediction failed")
        raise HTTPException(status_code=500, detail=str(e))

    try:
        explanation = explain_features_split(df, MODEL)  # вернет {positive: {...}, negative: {...}}
    except Exception as e:
        logger.warning(f"SHAP explanation failed: {e}")
        explanation = None

    return PredictionResponse(
        prediction=float(prediction_value),
        explanation=explanation,
        uid=uid,
    )


@router.get("/health", response_model=HealthResponse)
async def health_check():
    return HealthResponse(
        status="healthy" if IS_LOADED else "unhealthy",
        service=settings.APP_NAME,
        is_model_loaded=IS_LOADED,
        model_version=settings.APP_VERSION,
    )


@router.get("/model-info", response_model=ModelInfoResponse)
async def model_info():
    if not IS_LOADED:
        raise HTTPException(status_code=503, detail="Model not loaded")

    if hasattr(MODEL, "feature_names_in_"):
        features = MODEL.feature_names_in_.tolist()
    else:
        features = []

    return ModelInfoResponse(
        model_type=type(MODEL).__name__,
        model_version=settings.APP_VERSION,
        features=features,
        model_loaded=True,
    )


@router.post("/reload-model")
async def reload_model():
    load_assets()
    return {"status": "success"}
