"use client";

import { useState } from "react";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Separator } from "@/components/ui/separator";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Drawer } from "@/components/ui/drawer";
import { useAISettings } from "../hooks/useAISettings";
import {
  Loader2,
  Shield,
  Database,
  Key,
  Settings as SettingsIcon,
  Info,
  MapPin,
  TrendingUp,
  Package,
  Users,
  BarChart3,
} from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { ScrollArea } from "@/components/ui/scroll-area";
import { PRIVACY_GROUPS } from "../data/domain-registry";
import type { AIDataPrivacySettings } from "../types";

// Common timezones for selection
const TIMEZONES = [
  { value: "Asia/Jakarta", label: "Asia/Jakarta (GMT+7)" },
  { value: "Asia/Singapore", label: "Asia/Singapore (GMT+8)" },
  { value: "Asia/Bangkok", label: "Asia/Bangkok (GMT+7)" },
  { value: "Asia/Manila", label: "Asia/Manila (GMT+8)" },
  { value: "Asia/Kuala_Lumpur", label: "Asia/Kuala Lumpur (GMT+8)" },
  { value: "UTC", label: "UTC (GMT+0)" },
  { value: "America/New_York", label: "America/New York (GMT-5)" },
  { value: "America/Los_Angeles", label: "America/Los Angeles (GMT-8)" },
  { value: "Europe/London", label: "Europe/London (GMT+0)" },
  { value: "Europe/Paris", label: "Europe/Paris (GMT+1)" },
  { value: "Asia/Tokyo", label: "Asia/Tokyo (GMT+9)" },
  { value: "Asia/Shanghai", label: "Asia/Shanghai (GMT+8)" },
];

/** Map icon strings from domain registry to Lucide components */
const DOMAIN_ICONS: Record<string, React.ElementType> = {
  MapPin,
  TrendingUp,
  Package,
  Users,
  BarChart3,
  Settings: SettingsIcon,
};

interface AISettingsDrawerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  readOnly?: boolean;
}

