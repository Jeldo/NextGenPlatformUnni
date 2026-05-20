from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    database_url: str = "postgresql+asyncpg://calendar:calendar@localhost:5432/treatment_calendar"
    port: int = 8081
    aws_region: str = "us-east-1"
    bedrock_model_id: str = "us.anthropic.claude-opus-4-1-20250805-v1:0"

    model_config = SettingsConfigDict(env_file=".env")


settings = Settings()
