#!/bin/sh

# =========================================
# Project Initializer — starter-wahcah-be
# Jalankan sekali setelah clone:
#   sh init-project.sh
# =========================================

OLD_NAME="starter-wahcah-be"

# Ambil nama baru dari argumen atau tanya interaktif
if [ -n "$1" ]; then
    NEW_NAME="$1"
else
    printf "Nama project baru (contoh: my-awesome-api): "
    read NEW_NAME
fi

# Validasi tidak kosong
if [ -z "$NEW_NAME" ]; then
    echo "Error: nama project tidak boleh kosong."
    exit 1
fi

# Validasi format (hanya huruf, angka, dan tanda hubung)
case "$NEW_NAME" in
    *[!a-zA-Z0-9-_]*)
        echo "Error: nama project hanya boleh mengandung huruf, angka, - dan _"
        exit 1
        ;;
esac

echo ""
echo "Mengganti '$OLD_NAME' -> '$NEW_NAME' ..."
echo ""

# Replace di semua file Go, mod, env, yaml, dan markdown
find . \
    -type f \
    \( -name "*.go" -o -name "go.mod" -o -name "*.env" -o -name ".env" -o -name "*.yaml" -o -name "*.yml" -o -name "*.md" \) \
    -not -path "./.git/*" \
    -not -path "./vendor/*" \
    | while read FILE; do
        # Cek apakah file mengandung OLD_NAME sebelum replace
        if grep -q "$OLD_NAME" "$FILE"; then
            sed -i "s|$OLD_NAME|$NEW_NAME|g" "$FILE"
            echo "  updated: $FILE"
        fi
    done

# Rename folder bin jika ada
if [ -d "bin/$OLD_NAME" ]; then
    mv "bin/$OLD_NAME" "bin/$NEW_NAME"
    echo "  renamed: bin/$OLD_NAME -> bin/$NEW_NAME"
fi

echo ""
echo "Selesai! Project berhasil diinisialisasi sebagai '$NEW_NAME'."
echo ""
echo "Langkah selanjutnya:"
echo "  1. go mod tidy"
echo "  2. make dev"