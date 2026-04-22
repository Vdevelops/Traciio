"use client";

import { useState } from "react";
import { Calendar as CalendarIcon, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

interface DatePickerProps {
  readonly date: Date | undefined;
  readonly onDateChange: (date: Date | undefined) => void;
  readonly disabled?: boolean;
  readonly placeholder?: string;
  readonly className?: string;
}

// Helper function to format date
const formatDate = (date: Date): string => {
  const day = date.getDate();
  const month = date.getMonth();
  const year = date.getFullYear();
  const monthNames = [
    "January", "February", "March", "April", "May", "June",
    "July", "August", "September", "October", "November", "December"
  ];
  return `${monthNames[month]} ${day}, ${year}`;
};

export function DatePicker({
  date,
  onDateChange,
  disabled = false,
  placeholder = "Pick a date",
  className,
}: DatePickerProps) {
  const [open, setOpen] = useState(false);

  const handleClear = () => {
    onDateChange(undefined);
    setOpen(false);
  };

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn(
            "w-full justify-start text-left font-normal",
            !date && "text-muted-foreground",
            className
          )}
          disabled={disabled}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {date ? formatDate(date) : <span>{placeholder}</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <div className="flex flex-col">
          <Calendar
            mode="single"
            selected={date}
            onSelect={(selectedDate: Date | undefined) => {
              onDateChange(selectedDate);
              setOpen(false);
            }}
          />
          {date && (
            <>
              <div className="border-t" />
              <Button
                variant="ghost"
                size="sm"
                className="w-full justify-start text-destructive hover:text-destructive hover:bg-destructive/10 cursor-pointer rounded-none"
                onClick={handleClear}
              >
                <X className="mr-2 h-4 w-4" />
                Clear
              </Button>
            </>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
