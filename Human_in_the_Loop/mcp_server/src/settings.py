import os
from dotenv import load_dotenv

load_dotenv()


class Settings:
    def __init__(self):
        self.db_port = os.getenv("DB_PORT", "5432")
        self.db_host = os.getenv("DB_HOST","localhost")
        self.db_name = os.getenv("DB_NAME", "project")
        self.db_username = os.getenv("DB_USERNAME","postgres")
        self.db_password = os.getenv("DB_PASSWORD", "Admin123")

        self.port = int(os.getenv("PORT", "8080"))
        self.host = os.getenv("HOST", "0.0.0.0")
        self.log_level = os.getenv("LOG_LEVEL", "info")

        self.delivery_fee = float(os.getenv("DELIVERY_FEE", "40.0"))
        self.free_delivery_threshold  = float(os.getenv("FREE_DELIVERY_THRESHOLD","5000.0"))


settings = Settings()
