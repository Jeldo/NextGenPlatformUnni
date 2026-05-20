from datetime import datetime, timezone

from sqlalchemy import DateTime, ForeignKey, Integer, String, func
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship


class Base(DeclarativeBase):
    pass


class Category(Base):
    __tablename__ = "categories"

    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    name: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)

    treatments: Mapped[list["Treatment"]] = relationship(
        back_populates="category", cascade="all, delete-orphan"
    )


class Treatment(Base):
    __tablename__ = "treatments"

    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    category_id: Mapped[str] = mapped_column(ForeignKey("categories.id", ondelete="CASCADE"), nullable=False)
    name: Mapped[str] = mapped_column(String(200), nullable=False)

    category: Mapped["Category"] = relationship(back_populates="treatments")
    dosage_types: Mapped[list["DosageType"]] = relationship(
        back_populates="treatment", cascade="all, delete-orphan"
    )
    cycle_rule: Mapped["CycleRule | None"] = relationship(
        back_populates="treatment", cascade="all, delete-orphan", uselist=False
    )


class DosageType(Base):
    __tablename__ = "dosage_types"

    id: Mapped[str] = mapped_column(String(36), primary_key=True)
    treatment_id: Mapped[str] = mapped_column(ForeignKey("treatments.id", ondelete="CASCADE"), nullable=False)
    unit: Mapped[str] = mapped_column(String(20), nullable=False)

    treatment: Mapped["Treatment"] = relationship(back_populates="dosage_types")


class CycleRule(Base):
    __tablename__ = "cycle_rules"

    treatment_id: Mapped[str] = mapped_column(
        ForeignKey("treatments.id", ondelete="CASCADE"), primary_key=True
    )
    cycle_days: Mapped[int] = mapped_column(Integer, nullable=False)
    description: Mapped[str | None] = mapped_column(String(500))
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        default=lambda: datetime.now(timezone.utc),
        onupdate=lambda: datetime.now(timezone.utc),
        server_default=func.now(),
    )

    treatment: Mapped["Treatment"] = relationship(back_populates="cycle_rule")
