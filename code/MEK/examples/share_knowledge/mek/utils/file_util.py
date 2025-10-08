import json
from collections import defaultdict
from pathlib import Path
from typing import Any

from pydantic import BaseModel


class CustomEncoder(json.JSONEncoder):
    def default(self, o: Any) -> Any:
        # 如果对象是 Pydantic 模型，则调用其 .model_dump() 方法
        if isinstance(o, BaseModel):
            # 这里的 o.model_dump() 会正确调用你自定义的方法
            return o.model_dump()
        # 如果对象是 defaultdict，则转换为普通字典
        if isinstance(o, defaultdict):
            return dict(o)
        # 对于其他类型，使用默认的编码器
        return super().default(o)


def get_root_path() -> str:
    """
    获取项目根目录的绝对路径
    
    基于file_util.py的固定位置计算项目根目录。
    file_util.py位于: {PROJECT_ROOT}/mek/utils/file_util.py
    因此项目根目录 = file_util.py的父目录的父目录的父目录
    
    Returns:
        str: 项目根目录的绝对路径
    """
    # file_util.py 的绝对路径
    current_file = Path(__file__).resolve()

    # 项目根目录 = file_util.py 的父目录的父目录的父目录
    # file_util.py -> utils -> mek -> PROJECT_ROOT
    root_path = current_file.parent.parent.parent

    return str(root_path)


def get_relative_path_from_root(target_path: str) -> str:
    """
    获取相对于项目根目录的路径
    
    Args:
        target_path: 目标路径（可以是绝对路径或相对路径）
        
    Returns:
        str: 相对于项目根目录的路径
    """
    root_path = get_root_path()
    target = Path(target_path)

    # 如果是相对路径，先转换为绝对路径
    if not target.is_absolute():
        target = Path(root_path) / target

    # 计算相对路径
    try:
        return str(target.relative_to(root_path))
    except ValueError:
        # 如果目标路径不在项目根目录下，返回绝对路径
        return str(target)


def ensure_dir_exists(dir_path: str) -> str:
    """
    确保目录存在，如果不存在则创建
    
    Args:
        dir_path: 目录路径（相对于项目根目录或绝对路径）
        
    Returns:
        str: 目录的绝对路径
    """
    path = Path(dir_path)

    # 如果是相对路径，基于项目根目录解析
    if not path.is_absolute():
        root_path = get_root_path()
        path = Path(root_path) / path

    path.mkdir(parents=True, exist_ok=True)
    return str(path.resolve())


def get_project_file_path(*path_parts: str) -> str:
    """
    获取项目内文件的绝对路径
    
    Args:
        *path_parts: 相对于项目根目录的路径部分
        
    Returns:
        str: 文件的绝对路径
        
    Example:
        get_project_file_path("task", "web_walker", "log")
        # 返回: /path/to/project/task/web_walker/log
    """
    root_path = get_root_path()
    return str(Path(root_path).joinpath(*path_parts))
