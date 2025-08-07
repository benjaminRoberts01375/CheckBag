#!/bin/bash

# Docker Volume Backup/Restore Script for CheckBag

VOLUMES=("valkey-data" "backend-data" "frontend-node-modules")
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKUP_DIR="$SCRIPT_DIR/checkbag-volume-backup"

backup_volumes() {
    echo "Creating backup directory..."
    mkdir -p "$BACKUP_DIR"
    cd "$BACKUP_DIR"

    echo "Backing up Docker volumes..."
    for volume in "${VOLUMES[@]}"; do
        echo "Backing up volume: $volume"
        docker run --rm \
            -v "$volume":/data \
            -v "$(pwd)":/backup \
            alpine tar czf "/backup/${volume}.tar.gz" -C /data .

        if [ $? -eq 0 ]; then
            echo "✓ Successfully backed up $volume"
        else
            echo "✗ Failed to backup $volume"
        fi
    done

    echo "Backup complete! Files are in: $BACKUP_DIR"
    ls -lh *.tar.gz
}

restore_volumes() {
    if [ ! -d "$BACKUP_DIR" ]; then
        echo "Error: Backup directory $BACKUP_DIR not found!"
        exit 1
    fi

    cd "$BACKUP_DIR"

    echo "Restoring Docker volumes..."
    for volume in "${VOLUMES[@]}"; do
        if [ -f "${volume}.tar.gz" ]; then
            echo "Creating volume: $volume"
            docker volume create "$volume" 2>/dev/null

            echo "Restoring volume: $volume"
            docker run --rm \
                -v "$volume":/data \
                -v "$(pwd)":/backup \
                alpine tar xzf "/backup/${volume}.tar.gz" -C /data

            if [ $? -eq 0 ]; then
                echo "✓ Successfully restored $volume"
            else
                echo "✗ Failed to restore $volume"
            fi
        else
            echo "⚠ Backup file ${volume}.tar.gz not found, skipping..."
        fi
    done

    echo "Restore complete!"
}

case "$1" in
    backup)
        backup_volumes
        ;;
    restore)
        restore_volumes
        ;;
    *)
        echo "Usage: $0 {backup|restore}"
        echo "  backup  - Create backups of Docker volumes"
        echo "  restore - Restore Docker volumes from backups"
        exit 1
        ;;
esac
