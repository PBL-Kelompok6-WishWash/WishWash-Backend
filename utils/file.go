package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Slugify cleans the string to be filename-friendly.
func Slugify(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	// Return a default if empty
	res := result.String()
	if res == "" {
		return "file"
	}
	return res
}

// BuildEntityFolder creates and returns the per-entity subfolder name.
// Format: "{id}_{slugified_name}" e.g. "1_budi_santoso"
// The id parameter should be the entity's primary key (e.g. id_pelanggan).
func BuildEntityFolder(entityID uint, entityName string) string {
	return fmt.Sprintf("%d_%s", entityID, Slugify(entityName))
}

// SaveBase64Image decodes a base64 image string and saves it to a per-entity subfolder.
// Structure: uploads/{category}/{entityID}_{entityName}/{filename}.{ext}
// Example:   uploads/pelanggan/1_budi_santoso/profile_pelanggan_budi_santoso.png
//
// Parameters:
//   - base64Str: the base64 string (or empty/existing path)
//   - category: top-level subfolder ("pelanggan", "karyawan", "layanan", "promo", "metode_bayar")
//   - entityFolder: the per-entity folder name from BuildEntityFolder()
//   - filenameWithoutExt: desired filename without extension
//
// If the input is not a base64 string (empty or already a path/URL), it returns it as-is.
func SaveBase64Image(base64Str string, category string, entityFolder string, filenameWithoutExt string) (string, error) {
	if base64Str == "" {
		return "", nil
	}

	// If it is already a URL or relative path, return it as is
	if strings.HasPrefix(base64Str, "/uploads/") || strings.HasPrefix(base64Str, "http://") || strings.HasPrefix(base64Str, "https://") {
		return base64Str, nil
	}

	// 1. Process base64 format (handles prefixes like "data:image/png;base64,")
	ext := ".png" // Default extension
	dataParts := strings.Split(base64Str, ",")
	rawBase64 := base64Str

	if len(dataParts) > 1 {
		prefix := dataParts[0]
		rawBase64 = dataParts[1]

		if strings.Contains(prefix, "image/jpeg") || strings.Contains(prefix, "image/jpg") {
			ext = ".jpg"
		} else if strings.Contains(prefix, "image/png") {
			ext = ".png"
		} else if strings.Contains(prefix, "image/gif") {
			ext = ".gif"
		} else if strings.Contains(prefix, "image/webp") {
			ext = ".webp"
		} else if strings.Contains(prefix, "image/svg+xml") {
			ext = ".svg"
		}
	}

	// Standard cleanups for base64 strings (remove spaces, newlines, tabs)
	rawBase64 = strings.ReplaceAll(rawBase64, " ", "")
	rawBase64 = strings.ReplaceAll(rawBase64, "\n", "")
	rawBase64 = strings.ReplaceAll(rawBase64, "\r", "")
	rawBase64 = strings.ReplaceAll(rawBase64, "\t", "")

	// 2. Decode the base64 string
	dec, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		// Fallback to RawStdEncoding (without padding)
		dec, err = base64.RawStdEncoding.DecodeString(rawBase64)
		if err != nil {
			return "", fmt.Errorf("gagal mendekode string base64: %v", err)
		}
	}

	// 3. Prepare target directory: uploads/{category}/{entityFolder}/
	uploadDir := filepath.Join("uploads", category, entityFolder)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("gagal membuat direktori: %v", err)
	}

	// 4. Clean filename and generate path
	cleanFilename := Slugify(filenameWithoutExt) + ext
	filePath := filepath.Join(uploadDir, cleanFilename)

	// 5. Write the file to disk (overwriting if it already exists)
	if err := os.WriteFile(filePath, dec, 0644); err != nil {
		return "", fmt.Errorf("gagal menyimpan gambar ke disk: %v", err)
	}

	// 6. Return relative URL path using forward slashes (cross-platform compatible)
	return fmt.Sprintf("/uploads/%s/%s/%s", category, entityFolder, cleanFilename), nil
}

// DeleteImageFile deletes an image file from the disk given its relative URL path.
func DeleteImageFile(relativeUrl string) error {
	if relativeUrl == "" {
		return nil
	}

	// Only process files inside /uploads/
	if !strings.HasPrefix(relativeUrl, "/uploads/") {
		return nil
	}

	// Convert URL path to local OS file path (e.g. "/uploads/pelanggan/1_budi/x.png" -> "uploads/pelanggan/1_budi/x.png")
	localPath := filepath.Clean(strings.TrimPrefix(relativeUrl, "/"))

	// Check if file exists
	info, err := os.Stat(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Already deleted or doesn't exist, which is fine
			return nil
		}
		return err
	}

	// Ensure it is not a directory
	if info.IsDir() {
		return fmt.Errorf("path adalah direktori, bukan file: %s", localPath)
	}

	// Delete file from disk
	if err := os.Remove(localPath); err != nil {
		return fmt.Errorf("gagal menghapus file: %v", err)
	}

	return nil
}

// DeleteImageFolder deletes the entire entity subfolder from disk.
// This is useful when deleting an entity entirely (e.g. deleting a pelanggan removes their folder).
// Path format: uploads/{category}/{entityFolder}/
func DeleteImageFolder(category string, entityFolder string) error {
	if category == "" || entityFolder == "" {
		return nil
	}

	folderPath := filepath.Join("uploads", category, entityFolder)

	// Check if folder exists
	info, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Already gone
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("diharapkan direktori tetapi mendapatkan file: %s", folderPath)
	}

	// Remove the entire folder and all its contents
	return os.RemoveAll(folderPath)
}
