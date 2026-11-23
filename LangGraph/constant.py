from ..settings.settings import settings

app_setting=settings

DB_URI = (
    f"postgresql://{settings.db_username}:"
    f"{settings.db_password}@"
    f"{settings.db_host}:"
    f"{settings.db_port}/"
    f"{settings.db_name_checkpoint}"
)

PYLINT_WORKSPACE="C:/Users/AkshayKrishn_eo05n92/Documents/review_temp.py"

LIMIT_RETRIES=3