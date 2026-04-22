"use client";

import { useEffect } from "react";
import { useSearchParams } from "next/navigation";
import { CheckCircle2, XCircle, Loader2 } from "lucide-react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

export default function GoogleCalendarCallbackPage() {
  const searchParams = useSearchParams();
  const code = searchParams.get("code");
  const state = searchParams.get("state");
  const error = searchParams.get("error");
  const success = searchParams.get("success");

  useEffect(() => {
    // Send success message to parent window if in popup
    // Handle both direct callback (with code) and redirect from backend (with success=true)
    if (window.opener && !error && (success === "true" || (code && state))) {
      window.opener.postMessage(
        {
          type: "GOOGLE_CALENDAR_CALLBACK_SUCCESS",
          code: code || "",
          state: state || "",
        },
        window.location.origin
      );
      // Close popup after short delay
      setTimeout(() => {
        if (!window.closed) {
          window.close();
        }
      }, 1000);
    }
  }, [code, state, error, success]);

  const handleClose = () => {
    if (window.opener) {
      window.close();
    } else {
      // If not in popup, redirect to profile
      window.location.href = "/profile";
    }
  };

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="flex justify-center mb-4">
              <XCircle className="h-16 w-16 text-destructive" />
            </div>
            <CardTitle>Connection Failed</CardTitle>
            <CardDescription>
              Failed to connect Google Calendar. Please try again.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground text-center">
              Error: {error}
            </p>
            <Button onClick={handleClose} className="w-full">
              Close
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Show success if success=true (from backend redirect) or if we have code and state (direct callback)
  if (success === "true" || (code && state)) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-background p-4">
        <Card className="w-full max-w-md">
          <CardHeader className="text-center">
            <div className="flex justify-center mb-4">
              <CheckCircle2 className="h-16 w-16 text-green-500" />
            </div>
            <CardTitle>Successfully Connected!</CardTitle>
            <CardDescription>
              Your Google Calendar has been connected successfully.
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground text-center">
              You can now sync your schedules to Google Calendar.
            </p>
            <Button onClick={handleClose} className="w-full">
              Close
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  // Show loading if we don't have success or error yet
  return (
    <div className="min-h-screen flex items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">
            <Loader2 className="h-16 w-16 text-primary animate-spin" />
          </div>
          <CardTitle>Processing...</CardTitle>
          <CardDescription>
            Please wait while we connect your Google Calendar.
          </CardDescription>
        </CardHeader>
      </Card>
    </div>
  );
}
