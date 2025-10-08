from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime
from enum import Enum
from typing import Optional


class RetrievalType(str, Enum):
    """检索类型枚举"""
    EXACT = "exact"
    SEMANTIC = "semantic"
    TEMPORAL = "temporal"
    ASSOCIATIVE = "associative"


@dataclass
class RetrievalOptions:
    """检索选项配置"""
    limit: int = 20
    importance_weight: float = 1.0
    since: Optional[datetime] = None
    until: Optional[datetime] = None 