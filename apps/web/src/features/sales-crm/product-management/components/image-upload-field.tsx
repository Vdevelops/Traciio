"use client";

import { useState, useRef, useEffect } from "react";
import { Upload, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Field, FieldLabel, FieldError } from "@/components/ui/field";
import { toast } from "sonner";
import { useTranslations } from "next-intl";
import apiClient from "@/lib/api-client";
import { ProductImage } from "./product-image";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

function normalizeUploadedUrl(url: string): string {
  // API may return a relative path like "/uploads/...".
  // The form schema validates URLs using `new URL(...)`, so ensure an absolute URL.
  if (url.startsWith("/uploads/")) {
    return `${API_BASE_URL}${url}`;
  }
  return url;
}

interface ImageUploadFieldProps {
  readonly value?: string;
  readonly onChange: (url: string) => void;
  readonly disabled?: boolean;
  readonly error?: string;
}

export function ImageUploadField({ value, onChange, disabled, error }: ImageUploadFieldProps) {
  const [preview, setPreview] = useState<string | undefined>(value);
  const [isUploading, setIsUploading] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const t = useTranslations("productManagement.form");

  // Sync preview when value changes from parent
  useEffect(() => {
    setPreview(value);
  }, [value]);

  const handleFileSelect = async (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Validate file type
    if (!file.type.startsWith("image/")) {
      toast.error(t("imageErrors.selectImageFile"));
      return;
    }

    // Validate file size (max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      toast.error(t("imageErrors.fileTooLarge"));
      return;
    }

    setIsUploading(true);

    try {
      // Create FormData and upload to server
      const formData = new FormData();
      formData.append("image", file);

      const response = await apiClient.post<{ success: boolean; data: { url: string } }>(
        "/upload/image",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );

      const uploadedUrl = response.data.data.url;
      const normalizedUrl = normalizeUploadedUrl(uploadedUrl);
      setPreview(normalizedUrl);
      onChange(normalizedUrl);
      toast.success(t("imageSuccess.uploaded"));
    } catch (error) {
      toast.error(t("imageErrors.uploadFailed"));
    } finally {
      setIsUploading(false);
    }
  };

  const handleRemoveImage = () => {
    setPreview(undefined);
    onChange("");
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  };

  return (
    <Field orientation="vertical">
      <FieldLabel>{t("imageLabel")}</FieldLabel>
      
      {preview ? (
        <div className="relative w-full">
          <ProductImage
            src={preview}
            alt="Product preview"
            className="w-full h-48 rounded-md"
            fallbackClassName="border"
            loading="eager"
          />
          {!disabled && (
            <Button
              type="button"
              variant="destructive"
              size="icon"
              className="absolute top-2 right-2"
              onClick={handleRemoveImage}
              disabled={isUploading}
            >
              <X className="h-4 w-4" />
            </Button>
          )}
        </div>
      ) : (
        <div className="flex items-center gap-2">
          <Input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={handleFileSelect}
            className="flex-1"
            disabled={disabled || isUploading}
          />
          <Button
            type="button"
            variant="outline"
            size="icon"
            onClick={() => fileInputRef.current?.click()}
            disabled={disabled || isUploading}
          >
            {isUploading ? (
              <div className="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
            ) : (
              <Upload className="h-4 w-4" />
            )}
          </Button>
        </div>
      )}

      {error && <FieldError>{error}</FieldError>}
      
      {!preview && !error && (
        <p className="text-xs text-muted-foreground">
          {t("imageHint")}
        </p>
      )}
    </Field>
  );
}
