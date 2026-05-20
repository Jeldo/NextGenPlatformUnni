"""initial schema

Revision ID: 001
Revises:
Create Date: 2026-05-20

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa

revision: str = "001"
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    op.create_table(
        "categories",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("name", sa.String(100), unique=True, nullable=False),
    )
    op.create_table(
        "treatments",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("category_id", sa.String(36), sa.ForeignKey("categories.id", ondelete="CASCADE"), nullable=False),
        sa.Column("name", sa.String(200), nullable=False),
    )
    op.create_table(
        "dosage_types",
        sa.Column("id", sa.String(36), primary_key=True),
        sa.Column("treatment_id", sa.String(36), sa.ForeignKey("treatments.id", ondelete="CASCADE"), nullable=False),
        sa.Column("unit", sa.String(20), nullable=False),
    )
    op.create_table(
        "cycle_rules",
        sa.Column("treatment_id", sa.String(36), sa.ForeignKey("treatments.id", ondelete="CASCADE"), primary_key=True),
        sa.Column("cycle_days", sa.Integer, nullable=False),
        sa.Column("description", sa.String(500)),
        sa.Column("updated_at", sa.DateTime(timezone=True), server_default=sa.func.now()),
    )


def downgrade() -> None:
    op.drop_table("cycle_rules")
    op.drop_table("dosage_types")
    op.drop_table("treatments")
    op.drop_table("categories")
