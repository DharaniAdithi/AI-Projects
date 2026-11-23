import logging
import sys
from functools import lru_cache

from src.constants import DEFAULT_LOG_LEVEL

class _SingletonLogger:
    """Configures root logger once across the entire app."""

    _configured: bool = False

    @classmethod
    def configure(cls, level: str = DEFAULT_LOG_LEVEL) -> None:
        """
        Configure root logger once.

        Args:
            level (str): Log level e.g. 'INFO', 'DEBUG'.
        """
        if cls._configured:
            return
        log_level = getattr(logging, level.upper(), logging.INFO)
        handler = logging.StreamHandler(sys.stdout)
        handler.setFormatter(
            logging.Formatter(
                fmt="%(asctime)s | %(levelname)s | %(name)s | %(message)s",
                datefmt="%Y-%m-%dT%H:%M:%S",
            )
        )
        root = logging.getLogger()
        root.setLevel(log_level)
        root.addHandler(handler)
        cls._configured = True


@lru_cache(maxsize=None)
def get_logger(name: str) -> logging.Logger:
    """
    Return a named logger (cached — Singleton per name).

    Args:
        name (str): Logger name, use __name__.

    Returns:
        logging.Logger: Configured logger.
    """
    _SingletonLogger.configure()
    return logging.getLogger(name)