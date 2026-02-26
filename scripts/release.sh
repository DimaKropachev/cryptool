#!/bin/bash
set -e

# Настройки
REPO="DimaKropachev/cryptool"
PROJECT_NAME="cryptool"
VERSION="v1.0.0"
BUILD_DIR="./build/$VERSION"

# Создаем директорию для сборки
mkdir -p "$BUILD_DIR"

# Платформы
PLATFORMS=("windows")
ARCHS=("amd64")

# Сборка
for OS in "${PLATFORMS[@]}"; do
  for ARCH in "${ARCHS[@]}"; do
    BIN_NAME="${PROJECT_NAME}-${OS}-${ARCH}-${VERSION}"

    if [ "$OS" == "windows" ]; then
      BIN_NAME="${BIN_NAME}.exe"
    fi

    echo "Собираем для $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "$BUILD_DIR/$BIN_NAME" ./cmd/cryptool
  done
done

# Проверка существования релиза
if gh release view "$VERSION" --repo "$REPO" > /dev/null 2>&1; then
  echo "Релиз $VERSION уже существует"
else
  echo "Создаю новый релиз $VERSION..."
  gh release create "$VERSION" \
    --repo "$REPO" \
    --title "Release $VERSION" \
    --notes "Автоматический релиз."
fi

# Загрузка файлов
echo "Загружаю бинарники..."
gh release upload "$VERSION" "$BUILD_DIR"/* \
  --repo "$REPO" \
  --clobber

echo "Готово! Релиз $VERSION опубликован."