import pandas as pd
import numpy as np
from .loader import MODEL, IS_LOADED


def preprocess(raw: dict) -> pd.DataFrame:
    return pd.DataFrame([raw])


def run_model(df: pd.DataFrame) -> float:
    if not IS_LOADED:
        raise RuntimeError("Model not loaded")
    pred = MODEL.predict(df)
    return float(np.squeeze(pred))


def predict(raw: dict) -> float:
    df = preprocess(raw)
    return run_model(df)


def explain(raw: dict, list_is_feature=None):
    """
    Возвращает топ-5 признаков с SHAP значениями и пояснениями.
    list_is_feature — словарь feature -> описание
    """
    pass
