import pandas as pd
import numpy as np
import shap
from .loader import MODEL, IS_LOADED
from .feature_descriptions import FEATURE_DESCRIPTIONS

def clean_features(features_dict: dict) -> dict:
    """Очищает features от NaN перед созданием DataFrame"""
    cleaned = {}
    for key, value in features_dict.items():
        if isinstance(value, float) and np.isnan(value):
            # Заменяем NaN на None (будет преобразовано в pandas.NA)
            cleaned[key] = None
        else:
            cleaned[key] = value
    return cleaned

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


def explain_features_split(df_sample, model, top_k=5):
    explainer = shap.Explainer(model)
    shap_values = explainer(df_sample)

    values = shap_values.values[0]
    features = df_sample.columns

    rows = []
    for feature, val in zip(features, values):
        rows.append({
            "feature": feature,
            "value": float(val),
            "abs": abs(val),
            "description": FEATURE_DESCRIPTIONS.get(feature, feature)
        })

    rows = sorted(rows, key=lambda x: x["abs"], reverse=True)[:top_k]

    positive = {}
    negative = {}

    for r in rows:
        if r["value"] > 0:
            positive[r["description"]] = r["value"]
        else:
            negative[r["description"]] = r["value"]

    return {
        "positive": positive,
        "negative": negative
    }
