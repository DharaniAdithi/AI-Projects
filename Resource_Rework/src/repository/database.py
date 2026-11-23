from sqlalchemy import create_engine, inspect, text
from sqlalchemy.orm import sessionmaker
from src.settings import get_settings
from urllib.parse import quote_plus
from src.utils.logger import get_logger

logger = get_logger(__name__)


class Database:
    """
    Singleton database connection manager.

    Responsibilities:
        - Create SQLAlchemy engine once.
        - Provide session factory.
        - Test DB connectivity.
    """

    _instance: "Database | None" = None

    def __new__(cls) -> "Database":
        """Ensure only one Database instance exists (Singleton)."""
        if cls._instance is None:
            cls._instance = super().__new__(cls)
            cls._instance._initialized = False
        return cls._instance

    def __init__(self) -> None:
        """Initialize engine and session factory once."""
        if self._initialized:
            return
        self.engine = self._create_engine()
        self.SessionLocal = sessionmaker(
            bind=self.engine,
            autoflush=False,
            autocommit=False,
        )
        self._initialized = True
        logger.info("Database initialized.")

    def _create_engine(self):
        """
        Build and return SQLAlchemy engine.

        Returns:
            Engine: Configured SQLAlchemy engine.
        """
        config = get_settings()
        password = quote_plus(config.db_password)
        db_url = (
            f"postgresql+psycopg2://{config.db_username}:"
            f"{password}@{config.db_host}:{config.db_port}/{config.db_name}"
        )
        return create_engine(db_url)

    def get_session(self):
        """
        Provide a database session with commit/rollback handling.

        Yields:
            Session: Active SQLAlchemy session.
        """
        try:
            session = self.SessionLocal()
            yield session
            session.commit()
        except Exception as e:
            session.rollback()
            raise e
        finally:
            session.close()

    def inspector(self, engine):
        """
        Return SQLAlchemy inspector for the given engine.

        Args:
            engine: SQLAlchemy engine instance.

        Returns:
            Inspector: SQLAlchemy inspector.
        """
        return inspect(engine)

    def test_connection(self) -> bool:
        """
        Test database connectivity.

        Returns:
            bool: True if connection succeeds, False otherwise.
        """
        try:
            with self.engine.connect() as connection:
                connection.execute(text("SELECT 1"))
            logger.info("Database connection test passed.")
            return True
        except Exception as e:
            logger.error("Database connection test failed: %s", e)
            return False