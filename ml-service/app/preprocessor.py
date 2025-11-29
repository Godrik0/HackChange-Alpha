import re
import numpy as np
import pandas as pd
from sklearn.base import BaseEstimator, TransformerMixin


# ----- 1. Препроцессор -----

def simplify_job(title):
    if pd.isna(title):
        return "unknown"
    words = re.split(r'[\s\-–—]', str(title))
    words = words[:3]
    words = [re.sub(r'[^а-яА-Яa-zA-Z]', '', w) for w in words]
    words = [w for w in words if w]
    if not words:
        return "unknown"
    words = sorted(words)
    return " ".join(words).lower()


class BankPreprocessor(BaseEstimator, TransformerMixin):
    def __init__(self, min_job_freq=10, verbose=False):
        self.min_job_freq = min_job_freq
        self.frequent_jobs_ = None
        self.verbose = verbose

    def fit(self, X, y=None):
        X = X.copy()
        if "dp_ewb_last_employment_position" in X.columns:
            tmp = X["dp_ewb_last_employment_position"].apply(simplify_job)
            counts = tmp.value_counts()
            self.frequent_jobs_ = set(counts[counts >= self.min_job_freq].index)
        else:
            self.frequent_jobs_ = set(["unknown"])
        return self

    def transform(self, X):
        X = X.copy()

        # job_simplified
        if "dp_ewb_last_employment_position" in X.columns:
            X["job_simplified"] = X["dp_ewb_last_employment_position"].apply(simplify_job)
            X["job_simplified"] = X["job_simplified"].apply(
                lambda x: x if x in self.frequent_jobs_ else "other"
            )

        # region = adminarea / addrref
        if "adminarea" in X.columns and "addrref" in X.columns:
            X["adminarea"] = X["adminarea"].replace("", np.nan)
            X["addrref"] = X["addrref"].replace("", np.nan)
            X["region"] = X["adminarea"].combine_first(X["addrref"])

        # удаляем ненужные столбцы (и id как фичу)
        drop_cols = [
            'adminarea', 'addrref', 'city_smart_name',
            'period_last_act_ad', 'dp_address_unique_regions',
            'dt', 'dp_ewb_last_employment_position', 'id'
        ]
        existing_drop = [c for c in drop_cols if c in X.columns]
        X = X.drop(columns=existing_drop)

        # object -> category
        obj_cols = X.select_dtypes(include=["object"]).columns
        for col in obj_cols:
            X[col] = X[col].astype("category")

        return X
