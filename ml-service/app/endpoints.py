from fastapi import APIRouter, HTTPException
from uuid import uuid4

from .schemas import PredictionRequest, PredictionResponse, HealthResponse, ModelInfoResponse
from .inference import predict, explain
from .loader import IS_LOADED, MODEL, load_assets
from .config import settings
from .utils import logger

router = APIRouter()


@router.post("/predict", response_model=PredictionResponse)
async def predict_income(request: PredictionRequest):
    if not IS_LOADED:
        raise HTTPException(status_code=503, detail="Model not loaded")

    request_id = request.user_id or str(uuid4())
    logger.info(f"Prediction request: {request_id}")

    try:
        result = predict(request.client_data)
    except Exception as e:
        logger.exception("Prediction failed")
        raise HTTPException(status_code=500, detail=str(e))

    explanation = None
    if request.return_explanation:
        explanation = explain(request.client_data, result)

    return PredictionResponse(
        prediction=result,
        model_version=settings.APP_VERSION,
        status="success",
        explanation=explanation,
        confidence=0.9,
        request_id=request_id,
        message="OK",
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
