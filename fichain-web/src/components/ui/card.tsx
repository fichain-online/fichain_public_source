// components/ui/card.tsx
import * as React from "react";
import { cn } from "@/lib/utils"; // Assuming you have this utility for class name merging

// Main Card Container
export const Card = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn(
      "rounded-2xl border bg-card text-card-foreground shadow-md", // Updated: bg-white to bg-card for theming (can be bg-white if preferred), removed p-4, adjusted shadow
      className
    )}
    {...props}
  />
));
Card.displayName = "Card";

// Card Header
export const CardHeader = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("flex flex-col space-y-1.5 p-6", className)} // Added padding
    {...props}
  />
));
CardHeader.displayName = "CardHeader";

// Card Title
export const CardTitle = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLHeadingElement>
>(({ className, ...props }, ref) => (
  <h3
    ref={ref}
    className={cn(
      "text-2xl font-semibold leading-none tracking-tight", // Standard title styling
      className
    )}
    {...props}
  />
));
CardTitle.displayName = "CardTitle";

// Card Description
export const CardDescription = React.forwardRef<
  HTMLParagraphElement,
  React.HTMLAttributes<HTMLParagraphElement>
>(({ className, ...props }, ref) => (
  <p
    ref={ref}
    className={cn("text-sm text-muted-foreground", className)} // Standard description styling
    {...props}
  />
));
CardDescription.displayName = "CardDescription";

// Card Content
export const CardContent = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    // Updated: Added padding. The original text styling is kept but can be easily overridden.
    // pt-0 can be added via className if it directly follows a CardHeader for tighter spacing.
    className={cn("p-6 text-sm text-gray-700", className)}
    {...props}
  />
));
CardContent.displayName = "CardContent";

// Card Footer
export const CardFooter = React.forwardRef<
  HTMLDivElement,
  React.HTMLAttributes<HTMLDivElement>
>(({ className, ...props }, ref) => (
  <div
    ref={ref}
    className={cn("flex items-center p-6 pt-0", className)} // Added padding, pt-0 assuming content or header above often has bottom padding
    {...props}
  />
));
CardFooter.displayName = "CardFooter";
