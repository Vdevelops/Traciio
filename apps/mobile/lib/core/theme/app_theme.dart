import 'package:flutter/material.dart';

class AppTheme {
  const AppTheme._();

  // ========== LIGHT THEME COLORS ==========
  // Kimia Farma Theme - Biru untuk background, Oranye #F39200 untuk primary
  // Converted from OKLCH to Flutter Color

  // Primary: oklch(0.73 0.19 55) ≈ #F39200 (Orange)
  static const Color primaryColor = Color(0xFFF39200);
  static const Color primaryColorLight = Color(0xFFFFA726); // Lighter orange
  static const Color primaryColorDark = Color(0xFFE67E00); // Darker orange

  // Background: oklch(0.98 0.01 240) ≈ #F8F8FC (Very light blue-ish)
  static const Color backgroundColor = Color(0xFFF8F8FC);

  // Card: oklch(1.0 0 0) = Pure white
  static const Color cardBackground = Color(0xFFFFFFFF);

  // Secondary: oklch(0.95 0.03 240) ≈ #F0F0F8 (Light blue-ish)
  static const Color secondaryBackground = Color(0xFFF0F0F8);

  // Muted: oklch(0.96 0.01 240) ≈ #F4F4F8 (Muted blue-ish)
  static const Color mutedBackground = Color(0xFFF4F4F8);

  // Accent: oklch(0.92 0.05 55) ≈ #FDE8CB (Light orange tint)
  static const Color accentColor = Color(0xFFFDE8CB);

  // Foreground: oklch(0.25 0.02 240) ≈ #3A3A45 (Dark blue-ish text)
  static const Color textPrimary = Color(0xFF3A3A45);

  // Muted Foreground: oklch(0.45 0.02 240) ≈ #717180 (Muted text)
  static const Color textSecondary = Color(0xFF717180);

  // Border: oklch(0.90 0.02 240) ≈ #E3E3EB (Light blue-ish border)
  static const Color borderColor = Color(0xFFE3E3EB);

  // Destructive: oklch(0.5386 0.1937 26.7249) ≈ #DC2626 (Red)
  static const Color destructiveColor = Color(0xFFDC2626);

  // ========== DARK THEME COLORS ==========
  // Kimia Farma Theme Dark Mode - Biru gelap untuk background, Oranye untuk primary

  // Background: oklch(0.20 0.02 240) ≈ #2A2A35 (Dark blue-ish)
  static const Color darkBackground = Color(0xFF2A2A35);

  // Card: oklch(0.25 0.02 240) ≈ #35353F (Slightly lighter dark blue)
  static const Color darkCard = Color(0xFF35353F);

  // Foreground: oklch(0.95 0.01 240) ≈ #F1F1F5 (Light text)
  static const Color darkForeground = Color(0xFFF1F1F5);

  // Primary: oklch(0.76 0.19 55) ≈ #FFA726 (Lighter orange for dark mode)
  static const Color darkPrimary = Color(0xFFFFA726);

  // Primary Foreground: oklch(0.15 0.05 240) ≈ #1F1F28 (Dark blue for text on primary)
  static const Color darkPrimaryForeground = Color(0xFF1F1F28);

  // Secondary: oklch(0.30 0.03 240) ≈ #404050 (Dark blue secondary)
  static const Color darkSecondary = Color(0xFF404050);

  // Muted: oklch(0.28 0.02 240) ≈ #3A3A48 (Muted dark blue)
  static const Color darkMuted = Color(0xFF3A3A48);

  // Muted Foreground: oklch(0.85 0.01 240) ≈ #D8D8E0 (Muted light text)
  static const Color darkMutedForeground = Color(0xFFD8D8E0);

  // Accent: oklch(0.72 0.16 55) ≈ #FFB84D (Light orange accent)
  static const Color darkAccent = Color(0xFFFFB84D);

  // Border: oklch(0.35 0.02 240) ≈ #4A4A58 (Dark border)
  static const Color darkBorder = Color(0xFF4A4A58);

