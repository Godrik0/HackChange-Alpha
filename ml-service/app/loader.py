import dill
import os
from time import time
from contextlib import asynccontextmanager

from .config import settings
from .utils import logger

# Используем словарь чтобы обойти проблему с импортом глобальных переменных
_state = {
    'model': None,
    'is_loaded': False,
    'load_time': None
}

MODEL = None
IS_LOADED = False
LOAD_TIME = None


def is_model_loaded():
    """Проверка, загружена ли модель"""
    return _state['is_loaded']


def get_model():
    """Получить загруженную модель"""
    return _state['model']


def load_assets(model_path=None):
    global MODEL, IS_LOADED, LOAD_TIME

    model_path = model_path or settings.MODEL_PATH
    start = time()

    try:
        if not os.path.exists(model_path):
            logger.error(f"Model file NOT FOUND: {model_path}")
            IS_LOADED = False
            return

        with open(model_path, 'rb') as f:
            MODEL = dill.load(f)
            _state['model'] = MODEL

        logger.info(f"Pipeline loaded from {model_path}")

        if hasattr(MODEL, 'steps'):
            logger.info(f"Pipeline steps: {MODEL.steps}")
        else:
            logger.info(f"Model type: {type(MODEL).__name__}")

        IS_LOADED = True
        _state['is_loaded'] = True
        LOAD_TIME = time() - start
        _state['load_time'] = LOAD_TIME
        logger.info(f"Assets loaded in {LOAD_TIME:.2f}s")

    except Exception as e:
        logger.exception(f"Failed to load pipeline: {e}")
        IS_LOADED = False
        _state['is_loaded'] = False


@asynccontextmanager
async def lifespan(app):
    logger.info("Starting ML Service...")
    load_assets()
    yield
    logger.info("ML Service shutting down...")
