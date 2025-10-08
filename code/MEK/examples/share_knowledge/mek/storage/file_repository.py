from __future__ import annotations

import json
from datetime import datetime
from pathlib import Path
from typing import Iterable, List, Optional

from ..core.models import MemoryItem
from ..core.repository import MemoryRepository
from ..core.types import MemoryType
from ..utils.file_util import get_project_file_path, ensure_dir_exists


class FileMemoryRepository(MemoryRepository):
    def __init__(self, base_dir: str | None = None) -> None:
        if base_dir is None:
            base_dir = get_project_file_path("data", "memory")
        ensure_dir_exists(base_dir)
        self.base_dir = Path(base_dir)
    
    @property
    def storage_path(self) -> str:
        return str(self.base_dir)
    
    @storage_path.setter
    def storage_path(self, path: str) -> None:
        ensure_dir_exists(path)
        self.base_dir = Path(path)

    def _path_for(self, memory_id: str) -> Path:
        return self.base_dir / f"{memory_id}.json"

    def _serialize(self, item: MemoryItem) -> dict:
        return {
            "memory_id": item.memory_id,
            "memory_type": item.memory_type.value,
            "created_at": item.created_at.isoformat(),
            "content": item.content,
            "importance": item.importance.value,
            "related_memory_ids": list(item.related_memory_ids),
        }

    def _deserialize(self, data: dict) -> MemoryItem:
        return MemoryItem(
            memory_id=str(data["memory_id"]),
            memory_type=MemoryType(data["memory_type"]),
            created_at=datetime.fromisoformat(data["created_at"]),
            content=dict(data.get("content", {})),
            importance=self._parse_importance(data.get("importance")),
            related_memory_ids=list(data.get("related_memory_ids", [])),
        )

    def _parse_importance(self, value: str):
        from ..core.types import ImportanceLevel

        try:
            return ImportanceLevel(value)
        except Exception:
            return ImportanceLevel.MEDIUM

    def add(self, item: MemoryItem) -> None:
        with self._path_for(item.memory_id).open("w", encoding="utf-8") as f:
            json.dump(self._serialize(item), f, ensure_ascii=False)

    def get(self, memory_id: str) -> Optional[MemoryItem]:
        p = self._path_for(memory_id)
        if not p.exists():
            return None
        with p.open("r", encoding="utf-8") as f:
            data = json.load(f)
        return self._deserialize(data)

    def update(self, memory_id: str, item: MemoryItem) -> None:
        if not self._path_for(memory_id).exists():
            raise KeyError(memory_id)
        self.add(item)

    def delete(self, memory_id: str) -> None:
        p = self._path_for(memory_id)
        if p.exists():
            p.unlink()

    def link_memories(self, source_id: str, target_id: str) -> None:
        source = self.get(source_id)
        target = self.get(target_id)
        if not source or not target:
            return
        if target_id not in source.related_memory_ids:
            source.related_memory_ids.append(target_id)
            self.update(source.memory_id, source)
        if source_id not in target.related_memory_ids:
            target.related_memory_ids.append(source_id)
            self.update(target.memory_id, target)

    def list_all(self) -> Iterable[MemoryItem]:
        for p in self.base_dir.glob("*.json"):
            try:
                with p.open("r", encoding="utf-8") as f:
                    data = json.load(f)
                yield self._deserialize(data)
            except Exception:
                continue 