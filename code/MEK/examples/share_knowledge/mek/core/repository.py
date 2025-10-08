from __future__ import annotations

from typing import Iterable, List, Optional

from .models import MemoryItem
from .types import MemoryType


class MemoryRepository:
    """Abstract repository for memory items.

    Implementations may be in-memory, database-backed, or hybrid. Methods are
    intentionally left unimplemented at this stage per scaffolding request.
    """

    def add(self, item: MemoryItem) -> None:
        raise NotImplementedError

    def get(self, memory_id: str) -> Optional[MemoryItem]:
        raise NotImplementedError

    def update(self, memory_id: str, item: MemoryItem) -> None:
        raise NotImplementedError

    def delete(self, memory_id: str) -> None:
        raise NotImplementedError

    def link_memories(self, source_id: str, target_id: str) -> None:
        raise NotImplementedError

    def list_all(self) -> Iterable[MemoryItem]:
        raise NotImplementedError 