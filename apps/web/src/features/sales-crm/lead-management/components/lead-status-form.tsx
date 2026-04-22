"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import type { LeadStatus, CreateLeadStatusRequest, UpdateLeadStatusRequest } from "../types/lead-status";

interface LeadStatusFormProps {
  initialData?: LeadStatus;
  onSubmit: (data: CreateLeadStatusRequest | UpdateLeadStatusRequest) => void;
  isSubmitting?: boolean;
  onCancel?: () => void;
}

export function LeadStatusForm({ initialData, onSubmit, isSubmitting, onCancel }: LeadStatusFormProps) {
  const [formData, setFormData] = useState({
    name: initialData?.name || "",
    code: initialData?.code || "",
    description: initialData?.description || "",
    score: initialData?.score || 0,
    color: initialData?.color || "#6366f1",
    order: initialData?.order || 0,
    is_active: initialData?.is_active ?? true,
  });

  const [errors, setErrors] = useState<Record<string, string>>({});

  const getScoreColor = (score: number) => {
    if (score >= 60) return "text-green-600";
    if (score >= 30) return "text-yellow-600";
    return "text-red-600";
  };

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

    if (formData.score < 0 || formData.score > 100) {
      newErrors.score = "Score must be between 0 and 100";
    }

    if (!/^#[0-9A-Fa-f]{6}$/.test(formData.color)) {
      newErrors.color = "Color must be a valid hex code (e.g., #6366f1)";
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
          placeholder="New Lead"
          value={formData.name}
          onChange={(e) => handleChange("name", e.target.value)}
        />
        {errors.name && <p className="text-sm text-red-500">{errors.name}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="code">Code *</Label>
        <Input
          id="code"
          placeholder="NEW_LEAD"
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
          placeholder="Description of this status"
          value={formData.description}
          onChange={(e) => handleChange("description", e.target.value)}
          rows={3}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="score">Score *</Label>
        <div className="space-y-2">
          <div className="flex items-center gap-4">
            <Input
              id="score"
              type="number"
              min={0}
              max={100}
              step={5}
              value={formData.score}
              onChange={(e) => handleChange("score", parseInt(e.target.value) || 0)}
              className="w-24"
            />
            <span className={`text-lg font-medium ${getScoreColor(formData.score)}`}>
              {formData.score}%
            </span>
            <input
              type="range"
              min={0}
              max={100}
              step={5}
              value={formData.score}
              onChange={(e) => handleChange("score", parseInt(e.target.value))}
              className="flex-1"
            />
          </div>
          <p className="text-xs text-muted-foreground">
            Score percentage for this status (0-100%)
          </p>
        </div>
        {errors.score && <p className="text-sm text-red-500">{errors.score}</p>}
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="color">Color *</Label>
          <div className="flex gap-2">
            <Input
              type="color"
              value={formData.color}
              onChange={(e) => handleChange("color", e.target.value)}
              className="w-14 h-10 p-1 cursor-pointer shrink-0"
            />
            <Input
              id="color"
              type="text"
              placeholder="#6366f1"
              value={formData.color}
              onChange={(e) => handleChange("color", e.target.value)}
              className="flex-1"
            />
          </div>
          {errors.color && <p className="text-sm text-red-500">{errors.color}</p>}
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
            Active statuses are available for selection in lead forms
          </p>
        </div>
      </div>

      <div className="flex flex-col-reverse sm:flex-row justify-end gap-2 pt-4">
        {onCancel && (
          <Button type="button" variant="outline" onClick={onCancel} disabled={isSubmitting} className="w-full sm:w-auto cursor-pointer">
            Cancel
          </Button>
        )}
        <Button type="submit" disabled={isSubmitting} className="w-full sm:w-auto cursor-pointer">
          {isSubmitting ? "Saving..." : initialData ? "Update Status" : "Create Status"}
        </Button>
      </div>
    </form>
  );
}
