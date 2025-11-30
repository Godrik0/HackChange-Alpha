import pandas as pd
import numpy as np
import shap
import logging
from .feature_descriptions import FEATURE_DESCRIPTIONS

logger = logging.getLogger(__name__)


def explain_features_split(df_sample, model, top_k=5):
    preprocessor = model.named_steps['preprocessor']
    xgb_model = model.named_steps['model']
    
    df_transformed = preprocessor.transform(df_sample)
    importance_dict = xgb_model.get_booster().get_score(importance_type='gain')
    
    total_importance = sum(importance_dict.values())
    if total_importance > 0:
        importance_dict = {k: v / total_importance for k, v in importance_dict.items()}
    
    sample_values = df_transformed.iloc[0].to_dict()
    
    rows = []
    for feature, importance in importance_dict.items():
        value = sample_values.get(feature, 0)
        if isinstance(value, (int, float)):
            contrib = importance * (1 if value > 0 else -1)
        else:
            contrib = importance
        
        rows.append({
            "feature": feature,
            "value": float(contrib),
            "abs": abs(contrib),
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
    
    return {"positive": positive, "negative": negative}
