import pandas as pd
import logging
from fastapi import APIRouter, HTTPException
from uuid import uuid4

from .schemas import PredictionRequest, PredictionResponse, HealthResponse, ModelInfoResponse
from .inference import predict, explain_features_split, clean_features
from .loader import is_model_loaded, get_model, load_assets
from .config import settings
from .utils import logger

logger = logging.getLogger(__name__)
router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
async def predict_income(request: PredictionRequest):
    logger.info(f"[PREDICT] Received request: model_version={request.model_version}, pipeline_version={request.pipeline_version}, user_id={request.user_id}")
    logger.info(f"[PREDICT] Features count: {len(request.features)}")
    
    if not is_model_loaded():
        logger.error("[PREDICT] Model not loaded!")
        raise HTTPException(status_code=503, detail="Model not loaded")

    # Создаем uid для ответа и логов
    uid = str(uuid4())
    logger.info(f"[PREDICT] Request uid={uid}, user_id={request.user_id}")

    try:
        df = pd.DataFrame([request.features])
        logger.info(f"[PREDICT] DataFrame created, shape={df.shape}, columns={list(df.columns)[:5]}...")
    except Exception as e:
        logger.exception("[PREDICT] Failed to convert features to DataFrame")
        logger.warning("[PREDICT] Returning mock data due to DataFrame creation failure")
        return PredictionResponse(
            prediction=50000.0,
            explanation={"education": 0.3, "experience": 0.25, "age": -0.1},
            uid=uid,
        )

    # Предсказание
    try:
        logger.info("[PREDICT] Starting model prediction...")
        model = get_model()
        prediction_value = model.predict(df)[0]  # pipeline уже включает препроцессинг
        logger.info(f"[PREDICT] Model prediction SUCCESS: {prediction_value}")
    except Exception as e:
        logger.exception("[PREDICT] Prediction failed")
        logger.warning("[PREDICT] Returning mock data due to prediction failure")
        return PredictionResponse(
            prediction=50000.0,
            explanation={"education": 0.3, "experience": 0.25, "age": -0.1},
            uid=uid,
        )

    try:
        logger.info("[PREDICT] Calculating SHAP explanation...")
        explanation = explain_features_split(df, model)  # вернет {positive: {...}, negative: {...}}
        logger.info(f"[PREDICT] SHAP explanation SUCCESS: {len(explanation.get('positive', {}))} positive, {len(explanation.get('negative', {}))} negative")
    except Exception as e:
        logger.warning(f"[PREDICT] SHAP explanation failed: {e}")
        explanation = {"education": 0.3, "experience": 0.25, "age": -0.1}

    logger.info(f"[PREDICT] Returning response: prediction={prediction_value}, uid={uid}")
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
    if hasattr(model, "feature_names_in_"):
        features = model.feature_names_in_.tolist()
    else:
        features = []

    # Get allowed categories if available
    allowed_categories = []
    
    return ModelInfoResponse(
        model_type=type(model).__name__,
        model_version=settings.APP_VERSION,
        features=features,
        allowed_categories=allowed_categories,
        model_loaded=True,
    )


@router.post("/reload-model")
async def reload_model():
    load_assets()
    return {"status": "success"}
