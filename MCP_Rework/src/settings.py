from dataclasses import dataclass
from functools import lru_cache
import os
from dotenv import load_dotenv

from src.constants import (
    DEFAULT_HOST,
    DEFAULT_PORT,
    DEFAULT_LOG_LEVEL,
    DEFAULT_REGION
)


load_dotenv()


@dataclass
class Config:
    """
    Application configuration loaded from environment variables.

    Attributes:
        db_port (str): Database port.
        db_host (str): Database host.
        db_name (str): Database name.
        db_username (str): Database username.
        db_password (str): Database password.
        port (int): App server port.
        host (str): App server host.
        log_level (str): Logging level.
    """

    db_port: str
    db_host: str
    db_name: str
    db_username: str
    db_password: str
    port: int
    host: str
    log_level: str
    aws_region:str
    model_id:str
    provider:str


@lru_cache()
def get_settings() -> Config:
    """
    Return cached Config instance (Singleton via lru_cache).

    Returns:
        Config: Application settings.
    """
    return Config(
        db_port=os.getenv("DB_PORT"),
        db_host=os.getenv("DB_HOST", "localhost"),
        db_name=os.getenv("DB_NAME"),
        db_username=os.getenv("DB_USERNAME", "postgres"),
        db_password=os.getenv("DB_PASSWORD"),
        port=int(os.getenv("PORT", DEFAULT_PORT)),
        host=os.getenv("HOST", DEFAULT_HOST),
        log_level=os.getenv("LOG_LEVEL", DEFAULT_LOG_LEVEL),
        aws_region=os.getenv("AWS_REGION", DEFAULT_REGION),
        model_id=os.getenv("MODEL_ID", ""),
        provider= os.getenv("PROVIDER", "amazon")
    )


# Keep your existing config for backward compatibility
config = get_settings()
