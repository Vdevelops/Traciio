"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import type { LeadSource, CreateLeadSourceRequest, UpdateLeadSourceRequest } from "../types/lead-source";

interface LeadSourceFormProps {
  initialData?: LeadSource;
  onSubmit: (data: CreateLeadSourceRequest | UpdateLeadSourceRequest) => void;
  isSubmitting?: boolean;
  onCancel?: () => void;
}

export function LeadSourceForm({ initialData, onSubmit, isSubmitting, onCancel }: LeadSourceFormProps) {
  const [formData, setFormData] = useState({
    name: initialData?.name || "",
    code: initialData?.code || "",
    description: initialData?.description || "",
    order: initialData?.order || 0,
    is_active: initialData?.is_active ?? true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const validateForm = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.name.trim()) {
      newErrors.name = "Name is required";
    } else if (formData.name.length > 100) {
      newErrors.name = "Name must be at most 100 characters";
    }

    if (!formData.code.trim()) {
      newErrors.code = "Code is required";
    } else if (!/^[A-Z0-9_]+$/.test(formData.code)) {
      newErrors.code = "Code must be uppercase letters, numbers, and underscores only";
    } else if (formData.code.length > 50) {
      newErrors.code = "Code must be at most 50 characters";
    }

    if (formData.order < 0) {
      newErrors.order = "Order must be at least 0";
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    onSubmit(formData);
  };

  const handleChange = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }));
    // Clear error for this field
    if (errors[field]) {
      setErrors(prev => {
        const newErrors = { ...prev };
        delete newErrors[field];
        return newErrors;
      });
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="name">Name *</Label>
        <Input
          id="name"
          placeholder="Website"
          value={formData.name}
          onChange={(e) => handleChange("name", e.target.value)}
        />
        {errors.name && <p className="text-sm text-red-500">{errors.name}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="code">Code *</Label>
        <Input
          id="code"
          placeholder="WEBSITE"
          value={formData.code}
          onChange={(e) => handleChange("code", e.target.value.toUpperCase())}
        />
        <p className="text-xs text-muted-foreground">
          Unique identifier (uppercase letters, numbers, and underscores only)
        </p>
        {errors.code && <p className="text-sm text-red-500">{errors.code}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="description">Description</Label>
        <Textarea
          id="description"
          placeholder="Description of this lead source"
          value={formData.description}
          onChange={(e) => handleChange("description", e.target.value)}
          rows={3}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="order">Order *</Label>
        <Input
          id="order"
          type="number"
          min={0}
          value={formData.order}
          onChange={(e) => handleChange("order", parseInt(e.target.value) || 0)}
        />
        <p className="text-xs text-muted-foreground">Display order</p>
        {errors.order && <p className="text-sm text-red-500">{errors.order}</p>}
      </div>

      <div className="flex items-center space-x-3 rounded-md border p-4">
        <Checkbox
          id="is_active"
          checked={formData.is_active}
          onCheckedChange={(checked) => handleChange("is_active", checked === true)}
        />
        <div className="space-y-1 leading-none">
          <Label htmlFor="is_active" className="cursor-pointer">
            Active Status
          </Label>
          <p className="text-xs text-muted-foreground">
            Active lead sources are available for selection in lead forms
          </p>
        </div>
      </div>

      <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 pt-4">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting} className="cursor-pointer w-full sm:w-auto">
            Cancel
          </Button>
        )}
        <Button type="submit" disabled={isSubmitting} className="cursor-pointer w-full sm:w-auto">
          {isSubmitting ? "Saving..." : initialData ? "Update Lead Source" : "Create Lead Source"}
        </Button>
      </div>
    </form>
  );
}

