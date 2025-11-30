import pandas as pd
import numpy as np
import shap
import logging
from .feature_descriptions import FEATURE_DESCRIPTIONS

logger = logging.getLogger(__name__)


def explain_features_split(df_sample, model, top_k=5):
    """
    Generate SHAP-based feature explanations for a prediction.
    Works with sklearn Pipeline containing preprocessor and XGBoost model.
    """
    try:
        preprocessor = model.named_steps['preprocessor']
        xgb_model = model.named_steps['model']
        
        df_transformed = preprocessor.transform(df_sample)
        
        explainer = shap.Explainer(xgb_model)
        shap_values = explainer(df_transformed)
        
        values = shap_values.values[0]
        features = df_transformed.columns
        
        rows = []
        for feature, val in zip(features, values):
            description = FEATURE_DESCRIPTIONS.get(feature, None)
            if description is None:
                continue
            rows.append({
                "feature": feature,
                "value": float(val),
                "abs": abs(val),
                "description": description
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
    
    except Exception as e:
        logger.error(f"SHAP explanation failed: {e}", exc_info=True)
        return {"positive": {}, "negative": {}}
