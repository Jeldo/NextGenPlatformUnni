import os


DATABASE_URL = os.getenv(
    "DATABASE_URL",
    "postgresql://calendar:calendar@localhost:5432/treatment_calendar",
)
PORT = int(os.getenv("PORT", "8081"))
