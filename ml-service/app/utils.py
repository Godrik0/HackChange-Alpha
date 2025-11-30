import logging
from datetime import datetime, timezone
import json

logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
)
logger = logging.getLogger("ml_service")


def utc_now():
    return datetime.now(timezone.utc)


def safe_serialize(obj):
    try:
        return json.dumps(obj)
    except Exception:
        return str(obj)