export function AISettingsDrawer({ open, onOpenChange, readOnly = false }: Readonly<AISettingsDrawerProps>) {
  const {
    settings,
    settingsResponse,
    isLoading,
    updateDataPrivacy,
    toggleEnabled,
    updateProvider,
    updateModel,
    updateAPIKey,
    updateTimezone,
    isUpdating,
  } = useAISettings();
  
  const [apiKey, setAPIKey] = useState("");
  const [showAPIKey, setShowAPIKey] = useState(false);

  const handleUpdateAPIKey = () => {
    if (apiKey.trim()) {
      updateAPIKey(apiKey);
      setAPIKey("");
      setShowAPIKey(false);
    }
  };

  /** Toggle all privacy keys within a domain group */
  const handleToggleDomainGroup = (keys: { key: keyof AIDataPrivacySettings }[], enabled: boolean) => {
    if (!settings) return;
    const updates: Partial<AIDataPrivacySettings> = {};
    for (const { key } of keys) {
      updates[key] = enabled;
    }
    // Batch update all keys in the group
    for (const [key, value] of Object.entries(updates)) {
      updateDataPrivacy(key, value);
    }
  };

  /** Check if all keys in a group are enabled */
  const isGroupEnabled = (keys: { key: keyof AIDataPrivacySettings }[]): boolean => {
    return keys.every((k) => settings.data_privacy[k.key]);
  };

  /** Check if some (but not all) keys in a group are enabled */
  const isGroupPartial = (keys: { key: keyof AIDataPrivacySettings }[]): boolean => {
    const enabledCount = keys.filter((k) => settings.data_privacy[k.key]).length;
    return enabledCount > 0 && enabledCount < keys.length;
  };

  return (
    <Drawer
      open={open}
      onOpenChange={onOpenChange}
      title="AI Settings"
      description="Configure AI assistant preferences and data privacy settings"
      side="right"
      defaultWidth={672}
      minWidth={400}
      maxWidth={800}
      resizable={true}
    >
      {isLoading ? (
        <div className="flex items-center justify-center p-8">
          <Loader2 className="h-6 w-6 animate-spin" />
        </div>
      ) : (
        <ScrollArea className="h-full px-6 py-4">
          <div className="space-y-6">
            {readOnly && (
              <Alert>
                <Info className="h-4 w-4" />
                <AlertDescription>
                  You are viewing these settings in read-only mode. Contact your administrator for edit permissions.
                </AlertDescription>
              </Alert>
            )}
            
            <Alert>
              <Shield className="h-4 w-4" />
              <AlertDescription>
                These settings control which data can be sent to the external AI service.
                Disable data types you don&apos;t want to share with the AI.
              </AlertDescription>
            </Alert>

            {/* AI Configuration */}
            <div className="space-y-4">
              <div className="flex items-center gap-2">
                <SettingsIcon className="h-5 w-5" />
                <h3 className="font-medium">AI Configuration</h3>
              </div>
              
              <div className="flex items-center justify-between">
                <div className="space-y-0.5">
                  <Label htmlFor="ai-enabled">Enable AI Assistant</Label>
                  <p className="text-sm text-muted-foreground">
                    Turn on AI assistant to get insights and chat support
                  </p>
                </div>
                <Switch
                  id="ai-enabled"
                  checked={settings.enabled}
                  onCheckedChange={toggleEnabled}
                  disabled={isUpdating || readOnly}
                />
              </div>

              <Separator />

              <div className="space-y-2">
                <Label htmlFor="ai-provider">AI Provider</Label>
                <Select
                  value={settings.provider}
                  onValueChange={updateProvider}
                  disabled={isUpdating || !settings.enabled || readOnly}
                >
                  <SelectTrigger id="ai-provider">
                    <SelectValue placeholder="Select provider" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="cerebras">Cerebras</SelectItem>
                    <SelectItem value="openai">OpenAI</SelectItem>
                    <SelectItem value="anthropic">Anthropic</SelectItem>
                  </SelectContent>
                </Select>
                <p className="text-xs text-muted-foreground">
                  Choose the AI service provider for processing requests
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="ai-model">AI Model</Label>
                <Select
                  value="gpt-oss-120b"
                    onValueChange={updateModel}
                    disabled={isUpdating || !settings.enabled || readOnly}
                  >
                    <SelectTrigger id="ai-model">
                      <SelectValue placeholder="Select model" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="gpt-oss-120b">GPT-OSS-120B</SelectItem>
                    </SelectContent>
                  </Select>
                  <p className="text-xs text-muted-foreground">
                    Select the AI model for generating responses
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="api-key">API Key</Label>
                  <div className="flex gap-2">
                    <Input
                      id="api-key"
                      type={showAPIKey ? "text" : "password"}
                      value={apiKey}
                      onChange={(e) => setAPIKey(e.target.value)}
                      placeholder="Enter your API key"
                      disabled={isUpdating || !settings.enabled || readOnly}
                    />
                    <Button
                      variant="outline"
                      onClick={() => setShowAPIKey(!showAPIKey)}
                      disabled={isUpdating || !settings.enabled || readOnly}
                    >
                      {showAPIKey ? "Hide" : "Show"}
                    </Button>
                    <Button
                      onClick={handleUpdateAPIKey}
                      disabled={isUpdating || !settings.enabled || !apiKey.trim() || readOnly}
                    >
                      <Key className="h-4 w-4 mr-2" />
                      Save
                    </Button>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    {settingsResponse?.api_key ? "API key is configured" : "No API key configured"}
                  </p>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="timezone">Timezone</Label>
                  <Select
                    value={settings.timezone}
                    onValueChange={updateTimezone}
                    disabled={isUpdating || !settings.enabled || readOnly}
                  >
                    <SelectTrigger id="timezone">
                      <SelectValue placeholder="Select timezone" />
                    </SelectTrigger>
                    <SelectContent>
                      {TIMEZONES.map((tz) => (
                        <SelectItem key={tz.value} value={tz.value}>
                          {tz.label}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                  <p className="text-xs text-muted-foreground">
                    Set your preferred timezone for AI responses
                  </p>
                </div>
              </div>

              <Separator />

              {/* Data Privacy Settings - Grouped by Domain */}
              <div className="space-y-4">
                <div className="flex items-center gap-2">
                  <Database className="h-5 w-5" />
                  <h3 className="font-medium">Data Privacy</h3>
                </div>
                <p className="text-sm text-muted-foreground">
                  Control which CRM modules the AI assistant can access. Settings are grouped by domain for modular privacy control.
                </p>

                <div className="space-y-6">
                  {PRIVACY_GROUPS.map((group) => {
                    const IconComponent = DOMAIN_ICONS[group.icon] || Database;
                    const allEnabled = isGroupEnabled(group.keys);
                    const partial = isGroupPartial(group.keys);

                    return (
                      <div key={group.domain} className="space-y-3">
                        {/* Domain group header with master toggle */}
                        <div className="flex items-center justify-between">
                          <div className="flex items-center gap-2">
                            <IconComponent className="h-4 w-4 text-muted-foreground" />
                            <h4 className="text-sm font-medium">{group.label}</h4>
                            {partial && (
                              <span className="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
                                Partial
                              </span>
                            )}
                          </div>
                          <Switch
                            checked={allEnabled}
                            onCheckedChange={(checked) =>
                              handleToggleDomainGroup(group.keys, checked)
                            }
                            disabled={isUpdating || !settings.enabled || readOnly}
                          />
                        </div>

                        {/* Individual privacy toggles */}
                        <div className="ml-6 space-y-3">
                          {group.keys.map((privacyKey) => (
                            <div
                              key={privacyKey.key}
                              className="flex items-center justify-between"
                            >
                              <div className="space-y-0.5">
                                <Label htmlFor={`share-${privacyKey.key}`}>
                                  {privacyKey.label}
                                </Label>
                                <p className="text-xs text-muted-foreground">
                                  {privacyKey.description}
                                </p>
                              </div>
                              <Switch
                                id={`share-${privacyKey.key}`}
                                checked={settings.data_privacy[privacyKey.key] ?? true}
                                onCheckedChange={(checked) =>
                                  updateDataPrivacy(privacyKey.key, checked)
                                }
                                disabled={isUpdating || !settings.enabled || readOnly}
                              />
                            </div>
                          ))}
                        </div>

                        <Separator className="mt-2" />
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          </ScrollArea>
        )}
      </Drawer>
    );
  }
