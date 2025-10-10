import os
import threading
import tomllib
from pathlib import Path
from typing import Dict

from pydantic import BaseModel, Field

from app.utils.logger import logger


def get_project_root() -> Path:
    """Get the project root directory"""
    return Path(__file__).resolve().parent.parent


PROJECT_ROOT = get_project_root()


class LLMSettings(BaseModel):
    model: str = Field(..., description="Model name")
    base_url: str = Field(..., description="API base URL")
    api_key: str = Field(..., description="API key")
    max_tokens: int = Field(4096, description="Maximum number of tokens per request")
    temperature: float = Field(1.0, description="Sampling temperature")
    api_type: str = Field(..., description="AzureOpenai or Openai")
    api_version: str = Field(..., description="Azure Openai version if AzureOpenai")
    http_proxy: str = Field(..., description="Proxy URL")
    is_internal: bool = Field(False, description="Whether this is an internal model (no proxy needed)")


class AppConfig(BaseModel):
    llm: Dict[str, LLMSettings]


class Config:
    _instance = None
    _lock = threading.Lock()
    _initialized = False

    def __new__(cls):
        if cls._instance is None:
            with cls._lock:
                if cls._instance is None:
                    cls._instance = super().__new__(cls)
        return cls._instance

    def __init__(self):
        # 每次都重新加载配置，以支持动态切换模式
        with self._lock:
            self._config = None
            self._load_initial_config()
            self._initialized = True

    @staticmethod
    def _get_config_path() -> Path:
        root = PROJECT_ROOT
        config_path = root / "config" / "config.toml"
        if config_path.exists():
            return config_path
        example_path = root / "config" / "config.example.toml"
        if example_path.exists():
            return example_path
        raise FileNotFoundError("No configuration file found in config directory")

    def _load_config(self) -> dict:
        config_path = self._get_config_path()
        with config_path.open("rb") as f:
            return tomllib.load(f)

    def _load_initial_config(self):
        raw_config = self._load_config()

        # 检查是否为开发模式
        is_dev_mode = os.getenv('ZUI_DEV_MODE', 'false').lower() == 'true'
        mode = 'development' if is_dev_mode else 'production'
        logger.debug(f"mode: {mode}")

        # 获取对应模式的配置
        llm_config = raw_config.get("llm", {})
        mode_config = llm_config.get(mode, {})
        vision_mode_config = llm_config.get("vision", {}).get(mode, {})

        # 构建默认配置
        default_settings = {
            "model": mode_config.get("model"),
            "base_url": mode_config.get("base_url"),
            "api_key": mode_config.get("api_key"),
            "max_tokens": mode_config.get("max_tokens", 4096),
            "temperature": mode_config.get("temperature", 1.0),
            "api_type": mode_config.get("api_type", "Openai"),
            "api_version": mode_config.get("api_version", ""),
            "http_proxy": mode_config.get("http_proxy", ""),
            "is_internal": mode_config.get("is_internal", False)
        }

        # 构建vision配置
        vision_settings = {
            "model": vision_mode_config.get("model"),
            "base_url": vision_mode_config.get("base_url"),
            "api_key": vision_mode_config.get("api_key"),
            "max_tokens": vision_mode_config.get("max_tokens", 4096),
            "temperature": vision_mode_config.get("temperature", 1.0),
            "api_type": vision_mode_config.get("api_type", "Openai"),
            "api_version": vision_mode_config.get("api_version", ""),
            "http_proxy": vision_mode_config.get("http_proxy", ""),
            "is_internal": vision_mode_config.get("is_internal", False)
        }

        config_dict = {
            "llm": {
                "default": default_settings,
                "vision": vision_settings,
            }
        }

        self._config = AppConfig(**config_dict)

    @property
    def llm(self) -> Dict[str, LLMSettings]:
        return self._config.llm


config = Config()
