import dill
import os
from time import time
from contextlib import asynccontextmanager

from .config import settings
from .utils import logger

MODEL = None
IS_LOADED = False
LOAD_TIME = None


def load_assets(model_path=None):
    global MODEL, IS_LOADED, LOAD_TIME

    model_path = model_path or settings.MODEL_PATH
    start = time()

    try:
        if not os.path.exists(model_path):
            logger.error(f"Model file NOT FOUND: {model_path}")
            IS_LOADED = False
            return

        MODEL = dill.load(model_path)
        logger.info(f"Pipeline loaded from {model_path}")
        logger.info(f"Pipeline steps: {MODEL.steps}")

        IS_LOADED = True
        LOAD_TIME = time() - start
        logger.info(f"Assets loaded in {LOAD_TIME:.2f}s")

    except Exception as e:
        logger.exception(f"Failed to load pipeline: {e}")
        IS_LOADED = False


@asynccontextmanager
async def lifespan(app):
    logger.info("Starting ML Service...")
    load_assets()
    yield
    logger.info("ML Service shutting down...")
