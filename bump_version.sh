#!/bin/sh

VERSION_FILE="version.txt"
MAIN_GO_FILE="main.go"

if [ ! -f "$VERSION_FILE" ]; then
  echo "Version file not found!"
  exit 1
fi

VERSION_TYPE=$1
MANUAL_VERSION=$2

if [ -z "$VERSION_TYPE" ]; then
  VERSION_TYPE="patch"
fi

if [ -n "$MANUAL_VERSION" ]; then
  NEW_VERSION="$MANUAL_VERSION"
else
  CURRENT_VERSION=$(cat "$VERSION_FILE")
  IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"

  MAJOR=${VERSION_PARTS[0]}
  MINOR=${VERSION_PARTS[1]}
  PATCH=${VERSION_PARTS[2]}

  case "$VERSION_TYPE" in
    major)
      MAJOR=$((MAJOR + 1))
      MINOR=0
      PATCH=0
      ;;
    minor)
      MINOR=$((MINOR + 1))
      PATCH=0
      ;;
    patch)
      PATCH=$((PATCH + 1))
      ;;
    *)
      echo "Invalid version type: $VERSION_TYPE"
      exit 1
      ;;
  esac

  NEW_VERSION="$MAJOR.$MINOR.$PATCH"
fi

echo "$NEW_VERSION" > "$VERSION_FILE"

# Update the appVersion in main.go
sed -i.bak "s/const appVersion string = \".*\"/const appVersion string = \"$NEW_VERSION\"/" "$MAIN_GO_FILE"
rm -f "$MAIN_GO_FILE.bak"

echo "Version bumped to $NEW_VERSION"