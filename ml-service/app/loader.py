import dill
import os
from time import time
from contextlib import asynccontextmanager

from .config import settings
from .utils import logger

_state = {
    'model': None,
    'is_loaded': False
}

def is_model_loaded():
    return _state['is_loaded']

def get_model():
    return _state['model']


def load_assets(model_path=None):
    model_path = model_path or settings.MODEL_PATH
    start = time()

    if not os.path.exists(model_path):
        logger.error(f"Model file NOT FOUND: {model_path}")
        return

    with open(model_path, 'rb') as f:
        model = dill.load(f)
        _state['model'] = model

    logger.info(f"Pipeline loaded from {model_path}")

    if hasattr(model, 'steps'):
        preprocessor = model.steps[0][1]
        xgb_model = model.steps[-1][1]
        if hasattr(xgb_model, 'get_booster'):
            booster = xgb_model.get_booster()
            preprocessor.expected_feature_order_ = list(booster.feature_names)
            logger.info(f"Injected {len(booster.feature_names)} features into preprocessor")

    _state['is_loaded'] = True
    logger.info(f"Assets loaded in {time() - start:.2f}s")


@asynccontextmanager
async def lifespan(app):
    logger.info("Starting ML Service...")
    load_assets()
    yield
    logger.info("ML Service shutting down...")
