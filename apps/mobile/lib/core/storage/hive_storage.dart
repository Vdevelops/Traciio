import 'package:hive_flutter/hive_flutter.dart';

/// Hive storage initialization and management
/// 
/// This class handles:
/// - Initializing Hive for Flutter
/// - Opening boxes for different data types
/// - Managing box lifecycle
class HiveStorage {
  static bool _initialized = false;
  
  /// Initialize Hive storage
  /// 
  /// Call this in main.dart before runApp()
  static Future<void> init() async {
    if (_initialized) return;
    
    await Hive.initFlutter();
    
    // Note: Hive adapters will be registered here once models are updated
    // For now, we'll use JSON serialization with Hive
    
    _initialized = true;
  }
  
  /// Open a box for storing data
  /// 
  /// Boxes are automatically created if they don't exist
  static Future<Box<T>> openBox<T>(String name) async {
    if (!_initialized) {
      await init();
    }
    
    if (Hive.isBoxOpen(name)) {
      return Hive.box<T>(name);
    }
    
    return await Hive.openBox<T>(name);
  }
  
  /// Close a box
  static Future<void> closeBox(String name) async {
    if (Hive.isBoxOpen(name)) {
      await Hive.box(name).close();
    }
  }
  
  /// Clear all data from a box
  static Future<void> clearBox(String name) async {
    if (Hive.isBoxOpen(name)) {
      await Hive.box(name).clear();
    }
  }
  
  /// Delete a box (removes from disk)
  static Future<void> deleteBox(String name) async {
    if (Hive.isBoxOpen(name)) {
      await Hive.box(name).deleteFromDisk();
    }
  }
  
  /// Check if a box is open
  static bool isBoxOpen(String name) {
    return Hive.isBoxOpen(name);
  }
}


