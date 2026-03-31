#!/usr/bin/env python3
"""
将旧的 config.yaml 配置迁移到新的格式
"""

import os
import yaml
import shutil

def migrate_config():
    old_path = "./data/config.yaml"
    backup_path = "./data/config.yaml.backup"
    
    if not os.path.exists(old_path):
        print("Old config not found at", old_path)
        print("Creating new default config...")
        return
    
    try:
        # 备份旧配置
        shutil.copy2(old_path, backup_path)
        print(f"Backup created: {backup_path}")
        
        # 读取旧配置
        with open(old_path, 'r', encoding='utf-8') as f:
            old_config = yaml.safe_load(f)
        
        print("Old config structure:", old_config.keys())
        
        # 创建新配置结构
        new_config = {
            "web": {
                "host": old_config.get("web", {}).get("host", "127.0.0.1"),
                "port": old_config.get("web", {}).get("port", "8080"),
                "tls": old_config.get("web", {}).get("tls", False),
                "certPrivatePath": old_config.get("web", {}).get("certPrivatePath", ""),
                "certPublicPath": old_config.get("web", {}).get("certPublicPath", ""),
            },
            "database": {
                "type": old_config.get("database", {}).get("type", "sqlite"),
                "path": old_config.get("database", {}).get("path", "data/database.db"),
                "host": old_config.get("database", {}).get("host", ""),
                "port": old_config.get("database", {}).get("port", ""),
                "name": old_config.get("database", {}).get("name", ""),
                "user": old_config.get("database", {}).get("user", ""),
                "pass": old_config.get("database", {}).get("pass", ""),
            },
            "user": {
                "cookieTimeout": old_config.get("user", {}).get("cookieTimeout", 604800),
                "passHashType": old_config.get("user", {}).get("passHashType", "md5"),
            },
            "file": {
                "maxSize": old_config.get("file", {}).get("maxSize", 52428800),
                "paths": old_config.get("file", {}).get("pahts", {
                    "avatar": "data/static/avatar/",
                    "image":  "data/upload/image/",
                    "video":  "data/upload/video/",
                    "music":  "data/upload/music/",
                    "pdf":    "data/upload/pdf/",
                    "other":  "data/upload/other/",
                }),
                "allowImageMime": old_config.get("file", {}).get("allowImageMime", {
                    "image/jpeg": ".jpeg",
                    "image/png": ".png",
                    "image/gif": ".gif",
                    "image/bmp": ".bmp",
                }),
                "allowVideoMime": old_config.get("file", {}).get("allowVideoMime", {
                    "video/mp4": ".mp4",
                    "video/x-msvideo": ".avi",
                    "video/quicktime": ".mov",
                    "video/x-flv": ".flv",
                    "video/mpeg": ".mpeg",
                }),
                "allowMusicMime": old_config.get("file", {}).get("allowMusicMime", {
                    "audio/mpeg": ".mpeg",
                    "audio/aac": ".aac",
                    "audio/wav": ".wav",
                    "audio/flac": ".flac",
                }),
                "allowPdfMime": old_config.get("file", {}).get("allowPdfMime", {
                    "application/pdf": ".pdf",
                }),
            }
        }
        
        # 写入新配置
        with open(old_path, 'w', encoding='utf-8') as f:
            yaml.dump(new_config, f, default_flow_style=False, allow_unicode=True, sort_keys=False)
        
        print("Config migrated successfully!")
        print(f"New config saved to {old_path}")
        
    except Exception as e:
        print(f"Migration failed: {e}")
        if os.path.exists(backup_path):
            print("Restoring backup...")
            shutil.copy2(backup_path, old_path)
            print("Backup restored")

if __name__ == "__main__":
    migrate_config()