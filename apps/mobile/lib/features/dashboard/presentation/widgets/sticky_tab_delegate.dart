import 'package:flutter/material.dart';

class StickyTabDelegate extends SliverPersistentHeaderDelegate {
  final Widget child;
  final double height; // Tinggi normal
  final Color backgroundColor;

  StickyTabDelegate({
    required this.child,
    this.height = 60,
    required this.backgroundColor,
  });

  @override
  double get minExtent => height;

  @override
  double get maxExtent => height;

  @override
  Widget build(
    BuildContext context,
    double shrinkOffset,
    bool overlapsContent,
  ) {
    return Container(
      color:
          backgroundColor, // Background color agar konten di bawah tidak tembus saat scroll
      alignment: Alignment.center,
      child: child,
    );
  }

  @override
  bool shouldRebuild(StickyTabDelegate oldDelegate) {
    return oldDelegate.child != child ||
        oldDelegate.backgroundColor != backgroundColor;
  }
}
