from __future__ import annotations

import os
import yaml
from dataclasses import dataclass
from typing import Dict, Any, Optional
from pathlib import Path


@dataclass
class LLMConfig:
    """LLM配置"""
    provider: str
    model_name: str
    base_url: str
    api_key: str


@dataclass 
class EmbeddingConfig:
    """Embedding配置"""
    provider: str
    model_name: str
    base_url: str
    api_key: str
    embedding_dim: int = 1024


@dataclass
class MEKConfig:
    """MEK系统配置"""
    memory_storage_path: str = "data/memory"
    default_retrieval_limit: int = 20
    default_importance_weight: float = 1.0


class ConfigManager:
    """配置管理器"""
    
    def __init__(self, config_path: Optional[str] = None):
        self.config_path = config_path or self._find_config_file()
        self._config_data = self._load_config()
        
    def _find_config_file(self) -> str:
        """查找配置文件"""
        possible_paths = [
            "config.yaml",
            "config.yml", 
            os.path.expanduser("~/.mek/config.yaml"),
            "/etc/mek/config.yaml"
        ]
        
        for path in possible_paths:
            if os.path.exists(path):
                return path
                
        # 如果找不到配置文件，返回默认路径
        return "config.yaml"
    
    def _load_config(self) -> Dict[str, Any]:
        """加载配置文件"""
        if not os.path.exists(self.config_path):
            raise FileNotFoundError(f"配置文件不存在: {self.config_path}")
            
        with open(self.config_path, 'r', encoding='utf-8') as f:
            return yaml.safe_load(f)
    
    def get_llm_config(self) -> LLMConfig:
        """获取LLM配置"""
        llm_config = self._config_data.get('llm', {})
        if not llm_config:
            raise ValueError("配置文件中缺少LLM配置")
        
        required_fields = ['provider', 'model_name', 'base_url', 'api_key']
        for field in required_fields:
            if not llm_config.get(field):
                raise ValueError(f"LLM配置中缺少必需字段: {field}")
        
        return LLMConfig(
            provider=llm_config['provider'],
            model_name=llm_config['model_name'],
            base_url=llm_config['base_url'],
            api_key=llm_config['api_key']
        )
    
    def get_embedding_config(self) -> EmbeddingConfig:
        """获取Embedding配置"""
        embedding_config = self._config_data.get('embedding', {})
        if not embedding_config:
            raise ValueError("配置文件中缺少Embedding配置")
        
        required_fields = ['provider', 'model_name', 'base_url', 'api_key']
        for field in required_fields:
            if not embedding_config.get(field):
                raise ValueError(f"Embedding配置中缺少必需字段: {field}")
        
        return EmbeddingConfig(
            provider=embedding_config['provider'],
            model_name=embedding_config['model_name'],
            base_url=embedding_config['base_url'],
            api_key=embedding_config['api_key'],
            embedding_dim=embedding_config.get('embedding_dim', 1024)  # 保留这个默认值，因为通常是固定的
        )
    
    def get_mek_config(self) -> MEKConfig:
        """获取MEK系统配置"""
        mek_config = self._config_data.get('mek', {})
        return MEKConfig(
            memory_storage_path=mek_config.get('memory_storage_path', 'data/memory'),
            default_retrieval_limit=mek_config.get('default_retrieval_limit', 20),
            default_importance_weight=mek_config.get('default_importance_weight', 1.0)
        )


# 全局配置管理器实例
_config_manager: Optional[ConfigManager] = None


def get_config_manager(config_path: Optional[str] = None) -> ConfigManager:
    """获取配置管理器实例"""
    global _config_manager
    if _config_manager is None:
        _config_manager = ConfigManager(config_path)
    return _config_manager


def reload_config(config_path: Optional[str] = None) -> ConfigManager:
    """重新加载配置"""
    global _config_manager
    _config_manager = ConfigManager(config_path)
    return _config_manager 