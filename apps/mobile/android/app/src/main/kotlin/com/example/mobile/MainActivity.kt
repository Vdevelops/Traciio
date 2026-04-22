package com.example.mobile

import android.util.Log
import io.flutter.embedding.android.FlutterActivity
import io.flutter.embedding.engine.FlutterEngine

class MainActivity : FlutterActivity() {

    override fun configureFlutterEngine(flutterEngine: FlutterEngine) {
        super.configureFlutterEngine(flutterEngine)
        Log.d("FlutterCrash", "Engine configured")
    }

    override fun onDestroy() {
        Log.d("FlutterCrash", "MainActivity onDestroy")
        super.onDestroy()
    }
}
