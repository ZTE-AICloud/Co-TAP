from __future__ import annotations

from enum import Enum


class MemoryType(str, Enum):
    PROFILE = "profile"
    WORKFLOW = "workflow"
    SEMANTIC = "semantic" 
    EPISODE = "episode"


class ImportanceLevel(str, Enum):
    HIGH = "high"
    MEDIUM = "medium"
    LOW = "low" 