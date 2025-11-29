from pydantic_settings import BaseSettings
from pathlib import Path


class Settings(BaseSettings):
    app_name: str = "Alpha Income ML Service"
    app_version: str = "1.0.0"
    HOST: str = "0.0.0.0"
    PORT: int = 8000
    LOG_LEVEL: str = "info"

    MODEL_DIR: Path = Path(__file__).parent.parent / "models"
    MODEL_PATH: Path = MODEL_DIR / "pipeline.pkl"

    APP_VERSION: str = "1.0.0"

    class Config:
        env_file = ".env"
        env_file_encoding = "utf-8"


settings = Settings()