  static ThemeData get light => ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.light(
      primary: primaryColor,
      onPrimary: Colors.white,
      secondary: secondaryBackground,
      onSecondary: textPrimary,
      surface: cardBackground,
      onSurface: textPrimary,
      error: destructiveColor,
      onError: Colors.white,
      outline: borderColor,
      surfaceContainerHighest: mutedBackground,
    ),
    scaffoldBackgroundColor: backgroundColor,
    cardColor: cardBackground,
    dividerColor: borderColor,
    hintColor: textSecondary,
    textTheme: const TextTheme(
      displayLarge: TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.bold,
        color: textPrimary,
        letterSpacing: -0.02,
      ),
      headlineMedium: TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.bold,
        color: textPrimary,
        letterSpacing: -0.02,
      ),
      titleLarge: TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: textPrimary,
      ),
      titleMedium: TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: textPrimary,
      ),
      bodyLarge: TextStyle(fontSize: 16, color: textPrimary),
      bodyMedium: TextStyle(fontSize: 14, color: textPrimary),
      bodySmall: TextStyle(fontSize: 12, color: textSecondary),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: cardBackground,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: borderColor),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: borderColor),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: primaryColor, width: 2),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: destructiveColor),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
    ),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        backgroundColor: primaryColor,
        foregroundColor: Colors.white,
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: primaryColor,
        foregroundColor: Colors.white,
        elevation: 2,
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: primaryColor,
        side: const BorderSide(color: borderColor, width: 1.5),
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
    ),
    appBarTheme: const AppBarTheme(
      backgroundColor: cardBackground,
      foregroundColor: textPrimary,
      elevation: 0,
      centerTitle: false,
      titleTextStyle: TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: textPrimary,
      ),
    ),
    cardTheme: CardThemeData(
      color: cardBackground,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
    ),
    floatingActionButtonTheme: const FloatingActionButtonThemeData(
      backgroundColor: primaryColor,
      foregroundColor: Colors.white,
    ),
  );

  static ThemeData get dark => ThemeData(
    useMaterial3: true,
    colorScheme: ColorScheme.dark(
      primary: darkPrimary,
      onPrimary: darkPrimaryForeground,
      secondary: darkSecondary,
      onSecondary: darkForeground,
      surface: darkCard,
      onSurface: darkForeground,
      error: const Color(0xFFEF4444),
      onError: Colors.white,
      outline: darkBorder,
      surfaceContainerHighest: darkMuted,
    ),
    scaffoldBackgroundColor: darkBackground,
    cardColor: darkCard,
    dividerColor: darkBorder,
    hintColor: darkMutedForeground,
    textTheme: TextTheme(
      displayLarge: const TextStyle(
        fontSize: 32,
        fontWeight: FontWeight.bold,
        color: darkForeground,
        letterSpacing: -0.02,
      ),
      headlineMedium: const TextStyle(
        fontSize: 24,
        fontWeight: FontWeight.bold,
        color: darkForeground,
        letterSpacing: -0.02,
      ),
      titleLarge: const TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: darkForeground,
      ),
      titleMedium: const TextStyle(
        fontSize: 16,
        fontWeight: FontWeight.w600,
        color: darkForeground,
      ),
      bodyLarge: const TextStyle(fontSize: 16, color: darkForeground),
      bodyMedium: const TextStyle(fontSize: 14, color: darkForeground),
      bodySmall: const TextStyle(fontSize: 12, color: darkMutedForeground),
    ),
    inputDecorationTheme: InputDecorationTheme(
      filled: true,
      fillColor: darkCard,
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: darkBorder),
      ),
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: darkBorder),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: darkPrimary, width: 2),
      ),
      errorBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(12),
        borderSide: const BorderSide(color: Color(0xFFEF4444)),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
    ),
    filledButtonTheme: FilledButtonThemeData(
      style: FilledButton.styleFrom(
        backgroundColor: darkPrimary,
        foregroundColor: darkPrimaryForeground,
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    ),
    elevatedButtonTheme: ElevatedButtonThemeData(
      style: ElevatedButton.styleFrom(
        backgroundColor: darkPrimary,
        foregroundColor: darkPrimaryForeground,
        elevation: 0,
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
    ),
    outlinedButtonTheme: OutlinedButtonThemeData(
      style: OutlinedButton.styleFrom(
        foregroundColor: darkPrimary,
        side: const BorderSide(color: darkBorder, width: 1.5),
        padding: const EdgeInsets.symmetric(vertical: 14, horizontal: 24),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      ),
    ),
    appBarTheme: AppBarTheme(
      backgroundColor: darkCard,
      foregroundColor: darkForeground,
      elevation: 0,
      centerTitle: false,
      titleTextStyle: const TextStyle(
        fontSize: 20,
        fontWeight: FontWeight.w600,
        color: darkForeground,
      ),
    ),
    cardTheme: CardThemeData(
      color: darkCard,
      elevation: 0,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
    ),
    floatingActionButtonTheme: FloatingActionButtonThemeData(
      backgroundColor: darkPrimary,
      foregroundColor: darkPrimaryForeground,
    ),
  );
}
