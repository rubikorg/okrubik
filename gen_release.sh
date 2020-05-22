#!/bin/sh

current_version="0.1.0"
base=$(pwd)
release_folder="$base/releases/$current_version"
target_folder="$base/target"

# gen darwin release
# recreate the folder
mkdir -p "$release_folder/darwin"
mkdir -p "$release_folder/linux"
echo "Creating target ..."
mkdir -p $target_folder
rm $(find $target_folder -name "*-$current_version.zip")

echo "Building for OSX ..."
env GOOS=darwin GOARCH=amd64 go build -o "$release_folder/darwin/Darwin-amd64" cmd/okrubik/main.go
echo "Building for Linux ..."
env GOOS=linux GOARCH=amd64 go build -o "$release_folder/linux/Linux-amd64" cmd/okrubik/main.go

echo "Archiving OSX ..."
cd "$release_folder/darwin"
zip -r "$target_folder/darwin-amd64-$current_version.zip" .

echo "Archiving Linux ..."
cd "$release_folder/linux"
zip -r "$target_folder/linux-amd64-$current_version.zip" .

echo "Final cleanup ..."
rm -rf "$base/releases/"

echo "Done. Check target folder for binary release."