# Настройки
REPO="DimaKropachev/cryptool"    
PROJECT_NAME="cryptool"             
VERSION="v1.0.0-beta"                
BUILD_DIR="./build"           

# Создаем директорию для сборки
mkdir -p "$BUILD_DIR/$VERSION"

# Определяем платформы и архитектуры
PLATFORMS=("windows")
ARCHS=("amd64")                         # Поддерживаем 64-битные платформы

# Собираем бинарники
for OS in "${PLATFORMS[@]}"; do
  for ARCH in "${ARCHS[@]}"; do
    BIN_NAME="${PROJECT_NAME}-${OS}-${ARCH}-${VERSION}"
    if [ "$OS" == "windows" ]; then
      BIN_NAME="${BIN_NAME}.exe"
    fi

    echo "Собираем для $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "$BUILD_DIR/$VERSION/$BIN_NAME" ./cmd/cryptool
  done
done

# Проверка существования релиза (опционально)
if gh release view "$VERSION" > /dev/null 2>&1; then
  echo "Релиз $VERSION уже существует, добавляем файлы..."
else
  echo "Создаю новый релиз $VERSION..."
  gh release create "$VERSION" --title "Release $VERSION" --notes "Автоматический релиз."
fi

# Загружаем файлы в релиз
for FILE in "$BUILD_DIR"/*; do
  echo "Загружаем $FILE..."
  gh release upload "$REPO" "$FILE" --clobber
done

echo "Готово! Релиз $VERSION опубликован."