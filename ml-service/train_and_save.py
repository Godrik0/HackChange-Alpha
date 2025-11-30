import pandas as pd
import numpy as np
import dill
import os
import sys
from sklearn.pipeline import Pipeline

import xgboost as xgb

sys.path.append(os.path.dirname(__file__))

from app.preprocessor import BankPreprocessor

df_train = pd.read_csv('hackathon_income_train.csv',
                       sep=";", engine="python", decimal=",")

y = df_train['target'].astype(float)
w = df_train['w'].astype(float)
X_raw = df_train.drop(columns=['target', 'w'])

mask_full = (
        (~np.isnan(y)) & (~np.isinf(y)) & (np.abs(y) <= 1e10) &
        (~np.isnan(w)) & (~np.isinf(w))
)
X_raw = X_raw[mask_full]
y = y[mask_full]
w = w[mask_full]

print(f"Данные загружены: {X_raw.shape[0]} samples, {X_raw.shape[1]} features")

pipeline = Pipeline([
    ('preprocessor', BankPreprocessor(min_job_freq=10, verbose=False)),
    ('model', xgb.XGBRegressor(
        random_state=42,
        objective='reg:absoluteerror',
        enable_categorical=True,
        tree_method='hist',
        n_jobs=-1,
        subsample=0.7,
        reg_lambda=1,
        reg_alpha=0.5,
        n_estimators=600,
        min_child_weight=3,
        max_depth=8,
        learning_rate=0.05,
        colsample_bytree=0.9
    ))
])

pipeline.fit(X_raw, y, model__sample_weight=w)

with open('models/pipeline.pkl', 'wb') as f:
    dill.dump(pipeline, f)
