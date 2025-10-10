import datetime
import sys
from pathlib import Path

from loguru import logger

from app.config import PROJECT_ROOT  # 确保这个路径正确

# 定义日志级别的优先级顺序
LOG_LEVELS = ["DEBUG", "INFO", "WARNING", "ERROR"]


def configure_logger(
        console_level="INFO",
        file_level="DEBUG",
        log_dir: Path = None,
        max_size: str = "20 MB",  # 暂未支持，仅支持每天分割
        retention: str = "7 days",
):
    """
    配置 Loguru 日志系统，支持：
    - 控制台日志级别
    - 文件日志级别
    - 按天 + 按大小自动分割日志
    - 日志保留时间

    :param console_level: 控制台输出的日志级别
    :param file_level: 写入文件的日志级别
    :param log_dir: 日志文件保存目录
    :param log_name: 日志文件前缀名（如 app）
    :param max_size: 单个日志文件最大大小（如 "20 MB"）
    :param retention: 日志保留时间（如 "7 days"）
    """
    log_dir = log_dir or PROJECT_ROOT / "logs"
    log_dir.mkdir(parents=True, exist_ok=True)

    logger.remove()  # 移除默认的 logger

    # 控制台输出
    logger.add(
        sys.stderr,
        level=console_level,
        format="<green>{time:YYYY-MM-DD HH:mm:ss}</green> | <level>{level}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>",
        colorize=True,
    )

    # 文件输出，按天 + 按大小自动分割
    def add_file_log(log_level: str):
        logger.add(
            f"{log_dir}/{{time:YYYYMMDD}}/{log_level.lower()}.log",  # 动态文件名
            level=log_level,
            format="{time:YYYY-MM-DD HH:mm:ss} | {level} | {name}:{function}:{line} - {message}",
            # rotation=max_size,  # 每 20MB 分割
            rotation=datetime.time(hour=0, minute=0, second=0),  # 每天凌晨00:00:00分割
            retention=retention,  # 保留天数
            # compression="gz",  # 可选，压缩旧日志文件为： "zip", "gz", "bz2", "xz", "lz4", "zst"
            enqueue=True,  # 多线程安全
        )

    # 获取用户设置的 level 在优先级列表中的位置
    try:
        level_index = LOG_LEVELS.index(file_level.upper())
    except ValueError:
        raise ValueError(f"Invalid file_level: {file_level}. Must be one of {LOG_LEVELS}")

    # 所有需要写入文件的 level
    selected_levels = LOG_LEVELS[level_index:]
    for level in selected_levels:
        add_file_log(level)

    return logger


# 初始化日志配置
logger = configure_logger(
    console_level="INFO",
    file_level="DEBUG",
    log_dir=PROJECT_ROOT / "logs",
    max_size="50 MB",
    retention="7 days",
)
